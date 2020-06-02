package pkg

import "fmt"

func GetJeedomNetworkConfig(apiKey string, url string, debugMode bool) {
	toto := GetApiResult(apiKey, url, "config::byKey('internalAddr')", debugMode)
	fmt.Println(toto)
}