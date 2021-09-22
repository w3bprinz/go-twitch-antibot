package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

type Bots struct {
	BotName [][]interface{} `json:"bots"`
	Total   int             `json:"_total"`
}

// Load cron and godotenv
var dotenverr = godotenv.Load()
var c = cron.New()

// Creates the client for Twitch
var client = twitch.NewClient(os.Getenv("BOT_NAME"), os.Getenv("BOT_OAUTH"))

func getOnlineBots() []string {
	//fmt.Println("Bin in getOnlineBots function")

	url := "https://api.twitchinsights.net/v1/bots/online"

	botClient := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "k4nzi-antibot")

	res, getErr := botClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	//fmt.Println(res)

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//var result map[string]interface{}

	//json.Unmarshal([]byte(body), &result)

	var bots Bots

	json.Unmarshal([]byte(body), &bots)

	//var botNames []string
	//botNames := fmt.Sprintf("%v", bots.BotName)
	botNames := make([]string, len(bots.BotName))
	for i, v := range bots.BotName {
		botNames[i] = fmt.Sprint(v[0])
	}
	return botNames
}

// Function to check the user
func checkUser(channel string, friendlyBotsList []string) {
	botNames := getOnlineBots()
	userlist, userlisterr := client.Userlist(channel)
	fmt.Println(userlist)
	if userlisterr != nil {
		fmt.Println("Error getting Userlist from channel.")
	}
	for _, user := range userlist {
		// Check if user is in friendly list
		userInFriendlyBots := contains(friendlyBotsList, user)
		// If user is friendly do:
		if userInFriendlyBots {
			fmt.Printf("User %s was found in friendly list.\n", user)
		} else {
			// User is unfriendly do:
			fmt.Printf("User %s was not found in friendly list. Checking %s in TwitchInsights API...\n", user, user)
			userInBotList := contains(botNames, user)
			if userInBotList {
				// User is a Bot in TwitchInsights API.
				fmt.Printf("User %s is a bot. Checked in TwitchInsights API. Banning User...\n", user)
				// Ban the User in Twitch channel.
				banUser(channel, user)
			} else {
				// User is regular User.
				fmt.Printf("User %s is not a bot. Checked %s in TwitchInsights API.\n", user, user)
				//fmt.Println(user)
			}

		}
	}
}

func contains(bots []string, user string) bool {
	for _, bot := range bots {
		if bot == user {
			return true
		}
	}
	return false
}

func banUser(channel string, user string) {
	fmt.Printf("Banning the User: %s in channel #%s because: 'Bots are not allowed here.'\n", user, channel)
	client.Ban(channel, user, "Bots are not allowed here! Get rekt!")

}

func main() {
	if dotenverr != nil {
		log.Fatal("Error loading .env file")
	}

	friendlyBots := strings.Split(os.Getenv("FRIENDLY_BOTS"), ",")
	//getOnlineBots()
	//getFriendlyBots()

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	// client := twitch.NewClient("exilit", "oauth:ihzd9nu22ftyzl89fx5kxvbc0f5vol")

	// Function to check if connected and send Message to Console and Channel.
	client.OnConnect(func() {
		fmt.Println("Connected...")
		// Join Channel
		client.Join(os.Getenv("TWITCH_CHANNEL"))
		//client.Say("", "Bin da, wer noch?")

	})

	// Send Chat messages to console.
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Time, message.Channel, message.User.Name, message.Message)

		if message.Message == "?userlist" && message.User.ID == os.Getenv("BOT_OWNER_ID") {
			fmt.Println("Getting Userlist.")
			//client.Userlist("exilit")
			checkUser(message.Channel, friendlyBots)
			//fmt.Println(userlist)
		}
		if strings.Split(message.Message, " ")[0] == "?ban" && message.User.ID == os.Getenv("BOT_OWNER_ID") {
			//fmt.Println(strings.Split(message.Message, " ")[1])
			if strings.Split(message.Message, " ")[1] == "" {
				client.Say(message.Channel, "No Username provided.")
			} else {
				//Ban the given User in Twitch Channel.
				banUser(message.Channel, strings.Split(message.Message, " ")[1])
			}
		}
		if message.Message == "?startabb" && message.User.ID == os.Getenv("BOT_OWNER_ID") {
			client.Say(message.Channel, "Automatic Bot banning activated.")
			c.AddFunc("0 10 * * * *", func() { checkUser(message.Channel, friendlyBots) })
			c.Start()
		}
		if message.Message == "?stopabb" && message.User.ID == os.Getenv("BOT_OWNER_ID") {
			client.Say(message.Channel, "Automatic Bot banning de-activated.")
			c.Stop()
		}

	})

	err := client.Connect()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
