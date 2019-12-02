package ddosmanager

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
    "sort"

	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"

	clientv3 "github.com/coreos/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
	pcn_controllers "github.com/polycube-network/polycube/src/components/k8s/pcn_k8s/controllers"
	pcn_types "github.com/polycube-network/polycube/src/components/k8s/pcn_k8s/types"
	utils "github.com/polycube-network/polycube/src/components/k8s/utils"
	k8sddos "github.com/polycube-network/polycube/src/components/k8s/utils/k8sddos"
)

type DdosMitigator struct {
	Name         string
	BlacklistSrc []string
	BlacklistDst []string
	Dropcnt      uint64
}

type DdosMitigatorManager struct {
	localDM   map[string]DdosMitigator
	etcd_cli  *clientv3.Client
	dm_enable map[string]bool //dm_enalbe[default|"rel"="product","type"="security"] = enable
}

const (
	basePath              = "http://127.0.0.1:9000/polycube/v1"
	EtcdURLDefault string = "http://127.0.0.1:30901"
    baseCfgPath           = "/config/dfw/ddosmitigator/"
)

var k8sDdosAPI *k8sddos.DdosmitigatorApiService

func (m *DdosMitigatorManager) GetDdosAPI() *k8sddos.DdosmitigatorApiService {
	if k8sDdosAPI == nil {
		cfgK8sddos := k8sddos.Configuration{BasePath: basePath}
		srK8sddos := k8sddos.NewAPIClient(&cfgK8sddos)
		k8sDdosAPI = srK8sddos.DdosmitigatorApi

		if k8sDdosAPI == nil {
			panic("failed to create k8sDdosAPI")
		}
	}

	return k8sDdosAPI
}

func (m *DdosMitigatorManager) CreateDdosMitigator(name string) error {
	ddosApi := m.GetDdosAPI()
	if ddosApi == nil {
		log.Error("ddosApi is nil")
		return errors.New("ddosApi is nil")
	}

	log.Debugf("create ddosmitigator %s", name)
	if _, err := ddosApi.CreateDdosmitigatorByID(nil, name, k8sddos.Ddosmitigator{
		Name: name,
	}); err != nil {
		log.Errorf("create ddos %s failed, error: %s", name, err.Error())
		return err
	}

	return nil
}

func StartDdosMitigatorManager(node string) *DdosMitigatorManager {
	var err error
	manager := DdosMitigatorManager{
		localDM   : make(map[string]DdosMitigator),
        dm_enable : make(map[string]bool),
	}

	if manager.etcd_cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{EtcdURLDefault},
		DialTimeout: 10 * time.Second,
	}); err != nil {
		return nil
	}

	pcn_controllers.Pods().Subscribe(pcn_types.Update, nil, nil,
		&pcn_types.ObjectQuery{Name: node}, pcn_types.PodRunning, manager.PodUpdate)
	pcn_controllers.Pods().Subscribe(pcn_types.Delete, nil, nil,
		&pcn_types.ObjectQuery{Name: node}, pcn_types.PodAnyPhase, manager.PodDelete)

	return &manager
}

func (m *DdosMitigatorManager) PodUpdate(new, old *core_v1.Pod) {
	var new_name, old_name string
	if new != nil {
		new_name = new.Name
	} else {
		new_name = "nil"
	}
	if old != nil {
		old_name = old.Name
	} else {
		old_name = "nil"
	}
	log.Debugf("PodUpdate new pod: %s, old pod: %s", new_name, old_name)

	if m.CheckIfDMEnable(new) {
		name := "dm-" + new.Status.PodIP
		m.CreateDdosMitigator(name)
		dm := DdosMitigator{
			Name: name,
		}
		m.localDM[name] = dm
		log.Debugf("add dm %s", name)
	} else {
		log.Debug("delete dm")
	}
}

func (m *DdosMitigatorManager) PodDelete(new, old *core_v1.Pod) {
	var new_name, old_name string
	if new != nil {
		new_name = new.Name
	} else {
		new_name = "nil"
	}
	if old != nil {
		old_name = old.Name
	} else {
		old_name = "nil"
	}
	log.Debugf("PodDelete new pod: %s, old pod: %s", new_name, old_name)
}

func (m *DdosMitigatorManager) CheckIfDMEnable(pod *core_v1.Pod) bool {
    pod_ns := pod.Namespace
	pod_label := make(map[string]string)
	for k, v := range pod.Labels {
		pod_label[k] = v
	}

	//find the first label key string (already ordered) contained by pod label
	for labels, v := range m.dm_enable {
        str := strings.Split(labels, "|")
        ns := str[0]
        labels := str[1]

        if ns != pod_ns {
            continue
        }
		label_map := make(map[string]string)

		// convert the labels string to map
		label_set := strings.Split(labels, ",")
		for _, label := range label_set {
			pair := strings.Split(label, "=")
			label_map[pair[0]] = pair[1]
		}

		if utils.AreLabelsContained(label_map, pod_label) {
			return v
		}
	}

	return false
}

func (dmm DdosMitigatorManager) WatchDB() {
	rch := dmm.etcd_cli.Watch(context.Background(), baseCfgPath, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
            log.Debugf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
            if !strings.HasPrefix(string(ev.Kv.Key), baseCfgPath) {
                log.Warnf("key %q doesn't start with %q", ev.Kv.Key, baseCfgPath)
                return
            }
            str := strings.Split(strings.TrimPrefix(string(ev.Kv.Key), baseCfgPath), "/")
            ns := str[0]
            labels := strings.Split(str[1], ",")
            sort.Strings(labels)
            key := ns + "|"
            for _, label := range labels {
                key += label
            }
            log.Debugf("key: %q, value: %q", key, ev.Kv.Value)
			switch ev.Type {
			case mvccpb.PUT:
				log.Debug("put event")
                if _, ok := dmm.dm_enable[key]; !ok {
                    if string(ev.Kv.Value) == "enable" {
                        dmm.dm_enable[key] = true
                    } else {
                        dmm.dm_enable[key] = false
                    }
                } else {
                    log.Debugf("key %q exists, old value: %q, new value %q", key, dmm.dm_enable[key], ev.Kv.Value)
                }
				break
			case mvccpb.DELETE:
				fmt.Println("delete event")
                if _, ok := dmm.dm_enable[key]; !ok {
                    delete(dmm.dm_enable, key)
                }
				break
			default:
				fmt.Printf("invalid event type %d", ev.Type)
				break
			}
		}
	}
}
