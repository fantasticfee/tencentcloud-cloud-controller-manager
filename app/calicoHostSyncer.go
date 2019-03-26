package app

import (
	"github.com/projectcalico/felix/calc"
	"github.com/projectcalico/libcalico-go/lib/apiconfig"
	bapi "github.com/projectcalico/libcalico-go/lib/backend/api"
	glog "github.com/sirupsen/logrus"

	globalCfg "github.com/tencentcloud/tencentcloud-cloud-controller-manager/config"

	apiv3 "github.com/projectcalico/libcalico-go/lib/apis/v3"
	"github.com/projectcalico/libcalico-go/lib/backend/model"
	"github.com/projectcalico/libcalico-go/lib/backend/watchersyncer"

	"os"
)

type Startable interface {
	Start()
}

func RunCalicoHostSyncer(trigger chan struct{}) {
	if globalCfg.Cfg.Global.EtcdEndpoints == "" {
		globalCfg.Cfg.Global.EtcdEndpoints = "http://115.236.185.190:23799"
	}
	glog.Infof("edcdendpoint %v", globalCfg.Cfg.Global.EtcdEndpoints)

	cfg := &apiconfig.EtcdConfig{
		EtcdEndpoints:  globalCfg.Cfg.Global.EtcdEndpoints,
		EtcdKeyFile:    globalCfg.Cfg.Global.EtcdKeyFile,
		EtcdCertFile:   globalCfg.Cfg.Global.EtcdCertFile,
		EtcdCACertFile: globalCfg.Cfg.Global.EtcdCACertFile,
	}
	var backendClient bapi.Client
	var err error
	glog.SetOutput(os.Stdout)
	glog.SetLevel(glog.TraceLevel)

	if backendClient, err = NewHostClient(cfg); err != nil {
		glog.Errorf("err new client %v", err)
	}

	/*if backendClient, err = etcdv3.NewEtcdV3Client(cfg);err != nil {
		glog.Errorf("err new client %v",err)
	}*/

	var s Startable

	syncerToHandler := calc.NewSyncerCallbacksDecoupler()
	resourceTypes := []watchersyncer.ResourceType{
		{
			ListInterface: model.IpamResourceListOptions{Kind: apiv3.KindClusterInformation},
			//todo 实现
			UpdateProcessor: nil, //updateprocessors.NewClusterInfoUpdateProcessor(),
		},
	}

	s = watchersyncer.New(
		backendClient,
		resourceTypes,
		syncerToHandler,
	)
	s.Start()

	handler := NewHandler(trigger)
	go syncerToHandler.SendTo(handler)
}
