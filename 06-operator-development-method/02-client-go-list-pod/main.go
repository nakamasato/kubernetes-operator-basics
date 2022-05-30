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
	// String defines a string flag with specified name, default value, and usage string. The return value is the address of a string variable that stores the value of the flag.
	// https://pkg.go.dev/flag#String
	kubeconfig := flag.String("kubeconfig", defaultKubeConfigPath, "kubeconfig config file")
	flag.Parse()

	// Get kubeconfig
	// func BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error)
	// BuildConfigFromFlags is a helper function that builds configs from a master url or a kubeconfig filepath.
	config, _ := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	// func NewForConfig(c *rest.Config) (*Clientset, error)
	// NewForConfig creates a new Clientset for the given config.
	// https://pkg.go.dev/k8s.io/client-go/kubernetes#NewForConfig
	clientset, _ := kubernetes.NewForConfig(config)

	// Get list of Pod objects from all namespaces
	pods, _ := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})

	// Print Pod objects
	fmt.Println("INDEX\tNAMESPACE\tNAME")
	for i, pod := range pods.Items {
		fmt.Printf("%d\t%s\t%s\n", i, pod.GetNamespace(), pod.GetName())
	}
}
