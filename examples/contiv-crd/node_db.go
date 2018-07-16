package main

import (
	"sort"

	"github.com/ligato/cn-infra/logging"
	"github.com/pkg/errors"
)

//const name  =
type Node struct {
	ID       uint32
	IPAdr    string
	ManIPAdr string
	Name     string
	NodeInfo *NodeLiveness
}

type NodeLiveness struct {
	Build_version string `json:"build_version"`
	Build_date    string `json:"build_date"`
	State         uint32 `json:"state"`
	Start_time    uint32 `json:"start_time"`
	Last_change   uint32 `json:"last_change"`
	Last_update   uint32 `json:"last_update"`
	Commit_hash   string `json:"commit_hash"`
}

type NodeLivenessDTO struct {
	nodeName string
	NodeInfo *NodeLiveness
}

type Nodes interface {
	AddNode(ID uint32, nodeName, IPAdr, ManIPAdr string) error
	DeleteNode(key string) error
	GetNode(key string) (*Node, error)
	GetAllNodes() []*Node
	SetNodeInfo(name string, nL *NodeLiveness) error
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
	node.NodeInfo = nL
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
