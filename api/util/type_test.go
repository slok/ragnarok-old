package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/api/util"
)

func TestGetFullID(t *testing.T) {
	tests := []struct {
		name      string
		obj       api.Object
		expFullID string
	}{
		{
			name: "Node fullname",
			obj: &clusterv1.Node{
				TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
				Metadata: api.ObjectMeta{ID: "node1"},
			},
			expFullID: "cluster/v1/node/node1",
		},
		{

			name: "Experiment fullname",
			obj: &chaosv1.Experiment{
				TypeMeta: api.TypeMeta{Kind: chaosv1.ExperimentKind, Version: chaosv1.ExperimentVersion},
				Metadata: api.ObjectMeta{ID: "exp1"},
			},
			expFullID: "chaos/v1/experiment/exp1",
		},
		{
			name: "Failure fullname",
			obj: &chaosv1.Failure{
				TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
				Metadata: api.ObjectMeta{ID: "flr1"},
			},
			expFullID: "chaos/v1/failure/flr1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.expFullID, util.GetFullID(test.obj))
		})
	}
}

func TestGetFullIDFromType(t *testing.T) {
	tests := []struct {
		name      string
		typeMeta  api.TypeMeta
		id        string
		expFullID string
	}{
		{
			name:      "Node fullname",
			typeMeta:  api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
			id:        "node1",
			expFullID: "cluster/v1/node/node1",
		},
		{

			name:      "Experiment fullname",
			typeMeta:  api.TypeMeta{Kind: chaosv1.ExperimentKind, Version: chaosv1.ExperimentVersion},
			id:        "exp1",
			expFullID: "chaos/v1/experiment/exp1",
		},
		{
			name:      "Failure fullname",
			typeMeta:  api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
			id:        "flr1",
			expFullID: "chaos/v1/failure/flr1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.expFullID, util.GetFullIDFromType(test.typeMeta, test.id))
		})
	}
}

func TestGetFullType(t *testing.T) {
	tests := []struct {
		name      string
		obj       api.Object
		expFullID string
	}{
		{
			name: "Node fullname",
			obj: &clusterv1.Node{
				TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
				Metadata: api.ObjectMeta{ID: "node1"},
			},
			expFullID: "cluster/v1/node",
		},
		{

			name: "Experiment fullname",
			obj: &chaosv1.Experiment{
				TypeMeta: api.TypeMeta{Kind: chaosv1.ExperimentKind, Version: chaosv1.ExperimentVersion},
				Metadata: api.ObjectMeta{ID: "exp1"},
			},
			expFullID: "chaos/v1/experiment",
		},
		{
			name: "Failure fullname",
			obj: &chaosv1.Failure{
				TypeMeta: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
				Metadata: api.ObjectMeta{ID: "flr1"},
			},
			expFullID: "chaos/v1/failure",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.expFullID, util.GetFullType(test.obj))
		})
	}
}

func TestSplitFullID(t *testing.T) {
	tests := []struct {
		name    string
		fullID  string
		expType api.TypeMeta
		expID   string
	}{
		{
			name:    "Node type",
			fullID:  "cluster/v1/node/node1",
			expType: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
			expID:   "node1",
		},
		{
			name:    "Experiment type",
			fullID:  "chaos/v1/experiment/exp1",
			expType: api.TypeMeta{Kind: chaosv1.ExperimentKind, Version: chaosv1.ExperimentVersion},
			expID:   "exp1",
		},
		{
			name:    "Failure type",
			fullID:  "chaos/v1/failure/flr1",
			expType: api.TypeMeta{Kind: chaosv1.FailureKind, Version: chaosv1.FailureVersion},
			expID:   "flr1",
		},
		{
			name:    "Wrong type",
			fullID:  "chaos/v1/failure/exp1/wrong",
			expType: api.TypeMeta{},
			expID:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			tt, id := util.SplitFullID(test.fullID)
			assert.Equal(test.expType, tt)
			assert.Equal(test.expID, id)
		})
	}
}
