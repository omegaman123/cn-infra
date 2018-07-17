package main

import (
	"sort"

	"github.com/ligato/cn-infra/logging"
	"github.com/pkg/errors"
)

type Node struct {
	ID             uint32
	IPAdr          string
	ManIPAdr       string
	Name           string
	NodeLiveness   *NodeLiveness
	NodeInterfaces []NodeInterface
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
	NodeInfo *NodeLiveness
}

type NodeInterface struct {
	VppInternalName string   `json:"vpp_internal_name"`
	Name            string   `json:"name"`
	IfType          uint32   `json:"type,omitempty"`
	Enabled         bool     `json:"enabled"`
	PhysAddress     string   `json:"phys_address"`
	Mtu             uint32   `json:"mtu,omitempty"`
	IpAddresses     []string `json:"ip_addresses"`
	Tap             tap      `json:"tap,omitempty"`
}
type tap struct {
	Version    uint32 `json:"version"`
	HostIfName string `json:"host_if_name"`
}

type NodeInterfacesDTO struct {
	nodeName string
	nodeInfo map[int]NodeInterface
}

type Nodes interface {
	AddNode(ID uint32, nodeName, IPAdr, ManIPAdr string) error
	DeleteNode(key string) error
	GetNode(key string) (*Node, error)
	GetAllNodes() []*Node
	SetNodeInfo(name string, nL *NodeLiveness) error
	SetNodeInterfaces(name string, nInt []NodeInterface) error
}

type NodesDB struct {
	nMap   map[string]*Node
	logger logging.PluginLogger
}

func NewNodesDB(logger logging.PluginLogger) (n *NodesDB) {
	return &NodesDB{make(map[string]*Node), logger}
}

func (nDb *NodesDB) SetNodeInfo(name string, nL *NodeLiveness) error {
	node, err := nDb.GetNode(name)
	if err != nil {
		return err
	}
	node.NodeLiveness = nL
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
