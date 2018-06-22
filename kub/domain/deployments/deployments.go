package deployments

import "git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/domain/pods"

//easyjson
type Pods struct {
	Current uint `json:"current"`
	Running uint `json:"running"`
}

//easyjson
type ObjectMeta struct {
	Name string `json:"name"`
}

//easyjson
type Deployment struct {
	ObjectMeta ObjectMeta `json:"objectMeta"`
	Pods       Pods       `json:"pods"`
}

//easyjson
type Deployments struct {
	Deployments []Deployment `json:"deployments"`
	Status      pods.Status  `json:"status"`
}

//easyjson
type Response struct {
	DesiredReplicas uint `json:"desiredReplicas"`
	ActualReplicas  uint `json:"actualReplicas"`
}
