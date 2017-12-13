package informer

import (
	"fmt"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/util/queue"
	"github.com/slok/ragnarok/client/util/store"
	"github.com/slok/ragnarok/log"
)

// Informer is the one that will inform about the changes on resources so the controllers could apply their logic.
// The informer type will set the nature of how the controllers are informed.
type Informer interface {
	// Run will run the informer, this is: start handling the events and start informing. How it informs depends on the informer type.
	Run(stopCh chan struct{}) error
}

// WorkQueueInformerInterface is an informer type that uses a workqueue to inform to the controllers. In this case it will
// have a cache with the current state of the resources received as events and at the same time will push jobs on a queue
// so the controller can get them afterwards. This informer can only be used by a controller at the same time.
type WorkQueueInformerInterface interface {
	Informer
	GetStore() store.Store
	GetQueue() queue.Queue
}

// WorkQueueInformer implements the basic functionality of an informer using work queues.
// It will store the resources in an indexed store (so it can be retrieved afterwards). Then the resources
// that triggered the event will be added to a workqueue so they can be retrieved as jobs afterwards and
// get the latest state from the cache store.
type WorkQueueInformer struct {
	// store is where the received state from the client will be stored, in other words, this is a mini cache
	// system where all the event resources will be placed. When we want to access an event received object we will
	// look here its state. This is done 1 for performance and 2, when we have an error processing a resource event
	// we could get again the resource afterwards (and the state could change, for example deleted)
	store store.Store
	// lw is the way the informer will get all the events and resources.
	lw ListerWatcher
	// queue is where all the resource events resource will be add so they can be processed
	// in an ordered and ditributed way.
	queue queue.Queue

	handler ResourceEventHandler
	lwOpts  api.ListOptions
	indexer store.ObjectIndexKeyer
	logger  log.Logger
}

// NewWorkQueueInformer returns a WorkQueueInformer with t values.
func NewWorkQueueInformer(
	indexer store.ObjectIndexKeyer,
	queue queue.Queue,
	store store.Store,
	lwOpts api.ListOptions,
	lw ListerWatcher,
	logger log.Logger) *WorkQueueInformer {

	// Create a custom resource event handler that will update the cache store with the event received state and also
	// will add the resource to the job queue.
	reh := &ResourceEventHandlerFuncs{
		OnAddFunc: func(obj api.Object) {
			if err := store.Add(obj); err != nil {
				logger.Errorf("error adding object from add event to the cache store")
				return
			}
			key, err := indexer.GetKey(obj)
			if err != nil {
				logger.Errorf("error getting index key of object on add event to the cache store")
			}
			queue.Push(key)
		},
		OnUpdateFunc: func(oldObj, newObj api.Object) {
			if err := store.Update(newObj); err != nil {
				logger.Errorf("error adding object from update event to the cache store")
				return
			}
			key, err := indexer.GetKey(newObj)
			if err != nil {
				logger.Errorf("error getting index key of object on update event to the cache store")
			}
			queue.Push(key)
		},
		OnDeleteFunc: func(obj api.Object) {
			if err := store.Delete(obj); err != nil {
				logger.Errorf("error adding object from delete event to the cache store")
				return
			}
			key, err := indexer.GetKey(obj)
			if err != nil {
				logger.Errorf("error getting index key of object on delete event to the cache store")
			}
			queue.Push(key)
		},
	}

	return &WorkQueueInformer{
		indexer: indexer,
		store:   store,
		queue:   queue,
		handler: reh,
		lwOpts:  lwOpts,
		lw:      lw,
		logger:  logger,
	}
}

// setInitialState will get the current state of the resources and set on the cache store.
func (w *WorkQueueInformer) setInitialState() error {
	// Get the list of resources. (current state).
	objs, err := w.lw.List(w.lwOpts)
	if err != nil {
		return err
	}
	w.logger.Debugf("list got %d objects", len(objs))

	// Set state on cache.
	for _, obj := range objs {
		w.store.Add(obj)
		// Send a job to the queue so they check all the objects at the begginning.
		key, err := w.indexer.GetKey(obj)
		if err != nil {
			w.logger.Errorf("error getting key of object: %s", err)
		}
		if err := w.queue.Push(key); err != nil {
			w.logger.Errorf("error pushing job to queue: %s", err)
		}
	}
	return nil
}

// handleEvent will be set store on the cache store and push to the job queue.
func (w *WorkQueueInformer) handleEvent(ev watch.Event) error {
	switch ev.Type {
	case watch.AddedEvent:
		w.handler.OnAdd(ev.Object)
	case watch.UpdatedEvent:
		// Get the previous state from the cache store.
		old, exists, err := w.store.Get(ev.Object)
		if err != nil {
			return fmt.Errorf("Error retrieving old object for: %s", ev.Object)
		}
		// If doesn't exists previously, then is a creation not an update.
		if !exists {
			w.handler.OnAdd(ev.Object)
		} else {
			w.handler.OnUpdate(old, ev.Object)
		}
	case watch.DeletedEvent:
		w.handler.OnDelete(ev.Object)
	default:
		return fmt.Errorf("wrong type of event: %d", ev.Type)
	}
	return nil
}

// startWatcher will start the watcher that receives the events. and will handle the events
func (w *WorkQueueInformer) startWatcher(stopCh chan struct{}) error {
	watcher, err := w.lw.Watch(w.lwOpts)
	if err != nil {
		return err
	}

	// Start a gourutine to handle the events, this events
	eventChan := watcher.GetChan()
	for {
		select {
		case ev := <-eventChan:
			w.handleEvent(ev)
		case <-stopCh:
			w.logger.Infof("stoppping watching events")
			return nil
		}
	}
}

// Run will start the informer main event processing loop.
func (w *WorkQueueInformer) Run(stopCh chan struct{}) error {
	// TODO: Handle panics, cleanup...

	defer w.queue.ShutDown()
	w.logger.Info("starting informer")

	if err := w.setInitialState(); err != nil {
		return fmt.Errorf("could not set the initial state: %s", err)
	}

	w.logger.Info("started handling events")
	err := w.startWatcher(stopCh)

	w.logger.Info("ending informer run")
	return err
}

// GetStore satisfies WorkQueueInformerInterface.
func (w *WorkQueueInformer) GetStore() store.Store {
	return w.store
}

// GetQueue satisfies WorkQueueInformerInterface.
func (w *WorkQueueInformer) GetQueue() queue.Queue {
	return w.queue
}
