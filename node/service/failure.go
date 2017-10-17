package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/node/client"
)

const (
	stopTimeout = time.Second * 3
)

// FailureState is the interface all the services that want to process failures
// received from the master need to satisfy
type FailureState interface {
	// FailureStateHandler interface.
	client.FailureStateHandler
	// StartHandling will start handling the failures received from the master.
	StartHandling() error
	// StopHandling will stop handling the failures received from the master.
	StopHandling() error
}

// LogFailureState will process the failures received from the master and will only log them.
type LogFailureState struct {
	nodeID  string
	cli     client.Failure
	stopC   chan struct{}
	logger  log.Logger
	clock   clock.Clock
	running bool
	stMu    sync.Mutex // stMu is the status running mutex.
}

// NewLogFailureState returns a new Failurestate.
func NewLogFailureState(nodeID string, cli client.Failure, clock clock.Clock, logger log.Logger) *LogFailureState {
	logger = logger.WithField("kind", "log").WithField("service", "failureState")
	return &LogFailureState{
		nodeID: nodeID,
		cli:    cli,
		stopC:  make(chan struct{}),
		logger: logger,
		clock:  clock,
	}
}

// StartHandling satisfies FailureState interface.
func (l *LogFailureState) StartHandling() error {
	l.stMu.Lock()
	defer l.stMu.Unlock()

	if l.running {
		return fmt.Errorf("failure state handler already running")
	}

	l.logger.Infof("start handling failure status from master...")
	if err := l.cli.ProcessFailureStateStreaming(l.nodeID, l, l.stopC); err != nil {
		return err
	}
	l.running = true
	return nil
}

// StopHandling satisfies FailureState interface.
func (l *LogFailureState) StopHandling() error {
	l.stMu.Lock()
	defer l.stMu.Unlock()
	if !l.running {
		return fmt.Errorf("can't stop, failure state handler not running")
	}

	l.logger.Infof("stopping handling failure status from master...")
	select {
	case <-l.clock.After(stopTimeout):
		return fmt.Errorf("timeout stopping the handler of failure statuses from master")
	case l.stopC <- struct{}{}:
	}
	l.running = false

	return nil
}

// ProcessFailureStates implements client.FailureStateHandler
func (l *LogFailureState) ProcessFailureStates(failures []*v1.Failure) error {
	for _, fl := range failures {
		l.logger.WithField("failure", fl.Metadata.ID).Infof("%+v", fl)
	}
	return nil
}
