package events

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var stopCh chan struct{}

func Start(ctx context.Context, errCh chan<- error) error {
	// Use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		runtime.HandleError(err)
		return err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		runtime.HandleError(err)
		return err
	}
	// Create a new shared informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// Get the event informer from the factory
	eventInformer := informerFactory.Core().V1().Events().Informer()

	// Create a new event handler
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event := obj.(*corev1.Event)
			slog.Debug("add event", "kind", event.InvolvedObject.Kind, "name", event.Name, "namespace", event.Namespace, "reason", event.Reason)
			if event.InvolvedObject.Kind == "Deployment" || event.InvolvedObject.Kind == "ReplicaSet" || event.InvolvedObject.Kind == "StatefulSet" || event.InvolvedObject.Kind == "DaemonSet" || event.InvolvedObject.Kind == "Job" || event.InvolvedObject.Kind == "CronJob" {
				var artifactID string
				switch event.InvolvedObject.Kind {
				case "Deployment":
					deployment, err := clientset.AppsV1().Deployments(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = deployment.Annotations["hoeg.com/artifact-id"]

				case "ReplicaSet":
					replicaSet, err := clientset.AppsV1().ReplicaSets(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = replicaSet.Annotations["hoeg.com/artifact-id"]

				case "StatefulSet":
					statefulSet, err := clientset.AppsV1().StatefulSets(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = statefulSet.Annotations["hoeg.com/artifact-id"]

				case "DaemonSet":
					daemonSet, err := clientset.AppsV1().DaemonSets(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = daemonSet.Annotations["hoeg.com/artifact-id"]

				case "Job":
					job, err := clientset.BatchV1().Jobs(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = job.Annotations["hoeg.com/artifact-id"]

				case "CronJob":
					cronJob, err := clientset.BatchV1beta1().CronJobs(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
					if err != nil {
						runtime.HandleError(err)
						return
					}
					artifactID = cronJob.Annotations["hoeg.com/artifact-id"]

				default:
					slog.Error("Unsupported resource kind", "kind", event.InvolvedObject.Kind)
					return
				}
				// extract the commit sha from the artifactID
				if artifactID != "" {
					commitSha := extractCommitSha(artifactID)
					slog.Info("Found artifactID", "Pod", event.InvolvedObject.Name, "namespace", event.Namespace, "commitSha", commitSha)
				} else {
					slog.Debug("No artifact ID found", "Pod", event.InvolvedObject.Name, "namespace", event.Namespace)
				}
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			event := newObj.(*corev1.Event)
			slog.Debug("update event", "kind", event.InvolvedObject.Kind, "name", event.Name, "namespace", event.Namespace, "reason", event.Reason)
		},
		DeleteFunc: func(obj interface{}) {
			event := obj.(*corev1.Event)
			slog.Debug("delete event", "kind", event.InvolvedObject.Kind, "name", event.Name, "namespace", event.Namespace, "reason", event.Reason)
		},
	}

	// Add the event handler to the informer
	eventInformer.AddEventHandler(eventHandler)

	// Start the informer
	slog.Info("Starting the event informer")
	stopCh = make(chan struct{})
	go eventInformer.Run(stopCh)

	// Wait for the informer to sync
	if !cache.WaitForCacheSync(stopCh, eventInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("failed to sync informer cache"))
	}

	return nil
}

func Stop(ctx context.Context) error {
	// Stop the informer
	slog.Info("Stopping the event informer")
	close(stopCh)
	return nil
}

func extractCommitSha(artifactID string) string {
	// Find the last hyphen
	lastHyphenIndex := strings.LastIndex(artifactID, "-")
	if lastHyphenIndex == -1 {
		return ""
	}

	// Find the second last hyphen
	secondLastHyphenIndex := strings.LastIndex(artifactID[:lastHyphenIndex], "-")
	if secondLastHyphenIndex == -1 {
		return ""
	}

	// Extract the commit sha
	commitSha := artifactID[secondLastHyphenIndex+1 : lastHyphenIndex]

	return commitSha
}
