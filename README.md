# ansible-hcloud-inventory
A hcloud inventory for Ansible written in Go

## Features

* Pass the API token via environment variable
* Pass the API token via config file
* Pass the API token via subcommand call (pass, gopass, etc)
* Generate an ini file out of your inventory
* The generated ini file will have additional labels

## Groups
ansible-hcloud-inventory is grouping your hosts in the following groups:

| Group | Description |
| ---  | ---       |
| nbg1 | All servers running in location NBG1 |
| fsn1 | All servers running in location FSN1 |
| hel1 | All servers running in location HEL1 |

## How to use it

### Installation and usage

You need to do the following things to use this project:

* configure a valid reverse DNS entry for every server
* create a configuration file or pass the API token via environment variable
* Install the ansible-hcloud-inventory via `go install https://github.com/shibumi/ansible-hcloud-inventory`
* Make sure to set your go path accordingly: `export GOBIN=$GOPATH/bin`
* Create a config file or pass the API token via environment variable. For examples scroll down.

### Dynamic Ansible inventory

For using this project as your dynamic ansible inventory do the
following:

```sh
$ ansible -i ansible-hcloud-inventory all -m ping
```

### Parameters

`ansible-hcloud-inventory` supports two parameters right now:

* `--list` for listing all hosts including variables
* `--host <hostname>` for just listing all labels for a given hostname
* `--ini` for generating an INI inventory file and printing it to stdout

### Config file examples

You can specify the API token via multiple ways:

* Environment variable
* Plain value in the config file
* Invoking a command (useful for password stores)

If you want to use the environment variable, just do:

```sh
$ HETZNER_CLOUD_KEY=yourkey ansible-hcloud-inventory --list
```

For using a plain value do the following:
(Note: It's important to leave the command field empty in this case)

```json
{
	"command": "",
	"token": "yourkey"
}
```

For using your favourite command for retrieving the password automatically:

```json
{
	"command": "your command here",
	"token": ""
}
```

## Differences to existing solutions

### https://github.com/thannaske/hcloud-ansible-inv

* doesn't let you specify labels
* doesn't have an additional host parameter
* is not grouping your servers via location
* can't generate an ini file for you

### https://github.com/thetechnick/hcloud-ansible

* is no longer maintained.
* can't generate an ini file for you
* doesn't have an additional host parameter

### https://github.com/hg8496/ansible-hcloud-inventory

* is written in Python
* can't generate an ini file for you
* doesn't have an additional host parameter
* doesn't support all locations (according to their README)