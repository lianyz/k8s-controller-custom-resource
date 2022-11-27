/*
@Time : 2022/11/22 18:14
@Author : lianyz
@Description :
*/

package main

import (
	"flag"
	"github.com/golang/glog"
	clientset "github.com/lianyz/k8s-controller-custom-resource/pkg/client/clientset/versioned"
	informers "github.com/lianyz/k8s-controller-custom-resource/pkg/client/informers/externalversions"
	"github.com/lianyz/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

var (
	masterURL  string
	kubeconfig string
)

// init 在main函数调用之前被调用
func init() {
	flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any kubeconfig. Only required if out-of-cluster.")
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
}

func main() {
	flag.Parse()

	stopCh := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
	if err != nil {
		glog.Fatal("error building kubeconfig: %s", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		glog.Fatal("error building kubernetes clientset: %s", err.Error())
	}

	networkClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		glog.Fatal("error building example clientset: %s", err.Error())
	}

	networkInformerFactory := informers.NewSharedInformerFactory(networkClient, time.Second*30)

	controller := NewController(kubeClient, networkClient,
		networkInformerFactory.Samplecrd().V1().Networks())

	go networkInformerFactory.Start(stopCh)

	if err = controller.Run(2, stopCh); err != nil {
		glog.Fatal("error running controller: %s", err.Error())
	}
}
