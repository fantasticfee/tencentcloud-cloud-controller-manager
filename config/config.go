package config

import "os"

// CloudConfig wraps the settings for the Alicloud provider.
type CloudConfig struct {
	Global struct {
		KubernetesClusterTag string
		//UID                  string `json:"uid"`
		//VpcID                string `json:"vpcid"`
		//Region               string `json:"region"`
		//ZoneID               string `json:"zoneid"`
		//VswitchID            string `json:"vswitchid"`

		AccessKeyID     string `json:"accessKeyID"`
		AccessKeySecret string `json:"accessKeySecret"`
		EtcdEndpoints   string `json:"etcdEndpoints"`
		EtcdKeyFile     string `json:"etcdKeyFile"`
		EtcdCertFile    string `json:"etcdCertFile"`
		EtcdCACertFile  string `json:"etcdCACertFile"`
	}
}

var Cfg CloudConfig

func init() {
	Cfg.Global.AccessKeyID = os.Getenv("ACCESS_KEY_ID")
	Cfg.Global.AccessKeySecret = os.Getenv("ACCESS_KEY_SECRET")
	Cfg.Global.EtcdEndpoints = os.Getenv("ETCD")
	Cfg.Global.EtcdKeyFile = os.Getenv("ETCD_KEY_FILE")
	Cfg.Global.EtcdCertFile = os.Getenv("ETCD_CERT_FILE")
	Cfg.Global.EtcdCACertFile = os.Getenv("ETCD_CA_CERT_FILE")
}
