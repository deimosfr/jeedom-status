package pkg

import (
	"fmt"
	"os"
)

func GetJeedomGlobalStatus41(apiKey string, url string, debugMode bool) map[string]string {
	stringMap := make(map[string]string)
	result, _ := GetApiResult(apiKey, url, "summary::global", debugMode)

	for key, value := range result {
		if key == "error" {
			for message, content := range value.(map[string]interface{}) {
				if message == "message" {
					fmt.Printf("Error: %s", content)
					os.Exit(1)
				}
			}
		} else if key == "result" {
			for name, content := range result["result"].(map[string]interface{}) {
				resultMap := content.(map[string]interface{})
				stringMap[name] = fmt.Sprintf("%v", resultMap["value"])
			}
		}
	}

	return stringMap
}