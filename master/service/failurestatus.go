package service

import (
	"fmt"

	"github.com/slok/ragnarok/chaos/failure"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/types"
)

// FailureStatusService is how the master manages, enables, disables... attacks on the nodes.
type FailureStatusService interface {
	// GetNodeFailures returns all the failures of a given node.
	GetNodeFailures(nodeID string) []*failure.Failure
	// GetNodeExpectedEnabledFailures returns all the failures in enabled state of a given node.
	GetNodeExpectedEnabledFailures(nodeID string) []*failure.Failure
	// GetNodeExpectedDisabledFailures returns all the failures in disabled state of a given node.
	GetNodeExpectedDisabledFailures(nodeID string) []*failure.Failure
	// GetFailure returns an specific failure.
	GetFailure(id string) (*failure.Failure, error)
}

// FailureStatus is the implementation of failure status service.
type FailureStatus struct {
	repo   FailureRepository // Repo is the failure repository.
	logger log.Logger
}

// NewFailureStatus returns a new FailureStatus
func NewFailureStatus(repository FailureRepository, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		repo:   repository,
		logger: logger,
	}
}

// GetNodeFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeFailures(nodeID string) []*failure.Failure {
	return f.repo.GetAllByNode(nodeID)
}

// GetNodeExpectedEnabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedEnabledFailures(nodeID string) []*failure.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	fs := f.repo.GetAllByNode(nodeID)
	res := []*failure.Failure{}
	// Filter them by status.
	for _, flr := range fs {
		if flr.ExpectedState == types.EnabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetNodeExpectedDisabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedDisabledFailures(nodeID string) []*failure.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	fs := f.repo.GetAllByNode(nodeID)
	res := []*failure.Failure{}
	// Filter them by status.
	for _, flr := range fs {
		if flr.ExpectedState == types.DisabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetFailure implements FailureStatusService interface.
func (f *FailureStatus) GetFailure(id string) (*failure.Failure, error) {
	flr, ok := f.repo.Get(id)
	if !ok {
		return nil, fmt.Errorf("failure %s can't be retrieved", id)
	}
	return flr, nil
}
