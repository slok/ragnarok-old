package experiment_test

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/experiment"
	mclichaosv1 "github.com/slok/ragnarok/mocks/client/api/chaos/v1"
	mcliclusterv1 "github.com/slok/ragnarok/mocks/client/api/cluster/v1"
)

func TestEnsureFailures(t *testing.T) {
	mockCreationTime := time.Now()

	tests := []struct {
		name             string
		nodes            *clusterv1.NodeList
		failures         *chaosv1.FailureList
		experiment       *chaosv1.Experiment
		expDeletedFlrIDs []string
		expNewFlrs       []*chaosv1.Failure
		expErr           bool
	}{
		{
			name: "A new Experiment should create failures to all available nodes that match the selector",
			nodes: &clusterv1.NodeList{
				Items: []*clusterv1.Node{
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode0"},
					},
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode1"},
					},
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode2"},
					},
				},
			},
			failures: &chaosv1.FailureList{},
			experiment: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "exp-001",
					Annotations: map[string]string{
						"name":        "first experiment",
						"description": "first experiment is the first experiment :|",
					},
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{},
				},
			},
			expDeletedFlrIDs: []string{},
			expNewFlrs: []*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-0",
						Labels: map[string]string{
							"experiment": "exp-001",
							"node":       "testNode0",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-1",
						Labels: map[string]string{
							"experiment": "exp-001",
							"node":       "testNode1",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-2",
						Labels: map[string]string{
							"experiment": "exp-001",
							"node":       "testNode2",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
			},
			expErr: false,
		},
		{
			name: "An already created experiment should create failures only on the nodes that match the selector and that don't have already the failure",
			nodes: &clusterv1.NodeList{
				Items: []*clusterv1.Node{
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode0"},
					},
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode1"},
					},
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode2"},
					},
				},
			},
			failures: &chaosv1.FailureList{
				Items: []*chaosv1.Failure{
					&chaosv1.Failure{
						TypeMeta: chaosv1.FailureTypeMeta,
						Metadata: api.ObjectMeta{
							ID: "flrid-x",
							Labels: map[string]string{
								api.LabelExperiment: "exp-001",
								api.LabelNode:       "testNode0",
							},
						},
						Status: chaosv1.FailureStatus{
							CurrentState:  4,
							ExpectedState: 1,
							Creation:      mockCreationTime,
						},
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "exp-001",
					Annotations: map[string]string{
						"name":        "first experiment",
						"description": "first experiment is the first experiment :|",
					},
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{},
				},
			},
			expDeletedFlrIDs: []string{},
			expNewFlrs: []*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-0",
						Labels: map[string]string{
							api.LabelExperiment: "exp-001",
							api.LabelNode:       "testNode1",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-1",
						Labels: map[string]string{
							api.LabelExperiment: "exp-001",
							api.LabelNode:       "testNode2",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
			},
			expErr: false,
		},
		{
			name: "An already created experiment should create failures only on the nodes that match the selector and that don't have already the failure and delete the ones that don't have nodes",
			nodes: &clusterv1.NodeList{
				Items: []*clusterv1.Node{
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode0"},
					},
					&clusterv1.Node{
						Metadata: api.ObjectMeta{ID: "testNode1"},
					},
				},
			},
			failures: &chaosv1.FailureList{
				Items: []*chaosv1.Failure{
					&chaosv1.Failure{
						TypeMeta: chaosv1.FailureTypeMeta,
						Metadata: api.ObjectMeta{
							ID: "flrid-x",
							Labels: map[string]string{
								api.LabelExperiment: "exp-001",
								api.LabelNode:       "testNode0",
							},
						},
						Status: chaosv1.FailureStatus{
							CurrentState:  4,
							ExpectedState: 1,
							Creation:      mockCreationTime,
						},
					},
					&chaosv1.Failure{
						TypeMeta: chaosv1.FailureTypeMeta,
						Metadata: api.ObjectMeta{
							ID: "flrid-y",
							Labels: map[string]string{
								api.LabelExperiment: "exp-001",
								api.LabelNode:       "testNode2",
							},
						},
						Status: chaosv1.FailureStatus{
							CurrentState:  4,
							ExpectedState: 1,
							Creation:      mockCreationTime,
						},
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "exp-001",
					Annotations: map[string]string{
						"name":        "first experiment",
						"description": "first experiment is the first experiment :|",
					},
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{},
				},
			},
			expDeletedFlrIDs: []string{
				"chaos/v1/failure/flrid-y",
			},
			expNewFlrs: []*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flrid-0",
						Labels: map[string]string{
							api.LabelExperiment: "exp-001",
							api.LabelNode:       "testNode1",
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  4,
						ExpectedState: 1,
						Creation:      mockCreationTime,
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockFlrIDCtr := 0

			assert := assert.New(t)

			gotDeletedFlrIDs := []string{}
			gotNewFlrs := []*chaosv1.Failure{}

			// mocks.
			mnCli := &mcliclusterv1.NodeClientInterface{}
			mnCli.On("List", mock.Anything).Return(test.nodes, nil)
			mfCli := &mclichaosv1.FailureClientInterface{}
			mfCli.On("List", mock.Anything).Return(test.failures, nil)

			// Grab actions made on mocks to assert afterwards.
			mfCli.On("Delete", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				id := args.Get(0).(string)
				gotDeletedFlrIDs = append(gotDeletedFlrIDs, id)
			})
			mfCli.On("Create", mock.Anything).Return(nil, nil).Run(func(args mock.Arguments) {
				flr := args.Get(0).(*chaosv1.Failure)
				// Mock ids for equal assertion.
				flr.Metadata.ID = fmt.Sprintf("flrid-%d", mockFlrIDCtr)
				mockFlrIDCtr++
				flr.Status.Creation = mockCreationTime
				gotNewFlrs = append(gotNewFlrs, flr)
			})

			// Create the experiment manager and ensure the failures.
			sm := experiment.NewSimpleManager(mnCli, mfCli, log.Dummy)
			err := sm.EnsureFailures(test.experiment)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {

				// Sort before asserting them
				sort.Slice(gotDeletedFlrIDs, func(i, j int) bool {
					return gotDeletedFlrIDs[i] < gotDeletedFlrIDs[j]
				})

				sort.Slice(gotNewFlrs, func(i, j int) bool {
					return gotNewFlrs[i].Metadata.ID < gotNewFlrs[j].Metadata.ID
				})

				// Assert.
				assert.Equal(test.expDeletedFlrIDs, gotDeletedFlrIDs)
				assert.Equal(test.expNewFlrs, gotNewFlrs)
			}
		})
	}
}
