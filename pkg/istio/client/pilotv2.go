package client

import (
	"context"
	"fmt"
	"time"

	"github.com/ServiceComb/go-chassis/pkg/istio/util"
	envoy_api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	xds "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"google.golang.org/grpc"
	"github.com/ServiceComb/go-archaius/lager"
    "encoding/json"
)

// PilotClient is a interface for the client to communicate to pilot
type PilotClient interface {
	RDS
    	// TODO: add all xDS interface
	EDS
	CDS
	LDS
}

// RDS defines route discovery service interface
type RDS interface {
	GetAllRouteConfigurations() (*envoy_api.RouteConfiguration, error)
	GetRouteConfigurationsByPort(string) (*envoy_api.RouteConfiguration, error)
}

// EDS defines endpoint discovery service interface
type EDS interface{}

// CDS defines cluster discovery service interface
type CDS interface{
	GetAllClusterConfigurations() ([]envoy_api.Cluster, error)
}

// LDS defines listener discovery service interface
type LDS interface{}



type pilotClient struct {
	rawConn *grpc.ClientConn

	adsConn xds.AggregatedDiscoveryServiceClient
	edsConn envoy_api.EndpointDiscoveryServiceClient
}

// NewGRPCPilotClient returns new PilotClient from options
func NewGRPCPilotClient(cfg *PilotOptions) (PilotClient, error) {
	// TODO: credentials need to be added here
	// set dial options from config
	lager.Logger.Info(fmt.Sprintf("Raj the value of endponts ", cfg.Endpoints[0]))
	conn, err := grpc.Dial(cfg.Endpoints[0], grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("new grpc pilot client error: %v", err)
	}
	ads := xds.NewAggregatedDiscoveryServiceClient(conn)
	eds := envoy_api.NewEndpointDiscoveryServiceClient(conn)

	return &pilotClient{rawConn: conn,
		adsConn: ads, edsConn: eds,
	}, nil
}

func (c *pilotClient) GetAllRouteConfigurations() (*envoy_api.RouteConfiguration, error) {
	// TODO: this RDS stream can be reuse in all RDS request?
	rds, err := c.adsConn.StreamAggregatedResources(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[RDS] stream error: %v", err)
	}

	nodeID := util.BuildNodeID()
	err = rds.Send(&envoy_api.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node: &envoy_api_core.Node{
			Id: nodeID,
		},
		ResourceNames: []string{util.CDSHttpProxy},
		TypeUrl:       util.ClusterType})
	if err != nil {
		return nil, fmt.Errorf("[RDS] send req error for %s(%s): %v", util.RDSHttpProxy, nodeID, err)
	}

	res, err := rds.Recv()
	/*for e, value := range res.GetResources() {

		cla := &envoy_api.Cluster{}
		_  = cla.Unmarshal(value.Value)
		//lager.Logger.Info(fmt.Sprintf("%v %v", e,cla))
		//lager.Logger.Info(fmt.Sprintf("Raj : printing each resources %v %v", e, string(value.Value)))
		//lager.Logger.Info(fmt.Sprintf("####################################"))

	}*/
	//lager.Logger.Info(fmt.Sprintf("Raj : The response is %v and error is  %v", res, err))
	//result, err := json.Marshal(res)
	//lager.Logger.Info(fmt.Sprintf("Raj : The response after marshalling %v and err %v",string(result), err ))


	if err != nil {
		return nil, fmt.Errorf("[RDS] recv error for %s(%s): %v", util.RDSHttpProxy, nodeID, err)
	}
	return GetRouteConfiguration(res)
}

func (c *pilotClient) GetRouteConfigurationsByPort(port string) (*envoy_api.RouteConfiguration, error) {
	// TODO: this RDS stream can be reuse in all RDS request?
	rds, err := c.adsConn.StreamAggregatedResources(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[RDS] stream error: %v", err)
	}

	nodeID := util.BuildNodeID()
	err = rds.Send(&envoy_api.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node: &envoy_api_core.Node{
			Id: nodeID,
		},
		ResourceNames: []string{port},
		TypeUrl:       util.RouteType})
	if err != nil {
		return nil, fmt.Errorf("[RDS] send req error for %s(%s): %v", util.RDSHttpProxy, nodeID, err)
	}

	res, err := rds.Recv()
	if err != nil {
		return nil, fmt.Errorf("[RDS] recv error for %s(%s): %v", util.RDSHttpProxy, nodeID, err)
	}
	return GetRouteConfiguration(res)
}


func (c *pilotClient) GetAllClusterConfigurations() ([]envoy_api.Cluster, error) {

	lager.Logger.Info(fmt.Sprintf("Raj:  Inside GetAllClusterConfiguration"))

	cds, err := c.adsConn.StreamAggregatedResources(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[CDS] stream error: %v", err)
	}

	nodeID := util.BuildNodeID()
	err = cds.Send(&envoy_api.DiscoveryRequest{
		ResponseNonce: time.Now().String(),
		Node: &envoy_api_core.Node{
			Id: nodeID,
		},
		ResourceNames: []string{""},
		TypeUrl:       util.ClusterType})
	if err != nil {
		return nil, fmt.Errorf("[CDS] send req error for %s(%s): %v", util.ClusterType, nodeID, err)
	}

	res, err := cds.Recv()

	for e, value := range res.GetResources() {

	cla := &envoy_api.Cluster{}
	_  = cla.Unmarshal(value.Value)
	lager.Logger.Info(fmt.Sprintf("%v %v", e,cla))
	//lager.Logger.Info(fmt.Sprintf("Raj : printing each resources %v %v", e, string(value.Value)))
	lager.Logger.Info(fmt.Sprintf("####################################"))

}
	lager.Logger.Info(fmt.Sprintf("Raj : The response is %v and error is  %v", res, err))
	result, err := json.Marshal(res)
	lager.Logger.Info(fmt.Sprintf("Raj : The response after marshalling %v and err %v",string(result), err ))

	if err != nil {
		return nil, fmt.Errorf("[CDS] recv error for %s(%s): %v", util.ClusterType, nodeID, err)
	}
	return GetClusterConfiguration(res)
}
