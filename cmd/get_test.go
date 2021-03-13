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

func getTestGlobalStatus() pkg.JeedomCurrentStatus {
	return pkg.JeedomCurrentStatus{
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
}

func getTestAllBatteryNotification() pkg.JeedomEquipmentsBatteryStatus {
	return pkg.JeedomEquipmentsBatteryStatus{
		BatteryWarning: 1,
		BatteryDanger:  2,
	}
}

func TestTextBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "text"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "A 2D 4G 3H 5L 6M 7O 8P 2R 9S 1U 3W ① ② 1B 2B", mainLine)
}

func TestJeedomBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "jeedom"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "5\ue601 5\ue601 3\ue60a \ue60e 4\ue611 6\ue612 2\ue61d 7\ue61e 2\ue622 1\ue627 3\ue90f 8\uf0e7 ① ② 1\ue602 2\ue602", mainLine)
}
func TestNerdBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "nerd"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1S 3\ue373 \uf023 8\uf0e7 3\uf17a 7\uf1e6 2\uf2c7 4\uf834 5盛 6ﰌ 9ﲍ 2ﴘ ① ② 1\uf244 2\uf244", mainLine)
}
func TestEmojiBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "emoji"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1↕ 8⚡ 2🌡 6🏃 4💡 3💧 5🔆 7🔌 🔒 3🖼 9🚨 2🚪 12 1210🔋10🔋", mainLine)
}