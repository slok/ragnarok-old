package failure

import (
	"errors"
	"testing"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/stretchr/testify/assert"
)

func TestGoodReadConfig(t *testing.T) {
	assert := assert.New(t)
	config := `
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

	c, err := ReadConfig([]byte(config))
	if assert.NoError(err, "YAML unmarshalling shouldn't return an error") {
		expectedC := Config{
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
		assert.EqualValues(expectedC, c, "Configuration values should be equal after loading YAML definition")
	}
}

func TestBadReadConfig(t *testing.T) {
	assert := assert.New(t)
	config := `
timeout: 1h
attacks:
  - attack1:
      size: 524288000
  back-attack:
  	something: 12
`
	_, err := ReadConfig([]byte(config))
	assert.Error(err, "YAML unmarshalling should return an error")
}

func TestMultipleAttacksOnMapReadConfig(t *testing.T) {
	assert := assert.New(t)
	config := `
timeout: 1h
attacks:
  - attack1:
      size: 524288000
    bad-attack:
      key: value
  - attack2:
      something: 12
`
	_, err := ReadConfig([]byte(config))
	if assert.Error(err, "YAML unmarshalling should return an attack format error") {
		assert.Equal(err, errors.New("attacks format error, tip: check identantion and '-' indicator"))
	}

}

func TestGoodRenderConfig(t *testing.T) {
	assert := assert.New(t)

	c := Config{
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
		assert.EqualValues(expectY, string(d), "YAML values should be equal after redering config")
	}
}
