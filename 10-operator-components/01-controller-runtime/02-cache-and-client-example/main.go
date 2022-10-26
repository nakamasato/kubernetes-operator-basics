package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	// Get a kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create Cache
	cache, err := cache.New(cfg, cache.Options{}) // &informerCache{InformersMap: im}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[cache] created")

	ctx := context.Background()

	// Start Cache
	go func() {
		if err := cache.Start(ctx); err != nil { // func (m *InformersMap) Start(ctx context.Context) error {
			log.Fatal(err)
		}
	}()
	fmt.Println("[cache] started")

	// Wait for Cache Synced
	if isSynced := cache.WaitForCacheSync(ctx); !isSynced {
		log.Fatal("[cache] failed to sync")
	}
	fmt.Println("[cache] synced")

	// Get Nginx Pod from Cache -> Expect to have NotFound error
	pod := &v1.Pod{}
	err = cache.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "nginx",
	}, pod)
	if err != nil {
		fmt.Printf("[cache] failed to get nginx Pod (err: %v)\n", err)
	} else {
		fmt.Println("[cache] successfully got nginx Pod")
	}

	// Create Client
	cli, err := client.New(cfg, client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	// Create Nginx Pod
	err = cli.Create(ctx, &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "nginx",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}, &client.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[client] created nginx Pod")

	time.Sleep(100 * time.Millisecond) // wait cache is updated

	// Get Nginx Pod from Cache -> Expect to successfuly get the Pod
	err = cache.Get(ctx, client.ObjectKey{
		Namespace: "default",
		Name:      "nginx",
	}, pod)
	if err != nil {
		fmt.Printf("[cache] failed to get nginx Pod (err: %v)\n", err)
	} else {
		fmt.Println("[cache] successfully got nginx Pod")
	}

	// Delete Nginx Pod
	err = cli.Delete(ctx, &v1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "nginx",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}, &client.DeleteOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[client] deleted nginx Pod")
}
