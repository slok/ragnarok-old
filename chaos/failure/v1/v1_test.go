package v1_test

import (
	"errors"
	"testing"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/chaos/failure/v1"
	"github.com/stretchr/testify/assert"
)

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

	f, err := v1.ReadFailure([]byte(spec))
	if assert.NoError(err, "YAML unmarshalling shouldn't return an error") {
		expectedF := v1.Failure{
			Spec: v1.Spec{
				Timeout: 1 * time.Hour,
				Attacks: []v1.AttackMap{
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
	_, err := v1.ReadFailure([]byte(spec))
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
	_, err := v1.ReadFailure([]byte(spec))
	if assert.Error(err, "YAML unmarshalling should return an attack format error") {
		assert.Equal(err, errors.New("attacks format error, tip: check identantion and '-' indicator"))
	}

}

func TestGoodRenderDefinition(t *testing.T) {
	assert := assert.New(t)

	f := v1.Failure{
		Spec: v1.Spec{
			Timeout: 30 * time.Second,
			Attacks: []v1.AttackMap{
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

	expectFSpec := `spec:
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
	f := v1.Failure{
		Spec: v1.Spec{
			Timeout: 30 * time.Second,
			Attacks: []v1.AttackMap{
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
