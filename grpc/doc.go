/*
Package grpc will have all the proto and the generated go files of the required services
by the nodes and the masters to implement the servers and the clients
*/
package grpc // import "github.com/slok/ragnarok/grpc"

// Attack mocks
//go:generate protoc nodestatus/nodestatus.proto --go_out=plugins=grpc:.
