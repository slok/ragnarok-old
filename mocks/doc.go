/*
Package mocks will have all the mocks of the application, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/ragnarok/mocks"

// Attack mocks
//go:generate mockery -output . -dir ../attack -name Registry
//go:generate mockery -output . -dir ../attack -name Creater
//go:generate mockery -output . -dir ../attack -name Attacker

// Clock mocks
//go:generate mockery -output . -dir ../clock -name Clock

// Logger mock
//go:generate mockery -output . -dir ../log -name Logger

// Node mocks
//go:generate mockery -output ./node/client -dir ../node/client -name Status

// Services mocks
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name NodeStatusService
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name NodeRepository

// Types mocks
//go:generate mockery -output ./types -outpkg types -dir ../types -name NodeStateParser

// GRPC proto clients
//go:generate mockery -output ./grpc -outpkg grpc -dir ../grpc/nodestatus -name NodeStatusClient
