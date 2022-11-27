/*
@Time : 2022/11/22 18:13
@Author : lianyz
@Description :
*/

package main

import (
	"fmt"
	"github.com/golang/glog"
	samplecrdv1 "github.com/lianyz/k8s-controller-custom-resource/pkg/apis/samplecrd/v1"
	clientset "github.com/lianyz/k8s-controller-custom-resource/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"time"

	networkscheme "github.com/lianyz/k8s-controller-custom-resource/pkg/client/clientset/versioned/scheme"
	informers "github.com/lianyz/k8s-controller-custom-resource/pkg/client/informers/externalversions/samplecrd/v1"
	listers "github.com/lianyz/k8s-controller-custom-resource/pkg/client/listers/samplecrd/v1"
)

const (
	controllerAgentName = "network-controller"

	SucceedSynced = "Synced"

	MessageResourceSynced = "Network synced successfully"
)

type Controller struct {
	kubeClientset kubernetes.Interface

	networkClientset clientset.Interface

	deployInformer kubeinformers.DeploymentInformer

	networksLister listers.NetworkLister
	networkSynced  cache.InformerSynced

	workQueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

// NewController return a new network controller
func NewController(
	kubeClientset kubernetes.Interface,
	networkClientset clientset.Interface,
	deployInformer kubeinformers.DeploymentInformer,
	networkInformer informers.NetworkInformer) *Controller {

	utilruntime.Must(networkscheme.AddToScheme(scheme.Scheme))
	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeClientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeClientset:    kubeClientset,
		networkClientset: networkClientset,
		deployInformer:   deployInformer,
		networksLister:   networkInformer.Lister(),
		networkSynced:    networkInformer.Informer().HasSynced,
		workQueue:        workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks"),
		recorder:         recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Network resources change
	networkInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.enqueueNetwork,
		UpdateFunc: controller.enqueueNetworkForUpdate,
		DeleteFunc: controller.enqueueNetworkForDelete,
	})

	return controller
}

// Run will set up the event handlers for types we are interested, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shut down the workQueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.workQueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Network control loop")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.networkSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workQueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workQueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workQueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so that we can defer c.workQueue Done.
	err := func(obj interface{}) error {
		defer c.workQueue.Done(obj)

		key, ok := obj.(string)
		if !ok {
			c.workQueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}

		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}

		c.workQueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	network, err := c.networksLister.Networks(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			glog.Warningf("Network %s/%s does not exist in local cache, will delete it from Neutron...",
				namespace, name)
			glog.Infof("[Neutron] Deleting network: %s/%s ...", namespace, name)

			// todo call Neutron API to delete this network by name.
			//
			// neutron.Delete(namespace, name)
			return nil
		}
		runtime.HandleError(fmt.Errorf("failed to list network by: %s/%s", namespace, name))

		return err
	}

	glog.Infof("[Neutron] Try to process network: %#v ...", network)

	// todo Do diff()
	//
	// actualNetwork, exists := neutron.Get(namespace, name)
	//
	// if !exists {
	//   neutron.Create(namespace, name)
	// } else if !reflect.DeepEqual(actualNetwork, network) {
	//   neutron.Update(namespace, name)
	// }

	c.recorder.Event(network, corev1.EventTypeNormal, SucceedSynced, MessageResourceSynced)
	return nil
}

// enqueueNetwork takes a Network resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should not be
// passed resources of any type other than Network.
func (c *Controller) enqueueNetwork(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}
	c.workQueue.AddRateLimited(key)
}

// enqueueNetworkForDelete takes a deleted Network resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should not be
// passed resources of any type other than Network.
func (c *Controller) enqueueNetworkForUpdate(old, new interface{}) {
	oldNetwork := old.(*samplecrdv1.Network)
	newNetwork := new.(*samplecrdv1.Network)
	if oldNetwork.ResourceVersion == newNetwork.ResourceVersion {
		// Periodic resync will send update events for all known Networks.
		// Two different versions of the same Network will always have different ResourceVersions.
		return
	}
	c.enqueueNetwork(new)
}

func (c *Controller) enqueueNetworkForDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
		return
	}

	c.workQueue.AddRateLimited(key)
}
