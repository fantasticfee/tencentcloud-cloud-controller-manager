package tencentcloud

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func listRoutes(config Config) (*vpc.DescribeRouteTablesResponse, error) {
	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, config.Region, cpf)

	request := vpc.NewDescribeRouteTablesRequest()

	//params := `{"RouteTableIds":["rtb-oo956y0f"]}`
	params := `{"RouteTableIds":["` + config.ClusterRouteTableId + `"]}`
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.DescribeRouteTables(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	return response, nil
}

func createRoutes(config Config, route *cloudprovider.Route) (*vpc.CreateRoutesResponse, error) {
	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, config.Region, cpf)

	request := vpc.NewCreateRoutesRequest()

	params := `{"RouteTableId":"` + config.ClusterRouteTableId + `","Routes":[{"DestinationCidrBlock":"` + route.DestinationCIDR + `","GatewayType":"NORMAL_CVM","GatewayId":"` + string(route.TargetNode) + `","Enabled":"true","RouteType":"U"}]}`
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.CreateRoutes(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

func deleteRoutes(config Config, route *cloudprovider.Route) (*vpc.DeleteRoutesResponse, error) {
	credential := common.NewCredential(
		config.SecretId,
		config.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "vpc.tencentcloudapi.com"
	client, _ := vpc.NewClient(credential, config.Region, cpf)
	request := vpc.NewDeleteRoutesRequest()

	//params := `{"RouteTableId":"rtb-oo956y0f","Routes":[{"DestinationCidrBlock":"10.11.21.0/24","GatewayType":"NORMAL_CVM","GatewayId":"192.168.252.3","Enabled":"TRUE","RouteType":"U"}]}`
	params := `{"RouteTableId":"` + config.ClusterRouteTableId + `","Routes":[{"DestinationCidrBlock":"` + route.DestinationCIDR + `","GatewayType":"NORMAL_CVM","GatewayId":"` + string(route.TargetNode) + `","Enabled":"true","RouteType":"U"}]}`
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.DeleteRoutes(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ListRoutes lists all managed routes that belong to the specified clusterName
func (cloud *Cloud) ListRoutes(ctx context.Context, clusterName string) ([]*cloudprovider.Route, error) {
	/*
		cloudRoutes, err := cloud.ccs.DescribeClusterRoute(&ccs.DescribeClusterRouteArgs{RouteTableName: cloud.config.ClusterRouteTable})
		if err != nil {
			return []*cloudprovider.Route{}, err
		}*/

	cloudRoutes, err := listRoutes(cloud.config)
	if err != nil {
		return []*cloudprovider.Route{}, err
	}
	routes := make([]*cloudprovider.Route, len(cloudRoutes.Response.RouteTableSet[0].RouteSet))

	for idx, route := range cloudRoutes.Response.RouteTableSet[0].RouteSet {
		routes[idx] = &cloudprovider.Route{Name: *route.GatewayId, TargetNode: types.NodeName(*route.GatewayId), DestinationCIDR: *route.DestinationCidrBlock}
	}
	return routes, nil
}

// CreateRoute creates the described managed route
// route.Name will be ignored, although the cloud-provider may use nameHint
// to create a more user-meaningful name.
func (cloud *Cloud) CreateRoute(ctx context.Context, clusterName string, nameHint string, route *cloudprovider.Route) error {
	/*
		_, err := cloud.ccs.CreateClusterRoute(&ccs.CreateClusterRouteArgs{
			RouteTableName:       cloud.config.ClusterRouteTable,
			GatewayIp:            string(route.TargetNode),
			DestinationCidrBlock: route.DestinationCIDR,
		})
	*/
	_, err := createRoutes(cloud.config, route)

	return err
}

// DeleteRoute deletes the specified managed route
// Route should be as returned by ListRoutes
func (cloud *Cloud) DeleteRoute(ctx context.Context, clusterName string, route *cloudprovider.Route) error {
	/*
		_, err := cloud.ccs.DeleteClusterRoute(&ccs.DeleteClusterRouteArgs{
			RouteTableName:       cloud.config.ClusterRouteTable,
			GatewayIp:            string(route.TargetNode),
			DestinationCidrBlock: route.DestinationCIDR,
		})
	*/
	_, err := deleteRoutes(cloud.config, route)
	return err
}
