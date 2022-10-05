package main

import (
	"context"
	"flag"

	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("source-examples")

func main() {
	// Prepare log
	opts := zap.Options{
		Development: true,
		TimeEncoder: zapcore.ISO8601TimeEncoder,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	log.Info("source start")

	// Get a kubeconfig
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
	}

	// Create a Cache
	cache, err := cache.New(cfg, cache.Options{}) // &informerCache{InformersMap: im}, nil
	if err != nil {
		log.Error(err, "")
	}
	log.Info("cache is created")

	ctx := context.Background()
	pod := &v1.Pod{}
	cache.Get(ctx, client.ObjectKeyFromObject(pod), pod)

	// Start Cache
	go func() {
		if err := cache.Start(ctx); err != nil { // func (m *InformersMap) Start(ctx context.Context) error {
			log.Error(err, "failed to start cache")
		}
	}()
	log.Info("cache is started")

	// Create a Kind (Source)
	kindWithCachePod := source.NewKindWithCache(pod, cache)

	// Create a Queue
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "test")

	// Create an EventHandler
	eventHandler := handler.Funcs{
		CreateFunc: func(e event.CreateEvent, q workqueue.RateLimitingInterface) {
			log.Info("CreateFunc is called", "object", e.Object.GetName())
			queue.Add(WorkQueueItem{Event: "Create", Name: e.Object.GetName()})
		},
		UpdateFunc: func(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
			log.Info("UpdateFunc is called", "objectNew", e.ObjectNew.GetName(), "objectOld", e.ObjectOld.GetName())
			queue.Add(WorkQueueItem{Event: "Update", Name: e.ObjectNew.GetName()})
		},
		DeleteFunc: func(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
			log.Info("DeleteFunc is called", "object", e.Object.GetName())
			queue.Add(WorkQueueItem{Event: "Delete", Name: e.Object.GetName()})
		},
	}

	// Start Kind (Source)
	if err := kindWithCachePod.Start(ctx, eventHandler, queue); err != nil { // Get informer and set eventHandler
		log.Error(err, "")
	}

	// Wait for Cache
	if err := kindWithCachePod.WaitForSync(ctx); err != nil {
		log.Error(err, "")
	}
	log.Info("kindWithCache is ready")

	// Get items from Queue
	for {
		item, shutdown := queue.Get()
		if shutdown {
			break
		}
		log.Info("got item", "item", item)
	}
}

type WorkQueueItem struct {
	Event string
	Name  string
}
