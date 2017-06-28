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
