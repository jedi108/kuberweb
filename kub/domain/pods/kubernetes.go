package pods

//easyjson
type ApiPod struct {
	Status Status `json:"status"`
	Pods   []Pods `json:"pods"`
}

//easyjson
type Status struct {
	Failed    int `json:"failed"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Succeeded int `json:"succeeded"`
}

//easyjson
type Pods struct {
	NodeName   string     `json:"nodeName"`
	ObjectMeta ObjectMeta `json:"objectMeta"`
}

//easyjson
type ObjectMeta struct {
	Name string `json:"name"`
}

type CumulativeMetrics struct{}

//Overview
//type ConfigMapList struct {}
//type CronJobList struct {}
//type DaemonSetList struct {}
//type DeploymentList struct {}
//type IngressList struct {}
//type JobList struct {}
//type PodList struct {}
//type ReplicaSetList struct {}
//type ReplicationControllerList struct {}
//type SecretList struct {}
//type ServiceList struct {}
//type StatefulSetList struct {}
type Errors struct{}
