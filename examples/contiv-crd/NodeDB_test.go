package main

import (
	"testing"
	"github.com/onsi/gomega"
)

func TestNodesDB_AddNode(t *testing.T) {
		gomega.RegisterTestingT(t)
		db := NewNodesDB(nil)
		db.AddNode(1,"k8s_master", "10", "20")
		node, ok := db.GetNode("k8s_master")
		gomega.Expect(ok).To(gomega.BeNil())
		gomega.Expect(node.IPAdr).To(gomega.Equal("10"))

		err := db.AddNode(2,"k8s_master","20","20")
		gomega.Expect(err).To(gomega.Not(gomega.BeNil()))


}

func TestNodesDB_GetNode(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1,"k8s_master","10","10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(1))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))


	nodeTwo, err := db.GetNode("NonExistentNode")
	gomega.Ω(err.Error()).Should(gomega.Equal("value with given key not found: NonExistentNode"))
	gomega.Expect(nodeTwo).To(gomega.BeNil())
}

func TestNodesDB_DeleteNode(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1,"k8s_master","10","10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))

	err := db.DeleteNode("k8s_master")
	gomega.Expect(err).To(gomega.BeNil())
	node,err = db.GetNode("k8s_master")
	gomega.Expect(node).To(gomega.BeNil())
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

	err = db.DeleteNode("k8s_master")
	gomega.Expect(err).To(gomega.Not(gomega.BeNil()))

}

func TestNodesDB_GetAllNodes(t *testing.T) {
	gomega.RegisterTestingT(t)
	db := NewNodesDB(nil)
	db.AddNode(1,"k8s_master","10","10")
	node, ok := db.GetNode("k8s_master")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master"))
	gomega.Expect(node.ID).To(gomega.Equal(1))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))


	db.AddNode(2,"k8s_master2","10","10")
	node, ok = db.GetNode("k8s_master2")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master2"))
	gomega.Expect(node.ID).To(gomega.Equal(2))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	db.AddNode(3,"k8s_master3","10","10")
	node, ok = db.GetNode("k8s_master3")
	gomega.Expect(ok).To(gomega.BeNil())
	gomega.Expect(node.IPAdr).To(gomega.Equal("10"))
	gomega.Expect(node.Name).To(gomega.Equal("k8s_master3"))
	gomega.Expect(node.ID).To(gomega.Equal(3))
	gomega.Expect(node.ManIPAdr).To(gomega.Equal("10"))

	nodeList := db.GetAllNodes()
	gomega.Expect(len(nodeList)).To(gomega.Equal(3))


}

