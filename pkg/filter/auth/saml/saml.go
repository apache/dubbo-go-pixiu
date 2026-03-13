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

package saml

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	stdHttp "net/http"
	"net/url"
	"os"
	"strings"

	samlcore "github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	pixiuhttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPAuthSamlFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct{}

	// FilterFactory is http filter instance factory.
	FilterFactory struct {
		cfg             *Config
		errMsg          []byte
		serviceProvider *samlcore.ServiceProvider
		middleware      *samlsp.Middleware
		metadataPath    string
		acsPath         string
	}

	// Filter is the actual filter instance.
	Filter struct {
		cfg             *Config
		errMsg          []byte
		serviceProvider *samlcore.ServiceProvider
		middleware      *samlsp.Middleware
		metadataPath    string
		acsPath         string
	}
)

func (p Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	if err := factory.cfg.Validate(); err != nil {
		return err
	}

	if factory.cfg.ErrMsg == "" {
		factory.cfg.ErrMsg = "SAML authentication failed"
	}
	factory.errMsg = []byte(factory.cfg.ErrMsg)

	acsURL, err := url.Parse(factory.cfg.AssertionConsumerURL)
	if err != nil {
		return fmt.Errorf("parse acs_url: %w", err)
	}
	metadataURL, err := url.Parse(factory.cfg.MetadataURL)
	if err != nil {
		return fmt.Errorf("parse metadata_url: %w", err)
	}
	if acsURL.Scheme == "" || acsURL.Host == "" {
		return fmt.Errorf("acs_url must be an absolute URL")
	}
	if metadataURL.Scheme == "" || metadataURL.Host == "" {
		return fmt.Errorf("metadata_url must be an absolute URL")
	}
	if metadataURL.Scheme != acsURL.Scheme || metadataURL.Host != acsURL.Host {
		return fmt.Errorf("metadata_url and acs_url must use the same scheme and host")
	}

	keyPair, err := tls.LoadX509KeyPair(factory.cfg.CertFile, factory.cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("load cert/key pair: %w", err)
	}
	if len(keyPair.Certificate) == 0 {
		return fmt.Errorf("certificate file %s does not contain a certificate", factory.cfg.CertFile)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return fmt.Errorf("parse certificate leaf: %w", err)
	}
	privateKey, ok := keyPair.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return fmt.Errorf("expected RSA private key, got %T", keyPair.PrivateKey)
	}

	idpMetadata, err := factory.loadIDPMetadata()
	if err != nil {
		return err
	}

	rootURL := url.URL{Scheme: acsURL.Scheme, Host: acsURL.Host}
	middleware, err := samlsp.New(samlsp.Options{
		EntityID:    factory.cfg.EntityID,
		URL:         rootURL,
		Key:         privateKey,
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})
	if err != nil {
		return fmt.Errorf("create saml middleware: %w", err)
	}

	middleware.ServiceProvider.MetadataURL = *metadataURL
	middleware.ServiceProvider.AcsURL = *acsURL
	factory.middleware = middleware
	factory.serviceProvider = &middleware.ServiceProvider
	factory.metadataPath = metadataURL.Path
	factory.acsPath = acsURL.Path

	logger.Infof("SAML filter initialized with entity_id=%s metadata_path=%s acs_path=%s", factory.cfg.EntityID, factory.metadataPath, factory.acsPath)
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *pixiuhttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		cfg:             factory.cfg,
		errMsg:          factory.errMsg,
		serviceProvider: factory.serviceProvider,
		middleware:      factory.middleware,
		metadataPath:    factory.metadataPath,
		acsPath:         factory.acsPath,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *pixiuhttp.HttpContext) filter.FilterStatus {
	path := ctx.Request.URL.Path

	matched := false
	for _, rule := range f.cfg.Rules {
		if strings.HasPrefix(path, rule.Match.Prefix) {
			matched = true
			break
		}
	}

	if !matched {
		return filter.Continue
	}

	logger.Warnf("SAML authentication not yet implemented for path: %s", path)
	ctx.SendLocalReply(stdHttp.StatusUnauthorized, f.errMsg)
	return filter.Stop
}

func (factory *FilterFactory) loadIDPMetadata() (*samlcore.EntityDescriptor, error) {
	if factory.cfg.IdPMetadataURL != "" {
		metadataURL, err := url.Parse(factory.cfg.IdPMetadataURL)
		if err != nil {
			return nil, fmt.Errorf("parse idp_metadata_url: %w", err)
		}
		idpMetadata, err := samlsp.FetchMetadata(context.Background(), stdHttp.DefaultClient, *metadataURL)
		if err != nil {
			return nil, fmt.Errorf("fetch idp metadata: %w", err)
		}
		return idpMetadata, nil
	}

	data, err := os.ReadFile(factory.cfg.IdPMetadataFile)
	if err != nil {
		return nil, fmt.Errorf("read idp_metadata_file: %w", err)
	}
	idpMetadata, err := samlsp.ParseMetadata(data)
	if err != nil {
		return nil, fmt.Errorf("parse idp metadata file: %w", err)
	}
	return idpMetadata, nil
}
