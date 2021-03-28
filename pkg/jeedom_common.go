package pkg

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"os"
	"reflect"
	"strconv"
)

type JeedomCurrentStatus struct {
	JeedomUrl			 string
	JeedomApiUrl         string
	JeedomVersion		 string
	JeedomApiKey         string
	JeedomGlobalStatus   map[string]string
	JeedomUpdates        int
	JeedomMessages       int
	BarsType             string
	Style                string
	DebugMode            bool
}

type JeedomEquipmentsBatteryStatus struct {
	BatteryWarning		int
	BatteryDanger		int
}

func GetJeedomUpdates(apiKey string, url string, debugMode bool) int {
	totalUpdates := 0
	result, _ := GetApiResult(apiKey, url, "update::all", debugMode)

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

func GetJeedomMessage(apiKey string, url string, debugMode bool) int {
	result, _ := GetApiResult(apiKey, url, "message::all", debugMode)

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

func GetJeedomBatteryInfo(apiKey string, url string, ignoreBatteryWarning bool, debugMode bool) JeedomEquipmentsBatteryStatus {
	result, _ := GetApiResult(apiKey, url, "eqLogic::all", debugMode)
	allBatteryNotification := JeedomEquipmentsBatteryStatus{
		BatteryWarning: 0,
		BatteryDanger:  0,
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
			pluginList := reflect.ValueOf(value)
			for i := 0; i < pluginList.Len(); i++ {
				// get current equipment info
				equipment := pluginList.Index(i).Interface().(map[string]interface{})
				if reflect.ValueOf(equipment["status"]).Len() == 0 {
					continue
				}

				// ensure status is not empty
				status := equipment["status"].(map[string]interface{})
				if len(status) == 0 {
					continue
				}

				// ensure all required fields are present
				batteryWarning, exists := status["batterywarning"]
				if !exists {
					continue
				}
				batteryWarningInt, err := strconv.Atoi(fmt.Sprintf("%v", batteryWarning))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				batteryDanger, exists := status["batterydanger"]
				if !exists {
					continue
				}
				batteryDangerInt, err := strconv.Atoi(fmt.Sprintf("%v", batteryDanger))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				// if not danger is returned and battery is lower or equal than 5, return an alert anyway
				batteryLevel, exists := status["battery"]
				if !exists {
					continue
				}
				currentBatteryLevel, err := strconv.Atoi(fmt.Sprintf("%v", batteryLevel))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				if batteryDangerInt == 0 && currentBatteryLevel <= 5 {
					batteryDangerInt = 1
				}

				if !ignoreBatteryWarning {
					allBatteryNotification.BatteryWarning += batteryWarningInt
				}
				allBatteryNotification.BatteryDanger += batteryDangerInt
			}
		}
	}

	return allBatteryNotification
}

func GetJeedomGlobalStatus(apiKey string, url string, jeedomVersion string, debugMode bool) map[string]string {
	jeedomVersion41, _ := version.NewVersion("4.1")
	jeedomVersionSemVer, _ := version.NewVersion(jeedomVersion)

	if jeedomVersionSemVer.LessThan(jeedomVersion41) {
		return GetJeedomGlobalStatus40(apiKey, url, debugMode)
	}
	return GetJeedomGlobalStatus41(apiKey, url, debugMode)
}

func BuildGlobalStatus(apiKey string, urlApi string, jeedomVersion string, selectedStyle string, selectedBarType string, debugMode bool) JeedomCurrentStatus {
	return JeedomCurrentStatus{
		JeedomApiUrl:         urlApi,
		JeedomApiKey:         apiKey,
		JeedomVersion:        jeedomVersion,
		JeedomGlobalStatus:   GetJeedomGlobalStatus(apiKey, urlApi, jeedomVersion, debugMode),
		JeedomUpdates:        GetJeedomUpdates(apiKey, urlApi, debugMode),
		JeedomMessages:       GetJeedomMessage(apiKey, urlApi, debugMode),
		BarsType:             selectedBarType,
		Style:                selectedStyle,
		DebugMode:            debugMode,
	}
}