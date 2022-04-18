package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var defaultKubeConfigPath string
	if home := homedir.HomeDir(); home != "" {
		// build kubeconfig path from $HOME dir
		defaultKubeConfigPath = filepath.Join(home, ".kube", "config")
	}

	// Set kubeconfig flag
	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	flag.Parse()

	// Retrieve kubeconfig
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// Get clientset for kubernetes resources
	clientset, _ := kubernetes.NewForConfig(config)

	// Get list of pod objects from all namespaces
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	// show pod object to stdout
	fmt.Println("INDEX\tNAMESPACE\tNAME")
	for i, pod := range pods.Items {
		fmt.Printf("%d\t%s\t%s\n", i, pod.GetNamespace(), pod.GetName())
	}
}
