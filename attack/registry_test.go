package attack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CreaterMock struct {
	mock.Mock
}

func (m *CreaterMock) Create(opts Opts) (Attacker, error) {
	m.Called(opts) // Track call.
	return nil, nil
}

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
		r := NewSimpleRegistry()

		err := r.Register(test.id, &CreaterMock{})
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
		r := SimpleRegistry(map[string]Creater{
			test.regID: &CreaterMock{},
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
		r := SimpleRegistry(map[string]Creater{})
		for _, id := range test.ids {
			r[id] = &CreaterMock{}
		}
		assert.Equal(test.want, r.Exists(test.checkID))
	}
}

func TestFactory(t *testing.T) {
	r := SimpleRegistry(map[string]Creater{})
	// Prepare 10 mocks on the registry
	mocks := make(map[string]*CreaterMock)
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("id%d", i)
		opts := Opts{"id": id, "idx": i}

		m := &CreaterMock{}
		m.On("Create", opts).Return(nil, nil)
		mocks[id] = m
		r[id] = m
	}

	// Use the factory and check it called the mocks
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("id%d", i)
		opts := Opts{"id": id, "idx": i}
		r.New(id, opts)
		mocks[id].AssertExpectations(t)
	}

}
