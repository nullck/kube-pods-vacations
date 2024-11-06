package kube_pods_vacations

import (
	"log"
	"os"
	"time"

	kubernetesmgmt "github.com/nullck/kube-pods-vacations/pkg/kubernetes_mgmt"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

//#kube_pods_vacations --annotationsscan true | --set-replicas true --replicas --resource-type --resource-name

func NewRootCommand() *cobra.Command {
	var setReplicas int
	var namespaceNames []string
	var resourceType, resourceName string
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var rootCmd = &cobra.Command{
		Use:     "kube-pods-vacations",
		Short:   "kube-pods-vacations",
		Version: "0.0.1",
		Long:    `kube-pods-vacations is an utility developed to reduce the number of pods for a specific time frame`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	var annotationScanCmd = &cobra.Command{
		Use:     "annotations-scan",
		Short:   "activate annotations scan",
		Version: rootCmd.Version,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, n := range namespaceNames {
				r := kubernetesmgmt.KubeMgmt{
					NamespaceName:    n,
					AnnotationPrefix: "kube-pods-vacations",
					Logger:           logger,
				}
				err := r.DeploymentAnnotationsBuilder()
				if err != nil {
					log.Printf("Error running annotations builder %v", err)
					return
				}
				time.Sleep(2 * time.Second)
			}

			logger.Info().Msg("worked")
		},
	}

	var setReplicasCmd = &cobra.Command{
		Use:     "set-replicas",
		Short:   "set replicas",
		Version: rootCmd.Version,
		Long:    `set the number of replicas for a resouce like Deployment or StatefulSet`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	annotationScanCmd.Flags().StringSliceVar(&namespaceNames, "namespaces", []string{}, "set a list of namespaces for annotations scan")
	setReplicasCmd.Flags().IntVar(&setReplicas, "set-replicas", 0, "set amount of pods replicas in a resource type")
	setReplicasCmd.Flags().StringVarP(&resourceType, "resource-type", "", "deployment", "set resource type name, e.g: deployment")
	setReplicasCmd.Flags().StringVar(&resourceName, "resource-name", "", "set resource name, e.g: deployment name")
	return annotationScanCmd
}
