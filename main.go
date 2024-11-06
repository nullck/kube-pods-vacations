package main

import (
	"fmt"
	"os"

	"github.com/nullck/kube-pods-vacations/cmd/kube_pods_vacations"
)

func main() {
	/*namespaceNames := []string{"default"}
	for _, n := range namespaceNames {
		r := kubernetesmgmt.KubeMgmt{
			NamespaceName:    n,
			AnnotationPrefix: "kube-pods-vacations",
		}
		r.DeploymentAnnotationsBuilder()
		time.Sleep(30 * time.Second)
	}
	*/
	cmd := kube_pods_vacations.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error running kube_pods_vacations %v", err)
		os.Exit(1)
	}
}
