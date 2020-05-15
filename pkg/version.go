package pkg

import (
	"errors"
	"net/http"
	"strings"
)

func GetCurrentVersion() string {
	return "0.6.0" // ci-version-check
}

func GetLatestOnlineVersionUrl() (string, error) {
	url := "https://github.com/deimosfr/jeedom-status/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("Can't reach GitHub website, please check your network connectivity")
	}
	return resp.Request.URL.Path, nil
}

func GetLatestOnlineVersionNumber() (string, error) {
	urlPath, err := GetLatestOnlineVersionUrl()
	if err != nil {
		return "", err
	}
	splitUrl := strings.Split(urlPath, "/v")
	return splitUrl[len(splitUrl)-1], nil
}

func CheckAvailableNewVersion() (bool, string) {
	latestOnlineVersion, err := GetLatestOnlineVersionNumber()
	if err != nil {
		return false, ""
	}
	if GetCurrentVersion() < latestOnlineVersion {
		return true, latestOnlineVersion
	}
	return false, latestOnlineVersion
}
