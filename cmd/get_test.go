package cmd

import (
	"github.com/deimosfr/jeedom-status/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

//func RandomNumberAsString() string {
//	return strconv.Itoa(rand.Intn(8)+1)
//}

//func getRandomJeedomGlobalStatus() map[string]string {
//	return map[string]string{
//		"alarm": "1",
//		"door": RandomNumberAsString(),
//		"humidity": RandomNumberAsString(),
//		"light": RandomNumberAsString(),
//		"luminosity": RandomNumberAsString(),
//		"motion": RandomNumberAsString(),
//		"outlet": RandomNumberAsString(),
//		"power": RandomNumberAsString(),
//		"security": RandomNumberAsString(),
//		"shutter": RandomNumberAsString(),
//		"temperature": RandomNumberAsString(),
//		"windows": RandomNumberAsString(),
//	}
//}

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
	assert.Equal(t, "A 2D 4G 3H 5L 6M 7O 8P 2R 9S 1U 3W ② ③ 1B 2B", mainLine)
}

func TestJeedomBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "jeedom"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "9 9 3  4 6 2 7 2 1 3 8 ② ③ 1 2", mainLine)
}
func TestNerdBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "nerd"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1S 3  8 3 7 2 4 5盛 6ﰌ 9ﲍ 2ﴘ ② ③ 1 2", mainLine)
}
func TestEmojiBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "emoji"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1↕ 8⚡ 2🌡 6🏃 4💡 3💧 5🔆 7🔌 🔒 3🖼 9🚨 2🚪 ② ③ 1🔋 2🔋", mainLine)
}

func TestEmojiMacBarOutput(t *testing.T) {
	currentGlobalStatus := getTestGlobalStatus()
	allBatteryNotification := getTestAllBatteryNotification()

	currentGlobalStatus.Style = "emoji"
	currentGlobalStatus.BarsType = "mac"
	mainLine := mainPrint(&currentGlobalStatus, &allBatteryNotification)
	assert.Equal(t, "1↕ 8⚡ 2🌡 6🏃 4💡 3💧 5🔆 7🔌 🔒 3🖼 9🚨 2🚪 [31m①[0m [33m②[0m[33m1[0m🔋[31m2[0m🔋", mainLine)
}