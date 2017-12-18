package util

import (
	"fmt"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	test "github.com/slok/ragnarok/test/api"
)

// NewObjectList is an util to abstract object list creator using interfaces.
// TODO: Redo using a registrator for the creators.
func NewObjectList(objs []api.Object, continueList string) (api.ObjectList, error) {
	if len(objs) == 0 {
		return nil, fmt.Errorf("object list can't be empty")
	}
	switch objs[0].(type) {
	case *clusterv1.Node:
		ns := make([]*clusterv1.Node, len(objs))
		for i, obj := range objs {
			ns[i] = obj.(*clusterv1.Node)
		}
		l := clusterv1.NewNodeList(ns, continueList)
		return &l, nil
	case *chaosv1.Failure:
		fs := make([]*chaosv1.Failure, len(objs))
		for i, obj := range objs {
			fs[i] = obj.(*chaosv1.Failure)
		}
		l := chaosv1.NewFailureList(fs, continueList)
		return &l, nil
	case *chaosv1.Experiment:
		es := make([]*chaosv1.Experiment, len(objs))
		for i, obj := range objs {
			es[i] = obj.(*chaosv1.Experiment)
		}
		l := chaosv1.NewExperimentList(es, continueList)
		return &l, nil
	// This is the test object used for some tests aroudn the app.
	// TODO: Rethink to not get this code here (registrators?)
	case *test.TestObj:
		ts := make([]*test.TestObj, len(objs))
		for i, obj := range objs {
			ts[i] = obj.(*test.TestObj)
		}
		l := test.NewTestObjList(ts, continueList)
		return &l, nil
	}
	return nil, fmt.Errorf("unknown object type")
}
