/*
Package mocks will have all the mocks of the application, we'll try to use mocking using blackbox
testing and integration tests whenever is possible.
*/
package mocks // import "github.com/slok/ragnarok/mocks"

// Attack mocks
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Registry
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Creater
//go:generate mockery -output ./attack -outpkg attack -dir ../attack -name Attacker

// Clock mocks
//go:generate mockery -output ./clock -outpkg clock -dir ../clock -name Clock

// Logger mock
//go:generate mockery -output ./log -outpkg log -dir ../log -name Logger

// Node mocks
//go:generate mockery -output ./node/client -outpkg client -dir ../node/client -name Status
//go:generate mockery -output ./node/client -outpkg client -dir ../node/client -name FailureStateHandler
//go:generate mockery -output ./node/client -outpkg client -dir ../node/client -name Failure
//go:generate mockery -output ./node/service -outpkg service -dir ../node/service -name FailureState
//go:generate mockery -output ./node/service -outpkg service -dir ../node/service -name Status

// Services mocks
//go:generate mockery -output ./master/service -outpkg service -dir ../master/service -name NodeStatusService
//go:generate mockery -output ./master/service -outpkg service -dir ../master/service -name FailureStatusService
//go:generate mockery -output ./master/service/scheduler -outpkg service -dir ../master/service/scheduler -name Scheduler

// GRPC proto clients
//go:generate mockery -output ./grpc/nodestatus -outpkg nodestatus -dir ../grpc/nodestatus -name NodeStatusClient
//go:generate mockery -output ./grpc/failurestatus -outpkg failurestatus -dir ../grpc/failurestatus -name FailureStatusClient
//go:generate mockery -output ./grpc/failurestatus -outpkg failurestatus -dir ../grpc/failurestatus -name FailureStatus_FailureStateListServer
//go:generate mockery -output ./grpc/failurestatus -outpkg failurestatus -dir ../grpc/failurestatus -name FailureStatus_FailureStateListClient

// apimachinery mocks
//go:generate mockery -output ./apimachinery/serializer -outpkg serializer -dir ../apimachinery/serializer -name Serializer
//go:generate mockery -output ./apimachinery/validator -outpkg validator -dir ../apimachinery/validator -name ObjectValidator
//go:generate mockery -output ./apimachinery/watch -outpkg watch -dir ../apimachinery/watch -name Multiplexer
//go:generate mockery -output ./apimachinery/watch -outpkg watch -dir ../apimachinery/watch -name MultiplexerFactory
//go:generate mockery -output ./apimachinery/watch -outpkg watch -dir ../apimachinery/watch -name Watcher

// client mocks
//go:generate mockery -output ./client/repository -outpkg v1 -dir ../client/repository -name Client
//go:generate mockery -output ./client/api/cluster/v1 -outpkg v1 -dir ../client/api/cluster/v1 -name NodeClientInterface
//go:generate mockery -output ./client/api/chaos/v1 -outpkg v1 -dir ../client/api/chaos/v1 -name FailureClientInterface
//go:generate mockery -output ./client/api/chaos/v1 -outpkg v1 -dir ../client/api/chaos/v1 -name ExperimentClientInterface
//go:generate mockery -output ./client/util/store -outpkg store -dir ../client/util/store -name ObjectIndexKeyer
//go:generate mockery -output ./client/util/store -outpkg store -dir ../client/util/store -name Store
//go:generate mockery -output ./client/util/queue -outpkg queue -dir ../client/util/queue -name Queue
//go:generate mockery -output ./client/informer -outpkg informer -dir ../client/informer -name ListerWatcher

// web mocks
//go:generate mockery -output ./master/web/handler -outpkg serializer -dir ../master/web/handler -name Handler
