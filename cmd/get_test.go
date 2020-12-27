package cmd

import (
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

func RandomNumberAsString() string {
	return strconv.Itoa(rand.Intn(8)+1)
}

func getRandomJeedomGlobalStatus() map[string]string {
	return map[string]string{
		"alarm": "1",
		"door": RandomNumberAsString(),
		"humidity": RandomNumberAsString(),
		"light": RandomNumberAsString(),
		"luminosity": RandomNumberAsString(),
		"motion": RandomNumberAsString(),
		"outlet": RandomNumberAsString(),
		"power": RandomNumberAsString(),
		"security": RandomNumberAsString(),
		"shutter": RandomNumberAsString(),
		"temperature": RandomNumberAsString(),
		"windows": RandomNumberAsString(),
	}
}

func getDefinedJeedomGlobalStatus() map[string]string {
	return map[string]string{
		"alarm": "1",
		"door": "2",
		"humidity": "3",
		"light": "4",
		"luminosity": "5",
		"motion": "6",
		"outlet": "7",
		"power": "8",
		"security": "9",
		"shutter": "1",
		"temperature": "2",
		"windows": "3",
	}
}

func TestBarOutput(t *testing.T) {
	currentGlobalStatus := pkg.JeedomCurrentStatus{
		JeedomApiUrl:         "",
		JeedomApiKey:         "",
		JeedomVersion: 		  "4.0",
		JeedomGlobalStatus:   getDefinedJeedomGlobalStatus(),
		JeedomUpdates:        1,
		JeedomMessages:       2,
		BarsType:             "none",
		Style:                "text",
		DebugMode:            false,
	}

	allBatteryNotification := pkg.JeedomEquipmentsBatteryStatus{
		BatteryWarning: 1,
		BatteryDanger:  2,
	}

	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "A 2D 4G 3H 5L 6M 7O 8P 2R 9S 1U 3W ① ② 1B 2B", mainLine)

	currentGlobalStatus.BarsType = "mac"
	currentGlobalStatus.Style = "emoji"
	mainLine = mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1↕ 8⚡ 2🌡 6🏃 4💡 3💧 5🔆 7🔌 🔒 3🖼 9🚨 2🚪 12 1210🔋10🔋", mainLine)
}