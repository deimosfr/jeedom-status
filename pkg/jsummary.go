package pkg

type JeedomSummary struct {
	Alarm string
	Door string
	Humidity string
	Light string
	Luminosity string
	Motion string
	Outlet string
	Power string
	Security string
	Shutter string
	Temperature string
	Windows string
}

func JeedomSummaryFontsIcons() JeedomSummary {
	// icons: https://www.nerdfonts.com/cheat-sheet
	return JeedomSummary{
		Alarm:       "\uF023",
		Door:        "\uFD18",
		Humidity:    "\uE373",
		Light:       "\uF834",
		Luminosity:  "\uFAA7",
		Motion:      "\uFC0C",
		Outlet:      "\uF1E6",
		Power:       "\uF0E7",
		Security:    "\uFC8D",
		Shutter:     "S",
		Temperature: "\uF2C7",
		Windows:     "\uF17A",
	}
}

func JeedomSummaryEmojiIcons() JeedomSummary {
	// emoji: https://unicode.org/emoji/charts/full-emoji-list.html
	return JeedomSummary{
		Alarm:       "\U0001F512",
		Door:        "\U0001F6AA",
		Humidity:    "\U0001F4A7",
		Light:       "\U0001F4A1",
		Luminosity:  "\U0001F506",
		Motion:      "\U0001F3C3",
		Outlet:      "\U0001F50C",
		Power:       "\u26A1",
		Security:    "\U0001F6A8",
		Shutter:     "S",
		Temperature: "\U0001F321",
		Windows:     "\U0001FA9F",
	}
}