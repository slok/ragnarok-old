package failure

import (
	"errors"
	"testing"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/stretchr/testify/assert"
)

func TestGoodReadDefinition(t *testing.T) {
	assert := assert.New(t)
	definition := `
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

	d, err := ReadDefinition([]byte(definition))
	if assert.NoError(err, "YAML unmarshalling shouldn't return an error") {
		expectedD := Definition{
			Timeout: 1 * time.Hour,
			Attacks: []AttackMap{
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
		}
		assert.EqualValues(expectedD, d, "Definitionuration values should be equal after loading YAML definition")
	}
}

func TestBadReadDefinition(t *testing.T) {
	assert := assert.New(t)
	definition := `
timeout: 1h
attacks:
  - attack1:
      size: 524288000
  back-attack:
  	something: 12
`
	_, err := ReadDefinition([]byte(definition))
	assert.Error(err, "YAML unmarshalling should return an error")
}

func TestMultipleAttacksOnMapReadDefinition(t *testing.T) {
	assert := assert.New(t)
	definition := `
timeout: 1h
attacks:
  - attack1:
      size: 524288000
    bad-attack:
      key: value
  - attack2:
      something: 12
`
	_, err := ReadDefinition([]byte(definition))
	if assert.Error(err, "YAML unmarshalling should return an attack format error") {
		assert.Equal(err, errors.New("attacks format error, tip: check identantion and '-' indicator"))
	}

}

func TestGoodRenderDefinition(t *testing.T) {
	assert := assert.New(t)

	c := Definition{
		Timeout: 30 * time.Second,
		Attacks: []AttackMap{
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
	}

	d, err := c.Render()

	expectY := `timeout: 30s
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
		assert.EqualValues(expectY, string(d), "YAML values should be equal after redering definition")
	}
}

func TestMultipleAttacksOnMapRenderDefinition(t *testing.T) {
	assert := assert.New(t)
	c := Definition{
		Timeout: 30 * time.Second,
		Attacks: []AttackMap{
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
	}

	_, err := c.Render()
	if assert.Error(err, "YAML marshalling should return an attack format error") {
		assert.Equal(err, errors.New("each attack map of the attack list needs to be a single map"))
	}

}
