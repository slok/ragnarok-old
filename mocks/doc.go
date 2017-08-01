/*
Package mocks will have all the mocks of the application, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/ragnarok/mocks"

// Attack mocks
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Registry
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Creater
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Attacker

// Failure mocks
//go:generate mockery -output ./failure -outpkg failure -dir ../failure -name Parser

// Clock mocks
//go:generate mockery -output ./clock -outpkg clock -dir ../clock -name Clock

// Logger mock
//go:generate mockery -output ./log -outpkg log -dir ../log -name Logger

// Node mocks
//go:generate mockery -output ./node/client -outpkg client -dir ../node/client -name Status

// Services mocks
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name NodeStatusService
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name NodeRepository
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name FailureStatusService
//go:generate mockery -output ./service -outpkg service -dir ../master/service -name FailureRepository

// Types mocks
//go:generate mockery -output ./types -outpkg types -dir ../types -name NodeStateParser

// GRPC proto clients
//go:generate mockery -output ./grpc/nodestatus -outpkg nodestatus -dir ../grpc/nodestatus -name NodeStatusClient
//go:generate mockery -output ./grpc/failurestatus -outpkg failurestatus -dir ../grpc/failurestatus -name FailureStatusClient
//go:generate mockery -output ./grpc/failurestatus -outpkg failurestatus -dir ../grpc/failurestatus -name FailureStatus_FailureStateListServer
