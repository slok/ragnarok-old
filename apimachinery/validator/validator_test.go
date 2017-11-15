package validator_test

import (
	"testing"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateNode(t *testing.T) {
	tests := []struct {
		name       string
		node       *clusterv1.Node
		expInvalid bool
	}{
		{
			name: "A correct node should not return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "node1",
					Labels: map[string]string{},
				},
			},
			expInvalid: false,
		},
		{
			name: "Not Id on a node should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID:     "",
					Labels: map[string]string{},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label key should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label key should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label key should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"my label !": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label value should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"something": "",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label value should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"node1": "qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label value should return an error.",
			node: &clusterv1.Node{
				Metadata: api.ObjectMeta{
					ID: "node1",
					Labels: map[string]string{
						"something": "my label !",
					},
				},
			},
			expInvalid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			ov := validator.NewObject()
			errs := ov.Validate(test.node)

			if test.expInvalid {
				assert.NotEmpty(errs)
			} else {
				assert.Empty(errs)
			}
		})
	}
}

func TestValidateFailure(t *testing.T) {
	tests := []struct {
		name       string
		failure    *chaosv1.Failure
		expInvalid bool
	}{
		{
			name: "A correct failure should not return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: false,
		},
		{
			name: "Not Id on a failure should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "",
					Labels: map[string]string{
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label key should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"":                  "something",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label key should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui": "something",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label key should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"my label !":        "something",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label value should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"something":         "",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label value should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"failure1":          "qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label value should return an error.",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"something":         "my label !",
						api.LabelExperiment: "exp1",
						api.LabelNode:       "node1",
					},
				},
			},
			expInvalid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			ov := validator.NewObject()
			errs := ov.Validate(test.failure)

			if test.expInvalid {
				assert.NotEmpty(errs)
			} else {
				assert.Empty(errs)
			}
		})
	}
}

func TestValidateExperiment(t *testing.T) {
	tests := []struct {
		name       string
		failure    *chaosv1.Experiment
		expInvalid bool
	}{
		{
			name: "A correct failure should not return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID:     "failure1",
					Labels: map[string]string{},
				},
			},
			expInvalid: false,
		},
		{
			name: "Not Id on a failure should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID:     "",
					Labels: map[string]string{},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label key should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label key should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label key should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"my label !": "something",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An empty label value should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"something": "",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "An big label value should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"failure1": "qwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyuiqwertyuiopqawertyuiopqwertyuiopqwertyuiopwertyuiopqwertyuiopwertyuiopwertyuioqwertyui",
					},
				},
			},
			expInvalid: true,
		},
		{
			name: "Invalid characters on a label value should return an error.",
			failure: &chaosv1.Experiment{
				Metadata: api.ObjectMeta{
					ID: "failure1",
					Labels: map[string]string{
						"something": "my label !",
					},
				},
			},
			expInvalid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			ov := validator.NewObject()
			errs := ov.Validate(test.failure)

			if test.expInvalid {
				assert.NotEmpty(errs)
			} else {
				assert.Empty(errs)
			}
		})
	}
}
