package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

type configuration struct {
	Token string `json:"token"`
}

func printHelp() {
	fmt.Println(`Usage:
    ansible-hcloud-inventory [option]

    --list    shows all groups including variables as JSON structure
    --host [hostname]    shows all variables for one host as JSON structure`)
	os.Exit(1)
}

func main() {
	// Try to read token from environment variable HETZNER_CLOUD_KEY otherwise use a config file
	token := os.Getenv("HETZNER_CLOUD_KEY")
	if token == "" {
		configFilePath, err := xdg.ConfigFile("ansible-hcloud-inventory/config.json")
		if err != nil {
			log.Fatalln(err)
		}
		file, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			log.Fatalln(err)
		}
		config := configuration{}
		if err = json.Unmarshal(file, &config); err != nil {
			log.Fatalln(err)
		}
		token = config.Token
	}

	// TODO: implement list and host functions
	client := hcloud.NewClient(hcloud.WithToken(token))
	servers, _ := client.Server.All(context.Background())
	for _, server := range servers {
		fmt.Println(server.Name)
	}
	// TODO: better flag handling
	args := os.Args
	if len(args) < 2 {
		printHelp()
	}
	switch args[1] {
	case "--list":
		fmt.Println("list")
		// TODO: get list
	case "--host":
		if len(args) != 3 {
			printHelp()
		}
		// TODO: get host variables
	default:
		printHelp()
	}
}
