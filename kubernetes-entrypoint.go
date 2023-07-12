package main

import (
	"flag"
	"os"

	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/config"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/container"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/customresource"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/daemonset"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/job"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/pod"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/service"
	_ "opendev.org/airship/kubernetes-entrypoint/dependencies/socket"
	entry "opendev.org/airship/kubernetes-entrypoint/entrypoint"
	"opendev.org/airship/kubernetes-entrypoint/logger"
	command "opendev.org/airship/kubernetes-entrypoint/util/command"
	"opendev.org/airship/kubernetes-entrypoint/util/env"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var comm []string
	var entrypoint *entry.Entrypoint
	var resetconfig *rest.Config
	var err error
	kc := flag.String("kubeconfig", "", "kubeconfig file path")
	flag.Parse()
	if kc != nil && *kc != "" {
		content, err := os.ReadFile(*kc)
		if err != nil {
			logger.Warning.Printf("read file %s failed: %s. use incluster config", *kc, err)
		} else {
			resetconfig, err = clientcmd.RESTConfigFromKubeConfig(content)
			if err != nil {
				logger.Warning.Printf("parse kubeconfig file failed: %s. use incluster config", err)
			}
		}
	}
	if entrypoint, err = entry.New(resetconfig); err != nil {
		logger.Error.Printf("Creating entrypoint failed: %v", err)
		os.Exit(1)
	}

	entrypoint.Resolve()

	if comm = env.SplitCommand(); len(comm) == 0 {
		// TODO(DTadrzak): we should consider other options to handle whether pod
		// is an init-container
		logger.Warning.Printf("COMMAND env is empty")
		os.Exit(0)
	}

	if err = command.Execute(comm); err != nil {
		logger.Error.Printf("Cannot execute command: %v", err)
		os.Exit(1)
	}
}
