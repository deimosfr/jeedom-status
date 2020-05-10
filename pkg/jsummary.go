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