package experiment

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	apiutil "github.com/slok/ragnarok/api/util"
	clichaosv1 "github.com/slok/ragnarok/client/api/chaos/v1"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/log"
)

const (
	fmtFailureID = "%s-flr-%s"
)

// Manager has the methods so it can set the correct state on
// the experiment resource.
type Manager interface {
	// EnsureFailures ensures the required failures should be on the cluster based on the experiment
	// this means cleaning up the ones that are not needed and creating the ones that need to be created.
	EnsureFailures(*chaosv1.Experiment) error
	// EnableFailures enables the failures of an experiment.
	EnableFailures(*chaosv1.Experiment) error
	// DisableFailures disables the failures of an experiment.
	DisableFailures(*chaosv1.Experiment) error
	// DeleteFailures deletes the failures of an experiment.
	DeleteFailures(*chaosv1.Experiment) error
}

// SimpleManager is the state manager that will use the controller to
// set the correct state on the experiments.
type SimpleManager struct {
	nodeCli    cliclusterv1.NodeClientInterface
	failureCli clichaosv1.FailureClientInterface
	logger     log.Logger
}

// NewSimpleManager returns a new experiment simple manager.
func NewSimpleManager(nodeCli cliclusterv1.NodeClientInterface, failureCli clichaosv1.FailureClientInterface, logger log.Logger) *SimpleManager {
	return &SimpleManager{
		nodeCli:    nodeCli,
		failureCli: failureCli,
		logger:     logger,
	}
}

// getNodes will get an experiment nodes based on the settings of the expriment, for example based on
// label selector.
func (s *SimpleManager) getNodes(exp *chaosv1.Experiment) (map[string]*clusterv1.Node, error) {
	res := map[string]*clusterv1.Node{}
	opts := api.ListOptions{
		LabelSelector: exp.Spec.Selector,
	}
	// TODO: Only get ready nodes.
	nodes, err := s.nodeCli.List(opts)
	if err != nil {
		return res, err
	}

	for _, node := range nodes {
		node := node
		res[node.Metadata.ID] = node
	}

	return res, nil
}

// getFailures gets the failures already created by an experiment.
func (s *SimpleManager) getFailures(exp *chaosv1.Experiment) ([]*chaosv1.Failure, error) {
	opts := api.ListOptions{
		LabelSelector: map[string]string{
			api.LabelExperiment: exp.Metadata.ID,
		},
	}

	return s.failureCli.List(opts)
}

// indexedFailuresByNode will get a list of failures and set in a map identified by the assigned node.
func (s *SimpleManager) indexedFailuresByNode(flrs []*chaosv1.Failure) map[string]*chaosv1.Failure {
	res := map[string]*chaosv1.Failure{}
	for _, flr := range flrs {
		flr := flr
		node, ok := flr.Metadata.Labels[api.LabelNode]
		if ok {
			res[node] = flr
		}
	}
	return res
}

// failureRandID creates a new failure random ID.
func (s *SimpleManager) failureRandID(exp *chaosv1.Experiment) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := fmt.Sprintf("%09d", r.Intn(999999999))
	return fmt.Sprintf(fmtFailureID, exp.Metadata.ID, id)
}

// createFailureFromExperiment is a helper function that creates a failure from an experiment and a node.
func (s *SimpleManager) createFailureFromExperiment(exp *chaosv1.Experiment, node *clusterv1.Node) *chaosv1.Failure {
	// TODO: better random ID
	flr := chaosv1.NewFailure()
	flr.Metadata.ID = s.failureRandID(exp)
	flr.Metadata.Labels = map[string]string{
		api.LabelExperiment: exp.Metadata.ID,
		api.LabelNode:       node.Metadata.ID,
	}
	flr.Spec = exp.Spec.Template.Spec
	flr.Status.CurrentState = chaosv1.DisabledFailureState
	flr.Status.ExpectedState = chaosv1.EnabledFailureState
	flr.Status.Creation = time.Now().UTC()

	return &flr
}

// gcFailures will delete the failures from the nodes that not longer exists on the cluster.
func (s *SimpleManager) gcFailures(exp *chaosv1.Experiment, flrsByNode map[string]*chaosv1.Failure, nodes map[string]*clusterv1.Node) error {
	logger := s.logger.With("experiment", exp.Metadata.ID)
	for nodeID, failure := range flrsByNode {
		// If not existent node then delete failure.
		if _, ok := nodes[nodeID]; !ok {
			logger.Debugf("node dissapeared %s: deleting failure %s", nodeID, failure.Metadata.ID)
			id := apiutil.GetFullID(failure)
			s.failureCli.Delete(id)
		}
	}

	return nil
}

// createAndScheduleFailures will create the required failures and schedule on the required nodes.
func (s *SimpleManager) createAndScheduleFailures(exp *chaosv1.Experiment, flrsByNode map[string]*chaosv1.Failure, nodes map[string]*clusterv1.Node) error {
	logger := s.logger.With("experiment", exp.Metadata.ID)
	for nodeID, node := range nodes {
		// if not present then create a failure on the node.
		if _, ok := flrsByNode[nodeID]; !ok {
			flr := s.createFailureFromExperiment(exp, node)
			if _, err := s.failureCli.Create(flr); err != nil {
				return fmt.Errorf("could not create failure: %s", err)
			}
			logger.Debugf("new failure %s created for node %s", flr.Metadata.ID, nodeID)
		}
	}

	return nil
}

// EnsureFailures will check that the actual experiment has the number of failures to be
// injected. This is based on the settings. At this moment the selector of node will mark how many
// failures will need to be made and schedule on the required nodes.
//
// This are the steps made to ensure a correct state:
// 1 - Check the failures that need to be deleted.
// 2 - Delete the ones that have asigned a failure on a non existent node.
// 3 - Get the difference between the desired number and the actual number.
// 4 - If is the desired ones are less this means that they need to create failures and assign the required node.
func (s *SimpleManager) EnsureFailures(exp *chaosv1.Experiment) error {
	// Get the selector and get the nodes.
	nodes, err := s.getNodes(exp)
	if err != nil {
		return err
	}
	// Get how many failures are already created form this experiment.
	flrs, err := s.getFailures(exp)
	if err != nil {
		return err
	}

	flrsByNode := s.indexedFailuresByNode(flrs)

	// Garbage collection of failures from this experiment.
	if err := s.gcFailures(exp, flrsByNode, nodes); err != nil {
		return fmt.Errorf("error on failure garbage collection: %s", err)
	}

	// Create the required failures and schedule them.
	if err := s.createAndScheduleFailures(exp, flrsByNode, nodes); err != nil {
		return fmt.Errorf("error on failure garbage collection: %s", err)
	}

	return nil
}
func (s *SimpleManager) EnableFailures(exp *chaosv1.Experiment) error {
	return fmt.Errorf("not implmeneted")
}
func (s *SimpleManager) DisableFailures(exp *chaosv1.Experiment) error {
	return fmt.Errorf("not implmeneted")
}
func (s *SimpleManager) DeleteFailures(exp *chaosv1.Experiment) error {
	return fmt.Errorf("not implmeneted")
}
