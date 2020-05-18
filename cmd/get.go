package cmd

import (
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"runtime"
	"strings"
)

type JeedomCurrentStatus struct {
	JeedomUrl			string
	JeedomApiKey		string
	JeedomGlobalStatus	map[string]string
	JeedomUpdates		int
	JeedomMessages		int
	BarsType			string
	Style				string
	DebugMode			bool
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Jeedom global summary",
	Run: func(cmd *cobra.Command, args []string) {
		var currentGlobalStatus JeedomCurrentStatus

		// Check args
		selectedStyle, _ := cmd.Flags().GetString("style")
		selectedBarType, _ := cmd.Flags().GetString("barType")
		debugMode, _ := cmd.Flags().GetBool("debug")
		if !pkg.CheckArgContent(selectedStyle, getStyles()) || !pkg.CheckArgContent(selectedBarType, getBarsTypes()) {
			os.Exit(1)
		}

		if res, _ := cmd.Flags().GetBool("fake"); res {
			currentGlobalStatus = JeedomCurrentStatus{
				JeedomUrl:          "",
				JeedomApiKey:       "",
				JeedomGlobalStatus: pkg.GetSampleJeedomGlobalStatus(),
				JeedomUpdates:      1,
				JeedomMessages:     2,
				BarsType:           selectedBarType,
				Style:              selectedStyle,
				DebugMode:          debugMode,
			}
		} else {
			apiKey, _ := cmd.Flags().GetString("apiKey")
			url, _ := cmd.Flags().GetString("url")
			url = url + "/core/api/jeeApi.php"
			currentGlobalStatus = JeedomCurrentStatus{
				JeedomUrl:          url,
				JeedomApiKey:       apiKey,
				JeedomGlobalStatus: getJeedomGlobalStatus(apiKey, url, debugMode),
				JeedomUpdates:      getJeedomUpdates(apiKey, url, debugMode),
				JeedomMessages:     getJeedomMessage(apiKey, url, debugMode),
				BarsType:           selectedBarType,
				Style:              selectedStyle,
				DebugMode:          debugMode,
			}
		}

		// Build lines
		mainLine := mainPrint(&currentGlobalStatus)
		additionalLines := additionalPrint(&currentGlobalStatus, mainLine)

		fmt.Println(mainLine)
		fmt.Println(additionalLines)
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
		fmt.Sprintf("Select the bar type: %s", strings.Join(getBarsTypes(), ", ")))

	getCmd.Flags().StringP("style", "s", "text",
		fmt.Sprintf("Choose output style: %s", strings.Join(getStyles(), ", ")))

	getCmd.Flags().BoolP("fake", "f", false,"Run a sample test (won't connect to Jeedom API)")
	getCmd.Flags().BoolP("debug", "d", false,"Run in debug mode")
}

func getStyles() []string {
	return []string{"text", "jeedom", "nerd", "emoji"}
}

func getBarsTypes() []string {
	return []string{"autodetect", "mac", "i3blocks", "none"}
}

func getJeedomGlobalStatus(apiKey string, url string, debugMode bool) map[string]string {
	stringMap := make(map[string]string)
	result := pkg.GetApiResult(apiKey, url, "summary::global", debugMode)

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

func getJeedomUpdates(apiKey string, url string, debugMode bool) int {
	totalUpdates := 0
	result := pkg.GetApiResult(apiKey, url, "update::all", debugMode)

	for key, value := range result {
		if key == "error" {
			for message, content := range value.(map[string]interface{}) {
				if message == "message" {
					fmt.Printf("Error: %s", content)
					os.Exit(1)
				}
			}
		} else if key == "result" {
			pluginList := reflect.ValueOf(value)
			for i := 0; i < pluginList.Len(); i++ {
				for name, content := range pluginList.Index(i).Interface().(map[string]interface{}) {
					if name == "status" && content != "ok" {
						totalUpdates++
						break
					}
				}
			}
		}
	}

	return totalUpdates
}

func getJeedomMessage(apiKey string, url string, debugMode bool) int {
	result := pkg.GetApiResult(apiKey, url, "message::all", debugMode)

	for key, value := range result {
		if key == "error" {
			for message, content := range value.(map[string]interface{}) {
				if message == "message" {
					fmt.Printf("Error: %s", content)
					os.Exit(1)
				}
			}
		} else if key == "result" {
			return reflect.ValueOf(value).Len()
		}
	}

	return 0
}

func mainPrint(jeedomCurrentInfos *JeedomCurrentStatus) string {
	var toPrint []string
	var icons pkg.JeedomSummary

	if jeedomCurrentInfos.DebugMode {
		fmt.Println(jeedomCurrentInfos.JeedomGlobalStatus)
	}

	icons = pkg.JeedomSummaryNoIcons()
	if jeedomCurrentInfos.Style == "nerd" {
		icons = pkg.JeedomSummaryNerdFontsIcons()
	} else if jeedomCurrentInfos.Style == "emoji" {
		icons = pkg.JeedomSummaryEmojiIcons()
	} else if jeedomCurrentInfos.Style == "jeedom" {
		icons = pkg.JeedomSummaryFontsIcons()
	}

	for key, value := range jeedomCurrentInfos.JeedomGlobalStatus {
		if key == "alarm" && value != "0" {
			toPrint = append(toPrint, icons.Alarm)
			continue
		} else if value == "<nil>" || value == "0" {
			continue
		} else if key == "security" {
			toPrint = append(toPrint, value+icons.Security)
		} else if key == "motion" {
			toPrint = append(toPrint, value+icons.Motion)
		} else if key == "windows" {
			toPrint = append(toPrint, value+icons.Windows)
		} else if key == "outlet" {
			toPrint = append(toPrint, value+icons.Outlet)
		} else if key == "humidity" {
			toPrint = append(toPrint, value+icons.Humidity)
		} else if key == "light" {
			toPrint = append(toPrint, value+icons.Light)
		} else if key == "luminosity" {
			toPrint = append(toPrint, value+icons.Luminosity)
		} else if key == "power" {
			toPrint = append(toPrint, value+icons.Power)
		} else if key == "door" {
			toPrint = append(toPrint, value+icons.Door)
		} else if key == "temperature" {
			toPrint = append(toPrint, value+icons.Temperature)
		} else if key == "shutter" {
			toPrint = append(toPrint, value+icons.Shutter)
		} else {
			toPrint = append(toPrint, value+key)
		}
	}

	// Print Jeedom if there is nothing else to print
	if toPrint == nil {
		if jeedomCurrentInfos.DebugMode {
			fmt.Println("Nothing to print, Global status is empty")
		}
		return "Jeedom"
	}

	lineGlobalStatus := strings.Join(toPrint, " ")

	return lineGlobalStatus
}

func additionalPrint(jeedomCurrentInfos *JeedomCurrentStatus, mainLine string) string {
	if jeedomCurrentInfos.BarsType == "mac" {
		return printMacBar(jeedomCurrentInfos)
	}

	if jeedomCurrentInfos.BarsType == "i3blocks" {
		return mainLine
	}

	if jeedomCurrentInfos.BarsType == "autodetect" {
		if runtime.GOOS == "darwin" {
			return printMacBar(jeedomCurrentInfos)
		}
	}
	return ""
}

func printMacBar(jeedomCurrentInfos *JeedomCurrentStatus) string {
	additionalInfo := "---\n"

	// Updates
	if jeedomCurrentInfos.JeedomUpdates > 0 {
		additionalInfo += fmt.Sprintf("Updates %d | color=red href=%s/index.php?v=d&p=update\n",
			jeedomCurrentInfos.JeedomUpdates,
			jeedomCurrentInfos.JeedomUrl)
	}
	// Messages
	if jeedomCurrentInfos.JeedomMessages > 0 {
		additionalInfo += fmt.Sprintf("Messages %d | color=yellow href=%s\n",
			jeedomCurrentInfos.JeedomMessages,
			jeedomCurrentInfos.JeedomUrl)
	}

	//Version / upgrade needed
	newAvailableVersion, version := pkg.GetLatestVersion()
	if newAvailableVersion {
		additionalInfo += fmt.Sprintf("New available version %s | bash=/usr/local/bin/brew param1=upgrade param2=jeedom-status terminal=false refresh=true\n", version)
	} else {
		additionalInfo += fmt.Sprintf("Current version %s\n", version)
	}
	return additionalInfo
}