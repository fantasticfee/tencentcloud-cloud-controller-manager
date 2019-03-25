package app

import (
	"strings"

	"github.com/cornelk/hashmap"
	"github.com/projectcalico/libcalico-go/lib/backend/api"
	"github.com/projectcalico/libcalico-go/lib/backend/model"
	log "github.com/sirupsen/logrus"
)

var IpamHostMap2Block hashmap.HashMap = hashmap.HashMap{}

type handler struct {
	triggerCh chan struct{}
}

func NewHandler(triggerCh chan struct{}) *handler {
	return &handler{
		triggerCh,
	}
}

func (v *handler) OnStatusUpdated(status api.SyncStatus) {
	// Pass through.
	log.Debug("handler OnStatusUpdated")
}

func (v *handler) OnUpdates(updates []api.Update) {
	log.Debug("handler OnUpdates")
	for _, update := range updates {
		if resourceKey, ok := update.Key.(model.ResourceKey); ok {
			if update.UpdateType == api.UpdateTypeKVNew {
				///nodeMap2Block[update.Key]
				params := strings.Split(resourceKey.Name, ":")
				if len(params) == 2 {
					if value, exist := IpamHostMap2Block.Get(params[0]); exist {
						if block, ok := value.([]string); ok {
							value = append(block, strings.Replace(params[1], "-", "/", 1))
						}
					} else {
						IpamHostMap2Block.Set(params[0], []string{strings.Replace(params[1], "-", "/", 1)})
					}
				}

				log.Infof("handler OnUpdates update %v", update.Key.String())
			} else if update.UpdateType == api.UpdateTypeKVDeleted {
				log.Infof("handler OnUpdates delete %v", update.Key.String())
				params := strings.Split(resourceKey.Name, ":")
				if len(params) == 2 {
					if value, exist := IpamHostMap2Block.Get(params[0]); exist {
						var temp []string
						if block, ok := value.([]string); ok {
							for _, v := range block {
								if v != params[1] {
									temp = append(temp, v)
								}
							}

						}
						IpamHostMap2Block.Set(params[0], temp)
					}
				}
			}
			if len(v.triggerCh) == 0 {
				v.triggerCh <- struct{}{}
			}
		}
	}
}
