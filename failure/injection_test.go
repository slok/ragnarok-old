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
	mattack "github.com/slok/ragnarok/mocks/attack"
	mclock "github.com/slok/ragnarok/mocks/clock"
	"github.com/slok/ragnarok/types"
)

func TestNewInjection(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
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
		},
	}
	at1 := &mattack.Attacker{}
	at2 := &mattack.Attacker{}
	at3 := &mattack.Attacker{}
	at4 := &mattack.Attacker{}
	// Mock registry.
	reg := &mattack.Registry{}
	reg.On("New", "attack1", f.Definition.Attacks[0]["attack1"]).Return(at1, nil)
	reg.On("New", "attack1", f.Definition.Attacks[1]["attack1"]).Return(at2, nil)
	reg.On("New", "attack2", f.Definition.Attacks[2]["attack2"]).Return(at3, nil)
	reg.On("New", "attack3", f.Definition.Attacks[3]["attack3"]).Return(at4, nil)

	// Test.
	ij, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.NoError(err) {
		assert.NotNil(ij, "A succesful creation shoudln't be an error")
		assert.Equal(types.EnabledFailureState, ij.CurrentState)
		reg.AssertExpectations(t)
	}

}

func TestNewInjectionError(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
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
		},
	}

	// Mock registry.
	reg := &mattack.Registry{}
	reg.On("New", "attack1", f.Definition.Attacks[0]["attack1"]).Return(nil, nil)
	reg.On("New", "attack2", f.Definition.Attacks[1]["attack2"]).Return(nil, errors.New("error test"))

	// Test.
	_, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	reg.AssertNotCalled(t, "New", "attack3", f.Definition.Attacks[2]["attack3"])
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}

}

func TestNewInjectionMultipleAttacksOnBlock(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
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
		},
	}

	// Mock registry.
	reg := &mattack.Registry{}
	reg.AssertNotCalled(t, "New")
	// Test.
	_, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.Error(err) {
		reg.AssertExpectations(t)
	}
}

func TestInjectionFailState(t *testing.T) {
	assert := assert.New(t)
	expectedErr := fmt.Errorf("invalid state. The only valid state for execution is: %s", types.EnabledFailureState)
	tests := []struct {
		state         types.FailureState
		expectedErr   error
		expectedState types.FailureState
	}{
		{
			state:         types.EnabledFailureState,
			expectedErr:   nil,
			expectedState: types.ExecutingFailureState,
		},
		{
			state:         types.ExecutingFailureState,
			expectedErr:   expectedErr,
			expectedState: types.ExecutingFailureState,
		},
		{
			state:         types.DisabledFailureState,
			expectedErr:   expectedErr,
			expectedState: types.DisabledFailureState,
		},
		{
			state:         types.ErroredFailureState,
			expectedErr:   expectedErr,
			expectedState: types.ErroredFailureState,
		},
		{
			state:         types.ErroredRevertingFailureState,
			expectedErr:   expectedErr,
			expectedState: types.ErroredRevertingFailureState,
		},
		{
			state:         types.UnknownFailureState,
			expectedErr:   expectedErr,
			expectedState: types.UnknownFailureState,
		},
	}

	for _, test := range tests {
		in, err := failure.NewInjection(&failure.Failure{}, nil, nil)
		in.CurrentState = test.state
		if assert.NoError(err) {
			err = in.Fail()
			assert.Equal(test.expectedErr, err)
			assert.Equal(test.expectedState, in.CurrentState, "Expected state should be '%s', got: '%s'", test.expectedState, in.CurrentState)
		}
	}
}

func TestSystemFailureAttacksOK(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
				{"attack2": attack.Opts{}},
				{"attack3": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at := &mattack.Attacker{}
	at.On("Apply", ctxMatcher).Times(3).Return(nil)

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at, nil)

	in, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.NoError(err) {
		if assert.NoError(in.Fail()) {
			assert.Equal(types.ExecutingFailureState, in.CurrentState)
			at.AssertExpectations(t)
		}
	}
}

func TestSystemFailureAttacksOKRevertOK(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
				{"attack2": attack.Opts{}},
				{"attack3": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at := &mattack.Attacker{}
	at.On("Apply", ctxMatcher).Times(3).Return(nil)
	at.On("Revert").Times(3).Return(nil)

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at, nil)

	in, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.NoError(err) {
		if assert.NoError(in.Fail()) {
			assert.Equal(types.ExecutingFailureState, in.CurrentState)
			if assert.NoError(in.Revert()) {
				assert.Equal(types.DisabledFailureState, in.CurrentState)
				at.AssertExpectations(t)
			}
		}
	}
}

func TestSystemFailureFailAttacksErrorAutoRevertOK(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
				{"attack2": attack.Opts{}},
				{"attack3": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1, at2, at3 := &mattack.Attacker{}, &mattack.Attacker{}, &mattack.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(errors.New("error1"))
	at2.On("Apply", ctxMatcher).Once().Return(nil)
	at3.On("Apply", ctxMatcher).Once().Return(errors.New("error3"))
	// the only one that needs to be reverted is the 2nd one because is the only one that applied correctly
	at2.On("Revert").Once().Return(nil)

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at2, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at3, nil)

	in, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.NoError(err) {
		err = in.Fail()
		if assert.Error(err) {
			assert.Equal(errors.New("error aplying failure"), err)
			assert.Equal(types.ErroredFailureState, in.CurrentState)
			at1.AssertExpectations(t)
			at2.AssertExpectations(t)
			at3.AssertExpectations(t)
		}
	}
}

func TestSystemFailureFailAttacksErrorAutoRevertError(t *testing.T) {
	assert := assert.New(t)
	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
				{"attack2": attack.Opts{}},
				{"attack3": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1, at2, at3 := &mattack.Attacker{}, &mattack.Attacker{}, &mattack.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	at2.On("Apply", ctxMatcher).Once().Return(errors.New("error2"))
	at3.On("Apply", ctxMatcher).Once().Return(nil)
	at1.On("Revert").Once().Return(nil)
	at3.On("Revert").Once().Return(errors.New("revert_error3"))

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)
	reg.On("New", "attack2", attack.Opts{}).Return(at2, nil)
	reg.On("New", "attack3", attack.Opts{}).Return(at3, nil)

	in, err := failure.NewInjectionFromReg(f, reg, nil, nil)
	if assert.NoError(err) {
		err = in.Fail()
		if assert.Error(err) {
			assert.Equal(errors.New("error aplying failure & error when trying to revert the applied ones"), err)
			assert.Equal(types.ErroredRevertingFailureState, in.CurrentState)
			at1.AssertExpectations(t)
			at2.AssertExpectations(t)
			at3.AssertExpectations(t)
		}
	}
}

func TestSystemFailureFailAttacksFinishWithTimeout(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1 := &mattack.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	reverted := make(chan struct{})
	at1.On("Revert").Once().Return(nil).Run(func(args mock.Arguments) {
		reverted <- struct{}{}
	})

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)

	// Mock clock
	cl := &mclock.Clock{}
	cl.On("After", f.Definition.Timeout).Return(time.After(0))
	cl.On("Now").Return(time.Now())

	in, err := failure.NewInjectionFromReg(f, reg, nil, cl)
	if assert.NoError(err) {
		err = in.Fail()
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

	f := &failure.Failure{
		Definition: failure.Definition{
			Timeout: 1 * time.Hour,
			Attacks: []failure.AttackMap{
				{"attack1": attack.Opts{}},
			},
		},
	}

	// Mock attackers
	ctxMatcher := mock.MatchedBy(func(ctx context.Context) bool { return true })
	at1 := &mattack.Attacker{}
	at1.On("Apply", ctxMatcher).Once().Return(nil)
	at1.On("Revert").Once().Return(nil)

	// Mock Registry
	reg := &mattack.Registry{}
	reg.On("New", "attack1", attack.Opts{}).Return(at1, nil)

	// Mock clock
	cl := &mclock.Clock{}
	cl.On("After", f.Definition.Timeout).Return(time.After(9999 * time.Hour)) // Never
	cl.On("Now").Return(time.Now())

	in, err := failure.NewInjectionFromReg(f, reg, nil, cl)
	if assert.NoError(err) {
		err = in.Fail()
		in.Revert()
		at1.AssertExpectations(t)
	}
}
