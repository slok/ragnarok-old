package validator

import (
	"fmt"
	"strings"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
)

const (
	labelMaxLength = 250
)

var (
	requiredFailureLabels = []string{api.LabelExperiment, api.LabelNode}
)

// ErrorList is a group of errors. Handy to validate once.
type ErrorList []error

func (e ErrorList) String() string {
	errStrs := make([]string, len(e))
	for i, err := range e {
		errStrs[i] = fmt.Sprintf("error: %s", err)
	}
	return strings.Join(errStrs, ",")
}

// ObjectValidator is the interface that should validate any kind of object based on the type is.
type ObjectValidator interface {
	// Calidate will return an error list if any errors arise when validating the object.
	Validate(obj api.Object) ErrorList
}

// Object inplements the validation of the objects
type Object struct{}

// DefaultObject is the default object validator.
var DefaultObject = NewObject()

// NewObject returns a new object validator.
func NewObject() *Object {
	return &Object{}
}

func (o *Object) validateObjectMeta(meta api.ObjectMeta) ErrorList {
	errors := []error{}

	// Check ID
	if err := errorIfEmptyStr(meta.ID); err != nil {
		errors = append(errors, fmt.Errorf("id error: %s", err))
	}

	// Check labels
	for k, v := range meta.Labels {
		// Check label length.
		if err := errorIfBiggerStr(k, labelMaxLength); err != nil {
			errors = append(errors, fmt.Errorf("label key error: %s", err))
		} else if err := errorIfEmptyStr(k); err != nil {
			errors = append(errors, fmt.Errorf("label value error: %s", err))
		}

		if err := errorIfBiggerStr(v, labelMaxLength); err != nil {
			errors = append(errors, fmt.Errorf("label key error: %s", err))
		} else if err := errorIfEmptyStr(v); err != nil {
			errors = append(errors, fmt.Errorf("label value error: %s", err))
		}

		// Check label valid string.
		if err := errorIfvalidLabelStr(k); err != nil {
			errors = append(errors, fmt.Errorf("label key error: %s", err))
		}
		if err := errorIfvalidLabelStr(v); err != nil {
			errors = append(errors, fmt.Errorf("label value error: %s", err))
		}
	}

	return errors
}

func (o *Object) validateNode(node *clusterv1.Node) ErrorList {
	errors := []error{}

	errors = append(errors, o.validateObjectMeta(node.Metadata)...)

	return errors
}

func (o *Object) validateFailure(flr *chaosv1.Failure) ErrorList {
	errors := []error{}

	errors = append(errors, o.validateObjectMeta(flr.Metadata)...)

	// Check failure labels correct.
	errors = append(errors, errorIfNoLabelKeys(flr.Metadata.Labels, requiredFailureLabels)...)

	return errors
}

func (o *Object) validateExperiment(exp *chaosv1.Experiment) ErrorList {
	errors := []error{}

	errors = append(errors, o.validateObjectMeta(exp.Metadata)...)

	return errors
}

func (o *Object) validateObject(obj api.Object) ErrorList {
	errs := []error{}
	if err := errorIfWrongTypeMeta(obj); err != nil {
		errs = append(errs, err)
	}

	return errs
}

// Validate satisfies ObjectValidator.
func (o *Object) Validate(obj api.Object) ErrorList {

	// Common object validation
	if errs := o.validateObject(obj); len(errs) > 0 {
		return errs
	}

	// Validation per type.
	switch v := obj.(type) {
	case *clusterv1.Node:
		return o.validateNode(v)
	case *chaosv1.Failure:
		return o.validateFailure(v)
	case *chaosv1.Experiment:
		return o.validateExperiment(v)
	default:
		return []error{fmt.Errorf("not a valid object to validate")}
	}
}
