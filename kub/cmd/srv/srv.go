package main

import (
	"fmt"
	"os"

	"time"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/clientKub"
	"git.betfavorit.cf/vadim.tsurkov/kuberweb/kub/kubService"
)

//this is test local
func main() {
	var ks *kubService.ServiceKubernetes
	ks = kubService.InitInstance(
		clientKub.NewRestClient(
			"https://10.200.38.101:32608",
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6InNhLWlsaWstY2x1c3Rlci1hZG1pbi10b2tlbi1mcmRybCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50Lm5hbWUiOiJzYS1pbGlrLWNsdXN0ZXItYWRtaW4iLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiJjNTNiN2ZlMS0wN2VlLTExZTgtOTc1Yi0wMDFhNGExNjAxMDAiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6ZGVmYXVsdDpzYS1pbGlrLWNsdXN0ZXItYWRtaW4ifQ.Y7eXC-VwsF5c0Qm2w8jj3nql9FB7pwNXdjg4xT48hxHZqz_yOLGzDYBhnqBXbLhkFLx5g52XSU_3o21sNirfzvOoR-FZUc6tEuCloRDVJ5CLDvj40hRH0Jy7m1T3tz-u_WeqGQXfGu5yQExbm8JEVa8z2XV6bQVLj6zs2J2V-FZldqF4-Uc7CTokf8ov09SdeqjCWPkQgiRaLVM5HivPkwhBomhdQHVwBhWpFszvjRkkEkNqyLc8pqu1Yjhuw4PRQ-c6euuJeDG1Oo0lfPaRafxrzclGL3xfe2GSVbxt1vepKtNkWybW7dOgvGEjkP8HoFvJ95hXBL-zbv4Wc59aSA",
			true),
	)
	go ks.Start()

	pods, err := ks.GetPods()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print((pods))

	{
		{
			{
				time.Sleep(100 * time.Second)
			}
		}
	}

}
