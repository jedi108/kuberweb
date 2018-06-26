package kubService

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"git.betfavorit.cf/backend/logger"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/clientKub"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/deployments"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/pods"
)

const (
	timeRefreshToken = 5
)

var (
	ch          = make(chan struct{})
	initialized uint32
	once        sync.Once
	instance    *ServiceKubernetes
)

type ServiceKubernetes struct {
	restClient *clientKub.RestClient
}

//reload auth if error data received
func ReloadAuth() {
	atomic.StoreUint32(&initialized, 2)
	fmt.Println("initialize:", atomic.LoadUint32(&initialized))
}

func InitInstance(restClient *clientKub.RestClient) *ServiceKubernetes {
	once.Do(func() {
		instance = &ServiceKubernetes{
			restClient: restClient,
		}
	})
	return instance
}

func GetInstance() *ServiceKubernetes {
	return instance
}

func (sk *ServiceKubernetes) GetRestClient() *clientKub.RestClient {
	return sk.restClient
}

func (sk *ServiceKubernetes) Start() {
	for {
		sk.auth()
		for {
			err := sk.refresh()
			if err != nil {
				break
			}

			if atomic.LoadUint32(&initialized) == 0 {
				ch <- struct{}{}
				atomic.StoreUint32(&initialized, 1)
			}
			if atomic.LoadUint32(&initialized) == 2 {
				fmt.Println("break initialize:", atomic.LoadUint32(&initialized))
				logger.Debug("break in service kubernetis", atomic.LoadUint32(&initialized))
				break
			}
			atomic.StoreUint32(&initialized, 1)
			time.Sleep(time.Minute * timeRefreshToken)
		}
		logger.Info("restart auth")
	}
	logger.Fatalf("exit from service kubernetis start")
}

func (sk *ServiceKubernetes) refresh() error {
	_, err := sk.restClient.Status()
	if err != nil {
		return err
	}

	err = sk.restClient.CsrfToken()
	if err != nil {
		return err
	}
	err = sk.restClient.UpdateRefreshToken()
	return err
}

func (sk *ServiceKubernetes) auth() {
	csrfToken, err := sk.restClient.CsrfLogin()
	if err != nil {
		logger.Errorf("auth error CsrfLogin %v", err)
	}
	logger.Debug("rest client auth")

	err = sk.restClient.Login(csrfToken)
	if err != nil {
		logger.Errorf("auth error Login %v", err)
	}
	logger.Debug("rest client login csrf")
}

func (sk *ServiceKubernetes) GetPods() (*pods.ApiPod, error) {
	if atomic.LoadUint32(&initialized) == 0 {
		<-ch
	}
	var apiPod *pods.ApiPod

	bufByte, err := sk.restClient.Pod("deploy")
	if err != nil {
		logger.Errorf("request getPods failed: %v", err)
		return apiPod, err
	}

	apiPod = &pods.ApiPod{}
	err = apiPod.UnmarshalJSON(bufByte)
	return apiPod, err
}

func (sk *ServiceKubernetes) GetDeployments() (*deployments.Deployments, error) {
	if atomic.LoadUint32(&initialized) == 0 {
		<-ch
	}
	var depls *deployments.Deployments

	bufByte, err := sk.restClient.Deployment("deploy")
	if err != nil {
		logger.Errorf("request getDeployments failed: %v", err)
		return depls, err
	}

	depls = &deployments.Deployments{}
	err = depls.UnmarshalJSON(bufByte)
	return depls, err
}

func (sk *ServiceKubernetes) ScaleBy(nameDep string, scaleBy uint64) (*deployments.Response, error) {
	if atomic.LoadUint32(&initialized) == 0 {
		<-ch
	}
	var depResp *deployments.Response
	bufByte, err := sk.restClient.Scale(nameDep, scaleBy)
	if err != nil {
		logger.Errorf("request scaleBy failed: %v", err)
		return depResp, err
	}

	depResp = &deployments.Response{}
	err = depResp.UnmarshalJSON(bufByte)
	return depResp, err
}
