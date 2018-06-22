package main

import (
	"flag"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"fmt"
	"os"
)

// optional - local kubeconfig for testing
var kubeconfig = flag.String("kubeconfig", "", "Path to a kubeconfig file")

func main() {

	// send logs to stderr so we can use 'kubectl logs'
	flag.Set("logtostderr", "true")
	flag.Set("v", "3")
	flag.Parse()

	//config, err := getConfig(*kubeconfig)
	//if err != nil {
	//	glog.Errorf("Failed to load client config: %v", err)
	//	return
	//}

	restConfig := &rest.Config{
		//Host:"10.200.38.101:32608",
		Host:"10.200.38.101:32608",
		//Host:"10.244.3.12",
		//Host:"k8s04.dc.betfavorit.cf",
		BearerToken:"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InNhLWlsaWstY2x1c3Rlci1hZG1pbi10b2tlbi1mcmRybCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJzYS1pbGlrLWNsdXN0ZXItYWRtaW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjNTNiN2ZlMS0wN2VlLTExZTgtOTc1Yi0wMDFhNGExNjAxMDAiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDpzYS1pbGlrLWNsdXN0ZXItYWRtaW4ifQ.Y7eXC-VwsF5c0Qm2w8jj3nql9FB7pwNXdjg4xT48hxHZqz_yOLGzDYBhnqBXbLhkFLx5g52XSU_3o21sNirfzvOoR-FZUc6tEuCloRDVJ5CLDvj40hRH0Jy7m1T3tz-u_WeqGQXfGu5yQExbm8JEVa8z2XV6bQVLj6zs2J2V-FZldqF4-Uc7CTokf8ov09SdeqjCWPkQgiRaLVM5HivPkwhBomhdQHVwBhWpFszvjRkkEkNqyLc8pqu1Yjhuw4PRQ-c6euuJeDG1Oo0lfPaRafxrzclGL3xfe2GSVbxt1vepKtNkWybW7dOgvGEjkP8HoFvJ95hXBL-zbv4Wc59aSA",
	}


	// build the Kubernetes client
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		glog.Errorf("Failed to create kubernetes client: %v", err)
		return
	}

	au := client.AuthorizationV1beta1()

	fmt.Println(au.LocalSubjectAccessReviews(""))
	fmt.Println(au.RESTClient().APIVersion())

	fmt.Println(au)

	poods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})

	//podds, err := client.CoreV1().Services("deploy").List(metav1.ListOptions{}) // Pods("deploy")
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
	fmt.Println(poods)
	os.Exit(0)


	//
	//z, err := client.AppsV1().Deployments("").List(metav1.ListOptions{})
	//if err !=nil {
	//	fmt.Println("ERROR: ", err)
	//}
	//
	//fmt.Println(z)


	//v1 := client.AuthenticationV1()
	//z := v1.TokenReviews()
	//z= z


	// list pods
	pods, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		glog.Errorf("Failed to retrieve pods: %v", err)
		return
	}

	for _, p := range pods.Items {
		glog.V(3).Infof("Found pods: %s/%s", p.Namespace, p.Name)
	}
}

func getConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return rest.InClusterConfig()
}

func z() {
	//clientcmd.AuthLoader()
}