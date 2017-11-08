package serializer

import (
	"bytes"
	"fmt"
	"strings"

	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/log"
)

// PBSerializer knows how to serialize objects back and forth using PB style.
type PBSerializer struct {
	serializer Serializer
	logger     log.Logger
}

// NewPBSerializer returns a new PBSerializer object.
func NewPBSerializer(logger log.Logger) *PBSerializer {
	return &PBSerializer{
		serializer: JSONSerializerDefault,
		logger:     logger,
	}
}

func (p *PBSerializer) encodeClusterV1Node(obj api.Object, out *clusterv1pb.Node) error {
	var b bytes.Buffer
	if err := p.serializer.Encode(obj, &b); err != nil {
		return err
	}
	out.SerializedData = strings.TrimSuffix(b.String(), "\n")
	return nil
}

func (p *PBSerializer) decodeClusterV1Node(in *clusterv1pb.Node) (api.Object, error) {
	return p.serializer.Decode([]byte(in.SerializedData))
}

// Encode satisfies Serializer interface.
func (p *PBSerializer) Encode(obj api.Object, out interface{}) error {
	var err error
	switch pb := out.(type) {
	case *clusterv1pb.Node:
		err = p.encodeClusterV1Node(obj, pb)
	default:
		err = fmt.Errorf("unknown pb type")
	}
	return err
}

// Decode satisfies Serializer interface.
func (p *PBSerializer) Decode(data interface{}) (api.Object, error) {
	var (
		err error
		obj api.Object
	)

	switch pb := data.(type) {
	case *clusterv1pb.Node:
		obj, err = p.decodeClusterV1Node(pb)
	default:
		err = fmt.Errorf("unknown pb type")
	}

	return obj, err
}
