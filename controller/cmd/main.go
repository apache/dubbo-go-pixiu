/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"github.com/apache/dubbo-go-pixiu/controller/pkg/controller"
	"github.com/apache/dubbo-go-pixiu/controller/pkg/proxy"
	"github.com/apache/dubbo-go-pixiu/controller/pkg/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

const IngressClassFlagHelpStr = "IngressClass for pixiu"

var (
	ingressClass string
)

func main() {
	// Define a command-line flag for ingress class, defaulting to "pixiu"
	flag.StringVar(&ingressClass, "IngressClass", "pixiu", IngressClassFlagHelpStr)
	flag.Parse()

	// Configure zerolog to output logs in a human-readable console format
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Function to get the user's home directory
	homeDir := func() string {
		if h := os.Getenv("HOME"); h != "" { // For Linux / Mac
			return h
		}
		return os.Getenv("USERPROFILE") // For Windows
	}

	// Try to load the in-cluster Kubernetes config (used when running inside a cluster)
	config, err := rest.InClusterConfig()
	if err != nil {
		// If not running inside a cluster, fall back to kubeconfig file in ~/.kube/config
		config, err = clientcmd.BuildConfigFromFlags("", filepath.Join(homeDir(), ".kube", "config"))
	}

	// Create a new Kubernetes client from the config
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get kubernetes client") // Exit if client creation fails
	}
	log.Info().Msg("get kubernetes client success")

	// Create a new Pixiu ingress controller for the specified ingress class
	ingressController := controller.NewIngressController(client, ingressClass)

	// Use a utility group to manage concurrent goroutines
	group := &utils.Group{}

	// Run the ingress controller in a separate goroutine
	group.Go(func() {
		ingressController.Run(5, nil) // 5 worker threads, no stop channel
	})

	// Run the Pixiu proxy in another goroutine
	group.Go(func() {
		err := proxy.Run(client, ingressClass)
		if err != nil {
			return
		}
	})

	// Wait for all goroutines to finish
	group.Wait()
}
