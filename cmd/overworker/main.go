package main

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/giongto35/cloud-game/pkg/util/logging"
	"github.com/giongto35/cloud-game/pkg/worker"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

func main() {
	// Due to the limit of libretro, one process can only run one game. For better performance, it's better run everything on one core
	runtime.GOMAXPROCS(1)
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := worker.NewDefaultConfig()
	cfg.AddFlags(pflag.CommandLine)

	logging.Init()
	defer logging.Flush()

	ctx, cancelCtx := context.WithCancel(context.Background())

	glog.Infof("Initializing worker server")
	glog.V(4).Infof("Worker configs %v", cfg)
	o := worker.New(ctx, cfg)
	if err := o.Run(); err != nil {
		glog.Errorf("Failed to run worker, reason %v", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	select {
	case <-stop:
		glog.Infoln("Received SIGTERM, Quiting Worker")
		o.Shutdown()
		cancelCtx()
	}
}
