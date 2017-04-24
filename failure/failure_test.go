package failure

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/mocks"
)

func TestNewSystemFailure(t *testing.T) {
	assert := assert.New(t)
	c := Config{
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
	at1 := &mocks.Attacker{}
	at2 := &mocks.Attacker{}
	at3 := &mocks.Attacker{}
	at4 := &mocks.Attacker{}
	// Mock registry.
	reg := &mocks.Registry{}
	reg.On("New", "attack1", c.Attacks[0]["attack1"]).Return(at1, nil)
	reg.On("New", "attack1", c.Attacks[1]["attack1"]).Return(at2, nil)
	reg.On("New", "attack2", c.Attacks[2]["attack2"]).Return(at3, nil)
	reg.On("New", "attack3", c.Attacks[3]["attack3"]).Return(at4, nil)

	// Test.
	f, err := NewSystemFailureFromReg(c, reg, nil)
	if assert.NoError(err) {
		assert.NotNil(f, "A succesful creation shoudln't be an error")
		reg.AssertExpectations(t)
		expectAtts := []attack.Attacker{at1, at2, at3, at4}
		assert.EqualValues(expectAtts, f.attacks, "Expected attacks are not the expected ones") // Check same values
		// Check same object in memory
		for i, a := range expectAtts {
			assert.True(a == f.attacks[i], "Attack pointer is not the same as expected")
		}
	}

}

func TestNewSystemFailureError(t *testing.T) {
	assert := assert.New(t)
	c := Config{
		Timeout: 1 * time.Hour,
		Attacks: []AttackMap{
			{
				"attack1": attack.Opts{
					"size": 524288000,
				},
			},
			{
				"attack2": nil,
			},
			{
				"attack3": nil,
			},
		},
	}

	// Mock registry.
	reg := &mocks.Registry{}
	reg.On("New", "attack1", c.Attacks[0]["attack1"]).Return(nil, nil)
	reg.On("New", "attack2", c.Attacks[1]["attack2"]).Return(nil, errors.New("error test"))

	// Test.
	_, err := NewSystemFailureFromReg(c, reg, nil)
	reg.AssertNotCalled(t, "New", "attack3", c.Attacks[2]["attack3"])
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}

}

func TestNewSystemFailureMultipleAttacksOnBlock(t *testing.T) {
	assert := assert.New(t)
	c := Config{
		Timeout: 1 * time.Hour,
		Attacks: []AttackMap{
			{
				"attack1": attack.Opts{
					"size": 524288000,
				},
				"attack2": nil,
			},
			{
				"attack3": nil,
			},
		},
	}

	// Mock registry.
	reg := &mocks.Registry{}
	reg.AssertNotCalled(t, "New")
	// Test.
	_, err := NewSystemFailureFromReg(c, reg, nil)
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}

}
