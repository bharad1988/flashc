package common

// This file defines all the structs and constants

//AgentPort is the tcp port where flashagent daemon runs
const AgentPort = "7979"

//ControllerPort is the tcp port where flashxontroller daemon runs
const ControllerPort = "8989"

//RepoPort serves imaging service
const RepoPort = "6969"

//JSONBodyType defines the json body type string
const JSONBodyType = "application/json; charset=UTF-8"

// MongoCtlrPort defines the port for mongo controller
const MongoCtlrPort = "27017"

// MongoAgentPort defines the port for mongo controller
const MongoAgentPort = "27017"

// ContMemSoftLimit is the soft limit for memory for each container (in ByteSize)
const ContMemSoftLimit = 26214400 // 25M

// ContMemLimit is the hard limit for memory (in ByteSize)
const ContMemLimit = 31457280 // 30M

// ConConfig is a superset of all configs related to flashc
// This structure is an example. Used to parse yaml structures
type ConConfig struct {
	//Container map[string]string
	Lxcpath    string
	Template   string
	Distro     string
	Release    string
	Arch       string
	Name       string
	Verbose    string
	Flush      string
	Validation string
	UUID       string
	State      string
	AgentUUID  string
	BatchUUID  string
	Backend    string
}

// MongoServer related data
type MongoServer struct {
	Server     string
	Port       string
	DB         string
	Collection string
}

// AgentInfo is the struct associated with agent node ( identifier )
type AgentInfo struct {
	AgentIP       string `bson:"agentip" json:"agentip"`
	UUID          string `bson:"uuid" json:"uuid"`
	Status        string `bson:"status" json:"status"`
	Hostname      string `bson:"hostname" json:"hostname"`
	TotalMemUsage string `bson:"totalmemusage" json:"totalmemusage"`
	//AgentPort string
}

// DestroyCon is a structure to provide details to destroy container
// mab be in future when count is also given
type DestroyCon struct {
	WithSnapShots string // set a value of true
}

// SuperCon is the structure to launch containers
// contains more meta information for container operations ( supplementary info. )
type SuperCon struct {
	Con   ConConfig
	Start string
	Count string
	DC    DestroyCon
	Snaps SnapShot
}

// Iface contains stats related to an interface
type Iface struct {
	Name string
	IP   string
	Rx   string
	Tx   string
}

// ContStat defines the stats related to a container
type ContStat struct {
	ContUUID string
	Intface  []Iface
	Memory   string
	Kmem     string
}

// SnapShot structure defines the structure for snapshot of a container . This is derived from lxc snapshot structure
type SnapShot struct {
	ContainerName string
	Name          string
	Path          string
	CommentPath   string
	TimeStamp     string
}

// SuperImage is the structure for container images from which new containers can be launched
type SuperImage struct {
	Image     ConImage
	ConPath   string
	ConUUID   string
	AgentUUID string
	Snap      string
}

// ConImage - struct about the image
type ConImage struct {
	User    string
	Name    string
	UUID    string
	Path    string
	Version string
}

// HTTPRetObj is a return object for REST calls, Returns the object along with a message
type HTTPRetObj struct {
	Obj     interface{}
	Message string
	Status  string
}

// RepoService Provides details about repo Server
type RepoService struct {
	Path string
}

// RepoUser defines the user name and his private key to login to image service
type RepoUser struct {
	User       string
	PrivateKey string
}

// Repo json file
type Repo struct {
	Path string
}
