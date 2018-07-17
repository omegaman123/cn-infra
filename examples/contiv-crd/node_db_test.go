package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
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

func TestNodesDB_SetNodeInfo(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1, "k8s_master", "10", "10.20.0.2")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(uint32(1)))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10.20.0.2"))
	res, err := http.Get("http://" + node.ManIPAdr + LivenessPort + LivessURL)
	gomega.Expect(err).To(gomega.BeNil())
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodeInfo := &NodeLiveness{}
	json.Unmarshal(b, nodeInfo)
	db.SetNodeInfo(node.Name, nodeInfo)
	//gomega.Expect(node.NodeLiveness.Build_version).To(gomega.Equal("v1.2-alpha-169-gcf4ac7e"))

	err = db.SetNodeInfo("NonExistantNode", nodeInfo)
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

}
