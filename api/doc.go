/*
Package api will have all the basic objects that will be used to work in the application.
The objects will be versioned by version and kind, also every object will have serialization
capabilities using the apimachinery/serializers package in YAML, JSON and Protobuffer formats.
*/
package api // import "github.com/slok/ragnarok/api"

// cluster/v1
//go:generate protoc -I. -I${GOOGLEPROTO_PATH} cluster/v1/pb/node.proto --go_out=plugins=grpc:.
