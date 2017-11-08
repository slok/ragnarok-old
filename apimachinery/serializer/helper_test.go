package serializer_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/apimachinery/serializer"
	pb "github.com/slok/ragnarok/apimachinery/serializer/testdata"
)

// NOTE: Serializers tests are on the api package (tests serialization per type).

func TestSafeTypeAsserterByteArray(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		expData []byte
		expErr  bool
	}{
		{
			name:    "Byte array assertion shouldn't return an error",
			data:    []byte("Peter Parker is Spiderman"),
			expData: []byte("Peter Parker is Spiderman"),
			expErr:  false,
		},
		{
			name:   "Int assertion should return an error",
			data:   5,
			expErr: true,
		},
		{
			name:   "String assertion should return an error",
			data:   "Bruce Wayne is Batman",
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := serializer.SafeTypeAsserter.ByteArray(test.data)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expData, got)
			}
		})
	}
}

func TestSafeTypeAsserterWriter(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		expData io.Writer
		expErr  bool
	}{
		{
			name:    "Buffer assertion shouldn't return an error",
			data:    &bytes.Buffer{},
			expData: &bytes.Buffer{},
			expErr:  false,
		},
		{
			name:    "Stdout assertion shouldn't return an error",
			data:    os.Stdout,
			expData: os.Stdout,
			expErr:  false,
		},
		{
			name:   "Int assertion should return an error",
			data:   5,
			expErr: true,
		},
		{
			name:   "String assertion should return an error",
			data:   "Bruce Wayne is Batman",
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := serializer.SafeTypeAsserter.Writer(test.data)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expData, got)
			}
		})
	}
}

func TestSafeTypeAsserterProtoMessage(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		expData proto.Message
		expErr  bool
	}{
		{
			name:    "PB Test assertion shouldn't return an error",
			data:    &pb.Test{},
			expData: &pb.Test{},
			expErr:  false,
		},
		{
			name:   "Int assertion should return an error",
			data:   5,
			expErr: true,
		},
		{
			name:   "String assertion should return an error",
			data:   "Bruce Wayne is Batman",
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			got, err := serializer.SafeTypeAsserter.ProtoMessage(test.data)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
				assert.Equal(test.expData, got)
			}
		})
	}
}
