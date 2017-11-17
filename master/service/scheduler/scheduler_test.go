package scheduler_test

import (
	"sort"
	"testing"
	"time"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/scheduler"
	mcliclusterv1 "github.com/slok/ragnarok/mocks/client/cluster/v1"
	mrepository "github.com/slok/ragnarok/mocks/master/service/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodeLabelsSchedule(t *testing.T) {

	tests := []struct {
		name        string
		nodes       []*clusterv1.Node
		experiment  *chaosv1.Experiment
		expFailures []*chaosv1.Failure
		expErr      bool
	}{
		{
			name:        "Scheduling on a missing node should return 0 failures",
			nodes:       []*clusterv1.Node{},
			experiment:  &chaosv1.Experiment{},
			expFailures: []*chaosv1.Failure{},
			expErr:      false,
		},
		{
			name: "Scheduling on a single node should return 1 failures",
			nodes: []*clusterv1.Node{
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node1",
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{ID: "exp-001"},
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
					Metadata: api.ObjectMeta{
						Labels: map[string]string{
							api.LabelNode:       "node1",
							api.LabelExperiment: "exp-001",
						},
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
			nodes: []*clusterv1.Node{
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node1",
					},
				},
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node2",
					},
				},
				&clusterv1.Node{
					Metadata: api.ObjectMeta{
						ID: "node3",
					},
				},
			},
			experiment: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{ID: "exp-001"},
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
					Metadata: api.ObjectMeta{
						Labels: map[string]string{
							api.LabelNode:       "node1",
							api.LabelExperiment: "exp-001",
						},
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
					Metadata: api.ObjectMeta{
						Labels: map[string]string{
							api.LabelNode:       "node2",
							api.LabelExperiment: "exp-001",
						},
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
					Metadata: api.ObjectMeta{
						Labels: map[string]string{
							api.LabelNode:       "node3",
							api.LabelExperiment: "exp-001",
						},
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
			mcn := &mcliclusterv1.Node{}
			mcn.On("List", mock.Anything).Return(test.nodes, nil)
			mfr := &mrepository.Failure{}
			mfr.On("Store", mock.Anything).Return(nil)

			s := scheduler.NewNodeLabels(mfr, mcn, log.Dummy)
			flrs, err := s.Schedule(test.experiment)

			if test.expErr {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					sort.Slice(flrs, func(i, j int) bool {
						return flrs[i].Metadata.Labels[api.LabelNode] < flrs[j].Metadata.Labels[api.LabelNode]
					})

					for i, expFlr := range test.expFailures {
						assert.Equal(expFlr.Spec, flrs[i].Spec)
						assert.Equal(expFlr.Metadata.Labels[api.LabelNode], flrs[i].Metadata.Labels[api.LabelNode])
						assert.Equal(expFlr.Status.CurrentState, flrs[i].Status.CurrentState)
						assert.Equal(expFlr.Status.ExpectedState, flrs[i].Status.ExpectedState)
					}
				}
			}
		})
	}
}
