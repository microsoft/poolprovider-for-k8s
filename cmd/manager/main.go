package main

import (
	"context"
	"fmt"
	"os"
	"log"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/microsoft/poolprovider-for-k8s/pkg/apis"
	"github.com/microsoft/poolprovider-for-k8s/pkg/controller"

	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Println("Failed to get watch namespace")
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Println(err, "")
		os.Exit(1)
	}

	ctx := context.TODO()
	// Become the leader before proceeding
	err = leader.Become(ctx, "k8s-poolprovider-lock")
	if err != nil {
		log.Println(err, "")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     restmapper.NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", "0.0.0.0", 8383),
	})
	if err != nil {
		log.Println(err, "")
		os.Exit(1)
	}

	log.Println("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Println(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Println(err, "")
		os.Exit(1)
	}

	log.Println("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Println(err, "Manager exited non-zero")
		os.Exit(1)
	}
}
