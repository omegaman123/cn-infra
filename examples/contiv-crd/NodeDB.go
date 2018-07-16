package main

import (
	"sort"
	"github.com/ligato/cn-infra/logging"
)
type Node struct {
	ID       int
	IPAdr    string
	ManIPAdr string
	Name     string
}

type Nodes interface {
	AddNode(ID int, nodeName, IPAdr, ManIPAdr string) error
	DeleteNode(key string) error
	GetNode(key string) (*Node, error)
	GetAllNodes()[]*Node
}

type NodesDB struct {
	nMap   map[string]Node
	logger logging.PluginLogger
}


func NewNodesDB(logger logging.PluginLogger) (n *NodesDB) {
	return &NodesDB{make(map[string]Node), logger}
}

//Returns a pointer to a node for the given key
//Returns an error if that key is not found.
func (nDB *NodesDB) GetNode(key string) (n *Node, err error) {
	node, ok := nDB.nMap[key]
	if !ok {
		return nil, err
	}
	return &node, nil
}

//Deletes a key with the given key
//Returns an error if the key is not found.
func (nDB *NodesDB) DeleteNode(key string) error {
	_, ok := nDB.nMap[key]
	if !ok {
		return nil
	}
	delete(nDB.nMap, key)
	return nil
}

//Adds a new node with the given information
//Returns an error if the node is already in the database
func (nDB *NodesDB) AddNode(ID int, nodeName, IPAdr, ManIPAdr string) error {
	n := Node{IPAdr: IPAdr, ManIPAdr: ManIPAdr, ID: ID, Name: nodeName}
	_, err := nDB.GetNode(nodeName)
	if err != nil {
		return err
	}
	nDB.nMap[nodeName] = n
	return nil
}

//Returns an ordered slice of all nodes in a database.
func (nDB *NodesDB) GetAllNodes() []*Node {
	var str []string
	for k := range nDB.nMap   {
		str = append(str, k)
	}
	var nList []*Node
	sort.Strings(str)
	for _ ,v := range str {
		n , _ := nDB.GetNode(v)
		nList = append(nList, n)
	}
	return nList
}
