package failure_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/failure"
	"github.com/slok/ragnarok/mocks"
)

func TestNewSystemFailure(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
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
	f, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.NoError(err) {
		assert.NotNil(f, "A succesful creation shoudln't be an error")
		assert.Equal(failure.Created, f.State)
		reg.AssertExpectations(t)
	}

}

func TestNewSystemFailureError(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
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
	_, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	reg.AssertNotCalled(t, "New", "attack3", c.Attacks[2]["attack3"])
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}

}

func TestNewSystemFailureMultipleAttacksOnBlock(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
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
	_, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}
}

func TestSystemFailureFailState(t *testing.T) {
	assert := assert.New(t)
	expectedErr := fmt.Errorf("invalid state. The only valid state for execution is: %s", failure.Created)
	tests := []struct {
		state         failure.State
		expectedErr   error
		expectedState failure.State
	}{
		{
			state:         failure.Created,
			expectedErr:   nil,
			expectedState: failure.Executing,
		},
		{
			state:         failure.Executing,
			expectedErr:   expectedErr,
			expectedState: failure.Executing,
		},
		{
			state:         failure.Reverted,
			expectedErr:   expectedErr,
			expectedState: failure.Reverted,
		},
		{
			state:         failure.Error,
			expectedErr:   expectedErr,
			expectedState: failure.Error,
		},
		{
			state:         failure.ErrorReverting,
			expectedErr:   expectedErr,
			expectedState: failure.ErrorReverting,
		},
		{
			state:         failure.Unknown,
			expectedErr:   expectedErr,
			expectedState: failure.Unknown,
		},
	}

	for _, test := range tests {
		f, err := failure.NewSystemFailure(failure.Config{}, nil, nil)
		f.State = test.state
		if assert.NoError(err) {
			err = f.Fail()
			assert.Equal(test.expectedErr, err)
			assert.Equal(test.expectedState, f.State, "Expected state should be '%s', got: '%s'", test.expectedState, f.State)
		}
	}
}

func TestSystemFailureAttacksOK(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
			{"attack2": attack.Opts{}},
			{"attack3": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at := &mocks.Attacker{}
	at.On("Apply", ctxMatcher).Times(3).Return(nil)

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at, nil)

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.NoError(err) {
		if assert.NoError(f.Fail()) {
			assert.Equal(failure.Executing, f.State)
			at.AssertExpectations(t)
		}
	}
}

func TestSystemFailureAttacksOKRevertOK(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
			{"attack2": attack.Opts{}},
			{"attack3": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at := &mocks.Attacker{}
	at.On("Apply", ctxMatcher).Times(3).Return(nil)
	at.On("Revert").Times(3).Return(nil)

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at, nil)

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.NoError(err) {
		if assert.NoError(f.Fail()) {
			assert.Equal(failure.Executing, f.State)
			if assert.NoError(f.Revert()) {
				assert.Equal(failure.Reverted, f.State)
				at.AssertExpectations(t)
			}
		}
	}
}

func TestSystemFailureFailAttacksErrorAutoRevertOK(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
			{"attack2": attack.Opts{}},
			{"attack3": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1, at2, at3 := &mocks.Attacker{}, &mocks.Attacker{}, &mocks.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(errors.New("error1"))
	at2.On("Apply", ctxMatcher).Once().Return(nil)
	at3.On("Apply", ctxMatcher).Once().Return(errors.New("error3"))
	// the only one that needs to be reverted is the 2nd one because is the only one that applied correctly
	at2.On("Revert").Once().Return(nil)

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at2, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at3, nil)

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.NoError(err) {
		err = f.Fail()
		if assert.Error(err) {
			assert.Equal(errors.New("error aplying failure"), err)
			assert.Equal(failure.Error, f.State)
			at1.AssertExpectations(t)
			at2.AssertExpectations(t)
			at3.AssertExpectations(t)
		}
	}
}

func TestSystemFailureFailAttacksErrorAutoRevertError(t *testing.T) {
	assert := assert.New(t)
	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
			{"attack2": attack.Opts{}},
			{"attack3": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1, at2, at3 := &mocks.Attacker{}, &mocks.Attacker{}, &mocks.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	at2.On("Apply", ctxMatcher).Once().Return(errors.New("error2"))
	at3.On("Apply", ctxMatcher).Once().Return(nil)
	at1.On("Revert").Once().Return(nil)
	at3.On("Revert").Once().Return(errors.New("revert_error3"))

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at2, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at3, nil)

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, nil)
	if assert.NoError(err) {
		err = f.Fail()
		if assert.Error(err) {
			assert.Equal(errors.New("error aplying failure & error when trying to revert the applied ones"), err)
			assert.Equal(failure.ErrorReverting, f.State)
			at1.AssertExpectations(t)
			at2.AssertExpectations(t)
			at3.AssertExpectations(t)
		}
	}
}

func TestSystemFailureFailAttacksFinishWithTimeout(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1 := &mocks.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	reverted := make(chan struct{})
	at1.On("Revert").Once().Return(nil).Run(func(args mock.Arguments) {
		reverted <- struct{}{}
	})

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)

	// Mock clock
	cl := &mocks.Clock{}
	cl.On("After", c.Timeout).Return(time.After(0))
	cl.On("Now").Return(time.Now())

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, cl)
	if assert.NoError(err) {
		err = f.Fail()
		// Wait until clock timeout to check revert called
		select {
		case <-clock.After(5 * time.Millisecond):
			require.Fail("Timeout calling revert after a timeout")
		case <-reverted:
			at1.AssertExpectations(t)
		}

	}
}

func TestSystemFailureFailAttacksFinishForced(t *testing.T) {
	assert := assert.New(t)

	c := failure.Config{
		Timeout: 1 * time.Hour,
		Attacks: []failure.AttackMap{
			{"attack1": attack.Opts{}},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1 := &mocks.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	at1.On("Revert").Once().Return(nil)

	// Mock Registry
	reg := &mocks.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)

	// Mock clock
	cl := &mocks.Clock{}
	cl.On("After", c.Timeout).Return(time.After(9999 * time.Hour)) // Never
	cl.On("Now").Return(time.Now())

	f, err := failure.NewSystemFailureFromReg(c, reg, nil, cl)
	if assert.NoError(err) {
		err = f.Fail()
		f.Revert()
		at1.AssertExpectations(t)
	}
}
