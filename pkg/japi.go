package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func GetApiResult(apiKey string, url string, method string, debugMode bool) map[string]interface{} {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      "1",
		"method":  method,
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
		if debugMode {
			fmt.Println(err.Error())
		}
		fmt.Println("Jeedom")
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

	return result
}