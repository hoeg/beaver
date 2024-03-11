package events

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func ListenForEvents() {

	// Use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		runtime.HandleError(err)
		return
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	// Create a new shared informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// Get the event informer from the factory
	eventInformer := informerFactory.Core().V1().Events().Informer()

	// Create a new event handler
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event := obj.(*corev1.Event)
			if event.InvolvedObject.Kind == "Pod" && (event.Reason == "Created" || event.Reason == "Updated") {
				fmt.Println("New event:", event.Name)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// Handle event updates if needed
			oldEvent := oldObj.(*corev1.Event)
			newEvent := newObj.(*corev1.Event)
			if oldEvent.InvolvedObject.Kind == "Pod" && newEvent.InvolvedObject.Kind == "Pod" && newEvent.Reason == "Updated" {
				fmt.Println("Updated event:", newEvent.Name)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// Handle event deletions if needed
		},
	}

	// Add the event handler to the informer
	eventInformer.AddEventHandler(eventHandler)

	// Start the informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	go eventInformer.Run(stopCh)

	// Wait for the informer to sync
	if !cache.WaitForCacheSync(stopCh, eventInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("failed to sync informer cache"))
	}

	// Run until interrupted
	select {}
}
