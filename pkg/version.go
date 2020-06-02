package pkg

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetCurrentVersion() string {
	return "0.7.0" // ci-version-check
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

func GetLatestVersion() (bool, string) {
	lastCheckVersionFile := "/tmp/jeedom-status"

	// If no file containing version exists
	_, err := os.Stat(lastCheckVersionFile)
	if os.IsNotExist(err) {
		newAvailableVersion, lastVersion := CheckAvailableNewVersion()
		StoreLastVersion(lastCheckVersionFile, lastVersion)
		return newAvailableVersion, lastVersion
	}

	// only check once every hour to speedup rendering
	_, minutes, _ := time.Now().Clock()
	if minutes == 0 {
		newAvailableVersion, lastVersion := CheckAvailableNewVersion()
		return newAvailableVersion, lastVersion
	}

	// Return saved content
	storedVersion := ReadLastCheckedVersion(lastCheckVersionFile)
	if GetCurrentVersion() != storedVersion {
		return true, storedVersion
	}
	return false, GetCurrentVersion()
}

func StoreLastVersion(lastCheckVersionFile string, version string) bool {
	f, err := os.Create(lastCheckVersionFile)
	if err != nil {
		return false
	}

	_, err = f.Write([]byte(version))
	if err != nil {
		err = f.Close()
		if err != nil {
			return false
		}
		return false
	}
	
	if f.Close() != nil {
		return false
	}
	return true
}

func ReadLastCheckedVersion(lastCheckVersionFile string) string {
	content, err := ioutil.ReadFile(lastCheckVersionFile)
	if err != nil {
		return ""
	}
	return string(content)
}
