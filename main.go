package main

import (
	"io"
	"log"
	"os"
	"strconv"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"gopkg.in/yaml.v3"
)

// Some config and other useful global vars
var (
	Token          string
	VK             *api.VK
	WeatherToken   string
	Debug          bool
	CommandStrings CommandString
)

// Get config params from env variables
func init() {
	var err error

	Token = getEnv("VK_TOKEN", "fallbacktoken")
	WeatherToken = getEnv("WEATHER_TOKEN", "fallbacktoken")
	if Debug, err = strconv.ParseBool(getEnv("DEBUG", "true")); err != nil {
		log.Fatalf("Could not parse DEBUG env variable: %v", err)
	}

	file, err := os.Open("strings.yml")
	if err != nil {
		log.Fatalf("Falied to open command strings config file: %v\n", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Falied to read from command strings config file: %v\n", err)
	}

	if err := yaml.Unmarshal(data, &CommandStrings); err != nil {
		log.Fatalf("Falied to unmarhal yaml from command strings config file: %v\n", err)
	}
}

// Main cycle
func main() {
	// Init new VK api
	VK = api.NewVK(Token)

	// Get information about the group
	group, err := VK.GroupsGetByID(nil)
	if err != nil {
		log.Fatalf("Could not obtain groups info: %v\n", err)
	}

	// Initializing Long Poll
	lp, err := longpoll.NewLongPoll(VK, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// Registering event handlers
	lp.MessageNew(OnMessageNew) // New message
	lp.MessageEvent(OnMessageEvent)

	// Running Bots Long Poll
	log.Println("Starting longpoll...\nBot in now online!")
	lp.Run()
}
