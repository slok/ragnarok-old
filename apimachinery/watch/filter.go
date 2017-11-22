package watch

import (
	"github.com/slok/ragnarok/api"
	apiutil "github.com/slok/ragnarok/api/util"
)

// NoFilter is handy filter to not filter anything.
var (
	NoFilter = FilterFunc(func(_ api.Object) bool { return false })
)

// ObjectFilter is a filter that should return true if the object should be filtered. This
// should be used to filter events based on the objects received.
type ObjectFilter interface {
	// Filter will return true if the object should be filtered (discarded).
	Filter(api.Object) bool
}

// FilterFunc is a helper and implementation for ObjectFilter interface.
type FilterFunc func(api.Object) bool

// Filter satisfies ObjectFilter interface.
func (f FilterFunc) Filter(obj api.Object) bool {
	return f(obj)
}

// TypeFilter will filter (filter returns true) if the object type is not the one
// selected when creating the type filter..
type TypeFilter struct {
	objFullType string
}

// NewTypeFilter returns a new Type filter.
func NewTypeFilter(obj api.Object) *TypeFilter {
	return &TypeFilter{
		objFullType: apiutil.GetFullType(obj),
	}
}

// Filter will return true if the object of the type filter is different to the
// object argument received. Satisfies ObjectFilter interface.
func (t *TypeFilter) Filter(obj api.Object) bool {
	objFT := apiutil.GetFullType(obj)
	return t.objFullType != objFT
}

// LabelFilter will filter (filter returns true) if the object doesn't have the
// required labels.
type LabelFilter struct {
	selector map[string]string
}

// NewLabelFilter returns a new Label filter.
func NewLabelFilter(selector map[string]string) *LabelFilter {
	return &LabelFilter{
		selector: selector,
	}
}

// selectorMatchesLabels will return true if the selector matches the labels.
func (l *LabelFilter) selectorMatchesLabels(labels map[string]string, selector map[string]string) bool {
	for lk, lv := range selector {
		// If one label does not match then we are done.
		if nv, ok := labels[lk]; !ok || (ok && nv != lv) {
			return false
		}
	}
	// All labels matched.
	return true
}

// Filter will return true if the received object does not satisfie the selector
// labels. If the selector doesn't have labels then all the received objects will not be ignored.
// Satisfies ObjectFilter interface.
func (l *LabelFilter) Filter(obj api.Object) bool {
	if len(l.selector) == 0 {
		return false
	}
	// If does not match then filter the object.
	return !l.selectorMatchesLabels(obj.GetObjectMetadata().Labels, l.selector)
}

// ListOptionsFilter will filter (filter returns true) based on list options.
type ListOptionsFilter struct {
	typeF  *TypeFilter
	labelF *LabelFilter
}

// NewListOptionsFilter returns a new ListOptions filter.
func NewListOptionsFilter(opts api.ListOptions) *ListOptionsFilter {
	return &ListOptionsFilter{
		typeF:  NewTypeFilter(opts),
		labelF: NewLabelFilter(opts.LabelSelector),
	}
}

// Filter will return true if received object does not satisfie the listoptions.
func (l *ListOptionsFilter) Filter(obj api.Object) bool {
	return l.typeF.Filter(obj) || l.labelF.Filter(obj)
}
