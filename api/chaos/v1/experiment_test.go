package v1_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/apimachinery"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/log"
)

func TestJSONEncodeChaosV1Experiment(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

	tests := []struct {
		name       string
		experiment *chaosv1.Experiment
		expEncNode string
		expErr     bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			experiment: &chaosv1.Experiment{
				Metadata: chaosv1.ExperimentMetadata{
					ID:          "exp-001",
					Name:        "first experiment",
					Description: " first experiment is the first experiment :|",
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"size": 524288000,
									},
								},
								{
									"attack2": nil,
								},
								{
									"attack3": attack.Opts{
										"target":   "myTarget",
										"quantity": 10,
										"pace":     "10m",
										"rest":     "30s",
									},
								},
							},
						},
					},
				},
				Status: chaosv1.ExperimentStatus{
					FailureIDs: []string{"node1", "node3", "node4"},
					Creation:   t1,
				},
			},
			expEncNode: `{"kind":"experiment","version":"chaos/v1","metadata":{"id":"exp-001","name":"first experiment","description":" first experiment is the first experiment :|"},"spec":{"selector":{"az":"eu-west-1a","kind":"master"},"template":{"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]}}},"status":{"failureIDs":["node1","node3","node4"],"creation":"2012-11-01T22:08:41Z"}}`,
			expErr:     false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			experiment: &chaosv1.Experiment{
				TypeMeta: api.TypeMeta{
					Kind:    chaosv1.ExperimentKind,
					Version: chaosv1.ExperimentVersion,
				},
				Metadata: chaosv1.ExperimentMetadata{
					ID:          "exp-001",
					Name:        "first experiment",
					Description: " first experiment is the first experiment :|",
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"size": 524288000,
									},
								},
								{
									"attack2": nil,
								},
								{
									"attack3": attack.Opts{
										"target":   "myTarget",
										"quantity": 10,
										"pace":     "10m",
										"rest":     "30s",
									},
								},
							},
						},
					},
				},
				Status: chaosv1.ExperimentStatus{
					FailureIDs: []string{"node1", "node3", "node4"},
					Creation:   t1,
				},
			},
			expEncNode: `{"kind":"experiment","version":"chaos/v1","metadata":{"id":"exp-001","name":"first experiment","description":" first experiment is the first experiment :|"},"spec":{"selector":{"az":"eu-west-1a","kind":"master"},"template":{"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]}}},"status":{"failureIDs":["node1","node3","node4"],"creation":"2012-11-01T22:08:41Z"}}`,
			expErr:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjTyper, apimachinery.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.experiment, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncNode, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestJSONDecodeChaosV1Experiment(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)

	tests := []struct {
		name           string
		experimentJSON string
		expExperiment  *chaosv1.Experiment
		expErr         bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			experimentJSON: `
{
   "kind":"experiment",
   "version":"chaos/v1",
   "metadata":{
      "id":"exp-001",
      "name":"first experiment",
      "description":" first experiment is the first experiment :|"
   },
   "spec":{
      "selector":{
         "az":"eu-west-1a",
         "kind":"master"
      },
      "template":{
         "spec":{
            "timeout":300000000000,
            "attacks":[
               {
                  "attack1":{
                     "size":524288000
                  }
               },
               {
                  "attack2":null
               },
               {
                  "attack3":{
                     "pace":"10m",
                     "quantity":10,
                     "rest":"30s",
                     "target":"myTarget"
                  }
               }
            ]
         }
      }
   },
   "status":{
      "failureIDs":[
         "node1",
         "node3",
         "node4"
      ],
      "creation":"2012-11-01T22:08:41Z"
   }
}`,
			expExperiment: &chaosv1.Experiment{
				TypeMeta: api.TypeMeta{
					Kind:    chaosv1.ExperimentKind,
					Version: chaosv1.ExperimentVersion,
				},
				Metadata: chaosv1.ExperimentMetadata{
					ID:          "exp-001",
					Name:        "first experiment",
					Description: " first experiment is the first experiment :|",
				},
				Spec: chaosv1.ExperimentSpec{
					Selector: map[string]string{"kind": "master", "az": "eu-west-1a"},
					Template: chaosv1.ExperimentFailureTemplate{
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"size": float64(524288000),
									},
								},
								{
									"attack2": nil,
								},
								{
									"attack3": attack.Opts{
										"target":   "myTarget",
										"quantity": float64(10),
										"pace":     "10m",
										"rest":     "30s",
									},
								},
							},
						},
					},
				},
				Status: chaosv1.ExperimentStatus{
					FailureIDs: []string{"node1", "node3", "node4"},
					Creation:   t1,
				},
			},
			expErr: false,
		},
		{
			name:           "Simple object decoding without kind or version should return an error",
			experimentJSON: ``,
			expExperiment:  &chaosv1.Experiment{},
			expErr:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjTyper, apimachinery.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.experimentJSON))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				experiment := obj.(*chaosv1.Experiment)
				assert.Equal(test.expExperiment, experiment)
			}
		})
	}
}