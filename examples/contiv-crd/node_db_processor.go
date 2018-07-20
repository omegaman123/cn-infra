package main

//Reads data sent by agent_client.go to the plugins Node Database channel.
//Decides how to process the data received based on the type of Data Transfer Object.
//After that, it orders the map of data into a list and updates the node Database with that ordered list.
func (plugin *Plugin) ProcessNodeData(nodeList []*Node) {
	for i := 0; i < NodeHTTPCalls *len(nodeList); i++ {
		data := <-plugin.nDBChannel
		switch data.(type) {
		case NodeLivenessDTO:
			nlDto := data.(NodeLivenessDTO)
			plugin.nodeDB.SetNodeLiveness(nlDto.nodeName, nlDto.nodeInfo)
		case NodeInterfacesDTO:
			niDto := data.(NodeInterfacesDTO)
			plugin.nodeDB.SetNodeInterfaces(niDto.nodeName, niDto.nodeInfo)
		case NodeBridgeDomainsDTO:
			nbdDto := data.(NodeBridgeDomainsDTO)
			plugin.nodeDB.SetNodeBridgeDomain(nbdDto.nodeName, nbdDto.nodeInfo)
		case NodeL2FibsDTO:
			nl2fDto := data.(NodeL2FibsDTO)
			plugin.nodeDB.SetNodeL2Fibs(nl2fDto.nodeName, nl2fDto.nodeInfo)
		case NodeTelemetryDTO:
			ntDto := data.(NodeTelemetryDTO)
			plugin.nodeDB.SetNodeTelemetry(ntDto.nodeName, ntDto.nodeInfo)
		case NodeIPArpDTO:
			nipaDto := data.(NodeIPArpDTO)
			plugin.nodeDB.SetNodeIPARPs(nipaDto.nodeName, nipaDto.nodeInfo)
		default:
			plugin.Log.Error("Unknown data type")
		}
	}
}
