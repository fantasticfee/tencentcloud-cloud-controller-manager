package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/projectcalico/libcalico-go/lib/apiconfig"
	"github.com/projectcalico/libcalico-go/lib/backend/api"
	"github.com/projectcalico/libcalico-go/lib/backend/model"
	cerrors "github.com/projectcalico/libcalico-go/lib/errors"
)

var (
	clientTimeout    = 10 * time.Second
	keepaliveTime    = 30 * time.Second
	keepaliveTimeout = 10 * time.Second
)

type HostClient struct {
	etcdClient *clientv3.Client
}

func NewHostClient(config *apiconfig.EtcdConfig) (api.Client, error) {
	// Split the endpoints into a location slice.
	etcdLocation := []string{}
	if config.EtcdEndpoints != "" {
		etcdLocation = strings.Split(config.EtcdEndpoints, ",")
	}

	if len(etcdLocation) == 0 {
		log.Warning("No etcd endpoints specified in etcdv3 API config")
		return nil, errors.New("no etcd endpoints specified")
	}

	// Create the etcd client
	tlsInfo := &transport.TLSInfo{
		CAFile:   config.EtcdCACertFile,
		CertFile: config.EtcdCertFile,
		KeyFile:  config.EtcdKeyFile,
	}

	tls, err := tlsInfo.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not initialize etcdv3 client: %+v", err)
	}

	// Build the etcdv3 config.
	cfg := clientv3.Config{
		Endpoints:            etcdLocation,
		TLS:                  tls,
		DialTimeout:          clientTimeout,
		DialKeepAliveTime:    keepaliveTime,
		DialKeepAliveTimeout: keepaliveTimeout,
	}

	// Plumb through the username and password if both are configured.
	if config.EtcdUsername != "" && config.EtcdPassword != "" {
		cfg.Username = config.EtcdUsername
		cfg.Password = config.EtcdPassword
	}

	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &HostClient{etcdClient: client}, nil
}

func (c *HostClient) Create(ctx context.Context, d *model.KVPair) (*model.KVPair, error) {

	return nil, nil
}

func (c *HostClient) Update(ctx context.Context, d *model.KVPair) (*model.KVPair, error) {

	return nil, nil
}

func (c *HostClient) Apply(ctx context.Context, d *model.KVPair) (*model.KVPair, error) {
	return nil, nil
}

func (c *HostClient) Delete(ctx context.Context, k model.Key, revision string) (*model.KVPair, error) {
	return nil, nil
}

func (c *HostClient) Get(ctx context.Context, k model.Key, revision string) (*model.KVPair, error) {
	return nil, nil
}

// List entries in the datastore.  This may return an empty list of there are
// no entries matching the request in the ListInterface.
func (c *HostClient) List(ctx context.Context, l model.ListInterface, revision string) (*model.KVPairList, error) {
	logCxt := log.WithFields(log.Fields{"list-interface": l, "rev": revision})
	logCxt.Debug("Processing List request")

	// To list entries, we enumerate from the common root based on the supplied
	// IDs, and then filter the results.
	key := model.ListOptionsToDefaultPathRoot(l)

	ops := []clientv3.OpOption{}

	ops = append(ops, clientv3.WithPrefix())

	logCxt = logCxt.WithField("etcdv3-etcdKey", key)

	// We may also need to perform a get based on a particular revision.
	if len(revision) != 0 {
		rev, err := parseRevision(revision)
		if err != nil {
			return nil, err
		}
		ops = append(ops, clientv3.WithRev(rev))
	}

	logCxt.Debug("Calling Get on etcdv3 client")
	resp, err := c.etcdClient.Get(ctx, key, ops...)
	if err != nil {
		logCxt.WithError(err).Debug("Error returned from etcdv3 client")
		return nil, cerrors.ErrorDatastoreError{Err: err}
	}
	logCxt.WithField("numResults", len(resp.Kvs)).Debug("Processing response from etcdv3")

	// Filter/process the results.
	list := []*model.KVPair{}
	for _, p := range resp.Kvs {
		if kv := convertListResponse(p, l); kv != nil {
			list = append(list, kv)
		}
	}

	return &model.KVPairList{
		KVPairs:  list,
		Revision: strconv.FormatInt(resp.Header.Revision, 10),
	}, nil
}

// EnsureInitialized makes sure that the etcd data is initialized for use by
// Calico.
func (c *HostClient) EnsureInitialized() error {
	//TODO - still need to worry about ready flag.
	return nil
}

// Clean removes all of the Calico data from the datastore.
func (c *HostClient) Clean() error {
	log.Warning("Cleaning etcdv3 datastore of all Calico data")
	_, err := c.etcdClient.Txn(context.Background()).If().Then(
		clientv3.OpDelete("/calico/", clientv3.WithPrefix()),
	).Commit()

	if err != nil {
		return cerrors.ErrorDatastoreError{Err: err}
	}
	return nil
}

// IsClean() returns true if there are no /calico/ prefixed entries in the
// datastore.  This is not part of the exposed API, but is public to allow
// direct consumers of the backend API to access this.
func (c *HostClient) IsClean() (bool, error) {
	log.Debug("Calling Get on etcdv3 client")
	resp, err := c.etcdClient.Get(context.Background(), "/calico/", clientv3.WithPrefix())
	if err != nil {
		log.WithError(err).Debug("Error returned from etcdv3 client")
		return false, cerrors.ErrorDatastoreError{Err: err}
	}

	// The datastore is clean if no results were enumerated.
	return len(resp.Kvs) == 0, nil
}

// getTTLOption returns a OpOption slice containing a Lease granted for the TTL.
func (c *HostClient) getTTLOption(ctx context.Context, d *model.KVPair) ([]clientv3.OpOption, error) {
	putOpts := []clientv3.OpOption{}

	if d.TTL != 0 {
		resp, err := c.etcdClient.Lease.Grant(ctx, int64(d.TTL.Seconds()))
		if err != nil {
			log.WithError(err).Error("Failed to grant a lease")
			return nil, cerrors.ErrorDatastoreError{Err: err}
		}

		putOpts = append(putOpts, clientv3.WithLease(resp.ID))
	}

	return putOpts, nil
}

// parseRevision parses the model.KVPair revision string and converts to the
// equivalent etcdv3 int64 value.
func parseRevision(revs string) (int64, error) {
	rev, err := strconv.ParseInt(revs, 10, 64)
	if err != nil {
		log.WithField("Revision", revs).Debug("Unable to parse Revision")
		return 0, cerrors.ErrorValidation{
			ErroredFields: []cerrors.ErroredField{
				{
					Name:  "ResourceVersion",
					Value: revs,
				},
			},
		}
	}
	return rev, nil
}

func convertListResponse(ekv *mvccpb.KeyValue, l model.ListInterface) *model.KVPair {
	log.WithField("etcdv3-etcdKey", string(ekv.Key)).Debug("Processing etcdv3 entry")
	if k := l.KeyFromDefaultPath(string(ekv.Key)); k != nil {
		log.WithField("model-etcdKey", k).Debug("Key is valid and converted to model-etcdKey")

		return &model.KVPair{Key: k, Value: "", Revision: strconv.FormatInt(ekv.ModRevision, 10)}
	}
	return nil
}
