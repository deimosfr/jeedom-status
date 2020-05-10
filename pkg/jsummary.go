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
		Motion:      "\uE373",
		Outlet:      "\uF834",
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
		Alarm:       "\u1F514",
		Door:        "\u1F6AA",
		Humidity:    "\u1F4A7",
		Light:       "\u1F4A1",
		Luminosity:  "\u1F506",
		Motion:      "\u1F3C3",
		Outlet:      "\u1F50C",
		Power:       "\u26A1",
		Security:    "\u1F512",
		Shutter:     "S",
		Temperature: "\u1F321",
		Windows:     "\u1FA9F",
	}
}