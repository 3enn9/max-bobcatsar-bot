package accountant

import (
	pb "bobcatsar-max-bot/internal/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	pb.AccountantServiceClient
	conn *grpc.ClientConn
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	return &Client{AccountantServiceClient: pb.NewAccountantServiceClient(conn), conn: conn}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
