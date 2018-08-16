package client

import (
	"errors"
	"fmt"

	"github.com/ServiceComb/go-chassis/pkg/istio/util"
	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
)

// GetRouteConfiguration returns routeconfiguration from discovery response
func GetRouteConfiguration(res *xdsapi.DiscoveryResponse) (*xdsapi.RouteConfiguration, error) {
	if res.TypeUrl != util.RouteType || res.Resources[0].TypeUrl != util.RouteType {
		return nil, errors.New("Invalid typeURL" + res.TypeUrl)
	}

	cla := &xdsapi.RouteConfiguration{}
	err := cla.Unmarshal(res.Resources[0].Value)
	if err != nil {
		return nil, err
	}
	return cla, nil
}

// GetClusterConfiguration returns cluster information from discovery response
func GetClusterConfiguration(res *xdsapi.DiscoveryResponse) ([]xdsapi.Cluster, error) {
	fmt.Println("Raj: valur of res.TypeUrl ", res.TypeUrl)
	if res.TypeUrl != util.ClusterType {
		fmt.Println("RaJ: return as type url not match")
		return nil, errors.New("Invalid typeURL" + res.TypeUrl)
	}

	var cluster []xdsapi.Cluster
	for _, value := range res.GetResources() {
        fmt.Println("Raj: the value of value is %v and %v", value, value.Value)
		cla := &xdsapi.Cluster{}
		err := cla.Unmarshal(value.Value)
		fmt.Println("Raj: is there a error %v", err)
		if err != nil {
            fmt.Println("Raj: calling append function")
			return nil, errors.New("unmarshall error")

		}
		cluster = append(cluster, *cla)

	}
	fmt.Println("Raj : inside getclusterconfiguration cluster%v ", cluster)
	return cluster, nil
}