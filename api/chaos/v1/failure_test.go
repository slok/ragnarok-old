package v1_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
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
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			expErr:        false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failure: &chaosv1.Failure{
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			expErr:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
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

func TestYAMLEncodeChaosV1Failure(t *testing.T) {
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
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: "kind: failure\nmetadata:\n  id: flr-001\nspec:\n  attacks:\n  - attack1:\n      size: 524288000\n  - attack2: null\n  - attack3:\n      pace: 10m\n      quantity: 10\n      rest: 30s\n      target: myTarget\n  timeout: 300000000000\nstatus:\n  creation: 2012-11-01T22:08:41Z\n  currentState: 1\n  executed: 2012-11-01T22:18:41Z\n  expectedState: 4\n  finished: 2012-11-01T22:28:41Z\nversion: chaos/v1",
			expErr:        false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failure: &chaosv1.Failure{
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: "kind: failure\nmetadata:\n  id: flr-001\nspec:\n  attacks:\n  - attack1:\n      size: 524288000\n  - attack2: null\n  - attack3:\n      pace: 10m\n      quantity: 10\n      rest: 30s\n      target: myTarget\n  timeout: 300000000000\nstatus:\n  creation: 2012-11-01T22:08:41Z\n  currentState: 1\n  executed: 2012-11-01T22:18:41Z\n  expectedState: 4\n  finished: 2012-11-01T22:28:41Z\nversion: chaos/v1",
			expErr:        false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
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
      "id":"flr-001"
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
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			name: "Simple object decoding without kind or version should return an error",
			failureJSON: `
{  
   "metadata":{  
      "id":"flr-001"
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
}
`,
			expFailure: &chaosv1.Failure{},
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
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

func TestYAMLDecodeChaosV1Failure(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t2s := "2012-11-01T22:18:41Z"
	t3s := "2012-11-01T22:28:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)
	t2, _ := time.Parse(time.RFC3339, t2s)
	t3, _ := time.Parse(time.RFC3339, t3s)

	tests := []struct {
		name        string
		failureYAML string
		expFailure  *chaosv1.Failure
		expErr      bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			failureYAML: `
kind: failure
version: chaos/v1
metadata:
  id: flr-001
spec:
  timeout: 300000000000
  attacks:
  - attack1:
      size: 524288000
  - attack2: null
  - attack3:
      pace: 10m
      quantity: 10
      rest: 30s
      target: myTarget
status:
  currentState: 1
  expectedState: 4
  creation: 2012-11-01T22:08:41Z
  executed: 2012-11-01T22:18:41Z
  finished: 2012-11-01T22:28:41Z`,
			expFailure: &chaosv1.Failure{
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			name: "Simple object decoding without kind or version should return an error",
			failureYAML: `
metadata:
  id: flr-001
spec:
  timeout: 300000000000
  attacks:
  - attack1:
      size: 524288000
  - attack2: null
  - attack3:
      pace: 10m
      quantity: 10
      rest: 30s
      target: myTarget
status:
  currentState: 1
  expectedState: 4
  creation: 2012-11-01T22:08:41Z
  executed: 2012-11-01T22:18:41Z
  finished: 2012-11-01T22:28:41Z
`,
			expFailure: &chaosv1.Failure{},
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.failureYAML))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				failure := obj.(*chaosv1.Failure)
				assert.Equal(test.expFailure, failure)
			}
		})
	}
}

func TestPBEncodeChaosV1Failure(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t2, _ := time.Parse(time.RFC3339, "2012-11-01T22:18:41+00:00")
	t3, _ := time.Parse(time.RFC3339, "2012-11-01T22:28:41+00:00")

	tests := []struct {
		name          string
		failure       *chaosv1.Failure
		expEncFailure *chaosv1pb.Failure
		expErr        bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			failure: &chaosv1.Failure{
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: &chaosv1pb.Failure{
				SerializedData: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			},

			expErr: false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failure: &chaosv1.Failure{
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			expEncFailure: &chaosv1pb.Failure{
				SerializedData: `{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}`,
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewPBSerializer(log.Dummy)
			pbflr := &chaosv1pb.Failure{}
			err := s.Encode(test.failure, pbflr)

			if test.expErr {
				assert.Error(err)
			} else {
				// Small fix for the \n
				pbflr.SerializedData = strings.TrimSuffix(pbflr.SerializedData, "\n")
				assert.Equal(test.expEncFailure, pbflr)
				assert.NoError(err)
			}
		})
	}
}

func TestPBDecodeChaosV1Failure(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t2s := "2012-11-01T22:18:41Z"
	t3s := "2012-11-01T22:28:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)
	t2, _ := time.Parse(time.RFC3339, t2s)
	t3, _ := time.Parse(time.RFC3339, t3s)

	tests := []struct {
		name       string
		failurePB  *chaosv1pb.Failure
		expFailure *chaosv1.Failure
		expErr     bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			failurePB: &chaosv1pb.Failure{
				SerializedData: `
{  
   "kind":"failure",
   "version":"chaos/v1",
   "metadata":{  
      "id":"flr-001"
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
			},
			expFailure: &chaosv1.Failure{
				TypeMeta: chaosv1.FailureTypeMeta,
				Metadata: api.ObjectMeta{
					ID: "flr-001",
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
			name:       "Simple object decoding without kind or version should return an error",
			failurePB:  &chaosv1pb.Failure{SerializedData: ``},
			expFailure: &chaosv1.Failure{},
			expErr:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewPBSerializer(log.Dummy)
			obj, err := s.Decode(test.failurePB)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				failure := obj.(*chaosv1.Failure)
				assert.Equal(test.expFailure, failure)
			}
		})
	}
}

func TestJSONEncodeChaosV1FailureList(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t2, _ := time.Parse(time.RFC3339, "2012-11-01T22:18:41+00:00")
	t3, _ := time.Parse(time.RFC3339, "2012-11-01T22:28:41+00:00")

	tests := []struct {
		name              string
		failureList       chaosv1.FailureList
		expEncFailureList string
		expErr            bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			failureList: chaosv1.FailureList{
				ListMetadata: api.ListMeta{
					Continue: "123454321",
				},
				Items: []*chaosv1.Failure{
					&chaosv1.Failure{
						Metadata: api.ObjectMeta{
							ID: "flr-001",
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
					&chaosv1.Failure{
						Metadata: api.ObjectMeta{
							ID: "flr-002",
						},
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"host":    "eu-west-1.aws.amazon.com",
										"timeout": "2m",
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
				},
			},
			expEncFailureList: `{"kind":"failureList","version":"chaos/v1","listMetadata":{"continue":"123454321"},"items":[{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}},{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-002"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"host":"eu-west-1.aws.amazon.com","timeout":"2m"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}]}`,
			expErr:            false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failureList: chaosv1.NewFailureList([]*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-001",
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
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-002",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"host":    "eu-west-1.aws.amazon.com",
									"timeout": "2m",
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
			},
				"123454321"),
			expEncFailureList: `{"kind":"failureList","version":"chaos/v1","listMetadata":{"continue":"123454321"},"items":[{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-001"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"size":524288000}},{"attack2":null},{"attack3":{"pace":"10m","quantity":10,"rest":"30s","target":"myTarget"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}},{"kind":"failure","version":"chaos/v1","metadata":{"id":"flr-002"},"spec":{"timeout":300000000000,"attacks":[{"attack1":{"host":"eu-west-1.aws.amazon.com","timeout":"2m"}}]},"status":{"currentState":1,"expectedState":4,"creation":"2012-11-01T22:08:41Z","executed":"2012-11-01T22:18:41Z","finished":"2012-11-01T22:28:41Z"}}]}`,
			expErr:            false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(&test.failureList, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncFailureList, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestJSONDecodeChaosV1FailureList(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t2s := "2012-11-01T22:18:41Z"
	t3s := "2012-11-01T22:28:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)
	t2, _ := time.Parse(time.RFC3339, t2s)
	t3, _ := time.Parse(time.RFC3339, t3s)

	tests := []struct {
		name            string
		failureListJSON string
		expFailureList  chaosv1.FailureList
		expErr          bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			failureListJSON: `
{
   "kind":"failureList",
   "version":"chaos/v1",
   "listMetadata":{
      "continue":"123454321"
   },
   "items":[
      {
         "kind":"failure",
         "version":"chaos/v1",
         "metadata":{
            "id":"flr-001"
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
      },
      {
         "kind":"failure",
         "version":"chaos/v1",
         "metadata":{
            "id":"flr-002"
         },
         "spec":{
            "timeout":300000000000,
            "attacks":[
               {
                  "attack1":{
                     "host":"eu-west-1.aws.amazon.com",
                     "timeout":"2m"
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
      }
   ]
}
`,
			expFailureList: chaosv1.NewFailureList([]*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-001",
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
									"target":   "myTarget",
									"quantity": float64(10),
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
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-002",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"host":    "eu-west-1.aws.amazon.com",
									"timeout": "2m",
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
			},
				"123454321"),
			expErr: false,
		},
		{
			name: "Simple object decoding without kind or version should return an error",
			failureListJSON: `
{
   "listMetadata":{
      "continue":"123454321"
   },
   "items":[
      {
         "kind":"failure",
         "version":"chaos/v1",
         "metadata":{
            "id":"flr-001"
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
      },
      {
         "kind":"failure",
         "version":"chaos/v1",
         "metadata":{
            "id":"flr-002"
         },
         "spec":{
            "timeout":300000000000,
            "attacks":[
               {
                  "attack1":{
                     "host":"eu-west-1.aws.amazon.com",
                     "timeout":"2m"
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
      }
   ]
}
`,
			expFailureList: chaosv1.FailureList{},
			expErr:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.failureListJSON))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				failureList := obj.(*chaosv1.FailureList)
				assert.Equal(&test.expFailureList, failureList)
			}
		})
	}
}

func TestYAMLEncodeChaosV1FailureList(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
	t2, _ := time.Parse(time.RFC3339, "2012-11-01T22:18:41+00:00")
	t3, _ := time.Parse(time.RFC3339, "2012-11-01T22:28:41+00:00")

	tests := []struct {
		name              string
		failureList       chaosv1.FailureList
		expEncFailureList string
		expErr            bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			failureList: chaosv1.FailureList{
				ListMetadata: api.ListMeta{
					Continue: "123454321",
				},
				Items: []*chaosv1.Failure{
					&chaosv1.Failure{
						Metadata: api.ObjectMeta{
							ID: "flr-001",
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
					&chaosv1.Failure{
						Metadata: api.ObjectMeta{
							ID: "flr-002",
						},
						Spec: chaosv1.FailureSpec{
							Timeout: 5 * time.Minute,
							Attacks: []chaosv1.AttackMap{
								{
									"attack1": attack.Opts{
										"host":    "eu-west-1.aws.amazon.com",
										"timeout": "2m",
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
				},
			},
			expEncFailureList: "items:\n- kind: failure\n  metadata:\n    id: flr-001\n  spec:\n    attacks:\n    - attack1:\n        size: 524288000\n    - attack2: null\n    - attack3:\n        pace: 10m\n        quantity: 10\n        rest: 30s\n        target: myTarget\n    timeout: 300000000000\n  status:\n    creation: 2012-11-01T22:08:41Z\n    currentState: 1\n    executed: 2012-11-01T22:18:41Z\n    expectedState: 4\n    finished: 2012-11-01T22:28:41Z\n  version: chaos/v1\n- kind: failure\n  metadata:\n    id: flr-002\n  spec:\n    attacks:\n    - attack1:\n        host: eu-west-1.aws.amazon.com\n        timeout: 2m\n    timeout: 300000000000\n  status:\n    creation: 2012-11-01T22:08:41Z\n    currentState: 1\n    executed: 2012-11-01T22:18:41Z\n    expectedState: 4\n    finished: 2012-11-01T22:28:41Z\n  version: chaos/v1\nkind: failureList\nlistMetadata:\n  continue: \"123454321\"\nversion: chaos/v1",
			expErr:            false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			failureList: chaosv1.NewFailureList([]*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-001",
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
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-002",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"host":    "eu-west-1.aws.amazon.com",
									"timeout": "2m",
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
			},
				"123454321"),
			expEncFailureList: "items:\n- kind: failure\n  metadata:\n    id: flr-001\n  spec:\n    attacks:\n    - attack1:\n        size: 524288000\n    - attack2: null\n    - attack3:\n        pace: 10m\n        quantity: 10\n        rest: 30s\n        target: myTarget\n    timeout: 300000000000\n  status:\n    creation: 2012-11-01T22:08:41Z\n    currentState: 1\n    executed: 2012-11-01T22:18:41Z\n    expectedState: 4\n    finished: 2012-11-01T22:28:41Z\n  version: chaos/v1\n- kind: failure\n  metadata:\n    id: flr-002\n  spec:\n    attacks:\n    - attack1:\n        host: eu-west-1.aws.amazon.com\n        timeout: 2m\n    timeout: 300000000000\n  status:\n    creation: 2012-11-01T22:08:41Z\n    currentState: 1\n    executed: 2012-11-01T22:18:41Z\n    expectedState: 4\n    finished: 2012-11-01T22:28:41Z\n  version: chaos/v1\nkind: failureList\nlistMetadata:\n  continue: \"123454321\"\nversion: chaos/v1",
			expErr:            false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(&test.failureList, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncFailureList, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestYAMLDecodeChaosV1FailureList(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t2s := "2012-11-01T22:18:41Z"
	t3s := "2012-11-01T22:28:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)
	t2, _ := time.Parse(time.RFC3339, t2s)
	t3, _ := time.Parse(time.RFC3339, t3s)

	tests := []struct {
		name            string
		failureListYAML string
		expFailureList  chaosv1.FailureList
		expErr          bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			failureListYAML: `
kind: failureList
version: chaos/v1
listMetadata:
  continue: "123454321"
items:
- kind: failure
  version: chaos/v1
  metadata:
    id: flr-001
  spec:
    attacks:
    - attack1:
        size: 524288000
    - attack2: null
    - attack3:
        pace: 10m
        quantity: 10
        rest: 30s
        target: myTarget
    timeout: 300000000000
  status:
    creation: 2012-11-01T22:08:41Z
    currentState: 1
    executed: 2012-11-01T22:18:41Z
    expectedState: 4
    finished: 2012-11-01T22:28:41Z
- kind: failure
  version: chaos/v1
  metadata:
    id: flr-002
  spec:
    attacks:
    - attack1:
        host: eu-west-1.aws.amazon.com
        timeout: 2m
    timeout: 300000000000
  status:
    creation: 2012-11-01T22:08:41Z
    currentState: 1
    executed: 2012-11-01T22:18:41Z
    expectedState: 4
    finished: 2012-11-01T22:28:41Z
`,
			expFailureList: chaosv1.NewFailureList([]*chaosv1.Failure{
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-001",
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
									"target":   "myTarget",
									"quantity": float64(10),
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
				&chaosv1.Failure{
					TypeMeta: chaosv1.FailureTypeMeta,
					Metadata: api.ObjectMeta{
						ID: "flr-002",
					},
					Spec: chaosv1.FailureSpec{
						Timeout: 5 * time.Minute,
						Attacks: []chaosv1.AttackMap{
							{
								"attack1": attack.Opts{
									"host":    "eu-west-1.aws.amazon.com",
									"timeout": "2m",
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
			},
				"123454321"),
			expErr: false,
		},
		{
			name: "Simple object decoding without kind or version should return an error",
			failureListYAML: `
listMetadata:
  continue: "123454321"
items:
- kind: failure
  version: chaos/v1
  metadata:
    id: flr-001
  spec:
    attacks:
    - attack1:
        size: 524288000
    - attack2: null
    - attack3:
        pace: 10m
        quantity: 10
        rest: 30s
        target: myTarget
    timeout: 300000000000
  status:
    creation: 2012-11-01T22:08:41Z
    currentState: 1
    executed: 2012-11-01T22:18:41Z
    expectedState: 4
    finished: 2012-11-01T22:28:41Z
- kind: failure
  version: chaos/v1
  metadata:
    id: flr-002
  spec:
    attacks:
    - attack1:
        host: eu-west-1.aws.amazon.com
        timeout: 2m
    timeout: 300000000000
  status:
    creation: 2012-11-01T22:08:41Z
    currentState: 1
    executed: 2012-11-01T22:18:41Z
    expectedState: 4
    finished: 2012-11-01T22:28:41Z
`,
			expFailureList: chaosv1.FailureList{},
			expErr:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.failureListYAML))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				failureList := obj.(*chaosv1.FailureList)
				assert.Equal(&test.expFailureList, failureList)
			}
		})
	}
}
