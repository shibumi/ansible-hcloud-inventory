package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/adrg/xdg"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

// configuration just holds a Token at the moment.
type configuration struct {
	Command string `json:"command"`
	Token   string `json:"token"`
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
	NBG1 struct {
		Hosts []string `json:"hosts"`
	} `json:"nbg1"`
	HEL1 struct {
		Hosts []string `json:"hosts"`
	} `json:"hel1"`
	FSN1 struct {
		Hosts []string `json:"hosts"`
	} `json:"fsn1"`
	Ungrouped struct {
		Hosts []string `json:"hosts"`
	} `json:"ungrouped"`
}

// printHelp just prints a formatted help.
func printHelp() {
	fmt.Println(`Usage:
    ansible-hcloud-inventory [option]

    --list    shows all groups including variables as JSON structure
    --host [hostname]    shows all variables for one host as JSON structure`)
	os.Exit(1)
}

// newInventory generates a new Hcloud Ansible inventory
func newInventory(token string) *inventory {
	inv := inventory{}
	client := hcloud.NewClient(hcloud.WithToken(token))
	servers, _ := client.Server.All(context.Background())
	// initialize All group
	inv.All.Children = []string{"nbg1", "hel1", "fsn1", "ungrouped"}
	inv.Meta.HostVars = make(map[string]map[string]interface{})
	inv.NBG1.Hosts = []string{}
	inv.FSN1.Hosts = []string{}
	inv.HEL1.Hosts = []string{}
	inv.Ungrouped.Hosts = []string{}
	for _, server := range servers {
		hostName := server.PublicNet.IPv4.DNSPtr
		inv.Meta.HostVars[hostName] = make(map[string]interface{})
		for k, v := range server.Labels {
			inv.Meta.HostVars[hostName][k] = v
		}
		switch server.Datacenter.Location.Name {
		case "nbg1":
			inv.NBG1.Hosts = append(inv.NBG1.Hosts, hostName)
		case "hel1":
			inv.HEL1.Hosts = append(inv.HEL1.Hosts, hostName)
		case "fsn1":
			inv.FSN1.Hosts = append(inv.FSN1.Hosts, hostName)
		default:
			inv.Ungrouped.Hosts = append(inv.Ungrouped.Hosts, hostName)
		}
	}
	return &inv
}

// list lists all servers including metadata
func (inv *inventory) list() {
	output, err := json.MarshalIndent(inv, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(output))
}

// host prints all labels for a given hostName (This has to be a RDNS pointer)
func (inv *inventory) host(hostName string) {
	for k, v := range inv.Meta.HostVars {
		if hostName == k {
			output, err := json.MarshalIndent(v, "", "    ")
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(output))
		}
	}
}

// generateIni generates an INI Ansible inventory and prints it to stdout
func (inv *inventory) generateIni() {
	t := template.New("inventory")
	t, err := t.Parse(`# This file has been generated via github.com/shibumi/ansible-hcloud-inventory
{{- $hostVars := .Meta.HostVars -}}
{{ range $server := .Ungrouped.Hosts -}}
{{$server}} {{ $index := index $hostVars $server}}{{ range $key, $value := $index}} {{$key}}={{$value}}{{end}}
{{end}}
[nbg1]
{{ range $server := .NBG1.Hosts -}}
{{$server}} {{ $index := index $hostVars $server}}{{ range $key, $value := $index}} {{$key}}={{$value}}{{end}}
{{end}}
[hel1]
{{ range $server := .HEL1.Hosts -}}
{{$server}} {{ $index := index $hostVars $server}}{{ range $key, $value := $index}} {{$key}}={{$value}}{{end}}
{{end}}
[fsn1]
{{ range $server := .FSN1.Hosts -}}
{{$server}} {{ $index := index $hostVars $server}}{{ range $key, $value := $index}} {{$key}}={{$value}}{{end}}
{{end}}
`)
	if err != nil {
		log.Fatalln(err)
	}
	if err = t.Execute(os.Stdout, inv); err != nil {
		log.Fatalln(err)
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
		if config.Command != "" {
			command := strings.Fields(config.Command)
			if out, err := exec.Command(command[0], command[1:]...).Output(); err != nil {
				log.Fatalln(err)
			} else {
				token = string(out)
			}
		} else {
			token = config.Token
		}
	}

	args := os.Args
	if len(args) < 2 {
		printHelp()
	}
	inv := newInventory(token)
	switch args[1] {
	case "--list":
		inv.list()
	case "--ini":
		inv.generateIni()
	case "--host":
		if len(args) != 3 {
			printHelp()
		}
		inv.host(args[2])
	default:
		printHelp()
	}
}
