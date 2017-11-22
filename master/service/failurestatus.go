package service

import (
	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clichaosv1 "github.com/slok/ragnarok/client/api/chaos/v1"
	"github.com/slok/ragnarok/log"
)

// FailureStatusService is how the master manages, enables, disables... attacks on the nodes.
type FailureStatusService interface {
	// GetNodeFailures returns all the failures of a given node.
	GetNodeFailures(nodeID string) []*chaosv1.Failure
	// GetNodeExpectedEnabledFailures returns all the failures in enabled state of a given node.
	GetNodeExpectedEnabledFailures(nodeID string) []*chaosv1.Failure
	// GetNodeExpectedDisabledFailures returns all the failures in disabled state of a given node.
	GetNodeExpectedDisabledFailures(nodeID string) []*chaosv1.Failure
	// GetFailure returns an specific failure.
	GetFailure(id string) (*chaosv1.Failure, error)
}

// FailureStatus is the implementation of failure status service.
type FailureStatus struct {
	client clichaosv1.FailureClientInterface // client is the client to manage failure objects.
	logger log.Logger
}

// NewFailureStatus returns a new FailureStatus
func NewFailureStatus(client clichaosv1.FailureClientInterface, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		client: client,
		logger: logger,
	}
}

func (f *FailureStatus) listFailuresByNode(nodeID string) ([]*chaosv1.Failure, error) {
	opts := api.ListOptions{
		LabelSelector: map[string]string{
			api.LabelNode: nodeID,
		},
	}
	return f.client.List(opts)
}

// GetNodeFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeFailures(nodeID string) []*chaosv1.Failure {
	flrs, err := f.listFailuresByNode(nodeID)
	if err != nil {
		f.logger.Errorf("error retrieving failures of node %s", nodeID)
	}

	return flrs
}

// GetNodeExpectedEnabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedEnabledFailures(nodeID string) []*chaosv1.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	flrs, err := f.listFailuresByNode(nodeID)
	if err != nil {
		f.logger.Errorf("error retrieving failures of node %s", nodeID)
	}
	res := []*chaosv1.Failure{}
	// Filter them by status.
	for _, flr := range flrs {
		if flr.Status.ExpectedState == chaosv1.EnabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetNodeExpectedDisabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedDisabledFailures(nodeID string) []*chaosv1.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	flrs, err := f.listFailuresByNode(nodeID)
	if err != nil {
		f.logger.Errorf("error retrieving failures of node %s", nodeID)
	}
	res := []*chaosv1.Failure{}
	// Filter them by status.
	for _, flr := range flrs {
		if flr.Status.ExpectedState == chaosv1.DisabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetFailure implements FailureStatusService interface.
func (f *FailureStatus) GetFailure(id string) (*chaosv1.Failure, error) {
	flr, err := f.client.Get(id)
	if err != nil {
		return nil, err
	}
	return flr, nil
}
