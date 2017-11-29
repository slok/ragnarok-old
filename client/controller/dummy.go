package controller

import (
	"github.com/slok/ragnarok/client/informer"
	"github.com/slok/ragnarok/log"
)

// DummyWorkQueueController will print the jobs that will be adding to the queue.
type DummyWorkQueueController struct {
	informer informer.WorkQueueInformerInterface
	logger   log.Logger
	stopC    chan struct{}
}

// NewDummyWorkQueueController returns a new DummyWorkQueueController.
func NewDummyWorkQueueController(informer informer.WorkQueueInformerInterface, logger log.Logger) *DummyWorkQueueController {
	return &DummyWorkQueueController{
		informer: informer,
		logger:   logger,
		stopC:    make(chan struct{}),
	}
}

func (d *DummyWorkQueueController) processOne(job interface{}) error {
	store := d.informer.GetStore()
	jobStr := job.(string)
	obj, exists, err := store.GetByKey(jobStr)
	if err != nil {
		d.logger.Warnf("error processing job: %s", jobStr)
		return err
	}
	if !exists {
		d.logger.Infof("Job doesn't exists, object deleted: %s", jobStr)
	} else {
		d.logger.Infof("Job exists: %#v", obj)
	}
	return nil
}

// ProcessingLoop will be the loop that processes all the job queue.
func (d *DummyWorkQueueController) processingLoop() error {
	q := d.informer.GetQueue()
	for {
		job, shutDown := q.Pop()
		if shutDown {
			return nil
		}
		// If error processing then requeue.
		if err := d.processOne(job); err != nil {
			q.Push(job)
		}
	}
}

// Run satisfies Controller interface.
func (d *DummyWorkQueueController) Run() error {
	// First start the informer to start handling the events.
	go d.informer.Run(d.stopC)

	// Start handling jobs from the informer.
	d.processingLoop()

	return nil
}

// Stop satisfies Controller interface.
func (d *DummyWorkQueueController) Stop() error {
	d.stopC <- struct{}{}
	return nil
}
