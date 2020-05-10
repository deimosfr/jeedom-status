package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Jeedom global summary",
	Run: func(cmd *cobra.Command, args []string) {
		result := make(map[string]string)

		// Check args
		selectedStyle, _ := cmd.Flags().GetString("style")
		_, found := pkg.Find(getStyles(), selectedStyle)
		if !found {
			fmt.Println(
				"Value %s is not a valid style, allowed values are: %s",
				selectedStyle,
				strings.Join(getStyles(), " "),
			)
			os.Exit(1)
		}

		if res, _ := cmd.Flags().GetBool("fake"); res {
			result = getSampleJeedomGlobalStatus()
		} else {
			apiKey, _ := cmd.Flags().GetString("apiKey")
			url, _ := cmd.Flags().GetString("url")
			result = getJeedomGlobalStatus(apiKey, url)
		}

		debugMode, _ := cmd.Flags().GetBool("debug")
		prettyPrint(result, selectedStyle, debugMode)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringP("url", "u", "", "Jeedom API URL (required)")
	getCmd.MarkFlagRequired("url")
	getCmd.Flags().StringP("apiKey", "k", "", "Jeedom API key or User Hash Key (required)")
	getCmd.MarkFlagRequired("apiKey")

	getCmd.Flags().StringP("style", "s", "fonts",
		fmt.Sprintf("Choose output style: %s", strings.Join(getStyles(), ", ")))

	getCmd.Flags().BoolP("fake", "f", false,"Run a sample test (won't connect to Jeedom API)")
	getCmd.Flags().BoolP("debug", "d", false,"Run in debug mode")
}

func getStyles() []string {
	return []string{"fonts", "emoji"}
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

func getSampleJeedomGlobalStatus() map[string]string {
	return map[string]string{
		"alarm": "1", // ok
		"door": "0",
		"humidity": "1", // ok
		"light": "1", // ok
		"luminosity": "4", // ok
		"motion": "0",
		"outlet": "2", // ok
		"power": "2", // ok
		"security": "0",
		"shutter": "0",
		"temperature": "1", // ok
		"windows": "3", // ok
	}
}

func prettyPrint(jeedomMap map[string]string, iconStyle string, debugMode bool) {
	var toPrint []string
	var icons pkg.JeedomSummary

	if debugMode {
		fmt.Println(jeedomMap)
	}

	if iconStyle == "fonts" {
		icons = pkg.JeedomSummaryFontsIcons()
	} else if iconStyle == "emoji" {
		icons = pkg.JeedomSummaryEmojiIcons()
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

	fmt.Println(lineToPrint)
	fmt.Println(lineToPrint)
}
