// +build integration

package memory_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/client/repository"
	"github.com/slok/ragnarok/client/repository/memory"
	"github.com/slok/ragnarok/log"
)

// Main objects for the integration tests.
var (
	// Nodes
	node0MasterEU = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "0",
			Labels: map[string]string{"kind": "master", "region": "eu"},
		},
	}
	node1MasterEU = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "1",
			Labels: map[string]string{"kind": "master", "region": "eu"},
		},
	}
	node2MasterAP = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "2",
			Labels: map[string]string{"kind": "master", "region": "ap"},
		},
	}

	node3NodeEU = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "3",
			Labels: map[string]string{"kind": "node", "region": "eu"},
		},
	}
	node4MasterUS = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "4",
			Labels: map[string]string{"kind": "master", "region": "us"},
		},
	}

	node4MasterUSUpdated = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:          "4",
			Labels:      map[string]string{"kind": "master", "region": "us"},
			Annotations: map[string]string{"updated": "true"},
		},
	}
	node5NodeUS = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "5",
			Labels: map[string]string{"kind": "node", "region": "us"},
		},
	}
	node5NodeUSUpdated = clusterv1.Node{
		TypeMeta: clusterv1.NodeTypeMeta,
		Metadata: api.ObjectMeta{
			ID:          "5",
			Labels:      map[string]string{"kind": "node", "region": "us"},
			Annotations: map[string]string{"updated": "true"},
		},
	}

	// Experiments.
	experiment0 = chaosv1.Experiment{
		TypeMeta: chaosv1.ExperimentTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "0",
			Labels: map[string]string{"machine": "node", "memory": "true", "cpu": "true", "level": "critical"},
		},
		Spec: chaosv1.ExperimentSpec{
			Name:        "BurnCPUNode",
			Description: "This experiment will burn the CPU in spikes and fill the memory",
			Selector:    map[string]string{"kind": "master", "region": "us"},
			Template: chaosv1.ExperimentFailureTemplate{
				Spec: chaosv1.FailureSpec{
					Timeout: 5 * time.Minute,
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
						{
							"attack2": nil,
						},
						{
							"attack3": attack.Opts{
								"target":   "myTarget",
								"quantity": 10,
								"pace":     "10m",
								"rest":     "30s",
							},
						},
					},
				},
			},
		},
	}
	experiment1 = chaosv1.Experiment{
		TypeMeta: chaosv1.ExperimentTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "1",
			Labels: map[string]string{"machine": "node", "host": "true", "level": "warning"},
		},
		Spec: chaosv1.ExperimentSpec{
			Name:        "DisableHost",
			Description: "This experiment will disable a host",
			Selector:    map[string]string{"kind": "node"},
			Template: chaosv1.ExperimentFailureTemplate{
				Spec: chaosv1.FailureSpec{
					Timeout: 5 * time.Minute,
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"host": "example.com",
							},
						},
					},
				},
			},
		},
	}
	experiment2 = chaosv1.Experiment{
		TypeMeta: chaosv1.ExperimentTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "2",
			Labels: map[string]string{"machine": "node", "host": "true", "level": "critical"},
		},
		Spec: chaosv1.ExperimentSpec{
			Name:        "DisableHost",
			Description: "This experiment will disable a AWS host",
			Selector:    map[string]string{"kind": "node"},
			Template: chaosv1.ExperimentFailureTemplate{
				Spec: chaosv1.FailureSpec{
					Timeout: 10 * time.Hour,
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"host": "aws.com",
							},
						},
					},
				},
			},
		},
	}
	experiment2Updated = chaosv1.Experiment{
		TypeMeta: chaosv1.ExperimentTypeMeta,
		Metadata: api.ObjectMeta{
			ID:          "2",
			Labels:      map[string]string{"machine": "node", "host": "true", "level": "critical"},
			Annotations: map[string]string{"updated": "true"},
		},
		Spec: chaosv1.ExperimentSpec{
			Name:        "DisableHost",
			Description: "This experiment will disable a AWS & GC host",
			Selector:    map[string]string{"kind": "node"},
			Template: chaosv1.ExperimentFailureTemplate{
				Spec: chaosv1.FailureSpec{
					Timeout: 10 * time.Hour,
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"host": "aws.com",
							},
						},
						{
							"attack2": attack.Opts{
								"host": "cloud.google.com",
							},
						},
					},
				},
			},
		},
	}

	// Failures.
	failure0 = chaosv1.Failure{
		TypeMeta: chaosv1.FailureTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "flr-000",
			Labels: map[string]string{api.LabelExperiment: "1"},
		},
		Spec: chaosv1.FailureSpec{
			Timeout: 5 * time.Minute,
			Attacks: []chaosv1.AttackMap{
				{
					"attack1": attack.Opts{},
				},
			},
		},
	}
	failure0Updated = chaosv1.Failure{
		TypeMeta: chaosv1.FailureTypeMeta,
		Metadata: api.ObjectMeta{
			ID:          "flr-000",
			Labels:      map[string]string{api.LabelExperiment: "1"},
			Annotations: map[string]string{"updated": "true"},
		},
		Spec: chaosv1.FailureSpec{
			Timeout: 5 * time.Minute,
		},
	}
	failure1 = chaosv1.Failure{
		TypeMeta: chaosv1.FailureTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "flr-001",
			Labels: map[string]string{api.LabelExperiment: "2"},
		},
		Spec: chaosv1.FailureSpec{
			Timeout: 5 * time.Minute,
			Attacks: []chaosv1.AttackMap{
				{
					"attack1": attack.Opts{},
				},
			},
		},
	}
	failure2 = chaosv1.Failure{
		TypeMeta: chaosv1.FailureTypeMeta,
		Metadata: api.ObjectMeta{
			ID:     "flr-002",
			Labels: map[string]string{api.LabelExperiment: "2"},
		},
		Spec: chaosv1.FailureSpec{
			Timeout: 5 * time.Minute,
			Attacks: []chaosv1.AttackMap{
				{
					"attack1": attack.Opts{},
				},
			},
		},
	}
	failure2Updated = chaosv1.Failure{
		TypeMeta: chaosv1.FailureTypeMeta,
		Metadata: api.ObjectMeta{
			ID:          "flr-002",
			Labels:      map[string]string{api.LabelExperiment: "2"},
			Annotations: map[string]string{"updated": "true"},
		},
		Spec: chaosv1.FailureSpec{
			Timeout: 5 * time.Minute,
		},
	}
)

// Main add envents.
var (
	// Nodes.
	node0MasterEUEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node0MasterEU,
	}
	node1MasterEUEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node1MasterEU,
	}
	node2MasterAPEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node2MasterAP,
	}
	node3NodeEUEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node3NodeEU,
	}
	node4MasterUSEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node4MasterUS,
	}
	node4MasterUSEventUpdate = watch.Event{
		Type:   watch.UpdatedEvent,
		Object: &node4MasterUSUpdated,
	}
	node4MasterUSEventDelete = watch.Event{
		Type:   watch.DeletedEvent,
		Object: &node4MasterUS,
	}
	node5NodeUSEventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &node5NodeUS,
	}
	node5NodeUSEventUpdate = watch.Event{
		Type:   watch.UpdatedEvent,
		Object: &node5NodeUSUpdated,
	}

	// Experiments
	experiment0EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &experiment0,
	}
	experiment1EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &experiment1,
	}
	experiment2EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &experiment2,
	}
	experiment2EventDelete = watch.Event{
		Type:   watch.DeletedEvent,
		Object: &experiment2,
	}
	experiment2EventUpdate = watch.Event{
		Type:   watch.UpdatedEvent,
		Object: &experiment2Updated,
	}

	// Failures
	failure0EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &failure0,
	}
	failure0EventUpdate = watch.Event{
		Type:   watch.UpdatedEvent,
		Object: &failure0Updated,
	}
	failure1EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &failure1,
	}
	failure2EventAdd = watch.Event{
		Type:   watch.AddedEvent,
		Object: &failure2,
	}
	failure2EventDelete = watch.Event{
		Type:   watch.DeletedEvent,
		Object: &failure2,
	}
	failure2EventUpdate = watch.Event{
		Type:   watch.UpdatedEvent,
		Object: &failure2Updated,
	}
)

type testWatcher struct {
	opts      api.ListOptions
	result    *[]watch.Event
	expResult []watch.Event
}

func newTestWatcher(opts api.ListOptions, expResult []watch.Event, client repository.Client, t *testing.T) *testWatcher {
	result := []watch.Event{}
	watcher, err := client.Watch(opts)
	require.NoError(t, err)
	go func() {
		for ev := range watcher.GetChan() {
			result = append(result, ev)
		}
	}()

	return &testWatcher{
		opts:      opts,
		result:    &result,
		expResult: expResult,
	}
}

// TestMemoryAddEventPropagation will test that doing different add actions on the
// memory repository the add events are propagated to different watchers.
func TestMemoryAddEventPropagation(t *testing.T) {
	assert := assert.New(t)

	// Create the memory client.
	logger := log.Dummy
	watcherfactory := watch.NewDefaultBroadcasterFactory(logger)
	client := memory.NewDefaultClient(watcherfactory, logger)

	// Watcher for all node objects.
	nodeAllOpts := api.ListOptions{
		TypeMeta: clusterv1.NodeTypeMeta,
	}
	nodeAllExpResult := []watch.Event{node0MasterEUEventAdd, node1MasterEUEventAdd, node2MasterAPEventAdd, node3NodeEUEventAdd, node4MasterUSEventAdd, node5NodeUSEventAdd}
	nodeAll := newTestWatcher(nodeAllOpts, nodeAllExpResult, client, t)

	// Watcher for eu node objects.
	nodeEUOpts := api.ListOptions{
		TypeMeta:      clusterv1.NodeTypeMeta,
		LabelSelector: map[string]string{"region": "eu"},
	}
	nodeEUExpResult := []watch.Event{node0MasterEUEventAdd, node1MasterEUEventAdd, node3NodeEUEventAdd}
	nodeEU := newTestWatcher(nodeEUOpts, nodeEUExpResult, client, t)

	// Watcher for master node objects.
	nodeMasterUSOpts := api.ListOptions{
		TypeMeta:      clusterv1.NodeTypeMeta,
		LabelSelector: map[string]string{"kind": "master", "region": "us"},
	}
	nodeMasterUSExpResult := []watch.Event{node4MasterUSEventAdd}
	nodeMasterUS := newTestWatcher(nodeMasterUSOpts, nodeMasterUSExpResult, client, t)

	// Watcher for critical experiments.
	experimentCriticalOpts := api.ListOptions{
		TypeMeta:      chaosv1.ExperimentTypeMeta,
		LabelSelector: map[string]string{"level": "critical"},
	}
	experimentCriticalExpResult := []watch.Event{experiment0EventAdd, experiment2EventAdd}
	experimentCritical := newTestWatcher(experimentCriticalOpts, experimentCriticalExpResult, client, t)

	// Watcher for all failures.
	failureAllOpts := api.ListOptions{
		TypeMeta:      chaosv1.FailureTypeMeta,
		LabelSelector: map[string]string{},
	}
	failureAllExpResult := []watch.Event{failure0EventAdd, failure1EventAdd, failure2EventAdd}
	failureAll := newTestWatcher(failureAllOpts, failureAllExpResult, client, t)

	// Watcher for all experiment 2 failures.
	failureExp2Opts := api.ListOptions{
		TypeMeta:      chaosv1.FailureTypeMeta,
		LabelSelector: map[string]string{api.LabelExperiment: "2"},
	}
	failureExp2ExpResult := []watch.Event{failure1EventAdd, failure2EventAdd}
	failureExp2 := newTestWatcher(failureExp2Opts, failureExp2ExpResult, client, t)

	// Add all objects.
	objectsToAdd := []api.Object{
		&node0MasterEU, &node1MasterEU, &node2MasterAP, &node3NodeEU, &node4MasterUS, &node5NodeUS,
		&experiment0, &experiment1, &experiment2,
		&failure0, &failure1, &failure2,
	}

	// Create all the objects
	for _, obj := range objectsToAdd {
		_, err := client.Create(obj)
		assert.NoError(err)
	}

	// Update some objects.

	// Safety sleep to leave time to process all the events.
	time.Sleep(1 * time.Second)

	// Add the tests so we can check checker.
	tests := []*testWatcher{
		nodeAll, nodeEU, nodeMasterUS,
		experimentCritical,
		failureAll, failureExp2,
	}
	for _, test := range tests {
		assert.Equal(test.expResult, *test.result)
	}
}

// TestMemoryUpdatedEventPropagation will test that doing different add & update actions on the
// memory repository the update events are propagated to different watchers.
func TestMemoryUpdatedEventPropagation(t *testing.T) {
	assert := assert.New(t)

	// Create the memory client.
	logger := log.Dummy
	watcherfactory := watch.NewDefaultBroadcasterFactory(logger)
	client := memory.NewDefaultClient(watcherfactory, logger)

	// Watcher for all node objects.
	nodeAllOpts := api.ListOptions{
		TypeMeta: clusterv1.NodeTypeMeta,
	}
	nodeAllExpResult := []watch.Event{node0MasterEUEventAdd, node1MasterEUEventAdd, node2MasterAPEventAdd, node3NodeEUEventAdd,
		node4MasterUSEventAdd, node5NodeUSEventAdd, node4MasterUSEventUpdate, node5NodeUSEventUpdate}
	nodeAll := newTestWatcher(nodeAllOpts, nodeAllExpResult, client, t)

	// Watcher for master node objects.
	nodeMasterUSOpts := api.ListOptions{
		TypeMeta:      clusterv1.NodeTypeMeta,
		LabelSelector: map[string]string{"kind": "master", "region": "us"},
	}
	nodeMasterUSExpResult := []watch.Event{node4MasterUSEventAdd, node4MasterUSEventUpdate}
	nodeMasterUS := newTestWatcher(nodeMasterUSOpts, nodeMasterUSExpResult, client, t)

	// Watcher for critical experiments.
	experimentCriticalOpts := api.ListOptions{
		TypeMeta:      chaosv1.ExperimentTypeMeta,
		LabelSelector: map[string]string{"level": "critical"},
	}
	experimentCriticalExpResult := []watch.Event{experiment0EventAdd, experiment2EventAdd, experiment2EventUpdate}
	experimentCritical := newTestWatcher(experimentCriticalOpts, experimentCriticalExpResult, client, t)

	// Watcher for all experiment 2 failures.
	failureExp2Opts := api.ListOptions{
		TypeMeta:      chaosv1.FailureTypeMeta,
		LabelSelector: map[string]string{api.LabelExperiment: "2"},
	}
	failureExp2ExpResult := []watch.Event{failure1EventAdd, failure2EventAdd, failure2EventUpdate}
	failureExp2 := newTestWatcher(failureExp2Opts, failureExp2ExpResult, client, t)

	// Add all objects.
	objectsToAdd := []api.Object{
		&node0MasterEU, &node1MasterEU, &node2MasterAP, &node3NodeEU, &node4MasterUS, &node5NodeUS,
		&experiment0, &experiment1, &experiment2,
		&failure0, &failure1, &failure2,
	}

	// Create all the objects
	for _, obj := range objectsToAdd {
		_, err := client.Create(obj)
		assert.NoError(err)
	}

	objectsToUpdate := []api.Object{
		&node4MasterUSUpdated, &node5NodeUSUpdated,
		&experiment2Updated,
		&failure0Updated, &failure2Updated,
	}

	// Create all the objects
	for _, obj := range objectsToUpdate {
		_, err := client.Update(obj)
		assert.NoError(err)
	}

	// Safety sleep to leave time to process all the events.
	time.Sleep(1 * time.Second)

	// Add the tests so we can check checker.
	tests := []*testWatcher{nodeAll, nodeMasterUS, experimentCritical, failureExp2}
	for _, test := range tests {
		assert.Equal(test.expResult, *test.result)
	}
}

// TestMemoryDeletedEventPropagation will test that doing different add & delete actions on the
// memory repository the delete events are propagated to different watchers.
func TestMemoryDeletedEventPropagation(t *testing.T) {
	assert := assert.New(t)

	// Create the memory client.
	logger := log.Dummy
	watcherfactory := watch.NewDefaultBroadcasterFactory(logger)
	client := memory.NewDefaultClient(watcherfactory, logger)

	// Watcher for all node objects.
	nodeAllOpts := api.ListOptions{
		TypeMeta: clusterv1.NodeTypeMeta,
	}
	nodeAllExpResult := []watch.Event{node0MasterEUEventAdd, node1MasterEUEventAdd, node2MasterAPEventAdd, node3NodeEUEventAdd,
		node4MasterUSEventAdd, node5NodeUSEventAdd, node4MasterUSEventDelete}
	nodeAll := newTestWatcher(nodeAllOpts, nodeAllExpResult, client, t)

	// Watcher for master node objects.
	nodeMasterUSOpts := api.ListOptions{
		TypeMeta:      clusterv1.NodeTypeMeta,
		LabelSelector: map[string]string{"kind": "master", "region": "us"},
	}
	nodeMasterUSExpResult := []watch.Event{node4MasterUSEventAdd, node4MasterUSEventDelete}
	nodeMasterUS := newTestWatcher(nodeMasterUSOpts, nodeMasterUSExpResult, client, t)

	// Watcher for critical experiments.
	experimentCriticalOpts := api.ListOptions{
		TypeMeta:      chaosv1.ExperimentTypeMeta,
		LabelSelector: map[string]string{"level": "critical"},
	}
	experimentCriticalExpResult := []watch.Event{experiment0EventAdd, experiment2EventAdd, experiment2EventDelete}
	experimentCritical := newTestWatcher(experimentCriticalOpts, experimentCriticalExpResult, client, t)

	// Watcher for all experiment 2 failures.
	failureExp2Opts := api.ListOptions{
		TypeMeta:      chaosv1.FailureTypeMeta,
		LabelSelector: map[string]string{api.LabelExperiment: "2"},
	}
	failureExp2ExpResult := []watch.Event{failure1EventAdd, failure2EventAdd, failure2EventDelete}
	failureExp2 := newTestWatcher(failureExp2Opts, failureExp2ExpResult, client, t)

	// Add all objects.
	objectsToAdd := []api.Object{
		&node0MasterEU, &node1MasterEU, &node2MasterAP, &node3NodeEU, &node4MasterUS, &node5NodeUS,
		&experiment0, &experiment1, &experiment2,
		&failure0, &failure1, &failure2,
	}

	// Create all the objects
	for _, obj := range objectsToAdd {
		_, err := client.Create(obj)
		assert.NoError(err)
	}

	objectsToDelete := []api.Object{
		&node4MasterUS,
		&experiment2,
		&failure0, &failure2,
	}

	// Create all the objects
	for _, obj := range objectsToDelete {
		fID := apiutil.GetFullID(obj)
		err := client.Delete(fID)
		assert.NoError(err)
	}

	// Safety sleep to leave time to process all the events.
	time.Sleep(1 * time.Second)

	// Add the tests so we can check checker.
	tests := []*testWatcher{nodeAll, nodeMasterUS, experimentCritical, failureExp2}
	for _, test := range tests {
		assert.Equal(test.expResult, *test.result)
	}
}
