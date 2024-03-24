package events

import (
	"context"
	"fmt"
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
			fmt.Printf("New event add: %v\n", obj)
			event := obj.(*corev1.Event)
			if event.InvolvedObject.Kind == "Pod" && (event.Reason == "Created" || event.Reason == "Updated") {
				fmt.Println("New event:", event.Name)
				pod, err := clientset.CoreV1().Pods(event.Namespace).Get(context.Background(), event.InvolvedObject.Name, metav1.GetOptions{})
				if err != nil {
					runtime.HandleError(err)
					return
				}
				//from the pod, get the value of the annotation lunarway.com/artifact-id: 'atl-2779-derp-eb6283670d-5855a38e0c' and extract the eb6283670d part as the commit sha
				artifactID := pod.Annotations["lunarway.com/artifact-id"]
				// extract the substring eb6283670d from the artifactID, the commit sha is the last string after the hyphen before the last hyphen, we cannot garantee the length of the commit sha, so we cannot use a fixed length, so use a regex instead
				commitSha := extractCommitSha(artifactID)
				fmt.Println(commitSha)
				// fmt.Println(artifactID)

			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Printf("New event udpated: %v\n", newObj)
			// Handle event updates if needed
			oldEvent := oldObj.(*corev1.Event)
			newEvent := newObj.(*corev1.Event)
			if oldEvent.InvolvedObject.Kind == "Pod" && newEvent.InvolvedObject.Kind == "Pod" && newEvent.Reason == "Updated" {
				fmt.Println("Updated event:", newEvent.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Printf("New event delete: %v\n", obj)
		},
	}

	// Add the event handler to the informer
	eventInformer.AddEventHandler(eventHandler)

	// Start the informer
	fmt.Println("Starting the event informer")
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
	fmt.Println("Stopping the event informer")
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
