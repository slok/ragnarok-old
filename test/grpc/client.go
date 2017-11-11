package grpc

import (
	"google.golang.org/grpc"

	// TODO: Change when GRPC supports std library context
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
)

// TestClient is a GRPC client ready to test GRPC services
type TestClient struct {
	conn *grpc.ClientConn

	nsCli pbns.NodeStatusClient
	fsCli pbfs.FailureStatusClient
}

// NewTestClient creates and returns a new test client
func NewTestClient(addr string) (*TestClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &TestClient{
		conn:  conn,
		nsCli: pbns.NewNodeStatusClient(conn),
		fsCli: pbfs.NewFailureStatusClient(conn),
	}, nil
}

// Close closes the GRPC connection.
func (t *TestClient) Close() error {
	return t.conn.Close()
}

// NodeStatusRegister wraps the call to nodestatus service.
func (t *TestClient) NodeStatusRegister(ctx context.Context, n *clusterv1pb.Node) (*pbempty.Empty, error) {
	return t.nsCli.Register(ctx, n)
}

// NodeStatusHeartbeat wraps the call to nodestatus service.
func (t *TestClient) NodeStatusHeartbeat(ctx context.Context, n *clusterv1pb.Node) (*pbempty.Empty, error) {
	return t.nsCli.Heartbeat(ctx, n)
}

// FailureStatusGetFailure wraps the call to failurestatus service.
func (t *TestClient) FailureStatusGetFailure(ctx context.Context, fID *pbfs.FailureId) (*chaosv1pb.Failure, error) {
	return t.fsCli.GetFailure(ctx, fID)
}

// FailureStatusFailureStateList wraps the call to failurestatus service.
func (t *TestClient) FailureStatusFailureStateList(ctx context.Context, nID *pbfs.NodeId) (pbfs.FailureStatus_FailureStateListClient, error) {
	return t.fsCli.FailureStateList(ctx, nID)
}
