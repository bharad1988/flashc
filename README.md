# flashc
A management suite to orchestrate lxc containers
  A highly scalable container management system. There will be a central controller node to which all the nodes participating in the cluster will register to. Each participating node will be running an agent which will provide the interface to connect to the controller.
  The controller provides with REST apis as a part of SDK to manage and monitor. 
As of now, the cluster can be managed by command line interface written on top REST apis

The container engine relies on btrfs to provide with snapshots.

More details will be added....
 
  
