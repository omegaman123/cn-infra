package main

type Node struct {
	IPAdr 	 string
	ManIPAdr string
}

type Nodes interface {
	AddNode(nodeName, IPAdr, ManIPAdr string) error
	DeleteNode(key string) error
	GetNode(key string) (error, Node)
}

type NodesDB struct {
  nMap map[string]Node
}

func NewNodesDB()(n *NodesDB) {
	return &NodesDB{make(map[string]Node)}
}

func (nDB *NodesDB)GetNode(key string) (n *Node, err error) {
node,ok := nDB.nMap[key]
	if !ok {
		return nil, err
	}
	return &node,nil
}

func (nDB *NodesDB)DeleteNode(key string) error {
	_, ok := nDB.nMap[key]
	if !ok{
		return nil
	}
	delete(nDB.nMap,key)
	return nil
}

func (nDB *NodesDB)AddNode(nodeName, IPAdr, ManIPAdr string) error {
	n:= Node{IPAdr:IPAdr, ManIPAdr:ManIPAdr}
	nDB.nMap[nodeName] = n
	return nil
}
