# Jeedom Status ![Master](https://github.com/deimosfr/jeedom-status/workflows/Push/badge.svg?branch=master)

Jeedom is a third party tool for [Jeedom](https://jeedom.com/) (Home assistant).

It shows the Jeedom global status in the status bars in the global status bar of the operating systems. Here is an example of what can be seen:

![all_output](assets/output_all.png)

You can download the binary directly from the [release page](https://github.com/deimosfr/jeedom-status/releases). It's available for Mac, Windows and Linux.

# Prerequisites

To use jeedom-status, you need to have :
* Your user hash key. Go into Tools -> Preferences -> Security -> User Hash.
* The URL of your jeedom API like (replace "jeedom" with the name or IP of Jeedom endpoint): http://jeedom/core/api/jeeApi.php
* Specific fonts containing icons: https://github.com/ryanoasis/nerd-fonts.

# Installation and usage

## Mac OS X

You need to install [brew](https://brew.sh/). If you don't have this tool, install it this way:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

Once done, you have to install the [Nerd fonts](https://github.com/ryanoasis/nerd-fonts) to get the Fonts containing the icons:
```bash
brew tap homebrew/cask-fonts
brew cask install font-hack-nerd-font
```

Install [bitbar](https://getbitbar.com/) to be able to add jeedom-status in the status bar:
```bash
brew cask install bitbar
```

We can now install jeedom-status:
```bash
brew tap deimosfr/jeedom-status
brew install jeedom-status
```

Finally, the final steps are:
* Open the bitbar application and define a folder to store plugins.
* Download the file "[jeedom-status.1m.sh](https://github.com/deimosfr/jeedom-status/integration/bitbar/jeedom-status.1m.sh)" and add it to the bitbar plugins folder.
* Edit this "jeedom-status.1m.sh" file with a text editor and update "APIKEY" and "APIKEY" with your information:

```bash
APIKEY="YOUR API OR UER HASH KEY HERE"
APIKEY="YOUR JEEDOM URL HERE"
```

## Linux - i3 and i3blocks

![i3_desktop](assets/i3_desktop.png)

![i3_output](assets/i3_output.png)

Here is an example with [i3blocks](https://github.com/vivien/i3blocks) for [i3wm](https://i3wm.org/). Add this in your i3blocks.conf:

```ini
[jeedom]
command=~/.config/i3/i3blocks_bin/jeedom_status
markup=pango
interval=60
```

