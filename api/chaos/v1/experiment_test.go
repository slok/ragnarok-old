package v1_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/apimachinery"
	"github.com/slok/ragnarok/log"
)

func TestJSONEncodeChaosV1Experiment(t *testing.T) {
	//t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

	tests := []struct {
		name       string
		experiment *chaosv1.Experiment
		expEncNode string
		expErr     bool
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjTyper, apimachinery.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.experiment, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncNode, b.String())
				assert.NoError(err)
			}
		})
	}
}
