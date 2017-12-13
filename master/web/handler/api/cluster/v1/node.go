package v1

import (
	"bytes"
	"io/ioutil"
	"net/http"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/serializer"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/master/web/handler/util"
)

const (
	nodeRoute = "/api/cluster/v1/node"
)

// NodeHandler is the handler that handlers Node resources.
type NodeHandler struct {
	serializer serializer.Serializer
	nodeCli    cliclusterv1.NodeClientInterface
}

// NewNodeHandler returns a new NodeHandler.
func NewNodeHandler(serializer serializer.Serializer, nodeCli cliclusterv1.NodeClientInterface) *NodeHandler {
	return &NodeHandler{
		serializer: serializer,
		nodeCli:    nodeCli,
	}
}

// Create will create a new node.
func (n *NodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.SetJSONBadRequest(w, err.Error())
		return
	}

	nodeTmp, err := n.serializer.Decode(b)
	if err != nil {
		util.SetJSONBadRequest(w, err.Error())
		return
	}

	node, ok := nodeTmp.(*clusterv1.Node)
	if !ok {
		util.SetJSONBadRequest(w, "decoded object is not a node")
		return
	}

	newNode, err := n.nodeCli.Create(node)
	if err != nil {
		util.SetJSONInternalError(w, err.Error())
		return
	}

	// TODO return node created node.
	var nb bytes.Buffer
	n.serializer.Encode(newNode, &nb)
	util.SetJSONOK(w, nb.Bytes())
}

// Update updates a node resource.
func (n *NodeHandler) Update(w http.ResponseWriter, r *http.Request, id string) {
	util.SetJSONNotImplementedError(w)
}

// Delete deletes a node.
func (n *NodeHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	util.SetJSONNotImplementedError(w)
}

// Get gets a node.
func (n *NodeHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	util.SetJSONNotImplementedError(w)
}

// List lists nodes.
func (n *NodeHandler) List(w http.ResponseWriter, r *http.Request, opts map[string]string) {
	util.SetJSONNotImplementedError(w)
}

// Watch watches nodes.
func (n *NodeHandler) Watch(w http.ResponseWriter, r *http.Request, opts map[string]string) {
	util.SetJSONNotImplementedError(w)
}

// GetRoute returns the route where teh handlers of nodes will listen.
func (n *NodeHandler) GetRoute() string {
	return nodeRoute
}
