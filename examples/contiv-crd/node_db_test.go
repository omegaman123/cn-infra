package main

import (
	"testing"

	"github.com/onsi/gomega"
)

//Checks adding a new node.
//Checks expected error for adding duplicate node.
func TestNodesDB_AddNode(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "20")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))

	err := db.AddNode(2, "k8s_master", "20", "20")
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

}

//Checks adding a node and then looking it up.
//Checks looking up a non-existent key.
func TestNodesDB_GetNode(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nodeTwo, err := db.GetNode("NonExistentNode")
	gomega.Ω(err.Error()).Should(gomega.Equal("value with given key not found: NonExistentNode"))
	gomega.Expect(nodeTwo).To(gomega.BeNil())
}

//Checks adding a node and then deleting it.
//Checks whether expected error is returned when deleting non-existent key.
func TestNodesDB_DeleteNode(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))

	err := db.DeleteNode("k8s_master")
	gomega.Expect(err).To(gomega.BeNil())
	node, err = db.GetNode("k8s_master")
	gomega.Expect(node).To(gomega.BeNil())
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

	err = db.DeleteNode("k8s_master")
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

}

//Creates 3 new nodes and adds them to a database.
//Then, the list is checked to see if it is in order.
func TestNodesDB_GetAllNodes(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	db.AddNode(2, "k8s_master2", "10", "10")
	node, ok = db.GetNode("k8s_master2")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master2"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(2)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	db.AddNode(3, "Ak8s_master3", "10", "10")
	node, ok = db.GetNode("Ak8s_master3")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("Ak8s_master3"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(3)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nodeList := db.GetAllNodes()
	gomega.Expect(len(nodeList)).To(gomega.Equal(3))
	gomega.Expect(nodeList[0].Name).To(gomega.Equal("Ak8s_master3"))

}

func TestNodesDB_SetNodeInterfaces(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))
	var ips []string
	nodeIFs := make([]NodeInterface,0)
	nodeIF := NodeInterface{"Test","Testing",0,true,"",0,vxlan{},ips,tap{}}
	nodeIFs = append(nodeIFs,nodeIF)
	err := db.SetNodeInterfaces("NENODE",nodeIFs)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
	err = db.SetNodeInterfaces("k8s_master",nodeIFs)
	gomega.Expect(err).To(gomega.BeNil())
}

func TestNodesDB_SetNodeBridgeDomain(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))
	var ifs []bdinterfaces
	nodeBD := NodeBridgeDomains{ifs,"",false,}
	nodesBDs := make([]NodeBridgeDomains,0)
	nodesBDs = append(nodesBDs,nodeBD)

	err := db.SetNodeBridgeDomain("NENODE",nodesBDs)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
	err = db.SetNodeBridgeDomain("k8s_master",nodesBDs)
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(node.NodeBridgeDomains[0]).To(gomega.BeEquivalentTo(NodeBridgeDomains{ifs,"",false}))

}

func TestNodesDB_SetNodeIPARPs(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nodeiparps := make([]NodeIPArp,0)
	nodeiparp := NodeIPArp{1,"1.2.3.4","12:34:56:78",false}
	nodeiparps = append(nodeiparps,nodeiparp)

	err := db.SetNodeIPARPs("NENODE", nodeiparps)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
	err = db.SetNodeIPARPs("k8s_master",nodeiparps)
	gomega.Expect(err).To(gomega.BeNil())

	gomega.Expect(nodeiparps[0]).To(gomega.BeEquivalentTo(NodeIPArp{1,"1.2.3.4","12:34:56:78",false}))
}

func TestNodesDB_SetNodeLiveness(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nlive := NodeLiveness{"54321","12345",0,0,0,0,""}
	err := db.SetNodeLiveness("NENODE",&nlive)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
	err = db.SetNodeLiveness("k8s_master",&nlive)
	gomega.Expect(err).To(gomega.BeNil())

	gomega.Expect(node.NodeLiveness).To(gomega.BeEquivalentTo(&NodeLiveness{"54321","12345",0,0,0,0,""}))

}

func TestNodesDB_SetNodeL2Fibs(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nfib := NodeL2Fib{1,2,"test",true,false}
	var nfibs []NodeL2Fib
	nfibs = append(nfibs, nfib)

	err := db.SetNodeL2Fibs("NENODE",nfibs)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))
	err = db.SetNodeL2Fibs("k8s_master",nfibs)
	gomega.Expect(err).To(gomega.BeNil())

	gomega.Expect(node.NodeL2Fibs[0]).To(gomega.BeEquivalentTo(NodeL2Fib{1,2,"test",true,false}))

}


