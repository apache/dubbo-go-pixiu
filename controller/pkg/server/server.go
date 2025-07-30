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

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Server struct {
	Client *kubernetes.Clientset
	Host   string
	Name   string
	Port   string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	// Get all Ingress
	ingressList, err := s.Client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		http.Error(w, "Failed to list Ingresses", http.StatusInternalServerError)
		return
	}

	for _, ingress := range ingressList.Items {
		for _, rule := range ingress.Spec.Rules {
			if rule.Host != host {
				continue
			}

			for _, path := range rule.HTTP.Paths {
				if !strings.HasPrefix(r.URL.Path, path.Path) {
					continue
				}

				// Get Service
				svc, err := s.Client.CoreV1().Services(ingress.Namespace).Get(context.TODO(), path.Backend.Service.Name, metav1.GetOptions{})
				if err != nil {
					http.Error(w, "Service not found", http.StatusBadGateway)
					return
				}

				// Get Endpoints
				endpoints, err := s.Client.CoreV1().Endpoints(ingress.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
				if err != nil || len(endpoints.Subsets) == 0 {
					http.Error(w, "No backend Pod available", http.StatusBadGateway)
					return
				}

				// Find the first available Pod IP + port
				var targetIP string
				var targetPort int32
				var scheme string = "http" // default http

				found := false
				for _, subset := range endpoints.Subsets {
					for _, addr := range subset.Addresses {
						for _, port := range subset.Ports {
							targetIP = addr.IP
							targetPort = port.Port
							portName := strings.ToLower(port.Name)
							if port.Port == 443 || strings.Contains(portName, "https") {
								scheme = "https"
							}
							found = true
							break
						}
						if found {
							break
						}
					}
					if found {
						break
					}
				}

				if !found {
					http.Error(w, "No backend Pod available", http.StatusBadGateway)
					return
				}

				// Constructing the backend URL
				backendURL, _ := url.Parse(fmt.Sprintf("%s://%s:%d", scheme, targetIP, targetPort))
				log.Printf("[proxy] %s%s -> %s\n", r.Host, r.URL.Path, backendURL.String())

				// Building a reverse proxy
				proxy := httputil.NewSingleHostReverseProxy(backendURL)
				proxy.Director = func(req *http.Request) {
					req.URL.Scheme = backendURL.Scheme
					req.URL.Host = backendURL.Host
					req.Host = r.Host
				}
				proxy.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				}

				proxy.ServeHTTP(w, r)
				return
			}
		}
	}

	// If no match is found, return 404
	http.NotFound(w, r)
}
