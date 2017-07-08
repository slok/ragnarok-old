package grpc

import (
	"google.golang.org/grpc"

	// TODO: Change when GRPC supports std librarie context
	"golang.org/x/net/context"

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
)

// TestClient is a GRPC client ready to test GRPC services
type TestClient struct {
	conn *grpc.ClientConn

	nsCli pbns.NodeStatusClient
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
	}, nil
}

// Close closes the GRPC connection.
func (t *TestClient) Close() error {
	return t.conn.Close()
}

// NodeStatusRegister wraps the call to nodestatus service
func (t *TestClient) NodeStatusRegister(ctx context.Context, ni *pbns.NodeInfo) (*pbns.RegisteredResponse, error) {
	return t.nsCli.Register(ctx, ni)
}
