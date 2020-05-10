# Jeedom Status ![Master](https://github.com/deimosfr/jeedom-status/workflows/Push/badge.svg?branch=master)

Jeedom is a third party tool for [Jeedom](https://jeedom.com/) (Home assistant).

It shows the Jeedom global status in the status bars in the global status bar of the operating systems. Here is an example of what can be seen:

![all_output](assets/output_all.png)

You can download the binary directly from the [release page](https://github.com/deimosfr/jeedom-status/releases). It's available for Mac, Windows and Linux.

# Prerequisites

To use jeedom-status, you have to get :
* Your user hash key. Go into Jeedom web interface, then click on Tools -> Preferences -> Security -> User Hash.
* The URL of your jeedom API like (replace "jeedom" with the name or IP of Jeedom endpoint): http://jeedom/core/api/jeeApi.php

# Installation and usage

## Mac OS X

You need have [brew](https://brew.sh/) installed. If you don't have brew, install it this way:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

Install [bitbar](https://getbitbar.com/) to be able to add jeedom-status in the status bar:
```bash
brew cask install bitbar
```

We can now install jeedom-status and download the Bitbar jeedom-status plugin:
```bash
brew tap deimosfr/jeedom-status
brew install jeedom-status
curl -o ~/Downloads/jeedom-status.1m.sh https://raw.githubusercontent.com/deimosfr/jeedom-status/master/integration/bitbar/jeedom-status.1m.sh
chmod 755 ~/Downloads/jeedom-status.1m.sh
```

Finally, the last steps are:
* Open the bitbar application and define a folder to store plugins.
* Move the "jeedom-status.1m.sh" plugin file from your Downloads folder to the bitbar plugins folder you've just defined.
* Edit this "jeedom-status.1m.sh" plugin file with a text editor and update "APIKEY" and "APIKEY" with your information (as described in the [prerequisites section](#Prerequisites)):

```bash
APIKEY="YOUR API OR UER HASH KEY HERE"
APIKEY="YOUR JEEDOM URL HERE"
```

You're done, click on the Bitbar and "refresh all". You'll see your Jeedom global status appearing.

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

For this example, I used Nerd fonts, containing additional icons: https://github.com/ryanoasis/nerd-fonts.