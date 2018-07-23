package main

//ProcessNodeData reads data sent by agent_client.go to the plugins Node Database channel.
//It decides how to process the data received based on the type of Data Transfer Object.
//Then it updates the node with the name from the DTO with the specific data from the DTO.
func (Plugin *Plugin) ProcessNodeData(nodeList []*Node) {
	for i := 0; i < nodeHTTPCalls*len(nodeList); i++ {
		data := <-Plugin.nDBChannel
		switch data.(type) {
		case NodeLivenessDTO:
			nlDto := data.(NodeLivenessDTO)
			Plugin.nodeDB.SetNodeLiveness(nlDto.nodeName, nlDto.nodeInfo)
		case NodeInterfacesDTO:
			niDto := data.(NodeInterfacesDTO)
			Plugin.nodeDB.SetNodeInterfaces(niDto.nodeName, niDto.nodeInfo)
		case NodeBridgeDomainsDTO:
			nbdDto := data.(NodeBridgeDomainsDTO)
			Plugin.nodeDB.SetNodeBridgeDomain(nbdDto.nodeName, nbdDto.nodeInfo)
		case NodeL2FibsDTO:
			nl2fDto := data.(NodeL2FibsDTO)
			Plugin.nodeDB.SetNodeL2Fibs(nl2fDto.nodeName, nl2fDto.nodeInfo)
		case NodeTelemetryDTO:
			ntDto := data.(NodeTelemetryDTO)
			Plugin.nodeDB.SetNodeTelemetry(ntDto.nodeName, ntDto.nodeInfo)
		case NodeIPArpDTO:
			nipaDto := data.(NodeIPArpDTO)
			Plugin.nodeDB.SetNodeIPARPs(nipaDto.nodeName, nipaDto.nodeInfo)
		default:
			Plugin.Log.Error("Unknown data type")
		}
	}
}
