package app

import (
	apiv3 "github.com/projectcalico/libcalico-go/lib/apis/v3"

	"github.com/projectcalico/libcalico-go/lib/backend/api"
	"github.com/projectcalico/libcalico-go/lib/backend/model"
	"github.com/projectcalico/libcalico-go/lib/backend/watchersyncer"
)

func New(client api.Client, callbacks api.SyncerCallbacks) api.Syncer {

	resourceTypes := []watchersyncer.ResourceType{
		{
			ListInterface: model.IpamResourceListOptions{Kind: apiv3.KindClusterInformation},
			//todo 实现
			UpdateProcessor: nil, //updateprocessors.NewClusterInfoUpdateProcessor(),
		},
	}

	return watchersyncer.New(
		client,
		resourceTypes,
		callbacks,
	)
}
