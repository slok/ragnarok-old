package scheduler

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clichaosv1 "github.com/slok/ragnarok/client/api/chaos/v1"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/log"
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
	flr.Metadata.Labels = map[string]string{
		api.LabelNode:       node.Metadata.ID,
		api.LabelExperiment: experiment.Metadata.ID,
	}
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
	nodecli    cliclusterv1.NodeClientInterface
	failurecli clichaosv1.FailureClientInterface
	logger     log.Logger
}

// NewNodeLabels will return a new NodeLabels scheduler.
func NewNodeLabels(failurecli clichaosv1.FailureClientInterface, nodecli cliclusterv1.NodeClientInterface, logger log.Logger) *NodeLabels {
	return &NodeLabels{
		nodecli:    nodecli,
		failurecli: failurecli,
		logger:     logger,
	}
}

// Schedule satisfies Scheduler interface.
func (n *NodeLabels) Schedule(experiment *chaosv1.Experiment) ([]*chaosv1.Failure, error) {
	flrs := []*chaosv1.Failure{}

	// Get all the nodes of the experiment based on the experiment labels.
	opts := api.ListOptions{
		LabelSelector: experiment.Spec.Selector,
	}
	nodes, err := n.nodecli.List(opts)
	if err != nil {
		return flrs, err
	}

	// TODO: Check failure already running on node.
	for _, node := range nodes {
		// Create the node and save.
		flr := createFailureFromExperiment(experiment, node)
		flrs = append(flrs, flr)
		_, err := n.failurecli.Create(flr)
		if err != nil {
			return flrs, err
		}
	}
	n.logger.WithField("experiment", experiment.Metadata.ID).Infof("scheduled %d failures", len(flrs))

	// TODO: Save experiment, or update required fields.
	return flrs, nil
}
