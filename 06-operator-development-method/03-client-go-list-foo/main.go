package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var gvr = schema.GroupVersionResource{
	Group:    "example.com",
	Version:  "v1alpha1",
	Resource: "foos",
}

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	TestString string `json:"testString"`
	TestNum    int    `json:"testNum"`
}

type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Foo `json:"items"`
}

func listFoos(client dynamic.Interface, namespace string) (*FooList, error) {
	// list is *unstructured.UnstructuredList
	// https://github.com/kubernetes/client-go/blob/master/dynamic/simple.go#L272-L294
	list, err := client.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// MarshalJSON ensures that the unstructured object produces proper JSON when passed to Go's standard JSON library.
	// func (u *Unstructured) MarshalJSON() ([]byte, error)
	// https://pkg.go.dev/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured#Unstructured.MarshalJSON
	data, err := list.MarshalJSON()
	if err != nil {
		return nil, err
	}

	var fooList FooList
	// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
	// https://pkg.go.dev/encoding/json#Unmarshal
	if err := json.Unmarshal(data, &fooList); err != nil {
		return nil, err
	}
	return &fooList, nil
}

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

	// func NewForConfig(inConfig *rest.Config) (Interface, error)
	// NewForConfig creates a new dynamic client or returns an error.
	// https://pkg.go.dev/k8s.io/client-go/dynamic#NewForConfig
	client, _ := dynamic.NewForConfig(config)

	// Get list of Foo objects from all namespaces
	foos, _ := listFoos(client, "")

	// Print Foo objects
	fmt.Println("INDEX\tNAMESPACE\tNAME")
	for i, foo := range foos.Items {
		fmt.Printf("%d\t%s\t%s\n", i, foo.GetNamespace(), foo.GetName())
	}
}
