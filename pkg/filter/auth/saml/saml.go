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
	"encoding/xml"
	"errors"
	"fmt"
	stdHttp "net/http"
	"net/url"
	"os"
	"strings"
)

import (
	samlcore "github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
)

import (
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
	if err := factory.initConfig(); err != nil {
		return err
	}

	acsURL, metadataURL, err := factory.parseServiceProviderURLs()
	if err != nil {
		return err
	}

	certificate, privateKey, err := factory.loadSigningCredentials()
	if err != nil {
		return err
	}

	idpMetadata, err := factory.loadIDPMetadata()
	if err != nil {
		return err
	}

	middleware, err := factory.newMiddleware(acsURL, metadataURL, certificate, privateKey, idpMetadata)
	if err != nil {
		return err
	}

	factory.middleware = middleware
	factory.serviceProvider = &middleware.ServiceProvider
	factory.metadataPath = metadataURL.Path
	factory.acsPath = acsURL.Path

	logger.Infof("SAML filter initialized with entity_id=%s metadata_path=%s acs_path=%s", factory.cfg.EntityID, factory.metadataPath, factory.acsPath)
	return nil
}

func (factory *FilterFactory) initConfig() error {
	if err := factory.cfg.Validate(); err != nil {
		return err
	}

	if factory.cfg.ErrMsg == "" {
		factory.cfg.ErrMsg = "SAML authentication failed"
	}
	factory.errMsg = []byte(factory.cfg.ErrMsg)
	return nil
}

func (factory *FilterFactory) parseServiceProviderURLs() (*url.URL, *url.URL, error) {
	acsURL, err := url.Parse(factory.cfg.AssertionConsumerURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse acs_url: %w", err)
	}
	metadataURL, err := url.Parse(factory.cfg.MetadataURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse metadata_url: %w", err)
	}
	if acsURL.Scheme == "" || acsURL.Host == "" {
		return nil, nil, fmt.Errorf("acs_url must be an absolute URL")
	}
	if metadataURL.Scheme == "" || metadataURL.Host == "" {
		return nil, nil, fmt.Errorf("metadata_url must be an absolute URL")
	}
	if metadataURL.Scheme != acsURL.Scheme || metadataURL.Host != acsURL.Host {
		return nil, nil, fmt.Errorf("metadata_url and acs_url must use the same scheme and host")
	}
	return acsURL, metadataURL, nil
}

func (factory *FilterFactory) loadSigningCredentials() (*x509.Certificate, *rsa.PrivateKey, error) {
	keyPair, err := tls.LoadX509KeyPair(factory.cfg.CertFile, factory.cfg.KeyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("load cert/key pair: %w", err)
	}
	if len(keyPair.Certificate) == 0 {
		return nil, nil, fmt.Errorf("certificate file %s does not contain a certificate", factory.cfg.CertFile)
	}

	certificate, err := x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, nil, fmt.Errorf("parse certificate leaf: %w", err)
	}
	privateKey, ok := keyPair.PrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("expected RSA private key, got %T", keyPair.PrivateKey)
	}
	return certificate, privateKey, nil
}

func (factory *FilterFactory) newMiddleware(
	acsURL, metadataURL *url.URL,
	certificate *x509.Certificate,
	privateKey *rsa.PrivateKey,
	idpMetadata *samlcore.EntityDescriptor,
) (*samlsp.Middleware, error) {
	rootURL := url.URL{Scheme: acsURL.Scheme, Host: acsURL.Host}

	// SAML ACS receives a cross-site POST from the IdP. For the request-tracking
	// cookie to survive that cross-site POST, browsers need SameSite=None + Secure.
	// This only works over HTTPS; for plain HTTP dev/test, use AllowIDPInitiated.
	cookieSameSite := stdHttp.SameSiteNoneMode
	if acsURL.Scheme != "https" {
		cookieSameSite = stdHttp.SameSiteDefaultMode
	}

	middleware, err := samlsp.New(samlsp.Options{
		EntityID:       factory.cfg.EntityID,
		URL:            rootURL,
		Key:            privateKey,
		Certificate:    certificate,
		IDPMetadata:    idpMetadata,
		CookieSameSite: cookieSameSite,
	})
	if err != nil {
		return nil, fmt.Errorf("create saml middleware: %w", err)
	}

	middleware.ServiceProvider.MetadataURL = *metadataURL
	middleware.ServiceProvider.AcsURL = *acsURL
	middleware.ServiceProvider.AllowIDPInitiated = factory.cfg.AllowIDPInitiated
	return middleware, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *pixiuhttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		cfg:             factory.cfg.DeepCopy(),
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

	// 1. Metadata endpoint — returns SP metadata XML for IdP configuration.
	if path == f.metadataPath {
		return f.handleMetadata(ctx)
	}

	// 2. ACS endpoint — processes the SAML Response posted by IdP.
	if path == f.acsPath {
		return f.handleACS(ctx)
	}

	// 3. Check if the request path matches any protected rule.
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

	// 4. Protected path: verify session cookie.
	session, err := f.middleware.Session.GetSession(ctx.Request)
	if err != nil {
		if errors.Is(err, samlsp.ErrNoSession) {
			// No valid session — redirect the browser to the IdP login page.
			logger.Debugf("SAML: no valid session for %s, redirecting to IdP", path)
			f.middleware.HandleStartAuthFlow(ctx.Writer, ctx.Request)
			return filter.Stop
		}
		// Unexpected error (e.g. cookie decryption failure) — don't silently redirect.
		logger.Errorf("SAML: session error for %s: %v", path, err)
		ctx.SendLocalReply(stdHttp.StatusInternalServerError, []byte("internal authentication error"))
		return filter.Stop
	}

	// 5. Valid session — forward configured attributes and let the request through.
	f.forwardAttributes(ctx, session)
	return filter.Continue
}

// handleMetadata serves the SP metadata document as XML.
// IdP administrators import this URL to configure their side of the trust.
func (f *Filter) handleMetadata(ctx *pixiuhttp.HttpContext) filter.FilterStatus {
	buf, err := xml.MarshalIndent(f.serviceProvider.Metadata(), "", "  ")
	if err != nil {
		logger.Errorf("SAML: failed to marshal metadata: %v", err)
		ctx.SendLocalReply(stdHttp.StatusInternalServerError, []byte("failed to generate metadata"))
		return filter.Stop
	}

	w := ctx.Writer
	w.Header().Set("Content-Type", "application/samlmetadata+xml; charset=utf-8")
	w.WriteHeader(stdHttp.StatusOK)
	_, _ = w.Write(buf)
	return filter.Stop
}

// handleACS delegates the Assertion Consumer Service request to the SAML middleware.
// The middleware validates the assertion, creates a session cookie, and redirects
// the user back to the originally requested URL (RelayState).
func (f *Filter) handleACS(ctx *pixiuhttp.HttpContext) filter.FilterStatus {
	f.middleware.ServeHTTP(ctx.Writer, ctx.Request)
	return filter.Stop
}

// forwardAttributes reads SAML attributes from the session and injects them
// into the request headers so that backend services can identify the user.
// It first removes any client-supplied values for the configured headers
// to prevent header spoofing.
func (f *Filter) forwardAttributes(ctx *pixiuhttp.HttpContext, session samlsp.Session) {
	// Always strip client-supplied values for SAML-controlled headers,
	// even if the session doesn't carry attributes.
	for _, fa := range f.cfg.ForwardAttributes {
		ctx.Request.Header.Del(fa.Header)
	}

	sa, ok := session.(samlsp.SessionWithAttributes)
	if !ok {
		return
	}
	attrs := sa.GetAttributes()
	for _, fa := range f.cfg.ForwardAttributes {
		if val := attrs.Get(fa.SAMLAttribute); val != "" {
			ctx.Request.Header.Set(fa.Header, val)
		}
	}
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
