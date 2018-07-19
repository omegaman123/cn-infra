package main

import (
	"sort"

	"github.com/ligato/cn-infra/logging"
	"github.com/pkg/errors"
)

type Node struct {
	ID                uint32
	IPAdr             string
	ManIPAdr          string
	Name              string
	NodeLiveness      *NodeLiveness
	NodeInterfaces    []NodeInterface
	NodeBridgeDomains []NodeBridgeDomains
	NodeL2Fibs        []NodeL2Fib
	NodeTelemetry     []NodeTelemetry
}

type NodeLiveness struct {
	BuildVersion string `json:"build_version"`
	BuildDate    string `json:"build_date"`
	State        uint32 `json:"state"`
	StartTime    uint32 `json:"start_time"`
	LastChange   uint32 `json:"last_change"`
	LastUpdate   uint32 `json:"last_update"`
	CommitHash   string `json:"commit_hash"`
}

type NodeLivenessDTO struct {
	nodeName string
	nodeInfo *NodeLiveness
}

type NodeTelemetryDTO struct {
	nodeName string
	nodeInfo map[string]NodeTelemetry
}

type NodeTelemetry struct {
	Command string   `json:"command"`
	Output  []output `json:"output"`
}

type output struct {
	command string
	output  []outputEntry
}

type outputEntry struct {
	nodeName string
	count    int
	reason   string
}

type NodeL2Fib struct {
	BridgeDomainIdx          uint32 `json:"bridge_domain_idx"`
	OutgoingInterfaceSwIfIdx uint32 `json:"outgoing_interface_sw_if_idx"`
	PhysAddress              string `json:"phys_address"`
	StaticConfig             bool   `json:"static_config"`
	BridgedVirtualInterface  bool   `json:"bridged_virtual_interface"`
}

type NodeL2FibsDTO struct {
	nodeName string
	nodeInfo map[string]NodeL2Fib
}

type NodeInterface struct {
	VppInternalName string   `json:"vpp_internal_name"`
	Name            string   `json:"name"`
	IfType          uint32   `json:"type,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	PhysAddress     string   `json:"phys_address,omitempty"`
	Mtu             uint32   `json:"mtu,omitempty"`
	Vxlan           vxlan    `json:"vxlan,omitempty"`
	IpAddresses     []string `json:"ip_addresses,omitempty"`
	Tap             tap      `json:"tap,omitempty"`
}
type vxlan struct {
	SrcAddress string `json:"src_address"`
	DstAddress string `json:"dst_address"`
	Vni        uint32 `json:"vni"`
}

type NodeIPArp struct {
	Interface  uint32 `json:"interface"`
	IPAddress  string `json:"IPAddress"`
	MacAddress string `json:"MacAddress"`
	static     bool   `json:"Static"`
}

type NodeIPArpDTO struct {
	nodeInfo map[string]NodeIPArp
	nodeName string
}
type tap struct {
	Version    uint32 `json:"version"`
	HostIfName string `json:"host_if_name"`
}

type NodeInterfacesDTO struct {
	nodeName string
	nodeInfo map[int]NodeInterface
}
type NodeBridgeDomains struct {
	Interfaces []bdinterfaces `json:"interfaces"`
	Name       string         `json:"name"`
	Forward    bool           `json:"forward"`
}
type bdinterfaces struct {
	SwIfIndex uint32 `json:"sw_if_index"`
}

type NodeBridgeDomainsDTO struct {
	nodeName string
	nodeInfo map[int]NodeBridgeDomains
}

type Nodes interface {
	AddNode(ID uint32, nodeName, IPAdr, ManIPAdr string) error
	DeleteNode(key string) error
	GetNode(key string) (*Node, error)
	GetAllNodes() []*Node
	SetNodeLiveness(name string, nL *NodeLiveness) error
	SetNodeInterfaces(name string, nInt []NodeInterface) error
	SetNodeBridgeDomain(name string, nBridge []NodeBridgeDomains) error
	SetNodeL2Fibs(name string, nL2f []NodeL2Fib) error
	SetNodeTelemetry(name string, nTele []NodeTelemetry) error
}

type NodesDB struct {
	nMap   map[string]*Node
	logger logging.PluginLogger
}

//Returns a pointer to a new node Database
func NewNodesDB(logger logging.PluginLogger) (n *NodesDB) {
	return &NodesDB{make(map[string]*Node), logger}
}

func (nDb *NodesDB) SetNodeLiveness(name string, nLive *NodeLiveness) error {
	node, err := nDb.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeLiveness = nLive
	return nil
}

func (nDB *NodesDB) SetNodeInterfaces(name string, nInt []NodeInterface) error {
	node, err := nDB.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeInterfaces = nInt
	return nil

}
func (nDB *NodesDB) SetNodeBridgeDomain(name string, nBridge []NodeBridgeDomains) error {
	node, err := nDB.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeBridgeDomains = nBridge
	return nil
}

func (nDB *NodesDB) SetNodeL2Fibs(name string, nL2F []NodeL2Fib) error {
	node, err := nDB.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeL2Fibs = nL2F
	return nil
}

func (nDB *NodesDB) SetNodeTelemetry(name string, nTele []NodeTelemetry) error {
	node, err := nDB.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeTelemetry = nTele
	return nil
}

//Returns a pointer to a node for the given key.
//Returns an error if that key is not found.
func (nDB *NodesDB) GetNode(key string) (n *Node, err error) {
	if node, ok := nDB.nMap[key]; ok {
		return node, nil
	}
	err = errors.Errorf("value with given key not found: %s", key)
	return nil, err
}

//Deletes node with the given key.
//Returns an error if the key is not found.
func (nDB *NodesDB) DeleteNode(key string) error {
	_, err := nDB.GetNode(key)
	if err != nil {
		return err
	}
	delete(nDB.nMap, key)
	return nil
}

//Adds a new node with the given information.
//Returns an error if the node is already in the database
func (nDB *NodesDB) AddNode(ID uint32, nodeName, IPAdr, ManIPAdr string) error {
	n := &Node{IPAdr: IPAdr, ManIPAdr: ManIPAdr, ID: ID, Name: nodeName}
	_, err := nDB.GetNode(nodeName)
	if err == nil {
		err = errors.Errorf("duplicate key found: %s", nodeName)
		return err
	}
	nDB.nMap[nodeName] = n
	return nil
}

//Returns an ordered slice of all nodes in a database organized by name.
func (nDB *NodesDB) GetAllNodes() []*Node {
	var str []string
	for k := range nDB.nMap {
		str = append(str, k)
	}
	var nList []*Node
	sort.Strings(str)
	for _, v := range str {
		n, _ := nDB.GetNode(v)
		nList = append(nList, n)
	}
	return nList
}
