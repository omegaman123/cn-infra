package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	LivenessPort      = ":9999"
	LivessURL         = "/liveness"
	Timeout           = time.Second*10
	InterfacePort     = ":9999"
	InterfaceURL      = "/interfaces"
	BridgeDomainsPort = ":9999"
	BridgeDomainURL   = "/bridgedomains"
	L2FibsPort        = ":9999"
	L2FibsURL         = "/l2fibs"
	TelemetryPort     = ":9999"
	TelemetryURL      = "/telemetry"
	NodeHTTPCalls     = 6
)

//Gathers a number of data points for every node in the Node List

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

		//TODO: Implement getTelemetry correctly.
		//Does not parse information correctly
		go plugin.getTelemetryInfo(client, node)

		go plugin.getIPArpInfo(client, node)

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
	res, err := client.Get("http://" + node.ManIPAdr + TelemetryPort + TelemetryURL)
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

func (plugin *Plugin) getIPArpInfo(client http.Client, node *Node) {
	res, err := client.Get("http://" + node.ManIPAdr + TelemetryPort + TelemetryURL)
	if err != nil {
		plugin.Log.Error(err)
		plugin.nDBChannel <- NodeIPArpDTO{nodeName: node.Name, nodeInfo: nil}
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	b = []byte(b)
	nodeiparp := make(map[string]NodeIPArp)
	json.Unmarshal(b, &nodeiparp)
	plugin.nDBChannel <- NodeIPArpDTO{nodeName: node.Name, nodeInfo: nodeiparp}
}
