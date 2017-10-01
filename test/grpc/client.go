package grpc

import (
	"google.golang.org/grpc"

	// TODO: Change when GRPC supports std library context
	pbempty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"

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
func (t *TestClient) NodeStatusRegister(ctx context.Context, ni *pbns.Node) (*pbns.RegisteredResponse, error) {
	return t.nsCli.Register(ctx, ni)
}

// NodeStatusHeartbeat wraps the call to nodestatus service.
func (t *TestClient) NodeStatusHeartbeat(ctx context.Context, ns *pbns.NodeState) (*pbempty.Empty, error) {
	return t.nsCli.Heartbeat(ctx, ns)
}

// FailureStatusGetFailure wraps the call to failurestatus service.
func (t *TestClient) FailureStatusGetFailure(ctx context.Context, fID *pbfs.FailureId) (*pbfs.Failure, error) {
	return t.fsCli.GetFailure(ctx, fID)
}

// FailureStatusFailureStateList wraps the call to failurestatus service.
func (t *TestClient) FailureStatusFailureStateList(ctx context.Context, nID *pbfs.NodeId) (pbfs.FailureStatus_FailureStateListClient, error) {
	return t.fsCli.FailureStateList(ctx, nID)
}
