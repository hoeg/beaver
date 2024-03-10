package events

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

func ListenForEvents() {
	// Create a new shared informer factory
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*30)

	// Get the event informer from the factory
	eventInformer := informerFactory.Core().V1().Events().Informer()

	// Create a new event handler
	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event := obj.(*corev1.Event)
			fmt.Println("New event:", event.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// Handle event updates if needed
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
		return
	}

	// Run until interrupted
	select {}
}
