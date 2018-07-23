package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	livenessPort      = ":9999"
	livenessURL       = "/liveness"
	timeout           = 100000000000
	interfacePort     = ":9999"
	interfaceURL      = "/interfaces"
	bridgeDomainsPort = ":9999"
	bridgeDomainURL   = "/bridgedomains"
	l2FibsPort        = ":9999"
	l2FibsURL         = "/l2fibs"
	telemetryPort     = ":9999"
	telemetryURL      = "/telemetry"
	arpPort           = ":9999"
	arpURL            = "/arps"
	nodeHTTPCalls     = 5
)

//Gathers a number of data points for every node in the Node List

func (Plugin *Plugin) collectAgentInfo() {
	nodeList := Plugin.nodeDB.GetAllNodes()
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       timeout,
	}

	for _, node := range nodeList {

		go Plugin.getLivenessInfo(client, node)

		go Plugin.getInterfaceInfo(client, node)

		go Plugin.getBridgeDomainInfo(client, node)

		go Plugin.getL2FibInfo(client, node)

		//TODO: Implement getTelemetry correctly.
		//Does not parse information correctly
		//go Plugin.getTelemetryInfo(client, node)

		go Plugin.getIPArpInfo(client, node)

	}
}

/* Here are the several functions that run as goroutines to collect information
about a specific node using an http client. First, an http request is made to the
specific url and port of the desired information and the request received is read
and unmarshalled into a struct to contain that information. Then, a data transfer
object is created to hold the struct of information as well as the name and is sent
over the plugins node database channel to node_db_processor.go where it will be read,
processed, and added to the node database.
*/

func (Plugin *Plugin) getLivenessInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + livenessPort + livenessURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodeInfo := &NodeLiveness{}
	json.Unmarshal(b, nodeInfo)
	Plugin.nDBChannel <- NodeLivenessDTO{nodeName: node.Name, nodeInfo: nodeInfo}

}

func (Plugin *Plugin) getInterfaceInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + interfacePort + interfaceURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)

	nodeInterfaces := make(map[int]NodeInterface, 0)
	json.Unmarshal(b, &nodeInterfaces)
	Plugin.nDBChannel <- NodeInterfacesDTO{nodeName: node.Name, nodeInfo: nodeInterfaces}

}
func (Plugin *Plugin) getBridgeDomainInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + bridgeDomainsPort + bridgeDomainURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeBridgeDomainsDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)

	nodeBridgeDomains := make(map[int]NodeBridgeDomains)
	json.Unmarshal(b, &nodeBridgeDomains)
	Plugin.nDBChannel <- NodeBridgeDomainsDTO{nodeName: node.Name, nodeInfo: nodeBridgeDomains}

}

func (Plugin *Plugin) getL2FibInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + l2FibsPort + l2FibsURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeL2FibsDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodel2fibs := make(map[string]NodeL2Fib)
	json.Unmarshal(b, &nodel2fibs)
	Plugin.nDBChannel <- NodeL2FibsDTO{nodeName: node.Name, nodeInfo: nodel2fibs}

}

func (Plugin *Plugin) getTelemetryInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + telemetryPort + telemetryURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeTelemetryDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodetelemetry := make(map[string]NodeTelemetry)
	json.Unmarshal(b, &nodetelemetry)
	Plugin.nDBChannel <- NodeTelemetryDTO{nodeName: node.Name, nodeInfo: nodetelemetry}
}

func (Plugin *Plugin) getIPArpInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + arpPort + arpURL)
	if err != nil {
		Plugin.Log.Error(err)
		Plugin.nDBChannel <- NodeIPArpDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)

	b = []byte(b)
	nodeiparpslice := make([]NodeIPArp, 0)
	json.Unmarshal(b, &nodeiparpslice)
	Plugin.nDBChannel <- NodeIPArpDTO{nodeName: node.Name, nodeInfo: nodeiparpslice}
}
