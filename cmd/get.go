package cmd

import (
	"fmt"
	"github.com/deimosfr/jeedom-status/pkg"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strconv"
	"strings"
)

const jeedomApiUrlSuffix = "/core/api/jeeApi.php"

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Jeedom global summary",
	Run: func(cmd *cobra.Command, args []string) {
		// Check args
		selectedStyle, _ := cmd.Flags().GetString("style")
		selectedBarType, _ := cmd.Flags().GetString("barType")
		debugMode, _ := cmd.Flags().GetBool("debug")
		apiKey, _ := cmd.Flags().GetString("apiKey")
		url, _ := cmd.Flags().GetString("url")
		alternateUrl, _ := cmd.Flags().GetString("alternateUrl")
		ignoreBatteryWarning, _ := cmd.Flags().GetBool("ignore-battery-warning")
		if !pkg.CheckArgContent(selectedStyle, getStyles()) || !pkg.CheckArgContent(selectedBarType, getBarsTypes()) {
			os.Exit(1)
		}

		// select reachable URL
		urlApi := pkg.CheckConnectivity(apiKey, url, alternateUrl, jeedomApiUrlSuffix, debugMode)
		if urlApi == "" {
			fmt.Println("Jeedom N/A")
			os.Exit(1)
		}
		jeedomVersion, err := pkg.GetVersion(apiKey, urlApi, debugMode)
		if err != nil {
			fmt.Println("Can't determine Jeedom version")
			os.Exit(1)
		}

		currentGlobalStatus := pkg.BuildGlobalStatus(apiKey, urlApi, jeedomVersion, selectedStyle, selectedBarType, debugMode)
		batteryStatus := pkg.GetJeedomBatteryInfo(apiKey, urlApi, ignoreBatteryWarning, debugMode)

		// Build lines
		mainLine := mainPrint(&currentGlobalStatus, &batteryStatus)
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
	getCmd.Flags().StringP("alternateUrl", "a", "", "Jeedom alternate API URL, like http://jeedom")
	getCmd.Flags().StringP("apiKey", "k", "", "Jeedom API key or User Hash Key (required)")
	err = getCmd.MarkFlagRequired("apiKey")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	getCmd.Flags().Float32P("jeedomVersion", "v", 4.0, "Specify the version of Jeedom")

	getCmd.Flags().StringP("barType", "b", "mac",
		fmt.Sprintf("Select the bar type: %s", strings.Join(getBarsTypes(), ", ")))

	getCmd.Flags().StringP("style", "s", "text",
		fmt.Sprintf("Choose output style: %s", strings.Join(getStyles(), ", ")))

	getCmd.Flags().BoolP("ignore-battery-warning", "w", false, "Ignore battery waning report")
	getCmd.Flags().BoolP("debug", "d", false, "Run in debug mode")
}

func getStyles() []string {
	return []string{"text", "jeedom", "nerd", "emoji"}
}

func getBarsTypes() []string {
	return []string{"mac", "i3blocks", "none"}
}

func mainPrint(jeedomCurrentInfo *pkg.JeedomCurrentStatus, batteryStatus *pkg.JeedomEquipmentsBatteryStatus) string {
	var lineToPrint string
	var icons pkg.JeedomSummary
	var iconsToPrint []string
	currentJeedomStatus := make(map[string]string)

	if jeedomCurrentInfo.DebugMode {
		fmt.Println(jeedomCurrentInfo.JeedomGlobalStatus)
	}

	icons = pkg.JeedomSummaryNoIcons()
	if jeedomCurrentInfo.Style == "nerd" {
		icons = pkg.JeedomSummaryNerdFontsIcons()
	} else if jeedomCurrentInfo.Style == "emoji" {
		icons = pkg.JeedomSummaryEmojiIcons()
	} else if jeedomCurrentInfo.Style == "jeedom" {
		icons = pkg.JeedomSummaryFontsIcons()
	}

	for key, value := range jeedomCurrentInfo.JeedomGlobalStatus {
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
		if jeedomCurrentInfo.DebugMode {
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
	lineToPrint += notificationsPrint(jeedomCurrentInfo)

	// Add battery info
	if batteryStatus.BatteryWarning > 0 {
		batteryYellowLine := batteryColorize(jeedomCurrentInfo.BarsType, "yellow", batteryStatus.BatteryWarning, icons.Battery)
		lineToPrint += batteryYellowLine
	}
	if batteryStatus.BatteryDanger > 0 {
		batteryRedLine := batteryColorize(jeedomCurrentInfo.BarsType, "red", batteryStatus.BatteryDanger, icons.Battery)
		lineToPrint += batteryRedLine
	}

	return strings.Trim(lineToPrint, " ")
}

func notificationsPrint(jeedomCurrentInfo *pkg.JeedomCurrentStatus) string {
	var result []string
	color := ""
	updateAndMessageCounts := [2]int{jeedomCurrentInfo.JeedomUpdates, jeedomCurrentInfo.JeedomMessages}

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
				content += notificationColorize(jeedomCurrentInfo.BarsType, color, number)
			}
			result = append(result, content)
		}
	}

	return strings.Join(result, " ")
}

func batteryColorize(barType string, color string, number int, icon string) string {
	content := ""
	var coloredContent string

	if barType == "i3blocks" {
		content = "<span color='" + color + "'><span font='Jeedom'>"
	}

	if barType == "mac" {
		if color == "yellow" {
			coloredContent = fmt.Sprintf("%s", Yellow(strconv.Itoa(number)))
		} else {
			coloredContent = fmt.Sprintf("%s", Red(strconv.Itoa(number)))
		}
		content += coloredContent
	} else {
		content += " " + strconv.Itoa(number)
	}
	content += icon

	if barType == "i3blocks" {
		content += "</span></span>"
	}

	return content
}

func notificationColorize(barType string, color string, number int) string {
	content := ""
	var coloredContent string
	icons := [20]string{
		"\u2460",
		"\u2461",
		"\u2462",
		"\u2463",
		"\u2464",
		"\u2465",
		"\u2466",
		"\u2467",
		"\u2468",
		"\u2469",
		"\u246A",
		"\u246B",
		"\u246C",
		"\u246D",
		"\u246E",
		"\u246F",
		"\u2470",
		"\u2471",
		"\u2472",
		"\u2473",
	}

	if barType == "i3blocks" {
		content = "<span color='" + color + "'><span font='FontAwesome'>"
	}

	if barType == "mac" {
		if color == "yellow" {
			coloredContent = Yellow(icons[number]).String()
		} else {
			coloredContent = Red(icons[number]).String()
		}
		content += coloredContent
	} else {
		content += icons[number]
	}

	if barType == "i3blocks" {
		content += "</span></span>"
	}

	return content
}

func additionalPrint(jeedomCurrentInfos *pkg.JeedomCurrentStatus, mainLine string) string {
	if jeedomCurrentInfos.BarsType == "mac" {
		return printMacBar(jeedomCurrentInfos)
	}

	if jeedomCurrentInfos.BarsType == "i3blocks" {
		return mainLine
	}

	return ""
}

func printMacBar(jeedomCurrentInfos *pkg.JeedomCurrentStatus) string {
	additionalInfo := "---\n"

	// Updates
	if jeedomCurrentInfos.JeedomUpdates > 0 {
		additionalInfo += fmt.Sprintf("Updates %d | color=red href=%s/index.php?v=d&p=update\n",
			jeedomCurrentInfos.JeedomUpdates,
			jeedomCurrentInfos.JeedomApiUrl)
	}
	// Messages
	if jeedomCurrentInfos.JeedomMessages > 0 {
		additionalInfo += fmt.Sprintf("Messages %d | color=orange href=%s\n",
			jeedomCurrentInfos.JeedomMessages,
			jeedomCurrentInfos.JeedomApiUrl)
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
