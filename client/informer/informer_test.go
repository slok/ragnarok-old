package informer_test

import (
	"testing"
	"time"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/informer"
	"github.com/slok/ragnarok/log"
	mwatch "github.com/slok/ragnarok/mocks/apimachinery/watch"
	minformer "github.com/slok/ragnarok/mocks/client/informer"
	mqueue "github.com/slok/ragnarok/mocks/client/util/queue"
	mstore "github.com/slok/ragnarok/mocks/client/util/store"
	testapi "github.com/slok/ragnarok/test/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWorkQueueInformerInitialState(t *testing.T) {
	tests := []struct {
		name    string
		objList api.ObjectList
	}{
		{
			name: "Initial state should populate the cache store and push the first jobs.",
			objList: &testapi.TestObjList{
				Items: []*testapi.TestObj{
					&testapi.TestObj{ID: "test0"},
					&testapi.TestObj{ID: "test1"},
					&testapi.TestObj{ID: "test2"},
					&testapi.TestObj{ID: "test3"},
					&testapi.TestObj{ID: "test4"},
					&testapi.TestObj{ID: "test5"},
					&testapi.TestObj{ID: "test6"},
				},
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			watchC := make(chan watch.Event)
			// channel to know when we are ready to check the result.
			readyC := make(chan struct{})

			// mocks.
			mindex := &mstore.ObjectIndexKeyer{}
			mstore := &mstore.Store{}
			mqueue := &mqueue.Queue{}
			mwatcher := &mwatch.Watcher{}
			mwatcher.On("GetChan").Return((<-chan watch.Event)(watchC))
			mlw := &minformer.ListerWatcher{}
			mlw.On("List", mock.Anything).Return(test.objList, nil)
			mlw.On("Watch", mock.Anything).Return(mwatcher, nil).Run(func(_ mock.Arguments) {
				// If we have call the watch, the list should be made.
				readyC <- struct{}{}
			})

			// Mock calls of resources on the mocks.
			for _, obj := range test.objList.GetItems() {
				id := obj.GetObjectMetadata().ID
				mindex.On("GetKey", mock.Anything).Once().Return(id, nil)
				mqueue.On("Push", id).Once().Return(nil) // This kind of informer sets the index Key.
				mstore.On("Add", obj).Once().Return(nil)
			}

			lOpts := api.ListOptions{}
			inf := informer.NewWorkQueueInformer(mindex, mqueue, mstore, lOpts, mlw, log.Dummy)
			stopC := make(chan struct{})

			// Run everything
			go inf.Run(stopC)

			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("Timeout waiting for results.")
			case <-readyC:
				// Ready to check.
				mstore.AssertExpectations(t)
				mqueue.AssertExpectations(t)
			}
		})
	}
}

func TestWorkQueueInformerWatchEvents(t *testing.T) {
	lastObjectID := "testLast"
	tests := []struct {
		name                 string
		events               []watch.Event
		expAdds              []api.Object
		updatedObjectsExists bool // do the updated objects exist on the store before?
		expUpdate            []api.Object
		expDelete            []api.Object
	}{
		{
			name:                 "Adding new objects should store and push on workqueue.",
			updatedObjectsExists: true,
			events: []watch.Event{
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test0"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: lastObjectID}},
			},
			expAdds: []api.Object{
				&testapi.TestObj{ID: "test0"},
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: lastObjectID},
			},
			expUpdate: []api.Object{},
			expDelete: []api.Object{},
		},
		{
			name:                 "Updating new objects should store and push on workqueue as adding them.",
			updatedObjectsExists: false,
			events: []watch.Event{
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test0"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: lastObjectID}},
			},
			expAdds: []api.Object{
				&testapi.TestObj{ID: "test0"},
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: lastObjectID},
			},
			expUpdate: []api.Object{},
			expDelete: []api.Object{},
		},
		{
			name:                 "Updating old objects should store and push on workqueue as adding them.",
			updatedObjectsExists: true,
			events: []watch.Event{
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test0"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: lastObjectID}},
			},
			expAdds: []api.Object{},
			expUpdate: []api.Object{
				&testapi.TestObj{ID: "test0"},
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: lastObjectID},
			},
			expDelete: []api.Object{},
		},
		{
			name:                 "Deleting objects should delete from store store and push on workqueue as adding them.",
			updatedObjectsExists: true,
			events: []watch.Event{
				watch.Event{Type: watch.DeletedEvent, Object: &testapi.TestObj{ID: "test0"}},
				watch.Event{Type: watch.DeletedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.DeletedEvent, Object: &testapi.TestObj{ID: lastObjectID}},
			},
			expAdds:   []api.Object{},
			expUpdate: []api.Object{},
			expDelete: []api.Object{
				&testapi.TestObj{ID: "test0"},
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: lastObjectID},
			},
		},
		{
			name:                 "Watching events should populate the cache correctly and insert the jobs.",
			updatedObjectsExists: true,
			events: []watch.Event{
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test0"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test2"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test3"}},
				watch.Event{Type: watch.UpdatedEvent, Object: &testapi.TestObj{ID: "test3"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: "test4"}},
				watch.Event{Type: watch.DeletedEvent, Object: &testapi.TestObj{ID: "test1"}},
				watch.Event{Type: watch.AddedEvent, Object: &testapi.TestObj{ID: lastObjectID}},
			},
			expAdds: []api.Object{
				&testapi.TestObj{ID: "test0"},
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: "test2"},
				&testapi.TestObj{ID: "test3"},
				&testapi.TestObj{ID: "test4"},
				&testapi.TestObj{ID: lastObjectID},
			},
			expUpdate: []api.Object{
				&testapi.TestObj{ID: "test1"},
				&testapi.TestObj{ID: "test3"},
			},
			expDelete: []api.Object{
				&testapi.TestObj{ID: "test1"},
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			watchC := make(chan watch.Event)
			// channel to know when we are ready to check the result.
			readyC := make(chan struct{})

			// mocks.
			mindex := &mstore.ObjectIndexKeyer{}
			mstore := &mstore.Store{}
			mqueue := &mqueue.Queue{}
			mwatcher := &mwatch.Watcher{}
			mwatcher.On("GetChan").Return((<-chan watch.Event)(watchC))
			mlw := &minformer.ListerWatcher{}
			mlw.On("List", mock.Anything).Return(&testapi.TestObjList{}, nil)
			mlw.On("Watch", mock.Anything).Return(mwatcher, nil)

			// Mock calls of resources on the mocks.

			expGetObjects := test.expAdds
			if test.updatedObjectsExists {
				expGetObjects = test.expUpdate
			}
			for _, obj := range expGetObjects {
				// Before updating it checks the old object.
				mstore.On("Get", obj).Return(nil, test.updatedObjectsExists, nil)
			}

			for _, obj := range test.expAdds {
				mstore.On("Add", obj).Return(nil)
			}
			for _, obj := range test.expUpdate {
				mstore.On("Update", obj).Return(nil)
			}
			for _, obj := range test.expDelete {
				mstore.On("Delete", obj).Return(nil)
			}
			for _, event := range test.events {
				id := event.Object.GetObjectMetadata().ID
				mindex.On("GetKey", mock.Anything).Once().Return(id, nil)
				mqueue.On("Push", id).Once().Return(nil).Run(func(args mock.Arguments) {
					objID := args.Get(0).(string)
					// We have finished, so make the call of finished to assert the expectations.
					if objID == lastObjectID {
						readyC <- struct{}{}
					}
				})
			}

			lOpts := api.ListOptions{}
			inf := informer.NewWorkQueueInformer(mindex, mqueue, mstore, lOpts, mlw, log.Dummy)
			stopC := make(chan struct{})
			go inf.Run(stopC)

			// Send the events.
			for _, event := range test.events {
				watchC <- event
			}

			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("Timeout waiting for results.")
			case <-readyC:
				// Ready to check.
				mstore.AssertExpectations(t)
				mqueue.AssertExpectations(t)
			}
		})
	}
}

func TestWorkQueueInformerStop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Stopping the informer should shutdown the queue.",
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			watchC := make(chan watch.Event)

			// mocks.
			mindex := &mstore.ObjectIndexKeyer{}
			mstore := &mstore.Store{}
			mqueue := &mqueue.Queue{}
			mqueue.On("ShutDown").Once().Return(nil)
			mwatcher := &mwatch.Watcher{}
			mwatcher.On("GetChan").Return((<-chan watch.Event)(watchC))
			mlw := &minformer.ListerWatcher{}
			mlw.On("List", mock.Anything).Return(&testapi.TestObjList{}, nil)
			mlw.On("Watch", mock.Anything).Return(mwatcher, nil)

			lOpts := api.ListOptions{}
			inf := informer.NewWorkQueueInformer(mindex, mqueue, mstore, lOpts, mlw, log.Dummy)
			stopC := make(chan struct{})
			result := make(chan error)
			go func() {
				result <- inf.Run(stopC)
			}()

			// Stop the informer.
			stopC <- struct{}{}

			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("The informer should be stopped.")
			case err := <-result:
				if assert.NoError(err) {
					mqueue.AssertExpectations(t)
				}
			}
		})
	}
}
