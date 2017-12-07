package controller

import (
	"fmt"

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/client/informer"
	"github.com/slok/ragnarok/log"
	experiment "github.com/slok/ragnarok/master/service/experiment"
)

// Experiment is the controller that will manage the expriment creation, enable, disable
// and deletion.
type Experiment struct {
	informer informer.WorkQueueInformerInterface
	service  experiment.Manager
	stopC    chan struct{}
	logger   log.Logger
}

// NewExperiment returns a new Experiment controller.
func NewExperiment(informer informer.WorkQueueInformerInterface, service experiment.Manager, logger log.Logger) *Experiment {
	return &Experiment{
		informer: informer,
		service:  service,
		stopC:    make(chan struct{}),
		logger:   logger,
	}
}

func (e *Experiment) processOne(job interface{}) error {
	store := e.informer.GetStore()
	jobStr := job.(string)
	obj, exists, err := store.GetByKey(jobStr)
	if err != nil {
		e.logger.Warnf("error processing job: %s", jobStr)
		return err
	}

	if !exists {
		e.logger.Warnf("experiment doesn't exist, TODO")

	}

	exp, ok := obj.(*chaosv1.Experiment)
	if !ok {
		return fmt.Errorf("invalid type received job object")
	}

	// Ensure failure instances.
	if err := e.service.EnsureFailures(exp); err != nil {
		return err
	}

	return nil
}

// ProcessingLoop will be the loop that processes all the job queue.
func (e *Experiment) processingLoop() error {
	q := e.informer.GetQueue()
	for {
		job, shutDown := q.Pop()
		if shutDown {
			return nil
		}
		// If error processing then requeue.
		if err := e.processOne(job); err != nil {
			q.Push(job)
		}
	}
}

// Run satisfies Controller interface.
func (e *Experiment) Run() error {
	// First start the informer to start handling the events.
	go e.informer.Run(e.stopC)

	// Start handling jobs from the informer.
	e.processingLoop()

	return nil
}

// Stop satisfies Controller interface.
func (e *Experiment) Stop() error {
	e.stopC <- struct{}{}
	return nil
}
