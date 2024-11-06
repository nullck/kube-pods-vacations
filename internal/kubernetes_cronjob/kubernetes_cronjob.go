package kubernetes_cronjob

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeCronJob struct {
	Namespace               string
	CronName                string
	ReducedReplicasCronExpr string
	DesiredReplicasCronExpr string
	ResourceType            string
	ResourceName            string
	MinimalReplicas         string
	DesiredReplicas         string
	Logger                  zerolog.Logger
}

func NewKubeCronJob(namespace, cronName, reduceRepCronExpr, desiredRepCronExpr, resourceType, resourceName, minReplicas, desReplicas string, logger zerolog.Logger) *KubeCronJob {
	kubeCronJob := KubeCronJob{
		CronName:                cronName,
		Namespace:               namespace,
		ReducedReplicasCronExpr: reduceRepCronExpr,
		DesiredReplicasCronExpr: desiredRepCronExpr,
		ResourceType:            resourceType,
		ResourceName:            resourceName,
		MinimalReplicas:         minReplicas,
		DesiredReplicas:         desReplicas,
		Logger:                  logger,
	}
	return &kubeCronJob
}

func (kc *KubeCronJob) CreateReduceCronJob(clientset *kubernetes.Clientset) error {
	kc.Logger.Info().Msgf("creating cron job %s", kc.CronName)

	cronJobsClient := clientset.BatchV1().CronJobs(kc.Namespace)

	cronJob := &batchv1.CronJob{ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("%s-reduce-pods", kc.CronName)},
		Spec: batchv1.CronJobSpec{
			Schedule: kc.ReducedReplicasCronExpr,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: apiv1.PodTemplateSpec{
						Spec: apiv1.PodSpec{
							Containers: []apiv1.Container{
								{
									Name:    "hello",
									Image:   "busybox:1.28",
									Command: []string{"/bin/sh", "-c", fmt.Sprintf("echo set replicas %s %s %s", kc.MinimalReplicas, kc.ResourceType, kc.ResourceName)},
								},
							},
							RestartPolicy: apiv1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	_, err := cronJobsClient.Create(context.Background(), cronJob, metav1.CreateOptions{})
	if err != nil {
		kc.Logger.Error().Msgf("failed to create reduce cron job %v", err)
		return err
	}

	kc.Logger.Info().Msg("created reduce cron job successfully")
	return nil
}

func (kc *KubeCronJob) CreateDesiredCronJob(clientset *kubernetes.Clientset) error {
	kc.Logger.Info().Msgf("creating cron job %s\n", kc.CronName)

	cronJobsClient := clientset.BatchV1().CronJobs(kc.Namespace)

	cronJob := &batchv1.CronJob{ObjectMeta: metav1.ObjectMeta{
		Name: fmt.Sprintf("%s-desired-pods", kc.CronName)},
		Spec: batchv1.CronJobSpec{
			Schedule: kc.DesiredReplicasCronExpr,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: apiv1.PodTemplateSpec{
						Spec: apiv1.PodSpec{
							Containers: []apiv1.Container{
								{
									Name:    "hello",
									Image:   "busybox:1.28",
									Command: []string{"/bin/sh", "-c", fmt.Sprintf("echo set replicas %s %s %s", kc.MinimalReplicas, kc.ResourceType, kc.ResourceName)},
								},
							},
							RestartPolicy: apiv1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}

	_, err := cronJobsClient.Create(context.Background(), cronJob, metav1.CreateOptions{})
	if err != nil {
		kc.Logger.Error().Msgf("failed to create reduce cron job %v", err)
		return err
	}

	kc.Logger.Info().Msg("created reduce cron job successfully")
	return nil
}
