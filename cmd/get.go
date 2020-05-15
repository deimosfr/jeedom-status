package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"runtime"
	"strings"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Jeedom global summary",
	Run: func(cmd *cobra.Command, args []string) {
		var result map[string]string

		// Check args
		selectedStyle, _ := cmd.Flags().GetString("style")
		if !pkg.CheckArgContent(selectedStyle, getStyles()) {
			os.Exit(1)
		}

		selectedBarType, _ := cmd.Flags().GetString("barType")
		if !pkg.CheckArgContent(selectedBarType, getBarsTypes()) {
			os.Exit(1)
		}

		if res, _ := cmd.Flags().GetBool("fake"); res {
			result = pkg.GetSampleJeedomGlobalStatus()
		} else {
			apiKey, _ := cmd.Flags().GetString("apiKey")
			url, _ := cmd.Flags().GetString("url")
			url = url + "/core/api/jeeApi.php"
			result = getJeedomGlobalStatus(apiKey, url)
		}

		debugMode, _ := cmd.Flags().GetBool("debug")

		prettyPrint(result, selectedStyle, selectedBarType, debugMode)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("url", "u", "", "Jeedom API URL, like http://jeedom (required)")
	err := getCmd.MarkFlagRequired("url")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	getCmd.Flags().StringP("apiKey", "k", "", "Jeedom API key or User Hash Key (required)")
	err = getCmd.MarkFlagRequired("apiKey")
	if err != nil {
		println(err)
		os.Exit(1)
	}

	getCmd.Flags().StringP("barType", "b", "autodetect",
		fmt.Sprintf("Select the bar type: %s", strings.Join(getStyles(), ", ")))

	getCmd.Flags().StringP("style", "s", "text",
		fmt.Sprintf("Choose output style: %s", strings.Join(getBarsTypes(), ", ")))

	getCmd.Flags().BoolP("fake", "f", false,"Run a sample test (won't connect to Jeedom API)")
	getCmd.Flags().BoolP("debug", "d", false,"Run in debug mode")
}

func getStyles() []string {
	return []string{"text", "jeedom", "nerd", "emoji"}
}

func getBarsTypes() []string {
	return []string{"autodetect", "mac", "i3blocks", "none"}
}

func getJeedomGlobalStatus(apiKey string, url string) map[string]string {
	stringMap := make(map[string]string)
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  "summary::global",
		"params": map[string]string{
			"apikey": apiKey,
			"id":     "1",
		},
	}

	bytesRepresentation, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Error while trying to reach url %s: %s", url, resp.Status)
		os.Exit(1)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for key, value := range result {
		if key == "error" {
			for message, content := range value.(map[string]interface{}) {
				if message == "message" {
					fmt.Printf("Error: %s", content)
					os.Exit(1)
				}
			}
		} else if key == "result" {
			for name, number := range result["result"].(map[string]interface{}) {
				stringMap[name] = fmt.Sprintf("%v", number)
			}
		}
	}

	return stringMap
}

func prettyPrint(jeedomMap map[string]string, iconStyle string, barType string, debugMode bool) {
	var toPrint []string
	var icons pkg.JeedomSummary

	if debugMode {
		fmt.Println(jeedomMap)
	}

	icons = pkg.JeedomSummaryNoIcons()
	if iconStyle == "nerd" {
		icons = pkg.JeedomSummaryNerdFontsIcons()
	} else if iconStyle == "emoji" {
		icons = pkg.JeedomSummaryEmojiIcons()
	} else if iconStyle == "jeedom" {
		icons = pkg.JeedomSummaryFontsIcons()
	}

	for key, value := range jeedomMap {
		if key == "alarm" && value != "0" {
			toPrint = append(toPrint, icons.Alarm)
			continue
		} else if value == "<nil>" || value == "0" {
			continue
		} else if key == "security" {
			toPrint = append(toPrint, value + icons.Security)
		} else if key == "motion" {
			toPrint = append(toPrint, value + icons.Motion)
		} else if key == "windows" {
			toPrint = append(toPrint, value + icons.Windows)
		} else if key == "outlet" {
			toPrint = append(toPrint, value + icons.Outlet)
		} else if key == "humidity" {
			toPrint = append(toPrint, value + icons.Humidity)
		} else if key == "light" {
			toPrint = append(toPrint, value + icons.Light)
		} else if key == "luminosity" {
			toPrint = append(toPrint, value + icons.Luminosity)
		} else if key == "power" {
			toPrint = append(toPrint, value + icons.Power)
		} else if key == "door" {
			toPrint = append(toPrint, value + icons.Door)
		} else if key == "temperature" {
			toPrint = append(toPrint, value + icons.Temperature)
		} else if key == "shutter" {
			toPrint = append(toPrint, value + icons.Shutter)
		} else {
			toPrint = append(toPrint, value + key)
		}
	}

	lineToPrint := strings.Join(toPrint, " ")
	fmt.Println(formatBarType(barType, lineToPrint))
}

func formatBarType(barType string, lineToPrint string) string {
	if barType == "mac" {
		return printMacBar(lineToPrint)
	}

	if barType == "i3blocks" {
		return lineToPrint + "\n" + lineToPrint
	}

	if barType == "autodetect" {
		if runtime.GOOS == "Darwin" {
			return printMacBar(lineToPrint)
		}
	}
	return lineToPrint
}

func printMacBar(lineToPrint string) string {
	result := lineToPrint

	newVersion, version := pkg.CheckAvailableNewVersion()
	if version == "" {
		return "Can't check latest version"
	}

	if newVersion {
		messageVersion := fmt.Sprintf("New version available v%s ", version)
		result = result + "\n" + messageVersion
	}

	return result
}