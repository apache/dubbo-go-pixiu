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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	stdHttp "net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

import (
	samlcore "github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const testSPEntityID = "test-sp"

// =============================================================================
// Test Helpers
// =============================================================================

// generateTestCert creates a self-signed certificate and key in temp files.
func generateTestCert(t *testing.T) (certFile, keyFile string) {
	t.Helper()
	dir := t.TempDir()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: testSPEntityID},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)

	certFile = filepath.Join(dir, "sp.crt")
	keyFile = filepath.Join(dir, "sp.key")

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	require.NoError(t, os.WriteFile(certFile, certPEM, 0644))

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	require.NoError(t, os.WriteFile(keyFile, keyPEM, 0644))

	return certFile, keyFile
}

// startTestIdPServer starts an httptest server that serves minimal IdP metadata.
func startTestIdPServer(t *testing.T) *httptest.Server {
	t.Helper()

	idpKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	idpTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "test-idp"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}
	idpCertDER, err := x509.CreateCertificate(rand.Reader, idpTemplate, idpTemplate, &idpKey.PublicKey, idpKey)
	require.NoError(t, err)

	idpCertB64 := base64.StdEncoding.EncodeToString(idpCertDER)

	server := httptest.NewServer(stdHttp.HandlerFunc(func(w stdHttp.ResponseWriter, r *stdHttp.Request) {
		metadata := fmt.Sprintf(`<EntityDescriptor entityID="%s" xmlns="urn:oasis:names:tc:SAML:2.0:metadata">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <KeyDescriptor use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>%s</X509Certificate>
        </X509Data>
      </KeyInfo>
    </KeyDescriptor>
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="%s/sso"/>
  </IDPSSODescriptor>
</EntityDescriptor>`, r.Host, idpCertB64, "http://"+r.Host)
		w.Header().Set("Content-Type", "application/xml")
		_, _ = w.Write([]byte(metadata))
	}))
	t.Cleanup(server.Close)
	return server
}

// buildTestFactory creates a fully initialized FilterFactory for testing.
func buildTestFactory(t *testing.T) *FilterFactory {
	t.Helper()

	return buildTestFactoryWithScheme(t, "http")
}

func buildTestFactoryWithScheme(t *testing.T, scheme string) *FilterFactory {
	t.Helper()

	certFile, keyFile := generateTestCert(t)
	idpServer := startTestIdPServer(t)

	factory := &FilterFactory{cfg: &Config{
		EntityID:             testSPEntityID,
		AssertionConsumerURL: fmt.Sprintf("%s://localhost:8888/saml/acs", scheme),
		MetadataURL:          fmt.Sprintf("%s://localhost:8888/saml/metadata", scheme),
		IdPMetadataURL:       idpServer.URL,
		CertFile:             certFile,
		KeyFile:              keyFile,
		AllowIDPInitiated:    true,
		Rules: []Rule{
			{Match: Match{Prefix: "/api"}},
			{Match: Match{Prefix: "/admin"}},
		},
		ForwardAttributes: []ForwardAttribute{
			{SAMLAttribute: "email", Header: "X-User-Email"},
			{SAMLAttribute: "displayName", Header: "X-User-Name"},
		},
	}}

	err := factory.Apply()
	require.NoError(t, err)
	return factory
}

// buildTestFilter creates a Filter from a fully initialized factory.
func buildTestFilter(t *testing.T) *Filter {
	t.Helper()
	factory := buildTestFactory(t)
	return &Filter{
		cfg:             factory.cfg,
		errMsg:          factory.errMsg,
		serviceProvider: factory.serviceProvider,
		middleware:      factory.middleware,
		metadataPath:    factory.metadataPath,
		acsPath:         factory.acsPath,
	}
}

func newHttpContext(method, path string) (*contexthttp.HttpContext, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	return &contexthttp.HttpContext{Writer: rec, Request: req}, rec
}

// mockSessionProvider implements samlsp.SessionProvider for testing.
type mockSessionProvider struct {
	session samlsp.Session
	err     error
}

func (m *mockSessionProvider) GetSession(_ *stdHttp.Request) (samlsp.Session, error) {
	return m.session, m.err
}

func (m *mockSessionProvider) CreateSession(_ stdHttp.ResponseWriter, _ *stdHttp.Request, _ *samlcore.Assertion) error {
	return nil
}

func (m *mockSessionProvider) DeleteSession(_ stdHttp.ResponseWriter, _ *stdHttp.Request) error {
	return nil
}

// mockSession implements samlsp.SessionWithAttributes for testing.
type mockSession struct {
	attrs samlsp.Attributes
}

func (m *mockSession) GetAttributes() samlsp.Attributes {
	return m.attrs
}

// =============================================================================
// Plugin Tests
// =============================================================================

func TestPluginKind(t *testing.T) {
	p := &Plugin{}
	assert.Equal(t, "dgp.filter.http.auth.saml", p.Kind())
}

func TestCreateFilterFactory(t *testing.T) {
	p := &Plugin{}
	factory, err := p.CreateFilterFactory()
	assert.NoError(t, err)
	assert.NotNil(t, factory)
}

// =============================================================================
// Apply Tests
// =============================================================================

func TestApply_Success(t *testing.T) {
	factory := buildTestFactory(t)
	assert.NotNil(t, factory.middleware)
	assert.NotNil(t, factory.serviceProvider)
	assert.Equal(t, "/saml/metadata", factory.metadataPath)
	assert.Equal(t, "/saml/acs", factory.acsPath)
}

func TestApply_InvalidConfig(t *testing.T) {
	factory := &FilterFactory{cfg: &Config{}}
	err := factory.Apply()
	assert.Error(t, err)
}

func TestApply_InvalidCertFile(t *testing.T) {
	idpServer := startTestIdPServer(t)
	factory := &FilterFactory{cfg: &Config{
		EntityID:             testSPEntityID,
		AssertionConsumerURL: "http://localhost:8888/saml/acs",
		MetadataURL:          "http://localhost:8888/saml/metadata",
		IdPMetadataURL:       idpServer.URL,
		CertFile:             "/nonexistent/sp.crt",
		KeyFile:              "/nonexistent/sp.key",
	}}
	err := factory.Apply()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "load cert/key pair")
}

func TestApply_HTTPSCookieSameSite(t *testing.T) {
	factory := buildTestFactoryWithScheme(t, "https")

	sessionProvider, ok := factory.middleware.Session.(samlsp.CookieSessionProvider)
	require.True(t, ok)
	assert.Equal(t, stdHttp.SameSiteNoneMode, sessionProvider.SameSite)
	assert.True(t, sessionProvider.Secure)

	requestTracker, ok := factory.middleware.RequestTracker.(samlsp.CookieRequestTracker)
	require.True(t, ok)
	assert.Equal(t, stdHttp.SameSiteNoneMode, requestTracker.SameSite)
}

func TestApply_HTTPCookieSameSite(t *testing.T) {
	factory := buildTestFactoryWithScheme(t, "http")

	sessionProvider, ok := factory.middleware.Session.(samlsp.CookieSessionProvider)
	require.True(t, ok)
	assert.Equal(t, stdHttp.SameSiteDefaultMode, sessionProvider.SameSite)
	assert.False(t, sessionProvider.Secure)

	requestTracker, ok := factory.middleware.RequestTracker.(samlsp.CookieRequestTracker)
	require.True(t, ok)
	assert.Equal(t, stdHttp.SameSiteDefaultMode, requestTracker.SameSite)
}

// =============================================================================
// Decode Tests
// =============================================================================

func TestDecode_UnprotectedPath(t *testing.T) {
	f := buildTestFilter(t)

	ctx, _ := newHttpContext("GET", "/public/health")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Continue, result)
}

func TestDecode_MetadataEndpoint(t *testing.T) {
	f := buildTestFilter(t)

	ctx, rec := newHttpContext("GET", "/saml/metadata")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Stop, result)
	assert.Equal(t, stdHttp.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/samlmetadata+xml")

	body := rec.Body.String()
	assert.Contains(t, body, fmt.Sprintf(`entityID="%s"`, testSPEntityID))
	assert.Contains(t, body, "AssertionConsumerService")
}

func TestDecode_ProtectedPath_NoSession(t *testing.T) {
	f := buildTestFilter(t)
	// ErrNoSession means "user hasn't logged in yet" → should redirect to IdP
	f.middleware.Session = &mockSessionProvider{err: samlsp.ErrNoSession}

	ctx, rec := newHttpContext("GET", "/api/users")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Stop, result)
	// Should redirect to IdP (302) or return an auth flow response
	assert.True(t, rec.Code == stdHttp.StatusFound || rec.Code == stdHttp.StatusOK,
		"expected redirect or auth flow response, got %d", rec.Code)
}

func TestDecode_ProtectedPath_SessionError(t *testing.T) {
	f := buildTestFilter(t)
	// A non-ErrNoSession error (e.g. cookie decryption failure) → should return 500
	f.middleware.Session = &mockSessionProvider{err: errors.New("crypto: bad key")}

	ctx, rec := newHttpContext("GET", "/api/users")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Stop, result)
	assert.Equal(t, stdHttp.StatusInternalServerError, rec.Code)
}

func TestDecode_ProtectedPath_ValidSession(t *testing.T) {
	f := buildTestFilter(t)
	// Replace session provider with mock that returns a valid session
	f.middleware.Session = &mockSessionProvider{
		session: &mockSession{
			attrs: samlsp.Attributes{
				"email":       {"test@example.com"},
				"displayName": {"Test User"},
			},
		},
	}

	ctx, _ := newHttpContext("GET", "/api/users")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Continue, result)
	// Verify attributes were forwarded
	assert.Equal(t, "test@example.com", ctx.Request.Header.Get("X-User-Email"))
	assert.Equal(t, "Test User", ctx.Request.Header.Get("X-User-Name"))
}

func TestDecode_ProtectedPath_SessionWithoutAttributes(t *testing.T) {
	f := buildTestFilter(t)
	// Session that does not implement SessionWithAttributes
	f.middleware.Session = &mockSessionProvider{session: "plain-session"}

	ctx, _ := newHttpContext("GET", "/api/users")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Continue, result)
	// No attributes forwarded
	assert.Empty(t, ctx.Request.Header.Get("X-User-Email"))
}

func TestDecode_ACSEndpoint(t *testing.T) {
	f := buildTestFilter(t)

	ctx, rec := newHttpContext("POST", "/saml/acs")
	result := f.Decode(ctx)

	assert.Equal(t, filter.Stop, result)
	// ACS without a valid SAMLResponse should return an error status
	assert.True(t, rec.Code >= 400, "expected error response for empty ACS request, got %d", rec.Code)
}

func TestDecode_MultipleRules(t *testing.T) {
	f := buildTestFilter(t)
	f.middleware.Session = &mockSessionProvider{err: samlsp.ErrNoSession}

	// /admin should also be protected
	ctx, _ := newHttpContext("GET", "/admin/dashboard")
	result := f.Decode(ctx)
	assert.Equal(t, filter.Stop, result)

	// /public should not be protected
	ctx2, _ := newHttpContext("GET", "/public/docs")
	result2 := f.Decode(ctx2)
	assert.Equal(t, filter.Continue, result2)
}

// =============================================================================
// ForwardAttributes Tests
// =============================================================================

func TestForwardAttributes_WithAttributes(t *testing.T) {
	f := buildTestFilter(t)
	ctx, _ := newHttpContext("GET", "/api/test")

	session := &mockSession{
		attrs: samlsp.Attributes{
			"email":       {"user@example.com"},
			"displayName": {"Jane Doe"},
		},
	}
	f.forwardAttributes(ctx, session)

	assert.Equal(t, "user@example.com", ctx.Request.Header.Get("X-User-Email"))
	assert.Equal(t, "Jane Doe", ctx.Request.Header.Get("X-User-Name"))
}

func TestForwardAttributes_MissingAttribute(t *testing.T) {
	f := buildTestFilter(t)
	ctx, _ := newHttpContext("GET", "/api/test")

	// Session only has email, not displayName
	session := &mockSession{
		attrs: samlsp.Attributes{
			"email": {"user@example.com"},
		},
	}
	f.forwardAttributes(ctx, session)

	assert.Equal(t, "user@example.com", ctx.Request.Header.Get("X-User-Email"))
	assert.Empty(t, ctx.Request.Header.Get("X-User-Name"))
}

func TestForwardAttributes_EmptyConfig(t *testing.T) {
	f := buildTestFilter(t)
	f.cfg.ForwardAttributes = nil
	ctx, _ := newHttpContext("GET", "/api/test")

	session := &mockSession{
		attrs: samlsp.Attributes{
			"email": {"user@example.com"},
		},
	}
	f.forwardAttributes(ctx, session)

	// Nothing should be forwarded
	assert.Empty(t, ctx.Request.Header.Get("X-User-Email"))
}

func TestForwardAttributes_NonAttributeSession(t *testing.T) {
	f := buildTestFilter(t)
	ctx, _ := newHttpContext("GET", "/api/test")

	// Pass a session that doesn't implement SessionWithAttributes
	f.forwardAttributes(ctx, "plain-session")

	assert.Empty(t, ctx.Request.Header.Get("X-User-Email"))
}

func TestForwardAttributes_ClearsClientSpoofedHeaders(t *testing.T) {
	f := buildTestFilter(t)
	ctx, _ := newHttpContext("GET", "/api/test")

	// Simulate a malicious client setting SAML-controlled headers
	ctx.Request.Header.Set("X-User-Email", "attacker@evil.com")
	ctx.Request.Header.Set("X-User-Name", "Attacker")

	// Session has email but NOT displayName
	session := &mockSession{
		attrs: samlsp.Attributes{
			"email": {"real@example.com"},
		},
	}
	f.forwardAttributes(ctx, session)

	// email should be overwritten with the real value
	assert.Equal(t, "real@example.com", ctx.Request.Header.Get("X-User-Email"))
	// displayName was not in the assertion → spoofed header must be gone
	assert.Empty(t, ctx.Request.Header.Get("X-User-Name"))
}

func TestForwardAttributes_ClearsHeadersEvenWithoutAttributes(t *testing.T) {
	f := buildTestFilter(t)
	ctx, _ := newHttpContext("GET", "/api/test")

	// Client sends spoofed header
	ctx.Request.Header.Set("X-User-Email", "attacker@evil.com")

	// Session doesn't support attributes at all
	f.forwardAttributes(ctx, "plain-session")

	// Spoofed header must still be removed
	assert.Empty(t, ctx.Request.Header.Get("X-User-Email"))
}

// =============================================================================
// HandleMetadata Tests
// =============================================================================

func TestHandleMetadata_ContainsSPInfo(t *testing.T) {
	f := buildTestFilter(t)
	ctx, rec := newHttpContext("GET", "/saml/metadata")

	result := f.handleMetadata(ctx)

	assert.Equal(t, filter.Stop, result)
	body := rec.Body.String()
	assert.Contains(t, body, fmt.Sprintf(`entityID="%s"`, testSPEntityID))
	assert.Contains(t, body, "http://localhost:8888/saml/acs")
	assert.Contains(t, body, "SPSSODescriptor")
	assert.True(t, strings.HasPrefix(rec.Header().Get("Content-Type"), "application/samlmetadata+xml"))
}

// =============================================================================
// PrepareFilterChain Tests
// =============================================================================

type mockFilterChain struct {
	decodeFilters []filter.HttpDecodeFilter
}

func (m *mockFilterChain) AppendDecodeFilters(filters ...filter.HttpDecodeFilter) {
	m.decodeFilters = append(m.decodeFilters, filters...)
}

func (m *mockFilterChain) AppendEncodeFilters(_ ...filter.HttpEncodeFilter) {}

func (m *mockFilterChain) OnDecode(_ *contexthttp.HttpContext) {}

func (m *mockFilterChain) OnEncode(_ *contexthttp.HttpContext) {}

func TestPrepareFilterChain(t *testing.T) {
	factory := buildTestFactory(t)
	ctx, _ := newHttpContext("GET", "/")
	chain := &mockFilterChain{}

	err := factory.PrepareFilterChain(ctx, chain)

	assert.NoError(t, err)
	assert.Len(t, chain.decodeFilters, 1)
}

func TestPrepareFilterChain_ConfigIsolation(t *testing.T) {
	factory := buildTestFactory(t)
	ctx, _ := newHttpContext("GET", "/")
	chain := &mockFilterChain{}

	_ = factory.PrepareFilterChain(ctx, chain)
	f := chain.decodeFilters[0].(*Filter)

	// Mutating factory config should NOT affect the filter's copy
	factory.cfg.Rules = append(factory.cfg.Rules, Rule{Match: Match{Prefix: "/new"}})
	assert.NotEqual(t, len(factory.cfg.Rules), len(f.cfg.Rules),
		"filter config should be independent of factory config")
}
