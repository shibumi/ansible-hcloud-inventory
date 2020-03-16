# ansible-hcloud-inventory
A hcloud inventory for Ansible written in Go

## How to use it

### First preparations

You need to do the following things to use this project:

* configure a valid reverse DNS entry for every server
* create a configuration file or pass the API token via environment variable

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

* hcloud-ansible-inv doesn't let you specify labels
* hcloud-ansible-inv doesn't have an additional host parameter
* hcloud-ansible-inv is not grouping your servers via location

### https://github.com/thetechnick/hcloud-ansible

* thetechnick/hcloud-ansible is no longer maintained.

### https://github.com/hg8496/ansible-hcloud-inventory

* hg8496/ansible-hcloud-inventory is written in Python
