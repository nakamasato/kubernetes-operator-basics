package main

import (
	"context"
	"flag"

	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	log = logf.Log.WithName("source-examples")
)

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

	// Set a mapper
	mapper, err := func(c *rest.Config) (meta.RESTMapper, error) {
		return apiutil.NewDynamicRESTMapper(c)
	}(cfg)
	if err != nil {
		log.Error(err, "")
	}

	// Create a Cache
	cache, err := cache.New(cfg, cache.Options{Mapper: mapper}) // &informerCache{InformersMap: im}, nil
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

	kindWithCachePod := source.NewKindWithCache(pod, cache)

	// Prepare queue and eventHandler
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "test")

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

	// Start Source
	if err := kindWithCachePod.Start(ctx, eventHandler, queue); err != nil { // Get informer and set eventHandler
		log.Error(err, "")
	}

	// Wait for cache
	if err := kindWithCachePod.WaitForSync(ctx); err != nil {
		log.Error(err, "")
	}
	log.Info("kindWithCache is ready")

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
