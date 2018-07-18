package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (plugin *Plugin) collectAgentInfo() {
	nodeList := plugin.nodeDB.GetAllNodes()
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       Timeout,
	}
	for _, node := range nodeList {

		go plugin.getLivenessInfo(client, node)
		//go func(client http.Client) {
		//	res, err := client.Get("http://" + node.ManIPAdr + LivenessPort + LivessURL)
		//	if err != nil {
		//		plugin.Log.Error(err)
		//		plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nil}
		//		return
		//	}
		//	b, _ := ioutil.ReadAll(res.Body)
		//	b = []byte(b)
		//	nodeInfo := &NodeLiveness{}
		//	json.Unmarshal(b, nodeInfo)
		//	plugin.Log.Info(node.NodeLiveness)
		//	plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nodeInfo}
		//}(client)

		go plugin.getInterfaceInfo(client, node)
		//go func(client http.Client) {
		//	res, err := client.Get("http://" + node.ManIPAdr + InterfacePort + InterfaceURL)
		//	if err != nil {
		//		plugin.Log.Error(err)
		//		plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nil}
		//		return
		//	}
		//	b, _ := ioutil.ReadAll(res.Body)
		//	b = []byte(b)
		//
		//	nodeInterfaces := make(map[int]NodeInterface, 0)
		//	json.Unmarshal(b, &nodeInterfaces)
		//	plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nodeInterfaces}
		//}(client)
		go plugin.getBridgeDomainInfo(client, node)

	}
}
func (plugin *Plugin) getLivenessInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + LivenessPort + LivessURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodeInfo := &NodeLiveness{}
	json.Unmarshal(b, nodeInfo)
	plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nodeInfo}

}

func (plugin *Plugin) getInterfaceInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + InterfacePort + InterfaceURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)

	nodeInterfaces := make(map[int]NodeInterface, 0)
	json.Unmarshal(b, &nodeInterfaces)
	plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nodeInterfaces}

}
func (plugin *Plugin) getBridgeDomainInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + BridgeDomainsPort + BridgeDomainURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)

	nodeBridgeDomains := make(map[int]NodeBridgeDomains)
	json.Unmarshal(b, &nodeBridgeDomains)
	plugin.nDBChannel <- NodeBridgeDomainsDTO{nodeName: node.Name, nodeInfo: nodeBridgeDomains}

}
