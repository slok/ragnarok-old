package v1_test

import (
	"bytes"
	"errors"
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

func TestJSONEncodeChaosV1Failure(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t2, _ := time.Parse(time.RFC3339, "2012-11-01T22:18:41+00:00")
	t3, _ := time.Parse(time.RFC3339, "2012-11-01T22:28:41+00:00")

	tests := []struct {
		name          string
		failure       *chaosv1.Failure
		expEncFailure string
		expErr        bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			failure: &chaosv1.Failure{
				Metadata: chaosv1.FailureMetadata{
					ID:     "flr-001",
					NodeID: "nd-034",
				},
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
				Status: chaosv1.FailureStatus{
					CurrentState:  chaosv1.EnabledFailureState,
					ExpectedState: chaosv1.DisabledFailureState,
					Creation:      t1,
					Executed:      t2,
					Finished:      t3,
				},
			},
			expEncFailure: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001","nodeid":"nd-034"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			expErr:        false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failure: &chaosv1.Failure{
				TypeMeta: api.TypeMeta{
					Kind:    chaosv1.FailureKind,
					Version: chaosv1.FailureVersion,
				},
				Metadata: chaosv1.FailureMetadata{
					ID:     "flr-001",
					NodeID: "nd-034",
				},
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
				Status: chaosv1.FailureStatus{
					CurrentState:  chaosv1.EnabledFailureState,
					ExpectedState: chaosv1.DisabledFailureState,
					Creation:      t1,
					Executed:      t2,
					Finished:      t3,
				},
			},
			expEncFailure: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001","nodeid":"nd-034"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			expErr:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjTyper, apimachinery.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.failure, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncFailure, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestJSONDecodeChaosV1Failure(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t2s := "2012-11-01T22:18:41Z"
	t3s := "2012-11-01T22:28:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)
	t2, _ := time.Parse(time.RFC3339, t2s)
	t3, _ := time.Parse(time.RFC3339, t3s)

	tests := []struct {
		name        string
		failureJSON string
		expFailure  *chaosv1.Failure
		expErr      bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			failureJSON: `
{  
   "kind":"failure",
   "version":"chaos/v1",
   "metadata":{  
      "id":"flr-001",
      "nodeid":"nd-034"
   },
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
   },
   "status":{  
      "currentState":1,
      "expectedState":4,
      "creation":"2012-11-01T22:08:41Z",
      "executed":"2012-11-01T22:18:41Z",
      "finished":"2012-11-01T22:28:41Z"
   }
}`,
			expFailure: &chaosv1.Failure{
				TypeMeta: api.TypeMeta{
					Kind:    chaosv1.FailureKind,
					Version: chaosv1.FailureVersion,
				},
				Metadata: chaosv1.FailureMetadata{
					ID:     "flr-001",
					NodeID: "nd-034",
				},
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
								"pace":     "10m",
								"quantity": float64(10),
								"rest":     "30s",
								"target":   "myTarget",
							},
						},
					},
				},
				Status: chaosv1.FailureStatus{
					CurrentState:  chaosv1.EnabledFailureState,
					ExpectedState: chaosv1.DisabledFailureState,
					Creation:      t1,
					Executed:      t2,
					Finished:      t3,
				},
			},
			expErr: false,
		},
		{
			name:        "Simple object decoding without kind or version should return an error",
			failureJSON: ``,
			expFailure:  &chaosv1.Failure{},
			expErr:      true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjTyper, apimachinery.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.failureJSON))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				failure := obj.(*chaosv1.Failure)
				assert.Equal(test.expFailure, failure)
			}
		})
	}
}

// TODO: LEGACY CODE NEEDS NO BE MIGRATED

func TestFailureStateStringer(t *testing.T) {
	tests := []struct {
		st    chaosv1.FailureState
		expSt string
	}{
		{chaosv1.EnabledFailureState, "enabled"},
		{chaosv1.ExecutingFailureState, "executing"},
		{chaosv1.RevertingFailureState, "reverting"},
		{chaosv1.DisabledFailureState, "disabled"},
		{chaosv1.StaleFailureState, "stale"},
		{chaosv1.ErroredFailureState, "errored"},
		{chaosv1.ErroredRevertingFailureState, "erroredreverting"},
		{chaosv1.UnknownFailureState, "unknown"},
		{99999, "unknown"},
	}

	for _, test := range tests {
		t.Run(test.expSt, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.expSt, test.st.String())
		})
	}
}

func TestGoodReadFailure(t *testing.T) {
	assert := assert.New(t)
	spec := `
spec:
  timeout: 1h
  attacks:
    - attack1:
        size: 524288000
    - attack1:
        size: 100
    - attack2:
    - attack3:
        target: myTarget
        quantity: 10
        pace: 10m
        rest: 30s
`

	f, err := chaosv1.ReadFailure([]byte(spec))
	if assert.NoError(err, "YAML unmarshalling shouldn't return an error") {
		expectedF := chaosv1.Failure{
			TypeMeta: api.TypeMeta{},
			Spec: chaosv1.FailureSpec{
				Timeout: 1 * time.Hour,
				Attacks: []chaosv1.AttackMap{
					{
						"attack1": attack.Opts{
							"size": 524288000,
						},
					},
					{
						"attack1": attack.Opts{
							"size": 100,
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
		}
		assert.EqualValues(expectedF, f, "Failure values should be equal after loading YAML definition")
	}
}

func TestBadReadFailure(t *testing.T) {
	assert := assert.New(t)
	spec := `
spec:
  timeout: 1h
  attacks:
    - attack1:
        size: 524288000
    back-attack:
    	something: 12
`
	_, err := chaosv1.ReadFailure([]byte(spec))
	assert.Error(err, "YAML unmarshalling should return an error")
}

func TestMultipleAttacksOnMapReadFailure(t *testing.T) {
	assert := assert.New(t)
	spec := `
spec:
  timeout: 1h
  attacks:
    - attack1:
        size: 524288000
      bad-attack:
        key: value
    - attack2:
        something: 12
`
	_, err := chaosv1.ReadFailure([]byte(spec))
	if assert.Error(err, "YAML unmarshalling should return an attack format error") {
		assert.Equal(err, errors.New("attacks format error, tip: check identantion and '-' indicator"))
	}

}

func TestGoodRenderDefinition(t *testing.T) {
	assert := assert.New(t)

	f := chaosv1.Failure{
		Spec: chaosv1.FailureSpec{
			Timeout: 30 * time.Second,
			Attacks: []chaosv1.AttackMap{
				{
					"attack1": attack.Opts{
						"size": 524288000,
					},
				},
				{
					"attack1": attack.Opts{
						"size": 100,
					},
				},
				{
					"attack3": attack.Opts{
						"pace":     "10m",
						"quantity": 10,
						"rest":     "30s",
						"target":   "myTarget",
					},
				},
			},
		},
	}

	fSpec, err := f.Render()

	expectFSpec := `kind: ""
version: ""
spec:
  timeout: 30s
  attacks:
  - attack1:
      size: 524288000
  - attack1:
      size: 100
  - attack3:
      pace: 10m
      quantity: 10
      rest: 30s
      target: myTarget
`

	if assert.NoError(err, "YAML marshalling shouldn't return an error") {
		assert.EqualValues(expectFSpec, string(fSpec), "YAML values should be equal after redering definition")
	}
}

func TestMultipleAttacksOnMapRenderFailure(t *testing.T) {
	assert := assert.New(t)
	f := chaosv1.Failure{
		Spec: chaosv1.FailureSpec{
			Timeout: 30 * time.Second,
			Attacks: []chaosv1.AttackMap{
				{
					"attack1": attack.Opts{
						"size": 524288000,
					},
					"wrong-attack": attack.Opts{
						"size": 524288000,
					},
				},
				{
					"attack1": attack.Opts{
						"size": 100,
					},
				},
			},
		},
	}

	_, err := f.Render()
	if assert.Error(err, "YAML marshalling should return an attack format error") {
		assert.Equal(err, errors.New("each attack map of the attack list needs to be a single map"))
	}
}
