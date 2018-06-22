package kubService

import (
	"fmt"
	"os"
	"time"

	"sync"

	"sync/atomic"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/clientKub"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/deployments"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/pods"
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
	sk.auth()
	for {
		sk.refresh()
		time.Sleep(time.Second * 2)
		fmt.Print(".")

		if atomic.LoadUint32(&initialized) == 0 {
			atomic.StoreUint32(&initialized, 1)
			ch <- struct{}{}
		}
	}
}

func (sk *ServiceKubernetes) refresh() {
	_, err := sk.restClient.Status()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//token
	err = sk.restClient.CsrfToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = sk.restClient.UpdateRefreshToken()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//bytes, err := sk.restClient.Pod("deploy")
	//if err != nil {
	//	fmt.Println(err)
	//}

	//bytes = bytes
	//fmt.Println(string(bytes))

}

func (sk *ServiceKubernetes) auth() {
	csrfToken, err := sk.restClient.CsrfLogin()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = sk.restClient.Login(csrfToken)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func (sk *ServiceKubernetes) GetPods() (*pods.ApiPod, error) {
	if atomic.LoadUint32(&initialized) == 0 {
		<-ch
	}
	var apiPod *pods.ApiPod

	bufByte, err := sk.restClient.Pod("deploy")
	if err != nil {
		return apiPod, err
	}

	apiPod = &pods.ApiPod{}
	err = apiPod.UnmarshalJSON(bufByte)
	return apiPod, nil
}

func (sk *ServiceKubernetes) GetDeployments() (*deployments.Deployments, error) {
	if atomic.LoadUint32(&initialized) == 0 {
		<-ch
	}
	var depls *deployments.Deployments

	bufByte, err := sk.restClient.Deployment("deploy")
	if err != nil {
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
		return depResp, err
	}

	depResp = &deployments.Response{}
	err = depResp.UnmarshalJSON(bufByte)
	return depResp, err
}
