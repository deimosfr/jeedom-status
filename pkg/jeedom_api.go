package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func CheckConnectivity(apiKey string, url string, alternateUrl string, jeedomApiUrlSuffix string, debugMode bool) string {
	urlFull := url + jeedomApiUrlSuffix
	_, err := GetApiResult(apiKey, urlFull, "ping", debugMode)
	if err == nil {
		return urlFull
	}

	urlFull = alternateUrl + jeedomApiUrlSuffix
	_, err = GetApiResult(apiKey, urlFull, "ping", debugMode)
	if err == nil {
		return urlFull
	}

	return ""
}

func GetVersion(apiKey string, url string, debugMode bool) (string, error) {
	resp, err := GetApiResult(apiKey, url, "version", debugMode)
	if err != nil {
		return "", err
	}

	for key, value := range resp {
		if key == "result" {
			jeedomVersion := fmt.Sprintf("%v", value)
			return jeedomVersion, nil
		}
	}
	return "", errors.New("Wasn't able to get Jeedom version from the API")
}

func GetApiResult(apiKey string, url string, method string, debugMode bool) (map[string]interface{}, error) {
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
		return nil, err
	}

	resp, errPost := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if debugMode {
		fmt.Println(err.Error())
	}
	if errPost != nil {
		return nil, errPost
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error while trying to reach url %s: %s", url, resp.Status)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}