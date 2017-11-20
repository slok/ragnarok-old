package memory

import (
	"fmt"
	"sync"

	"github.com/slok/ragnarok/api"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/util"
	"github.com/slok/ragnarok/log"
)

// Client is an instance that saves objects in memory.
type Client struct {
	reg    map[string]map[string]api.Object
	logger log.Logger
	sync.Mutex
}

// NewDefaultClient returns a default object memory repository.
func NewDefaultClient(logger log.Logger) *Client {
	return NewClient(map[string]map[string]api.Object{}, logger)
}

// NewClient returns an object memory repository.
func NewClient(registry map[string]map[string]api.Object, logger log.Logger) *Client {
	logger = logger.WithField("repository", "memory")
	return &Client{
		logger: logger,
		reg:    registry,
	}
}

// safeGet will get an object using a fullID safely from the registry.
func (c *Client) safeGet(fullID string) api.Object {
	t, id := apiutil.SplitFullID(fullID)
	// Get object registry.
	if objReg, ok := c.reg[t.String()]; ok {
		// Get object if registry present.
		if obj, ok := objReg[id]; ok {
			return obj
		}
	}
	c.logger.Debugf("retrieved obj %s at %s", id, t.String())
	return nil
}

// safeSet will set an object safely on the registry.
func (c *Client) safeSet(obj api.Object) {
	ft := apiutil.GetFullType(obj)
	// Get object registry.
	objReg, ok := c.reg[ft]
	if !ok {
		// If no registry then create.
		objReg = map[string]api.Object{}
		c.reg[ft] = objReg
	}
	objReg[obj.GetObjectMetadata().ID] = obj
	c.logger.Debugf("stored obj %s at %s", obj.GetObjectMetadata().ID, ft)
}

// safeDelete will delete an object safely on the registry.
func (c *Client) safeDelete(fullID string) {
	t, id := apiutil.SplitFullID(fullID)
	// Get object registry.
	objReg, ok := c.reg[t.String()]
	if !ok {
		return
	}
	delete(objReg, id)
}

// Create will store an object in memory.
func (c *Client) Create(obj api.Object) (api.Object, error) {
	c.Lock()
	defer c.Unlock()
	fullID := apiutil.GetFullID(obj)
	if obj := c.safeGet(fullID); obj != nil {
		return nil, fmt.Errorf("node %s already present", obj.GetObjectMetadata().ID)
	}
	c.safeSet(obj)
	return obj, nil
}

// Update will update an object in memory.
func (c *Client) Update(obj api.Object) (api.Object, error) {
	c.Lock()
	defer c.Unlock()
	fullID := apiutil.GetFullID(obj)
	if o := c.safeGet(fullID); o == nil {
		return nil, fmt.Errorf("node %s not present", obj.GetObjectMetadata().ID)
	}
	c.safeSet(obj)
	return obj, nil
}

// Delete will delete an object from memory.
func (c *Client) Delete(fullID string) error {
	c.Lock()
	defer c.Unlock()
	c.safeDelete(fullID)
	return nil
}

// Get will retrieve an object from memory.
func (c *Client) Get(fullID string) (api.Object, error) {
	c.Lock()
	defer c.Unlock()
	o := c.safeGet(fullID)
	if o == nil {
		return nil, fmt.Errorf("node %s not present", fullID)
	}
	return o, nil
}

func (c *Client) listAll(opts api.ListOptions) ([]api.Object, error) {
	c.Lock()
	defer c.Unlock()
	ol := []api.Object{}

	ft := apiutil.GetFullType(opts)
	reg, ok := c.reg[ft]
	if ok {
		for _, obj := range reg {
			ol = append(ol, obj)
		}
	}

	return ol, nil
}

func (c *Client) listBySelector(opts api.ListOptions) ([]api.Object, error) {
	c.Lock()
	defer c.Unlock()
	ol := []api.Object{}

	ft := apiutil.GetFullType(opts)
	reg, ok := c.reg[ft]
	if ok {
		for _, obj := range reg {
			if util.SelectorMatchesLabels(obj.GetObjectMetadata().Labels, opts.LabelSelector) {
				ol = append(ol, obj)
			}
		}
	}
	return ol, nil
}

// List will retrieve a list of objects in memory.
func (c *Client) List(opts api.ListOptions) ([]api.Object, error) {

	// Return all
	if len(opts.LabelSelector) == 0 {
		return c.listAll(opts)
	}

	// Return filtered.
	return c.listBySelector(opts)
}

// Watch will watch object events from memory operations.
func (c *Client) Watch(opts api.ListOptions) (watch.Watcher, error) {
	return nil, fmt.Errorf("not implemented")
}
