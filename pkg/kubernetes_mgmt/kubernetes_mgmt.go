package kubernetesmgmt

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nullck/kube-pods-vacations/internal/kubernetes_cronjob"
	"github.com/rs/zerolog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string

type KubeMgmt struct {
	NamespaceName    string
	AnnotationPrefix string
	Clientset        *kubernetes.Clientset
	Logger           zerolog.Logger
}

func (kt KubeMgmt) DeploymentAnnotationsBuilder() error {
	kt.Clientset = kubeConnect()
	deploymentsClient := kt.Clientset.AppsV1().Deployments(kt.NamespaceName)
	list, err := deploymentsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	cronAnnotations := make(map[string]string)
	cronAnnotations["cron-namespace"] = kt.NamespaceName
	cronAnnotations["resource-type"] = "deployment"
	for _, d := range list.Items {
		cronAnnotations["cron-name"] = fmt.Sprintf("%s-cron", d.Name)
		cronAnnotations["resource-name"] = d.Name
		if strings.Contains(fmt.Sprint(d.Annotations), kt.AnnotationPrefix) {
			for k, v := range d.Annotations {
				if strings.Contains(fmt.Sprint(k), kt.AnnotationPrefix) {
					annotationName := strings.Replace(fmt.Sprint(k), fmt.Sprintf("%s/", kt.AnnotationPrefix), "", -1)
					cronAnnotations[annotationName] = v
				}
			}
			kt.CreateCronFromDeployment(cronAnnotations)
			return nil
		}
	}
	return nil
}

func (kt KubeMgmt) CreateCronFromDeployment(annotations map[string]string) {
	namespace := annotations["cron-namespace"]
	cronName := annotations["cron-name"]
	reduceRepCronExpr := annotations["reduced-cron-expr"]
	desiredRepCronExpr := annotations["desired-cron-expr"]
	resourceType := annotations["resource-type"]
	resourceName := annotations["resource-name"]
	minReplicas := annotations["reduced-replicas"]
	desReplicas := annotations["desired-replicas"]

	kubeCron := kubernetes_cronjob.NewKubeCronJob(namespace, cronName, reduceRepCronExpr, desiredRepCronExpr, resourceType, resourceName, minReplicas, desReplicas, kt.Logger)
	err := kubeCron.CreateReduceCronJob(kt.Clientset)

	if err != nil {
		kt.Logger.Info().Msgf("Error creating reduce cronjob %v ", err)
	}
	err = kubeCron.CreateDesiredCronJob(kt.Clientset)
	if err != nil {
		kt.Logger.Error().Msgf("Error creating reduce cronjob %v ", err)
	}
}

func kubeConnect() *kubernetes.Clientset {
	if os.Getenv("INCLUSTER") == "true" {
		return inClusterKubeconfig()
	} else {
		return outClusterKubeconfig()
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // For Windows OS
}

func outClusterKubeconfig() *kubernetes.Clientset {
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "(optional) absolute path to the kubeconfig file")
	} else {
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		}
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func inClusterKubeconfig() *kubernetes.Clientset {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
