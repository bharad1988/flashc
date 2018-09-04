# flashc
### under development...
## A management suite to orchestrate lxc containers
  A highly scalable container management system. There will be a central controller node to which all the nodes participating in the cluster will register to. Each participating node will be running an agent which will provide the interface to connect to the controller.
  The controller provides with REST apis as a part of SDK to manage and monitor. 
As of now, the cluster can be managed by command line interface written on top REST apis

The container engine relies on btrfs to provide with snapshots.
## Further documentation will be added later



## Command line interface 

##### flashc main command set
```
flashc
	This is a command line REST client to communicate with controller
	This command is used to manage the flashc cluster with cli tools
	Use flashc -h for more detailed help

Usage:
  flashc [command]

Available Commands:
  agent       A brief description of your command
  container   A brief description of your command
  image       A brief description of your command

Flags:
  -C, --controllernode string   Provide the IP/DNS name of Controller Node (default "127.0.0.1")
  -h, --help                    help for flashc

Use "flashc [command] --help" for more information about a command.
```

##### flashc agent command set
```
flashc agent -h
Gets the status from Agents. Status of the agents are maintained in database.
	On running the command, the information is fetched from DB

Usage:
  flashc agent [flags]
  flashc agent [command]

Available Commands:
  register    Register an Agent node with the controller
  status      A brief description of your command
  unregister  UnRegister an Agent node with the controller
  update      update agent info

Global Flags:
  -C, --controllernode string   Provide the IP/DNS name of Controller Node (default "127.0.0.1")

Use "flashc agent [command] --help" for more information about a command.

flashc  agent status
Agent Status
IP			UUID					status		Hostname
172.16.210.1		5fd57838-3d84-4fdc-a87d-2f30eb55b870	offline		ajay-ubuntu


flashc  agent register -a 127.0.0.1
http://127.0.0.1:8989/agent/register
127.0.0.1
flashagent already registered
```

##### flashc container command set
```
flashc container -h
A longer description that spans multiple lines and likely contains examples
and usage of using the command. For example:

Usage:
  flashc container [flags]
  flashc container [command]

Available Commands:
  create      A brief description of your command
  destroy     destroys the containers with the given UUID
  list        A brief description of your command
  snap        snapshots the containers with the given UUID
  start       starts the containers with the given UUID
  stat        stats the containers with the given UUID
  stop        stops the containers with the given UUID

Global Flags:
  -C, --controllernode string   Provide the IP/DNS name of Controller Node (default "127.0.0.1")

Use "flashc container [command] --help" for more information about a command.

```
