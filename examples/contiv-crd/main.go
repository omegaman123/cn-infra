//go:generate protoc --proto_path=./model --go_out=./model ./model/node.proto
package main

import (
	"log"
	"os"
	"os/signal"

	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/ligato/cn-infra/agent"
	"github.com/ligato/cn-infra/config"
	"github.com/ligato/cn-infra/datasync"
	"github.com/ligato/cn-infra/datasync/kvdbsync"
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/db/keyval/etcd"
	"github.com/ligato/cn-infra/examples/contiv-crd/model"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/servicelabel"
)

// *************************************************************************
// This example demonstrates the usage of datasync API with etcd
// as the data store.
// ExamplePlugin spawns a data publisher and a data consumer (watcher)
// as two separate go routines.
// The publisher executes two operations on the same key: CREATE + UPDATE.
// The consumer is notified with each change and reports the events into
// the log.
// ************************************************************************/

const (
	PluginName        = "contiv-crd"
)

func main() {

	// Init close channel used to stop the example.
	pluginFinished := make(chan struct{})

	etcdPlug := etcd.NewPlugin(
		etcd.UseConf(etcd.Config{
			Endpoints: []string{":32379"},
		}),
	)

	etcdDataSync := kvdbsync.NewPlugin(
		kvdbsync.UseDeps(kvdbsync.Deps{
			KvPlugin: etcdPlug,
		}),
	)

	p := &Plugin{
		Deps: Deps{
			Log:          logging.ForPlugin(PluginName),
			PluginConfig: config.ForPlugin(PluginName),
			ServiceLabel: servicelabel.DefaultPlugin,
			Getter:       etcdDataSync,
		},
		closeChannel: pluginFinished,
	}

	a := agent.NewAgent(
		agent.AllPlugins(p),
		agent.QuitOn(pluginFinished),
	)

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

type KeyProtoValBroker interface {
	// Put <data> to ETCD or to any other key-value based data source.
	Put(key string, data proto.Message, opts ...datasync.PutOption) error

	// Delete data under the <key> in ETCD or in any other key-value based data
	// source.
	Delete(key string, opts ...datasync.DelOption) (existed bool, err error)

	// GetValue reads a value from etcd stored under the given key.
	GetValue(key string, reqObj proto.Message) (found bool, revision int64, err error)

	// List values stored in etcd under the given prefix.
	ListValues(prefix string) (keyval.ProtoKeyValIterator, error)
}

// ExamplePlugin demonstrates the usage of datasync API.
type Plugin struct {
	Deps

	stopCh chan struct{}
	wg     sync.WaitGroup

	//k8sClientConfig *rest.Config
	//k8sClientset    *kubernetes.Clientset
	closeChannel chan struct{}
	broker       KeyProtoValBroker
	nodeDB       Nodes
	nDBChannel   chan interface{}
}

// Deps lists dependencies of ExamplePlugin.
type Deps struct {
	Log          logging.PluginLogger
	PluginConfig config.PluginConfig
	ServiceLabel servicelabel.ReaderAPI

	Getter *kvdbsync.Plugin
}

// Name implements PluginNamed
func (p *Plugin) Name() string {
	return PluginName
}

// Init starts the consumer.
func (plugin *Plugin) Init() error {
	// Initialize plugin fields.
	plugin.broker = plugin.Getter.KvPlugin.NewBroker("")
	plugin.Log.Info("Initialization of the custom plugin for the contiv-crd example is completed")
	plugin.nodeDB = NewNodesDB(plugin.Log)

	// Start the consumer (ETCD watcher).
	go plugin.consumer()

	return nil
}

// Close shutdowns both the publisher and the consumer.
// Channels used to propagate data resync and data change events are closed
// as well.
func (plugin *Plugin) Close() error {
	return nil
}

// AfterInit starts the publisher and prepares for the shutdown.
func (plugin *Plugin) AfterInit() error {

	go plugin.closePlugin()

	return nil
}

func (plugin *Plugin) closePlugin() {
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	log.Println("Program killed !")

	// do last actions and wait for all write operations to end

	os.Exit(0)
}

// consumer (watcher) is subscribed to watch on data store changes.
// Changes arrive via data change channel, get identified based on the key
// and printed into the log.
func (plugin *Plugin) consumer() {

	plugin.Log.Print("KeyValProtoGetter started")

	messageList, err := plugin.broker.ListValues("/vnf-agent/contiv-ksr/allocatedIDs")
	if err != nil {
		plugin.Log.Error("Error: ", err)
	}
	for {
		message, stop := messageList.GetNext()
		protoMessage := &node.NodeInfo{}
		if stop {
			plugin.Log.Info("No more data under: ", messageList)
			break
		}
		err = message.GetValue(protoMessage)
		if err != nil {
			plugin.Log.Error("Error in getting value of iterator: ", err)
			continue
		}
		plugin.Log.Infof("Getting data under %+v : %+v", messageList, protoMessage)
		plugin.nodeDB.AddNode(protoMessage.Id, protoMessage.Name, protoMessage.IpAddress, protoMessage.ManagementIpAddress)
	}
	//Rest client
	nodeList := plugin.nodeDB.GetAllNodes()

	plugin.nDBChannel = make(chan interface{})

	plugin.collectAgentInfo()

	plugin.ProcessNodeData(nodeList)

	for _, node := range nodeList {
		plugin.Log.Infof("Node Info: %+v ", node)
		plugin.Log.Infof("Node Liveness: %+v", node.NodeLiveness)
		plugin.Log.Infof("Node Interfaces: %+v", node.NodeInterfaces)
		plugin.Log.Infof("Node Bridge Domains: %+v", node.NodeBridgeDomains)
		plugin.Log.Infof("Node L2Fibs: %+v", node.NodeL2Fibs)
		plugin.Log.Infof("Node Telemetry: %+v", node.NodeTelemetry)
	}

}
