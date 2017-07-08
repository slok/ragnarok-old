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
//go:generate mockery -output . -dir ../node/client -name Status

// Master mocks
//go:generate mockery -output ./master -outpkg master -dir ../master/ -name NodeRegistry
//go:generate mockery -output ./master -outpkg master -dir ../master/ -name Master

// GRPC clients
//go:generate mockery -output ./grpc -outpkg grpc -dir ../grpc/nodestatus -name NodeStatusClient
