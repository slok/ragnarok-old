package main

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/apimachinery/validator"
	"github.com/slok/ragnarok/apimachinery/watch"
	clichaosv1 "github.com/slok/ragnarok/client/api/chaos/v1"
	cliclusterv1 "github.com/slok/ragnarok/client/api/cluster/v1"
	"github.com/slok/ragnarok/client/controller"
	"github.com/slok/ragnarok/client/informer"
	"github.com/slok/ragnarok/client/repository"
	memrepository "github.com/slok/ragnarok/client/repository/memory"
	"github.com/slok/ragnarok/client/util/queue"
	"github.com/slok/ragnarok/client/util/store"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/cmd/master/flags"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	controlleripm "github.com/slok/ragnarok/master/controller"
	"github.com/slok/ragnarok/master/server"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/master/service/experiment"
	"github.com/slok/ragnarok/master/web"
)

// master dependencies is a helper object to group all the app dependencies
type masterDependencies struct {
	experimentClient clichaosv1.ExperimentClientInterface
	nodeClient       cliclusterv1.NodeClientInterface
	failureClient    clichaosv1.FailureClientInterface
	logger           log.Logger
	nodeStatus       service.NodeStatusService
	failureStatus    service.FailureStatusService
	serializer       serializer.Serializer
}

func createGRPCServer(cfg config.Config, deps masterDependencies, logger log.Logger) (*server.MasterGRPCServiceServer, error) {

	// Create the GRPC service server
	l, err := net.Listen("tcp", cfg.RPCListenAddress)
	if err != nil {
		return nil, err
	}
	srvServer := server.NewMasterGRPCServiceServer(deps.failureStatus, deps.nodeStatus, l, clock.Base(), logger)
	return srvServer, nil
}

func createHTTPServer(cfg config.Config, deps masterDependencies, logger log.Logger) (web.Server, error) {
	// Create the GRPC service server
	l, err := net.Listen("tcp", cfg.HTTPListenAddress)
	if err != nil {
		return nil, err
	}

	return web.NewDefaultHTTPServer(deps.serializer, deps.nodeClient, l, logger)
}

// TODO: Debugging stuff, remove.
func createExperimentController(nodeCli cliclusterv1.NodeClientInterface, failureCli clichaosv1.FailureClientInterface, repository repository.Client, logger log.Logger) (controller.Controller, error) {
	indexer := store.ObjectIndexKeyerFunc(func(obj api.Object) (string, error) {
		return apiutil.GetFullID(obj), nil
	})
	cache := store.NewIndexedStore(indexer, &sync.Map{}, logger)
	queue := queue.NewSimpleQueue()
	lw := informer.NewRepositoryListerWatcher(repository)
	lwOpts := api.ListOptions{
		TypeMeta: chaosv1.ExperimentTypeMeta,
	}

	logger = logger.WithField("controller", "experiment")
	inf := informer.NewWorkQueueInformer(indexer, queue, cache, lwOpts, lw, logger)
	service := experiment.NewSimpleManager(nodeCli, failureCli, logger)
	c := controlleripm.NewExperiment(inf, service, logger)
	return c, nil
}

// Main run main logic.
func Main() error {
	logger := log.Base()

	// Get the command line arguments.
	cfg, err := flags.GetMasterConfig(os.Args[1:])
	if err != nil {
		logger.Error(err)
		return err
	}

	// Set debug mode
	if cfg.Debug {
		logger.Set("debug")
	}

	// Create dependencies
	eventMux := watch.NewDefaultBroadcasterFactory(logger)
	memoryRepoClient := memrepository.NewDefaultClient(eventMux, logger)
	validator := validator.DefaultObject
	nodeCli := cliclusterv1.NewNodeClient(validator, memoryRepoClient)
	failureCli := clichaosv1.NewFailureClient(validator, memoryRepoClient)
	experimentCli := clichaosv1.NewExperimentClient(validator, memoryRepoClient)

	deps := masterDependencies{
		nodeClient:       nodeCli,
		failureClient:    failureCli,
		experimentClient: experimentCli,
		nodeStatus:       service.NewNodeStatus(*cfg, nodeCli, logger),
		failureStatus:    service.NewFailureStatus(failureCli, logger),
		serializer:       serializer.DefaultSerializer,
	}

	// Start dummy controller.
	experimentCtl, err := createExperimentController(nodeCli, failureCli, memoryRepoClient, logger)
	if err != nil {
		return err
	}
	go func() {
		experimentCtl.Run()
	}()

	// TODO: Autoregister this node as a master node.
	grpcServer, err := createGRPCServer(*cfg, deps, logger)
	if err != nil {
		return err
	}
	go func() {
		grpcServer.Serve()
	}()

	httpServer, err := createHTTPServer(*cfg, deps, logger)
	if err != nil {
		return err
	}
	httpServer.Serve()

	return nil
}

func clean() {
	log.Debug("Cleaning...")
}

func main() {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	errC := make(chan error)

	// Run main program
	go func() {
		if err := Main(); err != nil {
			errC <- err
		}
		return
	}()

	// Wait until signal (ctr+c, SIGTERM...).
	var exitCode int

Waiter:
	for {
		select {
		// Wait for errors
		case err := <-errC:
			if err != nil {
				exitCode = 1
				break Waiter
			}
			// Wait for signal
		case <-sigC:
			break Waiter
		}
	}

	clean()
	os.Exit(exitCode)
}
