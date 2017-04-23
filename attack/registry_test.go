package attack_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/mocks"
)

func TestRegister(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		id      string
		wantErr bool
	}{
		{id: "id1", wantErr: false},
		{id: "968", wantErr: false},
		{id: "", wantErr: true},
	}

	for _, test := range tests {
		r := attack.NewSimpleRegistry()

		m := &mocks.Creater{}
		err := r.Register(test.id, m)
		if !test.wantErr {
			assert.NoError(err, "An error was't expected")
			assert.Contains(r, test.id, "%s should be registered but is missing", test.id)
		} else {
			assert.Error(err, "An error was expected")
		}
	}
}

func TestDeregister(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		regID   string
		deregID string
		wantErr bool
	}{
		{regID: "id1", deregID: "id1", wantErr: false},
		{regID: "968", deregID: "968", wantErr: false},
		{regID: "968", deregID: "not_correct", wantErr: true},
		{regID: "968", deregID: "", wantErr: true},
	}

	for _, test := range tests {
		// Setup registry.
		r := attack.SimpleRegistry(map[string]attack.Creater{
			test.regID: &mocks.Creater{},
		})

		// Check.
		err := r.Deregister(test.deregID)
		if !test.wantErr {
			assert.NoError(err, "An error was't expected")
			assert.NotContains(r, test.regID, "%s should be deregistered but is present", test.deregID)
		} else {
			assert.Error(err, "An error was expected")
		}
	}
}

func TestExists(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		ids     []string
		checkID string
		want    bool
	}{
		{ids: []string{"id1", "id2", "id3"}, checkID: "id1", want: true},
		{ids: []string{"id1", "id2", "id3"}, checkID: "id0", want: false},
		{ids: []string{}, checkID: "id1", want: false},
	}
	for _, test := range tests {
		// Setup registry.
		r := attack.SimpleRegistry(map[string]attack.Creater{})
		for _, id := range test.ids {
			r[id] = &mocks.Creater{}
		}
		assert.Equal(test.want, r.Exists(test.checkID))
	}
}

func TestFactory(t *testing.T) {
	r := attack.SimpleRegistry(map[string]attack.Creater{})
	// Prepare 10 mocks on the registry
	creaters := make(map[string]*mocks.Creater)
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("id%d", i)
		opts := attack.Opts{"id": id, "idx": i}

		m := &mocks.Creater{}
		m.On("Create", opts).Return(nil, nil)
		creaters[id] = m
		r[id] = m
	}

	// Use the factory and check it called the mocks
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("id%d", i)
		opts := attack.Opts{"id": id, "idx": i}
		r.New(id, opts)
		creaters[id].AssertExpectations(t)
	}

}
