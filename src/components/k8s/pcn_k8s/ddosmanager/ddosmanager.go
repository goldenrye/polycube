package ddosmanager

import (
	"errors"

	uuid "github.com/google/uuid"
	k8sddos "github.com/polycube-network/polycube/src/components/k8s/utils/k8sddos"
	log "github.com/sirupsen/logrus"
	core_v1 "k8s.io/api/core/v1"
)

type DdosManager struct {
	pod *core_v1.Pod
}

const (
	basePath = "http://127.0.0.1:9000/polycube/v1"
)

var k8sDdosAPI *k8sddos.DdosmitigatorApiService

func NewDdosManager(pod *core_v1.Pod) *DdosManager {
	d := DdosManager{
		pod: pod,
	}
	return &d
}

func (m *DdosManager) GetDdosAPI() *k8sddos.DdosmitigatorApiService {
	if k8sDdosAPI == nil {
		cfgK8sddos := k8sddos.Configuration{BasePath: basePath}
		srK8sddos := k8sddos.NewAPIClient(&cfgK8sddos)
		k8sDdosAPI = srK8sddos.DdosmitigatorApi
	}
	return k8sDdosAPI
}

func (m *DdosManager) CreateDdosMitigator() error {
	ddosApi := m.GetDdosAPI()
	if ddosApi == nil {
		log.Error("ddosApi is nil")
		return errors.New("ddosApi is nil")
	}

	name := m.pod.Namespace + m.pod.Name
	if _, err := ddosApi.CreateDdosmitigatorByID(nil, name, k8sddos.Ddosmitigator{
		Name: name,
		Uuid: uuid.New().String(),
	}); err != nil {
		log.Errorf("create ddos %s failed")
		return errors.New("create ddos failed")
	}

	return nil
}
