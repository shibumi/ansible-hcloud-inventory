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

// configuration just holds a Token at the moment.
// TODO: Gain token via any command. For example: `gopass api/hcloud`
type configuration struct {
	Token string `json:"token"`
}

// inventory is holding the format as specified in:
// https://docs.ansible.com/ansible/latest/dev_guide/developing_inventory.html#id16
type inventory struct {
	Meta struct {
		HostVars map[string]map[string]interface{} `json:"hostvars"`
	} `json:"_meta"`
	All struct {
		Children []string `json:"children"`
	} `json:"all"`
	Ungrouped struct {
		Hosts []string `json:"hosts"`
	}
}

// printHelp just prints a formatted help.
func printHelp() {
	fmt.Println(`Usage:
    ansible-hcloud-inventory [option]

    --list    shows all groups including variables as JSON structure
    --host [hostname]    shows all variables for one host as JSON structure`)
	os.Exit(1)
}

// list lists all servers including metadata
func (inv *inventory) list(token string) {
	client := hcloud.NewClient(hcloud.WithToken(token))
	servers, _ := client.Server.All(context.Background())
	for _, server := range servers {
		hostName := server.PublicNet.IPv4.DNSPtr
		inv.Meta.HostVars = make(map[string]map[string]interface{})
		inv.Meta.HostVars[hostName] = make(map[string]interface{})
		for k, v := range server.Labels {
			inv.Meta.HostVars[hostName][k] = v
		}
		inv.All.Children = append(inv.All.Children, "ungrouped")
		inv.Ungrouped.Hosts = append(inv.Ungrouped.Hosts, hostName)
	}
	output, err := json.MarshalIndent(inv, "", "    ")
	if err != nil {
		log.Println("Couldn't marshal inventory")
	}
	fmt.Println(string(output))
}

// host prints all labels for a given hostName (This has to be a RDNS pointer)
func host(token string, hostName string) {
	client := hcloud.NewClient(hcloud.WithToken(token))
	servers, _ := client.Server.All(context.Background())
	for _, server := range servers {
		if hostName == server.PublicNet.IPv4.DNSPtr {
			output, err := json.MarshalIndent(server.Labels, "", "    ")
			if err != nil {
				log.Fatal("Couldn't marshal label list")
			}
			fmt.Println(string(output))
		}
	}
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

	// TODO: better flag handling
	args := os.Args
	if len(args) < 2 {
		printHelp()
	}
	switch args[1] {
	case "--list":
		inv := inventory{}
		inv.list(token)
	case "--host":
		if len(args) != 3 {
			printHelp()
		}
		host(token, args[2])
	default:
		printHelp()
	}
}
