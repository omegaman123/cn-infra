package main

import "sort"

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
			var keyslice []int
			nodeinterfaces := make([]NodeInterface, 0)
			niDto := data.(NodeInterfacesDTO)
			for itfkey := range niDto.nodeInfo {
				keyslice = append(keyslice, itfkey)
			}
			sort.Ints(keyslice)
			for _, itfkey := range keyslice {
				nodeinterfaces = append(nodeinterfaces, niDto.nodeInfo[itfkey])
			}
			plugin.nodeDB.SetNodeInterfaces(niDto.nodeName, nodeinterfaces)
		case NodeBridgeDomainsDTO:
			nbdDto := data.(NodeBridgeDomainsDTO)
			var keyslice []int
			nodebridgedomains := make([]NodeBridgeDomains, 0)
			for bdkey := range nbdDto.nodeInfo {
				keyslice = append(keyslice, bdkey)
			}
			sort.Ints(keyslice)
			for _, bdkey := range keyslice {
				nodebridgedomains = append(nodebridgedomains, nbdDto.nodeInfo[bdkey])
			}
			plugin.nodeDB.SetNodeBridgeDomain(nbdDto.nodeName, nodebridgedomains)
		case NodeL2FibsDTO:
			nl2fDto := data.(NodeL2FibsDTO)
			var keyslice []string
			nodel2fibs := make([]NodeL2Fib, 0)
			for l2fkey := range nl2fDto.nodeInfo {
				keyslice = append(keyslice, l2fkey)
			}
			sort.Strings(keyslice)
			for _, l2fkey := range keyslice {
				nodel2fibs = append(nodel2fibs, nl2fDto.nodeInfo[l2fkey])
			}
			plugin.nodeDB.SetNodeL2Fibs(nl2fDto.nodeName, nodel2fibs)
		case NodeTelemetryDTO:
			ntDto := data.(NodeTelemetryDTO)
			var keyslice []string
			nodetelemetry := make([]NodeTelemetry, 0)
			for telekey := range ntDto.nodeInfo {
				keyslice = append(keyslice, telekey)
			}
			sort.Strings(keyslice)
			for _, telekey := range keyslice {
				nodetelemetry = append(nodetelemetry, ntDto.nodeInfo[telekey])
			}
			plugin.nodeDB.SetNodeTelemetry(ntDto.nodeName, nodetelemetry)
		case NodeIPArpDTO:
			nipaDto := data.(NodeIPArpDTO)
			var keyslice []string
			nodeiparp := make([]NodeIPArp, 0)
			for arpkey := range nipaDto.nodeInfo {
				keyslice = append(keyslice, arpkey)
			}
			sort.Strings(keyslice)
			for _, arpkey := range keyslice {
				nodeiparp = append(nodeiparp, nipaDto.nodeInfo[arpkey])
			}
			plugin.nodeDB.SetNodeIPARPs(nipaDto.nodeName, nodeiparp)
		default:
			plugin.Log.Error("Unknown data type")
		}
	}
}
