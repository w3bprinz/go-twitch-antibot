# go-twitch-antibot

## What is go-twitch-antibot?
---
**go-twitch-antibot** is a bot written in Go (golang) that fights against bots. Using the **TwitchInsights API** (https://twitchinsights.net), the bot checks for viewers on the Twitch channel and bans them if they are flagged as a bot in the TwitchInsights API.

The bot was developed using existing Go modules and uses the following modules:

- github.com/gempir/go-twitch-irc
- github.com/joho/godotenv
- github.com/robfig/cron

## How to use the bot?
---
### **Edit the .env File**

BOT_NAME={YOUR_BOT_NAME}\
BOT_OAUTH=oauth:{YOUR_OAUTH_TOKEN}\
BOT_OWNER_ID={YOUR_TWITCH_ID}\
TWITCH_CHANNEL={YOUR_TWITCH_CHANNEL}\
FRIENDLY_BOTS={YOUR_TWITCH_BOTS}

How to get a Twitch OAUTH Token?\
https://twitchapps.com/tmi/\
https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/


How to Convert Twitch Username to User ID?\
https://www.streamweasels.com/support/convert-twitch-username-to-user-id/

or just via the TwitchAPI:\
https://dev.twitch.tv/docs/api/reference#get-users
## Bot Commands
---
- ?userlist - Will get the current Userlist from the Twitch channel and will check the users if they are bots.
- ?startabb - Will start the automatic bot banning mode. (Runs every 10 minutes)
- ?stopabb - Will stop the automatic bot banning mode.

