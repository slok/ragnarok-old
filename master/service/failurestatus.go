package service

import (
	"fmt"

	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/log"
)

// FailureStatusService is how the master manages, enables, disables... attacks on the nodes.
type FailureStatusService interface {
	// GetNodeFailures returns all the failures of a given node.
	GetNodeFailures(nodeID string) []*v1.Failure
	// GetNodeExpectedEnabledFailures returns all the failures in enabled state of a given node.
	GetNodeExpectedEnabledFailures(nodeID string) []*v1.Failure
	// GetNodeExpectedDisabledFailures returns all the failures in disabled state of a given node.
	GetNodeExpectedDisabledFailures(nodeID string) []*v1.Failure
	// GetFailure returns an specific failure.
	GetFailure(id string) (*v1.Failure, error)
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
func (f *FailureStatus) GetNodeFailures(nodeID string) []*v1.Failure {
	return f.repo.GetAllByNode(nodeID)
}

// GetNodeExpectedEnabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedEnabledFailures(nodeID string) []*v1.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	fs := f.repo.GetAllByNode(nodeID)
	res := []*v1.Failure{}
	// Filter them by status.
	for _, flr := range fs {
		if flr.Status.ExpectedState == v1.EnabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetNodeExpectedDisabledFailures implements FailureStatusService interface.
func (f *FailureStatus) GetNodeExpectedDisabledFailures(nodeID string) []*v1.Failure {
	// Get all.
	// TODO: Ask filtered directly to the repository.
	fs := f.repo.GetAllByNode(nodeID)
	res := []*v1.Failure{}
	// Filter them by status.
	for _, flr := range fs {
		if flr.Status.ExpectedState == v1.DisabledFailureState {
			res = append(res, flr)
		}
	}
	return res
}

// GetFailure implements FailureStatusService interface.
func (f *FailureStatus) GetFailure(id string) (*v1.Failure, error) {
	flr, ok := f.repo.Get(id)
	if !ok {
		return nil, fmt.Errorf("failure %s can't be retrieved", id)
	}
	return flr, nil
}
