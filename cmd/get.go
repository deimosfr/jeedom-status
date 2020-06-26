package cmd

import (
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type JeedomCurrentStatus struct {
	JeedomApiUrl         string
	JeedomUrl            string
	JeedomAlternativeUrl string
	JeedomApiKey         string
	JeedomGlobalStatus   map[string]string
	JeedomUpdates        int
	JeedomMessages       int
	BarsType             string
	Style                string
	DebugMode            bool
}

const jeedomApiUrlSuffix = "/core/api/jeeApi.php"

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
				JeedomApiUrl:         "",
				JeedomUrl:            "",
				JeedomAlternativeUrl: "",
				JeedomApiKey:         "",
				JeedomGlobalStatus:   pkg.GetSampleJeedomGlobalStatus(),
				JeedomUpdates:        1,
				JeedomMessages:       2,
				BarsType:             selectedBarType,
				Style:                selectedStyle,
				DebugMode:            debugMode,
			}
		} else {
			apiKey, _ := cmd.Flags().GetString("apiKey")
			url, _ := cmd.Flags().GetString("url")
			alternateUrl, _ := cmd.Flags().GetString("alternateUrl")
			urlApi := pkg.CheckConnectivity(apiKey, url, alternateUrl, jeedomApiUrlSuffix, debugMode)
			if urlApi == "" {
				fmt.Println("Jeedom N/A")
				os.Exit(1)
			}

			currentGlobalStatus = JeedomCurrentStatus{
				JeedomApiUrl:         urlApi,
				JeedomUrl:            url,
				JeedomAlternativeUrl: alternateUrl,
				JeedomApiKey:         apiKey,
				JeedomGlobalStatus:   getJeedomGlobalStatus(apiKey, urlApi, debugMode),
				JeedomUpdates:        getJeedomUpdates(apiKey, urlApi, debugMode),
				JeedomMessages:       getJeedomMessage(apiKey, urlApi, debugMode),
				BarsType:             selectedBarType,
				Style:                selectedStyle,
				DebugMode:            debugMode,
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
	getCmd.Flags().StringP("alternateUrl", "a", "", "Jeedom lternate API URL, like http://jeedom")
	getCmd.Flags().StringP("apiKey", "k", "", "Jeedom API key or User Hash Key (required)")
	err = getCmd.MarkFlagRequired("apiKey")
	if err != nil {
		println(err)
		os.Exit(1)
	}

	getCmd.Flags().StringP("barType", "b", "mac",
		fmt.Sprintf("Select the bar type: %s", strings.Join(getBarsTypes(), ", ")))

	getCmd.Flags().StringP("style", "s", "text",
		fmt.Sprintf("Choose output style: %s", strings.Join(getStyles(), ", ")))

	getCmd.Flags().BoolP("fake", "f", false, "Run a sample test (won't connect to Jeedom API)")
	getCmd.Flags().BoolP("debug", "d", false, "Run in debug mode")
}

func getStyles() []string {
	return []string{"text", "jeedom", "nerd", "emoji"}
}

func getBarsTypes() []string {
	return []string{"mac", "i3blocks", "none"}
}

func getJeedomGlobalStatus(apiKey string, url string, debugMode bool) map[string]string {
	stringMap := make(map[string]string)
	result, _ := pkg.GetApiResult(apiKey, url, "summary::global", debugMode)

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
	result, _ := pkg.GetApiResult(apiKey, url, "update::all", debugMode)

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
	result, _ := pkg.GetApiResult(apiKey, url, "message::all", debugMode)

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
	var lineToPrint string
	var icons pkg.JeedomSummary
	var iconsToPrint []string
	currentJeedomStatus := make(map[string]string)

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
			currentJeedomStatus[icons.Alarm] = ""
			iconsToPrint = append(iconsToPrint, icons.Alarm)
			continue
		} else if value == "<nil>" || value == "0" {
			continue
		} else if key == "security" {
			currentJeedomStatus[icons.Security] = value
			iconsToPrint = append(iconsToPrint, icons.Security)
		} else if key == "motion" {
			currentJeedomStatus[icons.Motion] = value
			iconsToPrint = append(iconsToPrint, icons.Motion)
		} else if key == "windows" {
			currentJeedomStatus[icons.Windows] = value
			iconsToPrint = append(iconsToPrint, icons.Windows)
		} else if key == "outlet" {
			currentJeedomStatus[icons.Outlet] = value
			iconsToPrint = append(iconsToPrint, icons.Outlet)
		} else if key == "humidity" {
			currentJeedomStatus[icons.Humidity] = value
			iconsToPrint = append(iconsToPrint, icons.Humidity)
		} else if key == "light" {
			currentJeedomStatus[icons.Light] = value
			iconsToPrint = append(iconsToPrint, icons.Light)
		} else if key == "luminosity" {
			currentJeedomStatus[icons.Luminosity] = value
			iconsToPrint = append(iconsToPrint, icons.Luminosity)
		} else if key == "power" {
			currentJeedomStatus[icons.Power] = value
			iconsToPrint = append(iconsToPrint, icons.Power)
		} else if key == "door" {
			currentJeedomStatus[icons.Door] = value
			iconsToPrint = append(iconsToPrint, icons.Door)
		} else if key == "temperature" {
			currentJeedomStatus[icons.Temperature] = value
			iconsToPrint = append(iconsToPrint, icons.Temperature)
		} else if key == "shutter" {
			currentJeedomStatus[icons.Shutter] = value
			iconsToPrint = append(iconsToPrint, icons.Shutter)
		} else {
			currentJeedomStatus[key] = value
			iconsToPrint = append(iconsToPrint, key)
		}
	}

	// Print Jeedom if there is nothing else to print
	if iconsToPrint == nil {
		if jeedomCurrentInfos.DebugMode {
			fmt.Println("Nothing to print, Global status is empty")
		}
		return "Jeedom"
	}

	sort.Strings(iconsToPrint)
	for currentIcon := range iconsToPrint {
		for icon, value := range currentJeedomStatus {
			if iconsToPrint[currentIcon] == icon {
				lineToPrint += value + icon + " "
				continue
			}
		}
	}

	// Add notifications
	lineToPrint += notificationsPrint(jeedomCurrentInfos)

	return strings.Trim(lineToPrint, " ")
}

func notificationsPrint(jeedomCurrentInfos *JeedomCurrentStatus) string {
	var result []string
	color := ""
	updateAndMessageCounts := [2]int{jeedomCurrentInfos.JeedomUpdates, jeedomCurrentInfos.JeedomMessages}

	for kind, number := range updateAndMessageCounts {
		if number > 0 {
			content := ""

			color = "red"
			if kind == 1 {
				color = "orange"
			}

			if number > 20 {
				content += "+"
			}
			if number <= 20 {
				content += notificationColorize(jeedomCurrentInfos.BarsType, color, number)
			}
			result = append(result, content)
		}
	}

	return strings.Join(result, " ")
}

func notificationColorize(barType string, color string, number int) string {
	content := ""
	var coloredContent int
	icons := map[int]string{
		1:  "\u2460",
		2:  "\u2461",
		3:  "\u2462",
		4:  "\u2463",
		5:  "\u2464",
		6:  "\u2465",
		7:  "\u2466",
		8:  "\u2467",
		9:  "\u2468",
		10: "\u2469",
		11: "\u246A",
		12: "\u246B",
		13: "\u246C",
		14: "\u246D",
		15: "\u246E",
		16: "\u246F",
		17: "\u2470",
		18: "\u2471",
		19: "\u2472",
		20: "\u2473",
	}

	if barType == "i3blocks" {
		content = "<span color='" + color + "'><span font='FontAwesome'>"
	}

	if barType == "mac" {
		if color == "yellow" {
			coloredContent, _ = fmt.Printf("%s", Yellow(icons[number]))
		} else {
			coloredContent, _ = fmt.Printf("%s", Red(icons[number]))
		}
		content += strconv.Itoa(coloredContent)
	} else {
		content += icons[number]
	}

	if barType == "i3blocks" {
		content += "</span></span>"
	}

	return content
}

func additionalPrint(jeedomCurrentInfos *JeedomCurrentStatus, mainLine string) string {
	if jeedomCurrentInfos.BarsType == "mac" {
		return printMacBar(jeedomCurrentInfos)
	}

	if jeedomCurrentInfos.BarsType == "i3blocks" {
		return mainLine
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
		additionalInfo += fmt.Sprintf("Messages %d | color=orange href=%s\n",
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
