package scheduler

import (
	"fmt"
	"math/rand"
	"time"

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/repository"
)

const (
	fmtFailureID = "%s-%s-flr-%s"
)

func randID() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%09d", r.Intn(999999999))
}

// createFailureFromExperiment is a helper function that creates a failure from an experiment and a node.
func createFailureFromExperiment(experiment *chaosv1.Experiment, node *clusterv1.Node) *chaosv1.Failure {
	// TODO: better random ID

	flr := chaosv1.NewFailure()
	flr.Metadata.ID = fmt.Sprintf(fmtFailureID,
		experiment.Metadata.ID,
		node.Metadata.ID,
		randID(),
	)
	flr.Metadata.NodeID = node.Metadata.ID
	flr.Spec = experiment.Spec.Template.Spec
	flr.Status.CurrentState = chaosv1.DisabledFailureState
	flr.Status.ExpectedState = chaosv1.EnabledFailureState
	flr.Status.Creation = time.Now().UTC()

	return &flr
}

// Scheduler is the interface that knows how to schedule the failures in different
// nodes.
type Scheduler interface {
	// Schedule will schedule an experiment and create the appropiate failures. It will
	// return the created failures.
	Schedule(experiment *chaosv1.Experiment) ([]*chaosv1.Failure, error)
}

// NodeLabels is an scheduler that will schedule the failures based on the labels of an experiment.
// appart from returning the failures it will store them on the repository based on the required node.
type NodeLabels struct {
	nodeRepo    repository.Node
	failureRepo repository.Failure
	logger      log.Logger
}

// NewNodeLabels will return a new NodeLabels scheduler.
func NewNodeLabels(failureRepo repository.Failure, nodeRepo repository.Node, logger log.Logger) *NodeLabels {
	return &NodeLabels{
		nodeRepo:    nodeRepo,
		failureRepo: failureRepo,
		logger:      logger,
	}
}

// Schedule satisfies Scheduler interface.
func (n *NodeLabels) Schedule(experiment *chaosv1.Experiment) ([]*chaosv1.Failure, error) {
	flrs := []*chaosv1.Failure{}

	// Get all the nodes of the experiment based on the experiment labels.
	nodes := n.nodeRepo.GetNodesByLabels(experiment.Spec.Selector)

	// TODO: Check failure already running on node.
	for _, node := range nodes {
		// Create the node and save.
		flr := createFailureFromExperiment(experiment, node)
		flrs = append(flrs, flr)
		err := n.failureRepo.Store(flr)
		if err != nil {
			return flrs, err
		}
	}
	n.logger.WithField("experiment", experiment.Metadata.ID).Infof("scheduled %d failures", len(flrs))

	// TODO: Save experiment, or update required fields.
	return flrs, nil
}
