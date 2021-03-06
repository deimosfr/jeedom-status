package pkg

type JeedomSummary struct {
	Alarm       string
	Battery     string
	Door        string
	Humidity    string
	Light       string
	Luminosity  string
	Motion      string
	Outlet      string
	Power       string
	Security    string
	Shutter     string
	Temperature string
	Updates     string
	Windows     string
}

func JeedomSummaryNerdFontsIcons() JeedomSummary {
	// icons: https://www.nerdfonts.com/cheat-sheet
	return JeedomSummary{
		Alarm:       "\uF023",
		Battery:     "\uF244",
		Door:        "\uFD18", // nok
		Humidity:    "\uE373",
		Light:       "\uF834",
		Luminosity:  "\uFAA7",
		Motion:      "\uFC0C", // nok
		Outlet:      "\uF1E6",
		Power:       "\uF0E7",
		Security:    "\uFC8D", // nok
		Shutter:     "S",
		Temperature: "\uF2C7",
		Updates:     "\uF62E",
		Windows:     "\uF17A",
	}
}

func JeedomSummaryEmojiIcons() JeedomSummary {
	// emoji: https://unicode.org/emoji/charts/full-emoji-list.html
	return JeedomSummary{
		Alarm:       "\U0001F512",
		Battery:     "\U0001F50B",
		Door:        "\U0001F6AA",
		Humidity:    "\U0001F4A7",
		Light:       "\U0001F4A1",
		Luminosity:  "\U0001F506",
		Motion:      "\U0001F3C3",
		Outlet:      "\U0001F50C",
		Power:       "\u26A1",
		Security:    "\U0001F6A8",
		Shutter:     "\u2195",
		Temperature: "\U0001F321",
		Updates:     "\U0001F534",
		Windows:     "\U0001F5BC",
	}
}

func JeedomSummaryFontsIcons() JeedomSummary {
	// Load fonts with http://mathew-kurian.github.io/CharacterMap/
	return JeedomSummary{
		Alarm:       "\uE60E", //Jeedom font
		Battery: 	 "\uE602", //Jeedom font
		Door:        "\uE61D", //Jeedom font
		Humidity:    "\uE90F", //Jeedomapp font
		Light:       "\uE611", //Jeedom font
		Luminosity:  "\uE601", //Nature font
		Motion:      "\uE612", //Jeedom font
		Outlet:      "\uE61E", //Jeedom font
		Power:       "\uF0E7", //General font / fonts awesome
		Security:    "\uE601", //Jeedom font
		Shutter:     "\uE627", //Jeedom font
		Temperature: "\uE622", //Jeedom font
		Updates:     "\uE91D", // Jeedomapp font
		Windows:     "\uE60A", //Jeedom font
	}
}

func JeedomSummaryNoIcons() JeedomSummary {
	return JeedomSummary{
		Alarm:       "A",
		Battery: 	 "B",
		Door:        "D",
		Humidity:    "H",
		Light:       "G",
		Luminosity:  "L",
		Motion:      "M",
		Outlet:      "O",
		Power:       "P",
		Security:    "S",
		Shutter:     "U",
		Temperature: "R",
		Updates:     "U",
		Windows:     "W",
	}
}