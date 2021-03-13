#!/usr/bin/env sh
#
# Jeedom status, using jeedom-status (https://github.com/deimosfr/jeedom-status)
#
# <bitbar.title>Jeedom status</bitbar.title>
# <bitbar.version>v1.0</bitbar.version>
# <bitbar.author>Pierre Mavro (deimosfr)</bitbar.author>
# <bitbar.author.github>deimosfr</bitbar.author.github>
# <bitbar.desc>Jeedom global status for operating systems status bars</bitbar.desc>
# <bitbar.dependencies>jeedom-status</bitbar.dependencies>
# <bitbar.image>https://github.com/deimosfr/jeedom-status/blob/master/assets/mac_output.png?raw=true</bitbar.image>
#
# Dependencies:
#   Jeedom-status (https://github.com/deimosfr/jeedom-status)

export PATH=/opt/homebrew/bin:$PATH
APIKEY=
JEEDOM_URL=
STYLE="emoji"
jeedom-status get --apiKey $APIKEY --url $JEEDOM_URL -s $STYLE -b mac
