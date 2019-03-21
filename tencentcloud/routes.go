package tencentcloud

import (
	"context"
        "fmt"

	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
        //"github.com/sirupsen/logrus"
        //"github.com/dbdd4us/qcloudapi-sdk-go/common"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

var routeTimes = 0
// ListRoutes lists all managed routes that belong to the specified clusterName
func (cloud *Cloud) ListRoutes(ctx context.Context, clusterName string) ([]*cloudprovider.Route, error) {
        //routeTimes++
        //logger := logrus.New()
	//logger.SetLevel(logrus.ErrorLevel)
        //client, err := ccs.NewClient(common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey}, common.Opts{Logger: logger, Region: cloud.config.Region})
	//if err != nil {
	//	return []*cloudprovider.Route{}, err
	//}
        //fmt.Println("come in ListRoutes times", routeTimes, "cloud config:", cloud.config)
        //cloudRoutes, err := client.DescribeClusterRoute(&ccs.DescribeClusterRouteArgs{
	//		RouteTableName: cloud.config.ClusterRouteTable,
        //})
	cloudRoutes, err := cloud.ccs.DescribeClusterRoute(&ccs.DescribeClusterRouteArgs{RouteTableName: cloud.config.ClusterRouteTable})
	if err != nil {
        //fmt.Println("come in ListRoutes get Routes err:", err)
		return []*cloudprovider.Route{}, err
	}
        //fmt.Println("After ListRoutes cloudRoutes:", clusterName)

/*
        if len(cloudRoutes.Data.RouteSet) == 0 {
			fmt.Println("No route found")
			return []*cloudprovider.Route{},nil
        }
*/
	routes := make([]*cloudprovider.Route, len(cloudRoutes.Data.RouteSet))

	for idx, route := range cloudRoutes.Data.RouteSet {
		routes[idx] = &cloudprovider.Route{Name: route.GatewayIp, TargetNode: types.NodeName(route.GatewayIp), DestinationCIDR: route.DestinationCidrBlock}
	}
	return routes, nil
}

// CreateRoute creates the described managed route
// route.Name will be ignored, although the cloud-provider may use nameHint
// to create a more user-meaningful name.
func (cloud *Cloud) CreateRoute(ctx context.Context, clusterName string, nameHint string, route *cloudprovider.Route) error {
	fmt.Println("come in CreateRoute route:",*route)
	_, err := cloud.ccs.CreateClusterRoute(&ccs.CreateClusterRouteArgs{
		RouteTableName:       cloud.config.ClusterRouteTable,
		GatewayIp:            string(route.TargetNode),
		DestinationCidrBlock: route.DestinationCIDR,
	})

	return err
}

// DeleteRoute deletes the specified managed route
// Route should be as returned by ListRoutes
func (cloud *Cloud) DeleteRoute(ctx context.Context, clusterName string, route *cloudprovider.Route) error {
	_, err := cloud.ccs.DeleteClusterRoute(&ccs.DeleteClusterRouteArgs{
		RouteTableName:       cloud.config.ClusterRouteTable,
		GatewayIp:            string(route.TargetNode),
		DestinationCidrBlock: route.DestinationCIDR,
	})
	return err
}
