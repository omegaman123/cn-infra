package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//Gathers a number of data points for every node in the Node List
//using goroutines and a
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

		go plugin.getInterfaceInfo(client, node)

		go plugin.getBridgeDomainInfo(client, node)

		go plugin.getL2FibInfo(client, node)

		go plugin.getTelemetryInfo(client, node)

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
		plugin.nDBChannel <- NodeBridgeDomainsDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)

	nodeBridgeDomains := make(map[int]NodeBridgeDomains)
	json.Unmarshal(b, &nodeBridgeDomains)
	plugin.nDBChannel <- NodeBridgeDomainsDTO{nodeName: node.Name, nodeInfo: nodeBridgeDomains}

}

func (plugin *Plugin) getL2FibInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + L2FibsPort + L2FibsURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeL2FibsDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodel2fibs := make(map[string]NodeL2Fib)
	json.Unmarshal(b, &nodel2fibs)
	plugin.nDBChannel <- NodeL2FibsDTO{nodeName: node.Name, nodeInfo: nodel2fibs}

}

func (plugin *Plugin) getTelemetryInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + L2FibsPort + L2FibsURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeTelemetryDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodetelemetry := make(map[string]NodeTelemetry)
	json.Unmarshal(b, &nodetelemetry)
	plugin.nDBChannel <- NodeTelemetryDTO{nodeName: node.Name, nodeInfo: nodetelemetry}
}
