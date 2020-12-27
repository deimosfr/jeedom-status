package pkg

import (
	"fmt"
	"os"
)

func GetJeedomGlobalStatus40(apiKey string, url string, debugMode bool) map[string]string {
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
			for name, number := range result["result"].(map[string]interface{}) {
				stringMap[name] = fmt.Sprintf("%v", number)
			}
		}
	}

	return stringMap
}

