package scheduler_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/scheduler"
	mservice "github.com/slok/ragnarok/mocks/master/service"
)

func TestNodeLabelsSchedule(t *testing.T) {

	tests := []struct {
		name        string
		nodes       map[string]*clusterv1.Node
		experiment  *chaosv1.Experiment
		expFailures []*chaosv1.Failure
		expErr      bool
	}{
		{
			name:        "Scheduling on a missing node should return 0 failures",
			nodes:       map[string]*clusterv1.Node{},
			experiment:  &chaosv1.Experiment{},
			expFailures: []*chaosv1.Failure{},
			expErr:      false,
		},
		{
			name: "Scheduling on a single node should return 1 failures",
			nodes: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: clusterv1.NodeMetadata{ID: "node1"},
					Spec: clusterv1.NodeSpec{
						Labels: map[string]string{"ID": "node1"},
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: chaosv1.ExperimentMetadata{ID: "exp-001"},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"ID": "node1"},
					Template: chaosv1.ExperimentFailureTemplate{
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"size": 524288000,
									},
								},
							},
						},
					},
				},
			},
			expFailures: []*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
					Metadata: chaosv1.FailureMetadata{
						NodeID: "node1",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.DisabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
			},
			expErr: false,
		},
		{
			name: "Scheduling on a multiple nodes should return multiple failures",
			nodes: map[string]*clusterv1.Node{
				"node1": &clusterv1.Node{
					Metadata: clusterv1.NodeMetadata{ID: "node1"},
					Spec: clusterv1.NodeSpec{
						Labels: map[string]string{"ID": "node1"},
					},
				},
				"node2": &clusterv1.Node{
					Metadata: clusterv1.NodeMetadata{ID: "node2"},
					Spec: clusterv1.NodeSpec{
						Labels: map[string]string{"ID": "node2"},
					},
				},
				"node3": &clusterv1.Node{
					Metadata: clusterv1.NodeMetadata{ID: "node3"},
					Spec: clusterv1.NodeSpec{
						Labels: map[string]string{"ID": "node3"},
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: chaosv1.ExperimentMetadata{ID: "exp-001"},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"ID": "node1"},
					Template: chaosv1.ExperimentFailureTemplate{
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"size": 524288000,
									},
								},
							},
						},
					},
				},
			},
			expFailures: []*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
					Metadata: chaosv1.FailureMetadata{
						NodeID: "node1",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.DisabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
					Metadata: chaosv1.FailureMetadata{
						NodeID: "node2",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.DisabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
					Metadata: chaosv1.FailureMetadata{
						NodeID: "node3",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.DisabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Create mocks.
			mnr := &mservice.NodeRepository{}
			mnr.On("GetNodesByLabels", mock.Anything).Return(test.nodes)
			mfr := &mservice.FailureRepository{}
			mfr.On("Store", mock.Anything).Return(nil)

			s := scheduler.NewNodeLabels(mfr, mnr, log.Dummy)
			flrs, err := s.Schedule(test.experiment)

			if test.expErr {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					for i, expFlr := range test.expFailures {
						assert.Equal(expFlr.Spec, flrs[i].Spec)
						assert.Equal(expFlr.Metadata.NodeID, flrs[i].Metadata.NodeID)
						assert.Equal(expFlr.Status.CurrentState, flrs[i].Status.CurrentState)
						assert.Equal(expFlr.Status.ExpectedState, flrs[i].Status.ExpectedState)
					}
				}
			}
		})
	}
}
