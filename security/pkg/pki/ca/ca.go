// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ca

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

import (
	"istio.io/pkg/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/security/pkg/cmd"
	k8ssecret "github.com/apache/dubbo-go-pixiu/security/pkg/k8s/secret"
	caerror "github.com/apache/dubbo-go-pixiu/security/pkg/pki/error"
	"github.com/apache/dubbo-go-pixiu/security/pkg/pki/util"
	certutil "github.com/apache/dubbo-go-pixiu/security/pkg/util"
)

const (
	// istioCASecretType is the Istio secret annotation type.
	istioCASecretType = "istio.io/ca-root"

	// CACertFile is the CA certificate chain file.
	CACertFile = "ca-cert.pem"
	// CAPrivateKeyFile is the private key file of CA.
	CAPrivateKeyFile = "ca-key.pem"
	// CASecret stores the key/cert of self-signed CA for persistency purpose.
	CASecret = "istio-ca-secret"
	// CertChainFile is the ID/name for the certificate chain file.
	CertChainFile = "cert-chain.pem"
	// PrivateKeyFile is the ID/name for the private key file.
	PrivateKeyFile = "key.pem"
	// RootCertFile is the ID/name for the CA root certificate file.
	RootCertFile = "root-cert.pem"
	// TLSSecretCACertFile is the CA certificate file name as it exists in tls type k8s secret.
	TLSSecretCACertFile = "tls.crt"
	// TLSSecretCAPrivateKeyFile is the CA certificate key file name as it exists in tls type k8s secret.
	TLSSecretCAPrivateKeyFile = "tls.key"
	// TLSSecretRootCertFile is the root cert file name as it exists in tls type k8s secret.
	TLSSecretRootCertFile = "ca.crt"
	// The standard key size to use when generating an RSA private key
	rsaKeySize = 2048
)

// SigningCAFileBundle locations of the files used for the signing CA
type SigningCAFileBundle struct {
	RootCertFile    string
	CertChainFiles  []string
	SigningCertFile string
	SigningKeyFile  string
}

var pkiCaLog = log.RegisterScope("pkica", "Citadel CA log", 0)

// caTypes is the enum for the CA type.
type caTypes int

type CertOpts struct {
	// SubjectIDs are used for building the SAN extension for the certificate.
	SubjectIDs []string

	// TTL is the requested lifetime (Time to live) to be applied in the certificate.
	TTL time.Duration

	// ForCA indicates whether the signed certificate if for CA.
	// If true, the signed certificate is a CA certificate, otherwise, it is a workload certificate.
	ForCA bool

	// Cert Signer info
	CertSigner string
}

const (
	// selfSignedCA means the Istio CA uses a self signed certificate.
	selfSignedCA caTypes = iota
	// pluggedCertCA means the Istio CA uses a operator-specified key/cert.
	pluggedCertCA
)

// IstioCAOptions holds the configurations for creating an Istio CA.
// TODO(myidpt): remove IstioCAOptions.
type IstioCAOptions struct {
	CAType caTypes

	DefaultCertTTL time.Duration
	MaxCertTTL     time.Duration
	CARSAKeySize   int

	KeyCertBundle *util.KeyCertBundle

	// Config for creating self-signed root cert rotator.
	RotatorConfig *SelfSignedCARootCertRotatorConfig
}

// NewSelfSignedIstioCAOptions returns a new IstioCAOptions instance using self-signed certificate.
func NewSelfSignedIstioCAOptions(ctx context.Context,
	rootCertGracePeriodPercentile int, caCertTTL, rootCertCheckInverval, defaultCertTTL,
	maxCertTTL time.Duration, org string, dualUse bool, namespace string,
	readCertRetryInterval time.Duration, client corev1.CoreV1Interface,
	rootCertFile string, enableJitter bool, caRSAKeySize int) (caOpts *IstioCAOptions, err error) {
	// For the first time the CA is up, if readSigningCertOnly is unset,
	// it generates a self-signed key/cert pair and write it to CASecret.
	// For subsequent restart, CA will reads key/cert from CASecret.
	caSecret, scrtErr := client.Secrets(namespace).Get(context.TODO(), CASecret, metav1.GetOptions{})
	if scrtErr != nil && readCertRetryInterval > time.Duration(0) {
		pkiCaLog.Infof("Citadel in signing key/cert read only mode. Wait until secret %s:%s can be loaded...", namespace, CASecret)
		ticker := time.NewTicker(readCertRetryInterval)
		defer ticker.Stop()
		for scrtErr != nil {
			select {
			case <-ticker.C:
				if caSecret, scrtErr = client.Secrets(namespace).Get(context.TODO(), CASecret, metav1.GetOptions{}); scrtErr == nil {
					pkiCaLog.Infof("Citadel successfully loaded the secret.")
				}
			case <-ctx.Done():
				pkiCaLog.Errorf("Secret waiting thread is terminated.")
				return nil, fmt.Errorf("secret waiting thread is terminated")
			}
		}
	}

	caOpts = &IstioCAOptions{
		CAType:         selfSignedCA,
		DefaultCertTTL: defaultCertTTL,
		MaxCertTTL:     maxCertTTL,
		RotatorConfig: &SelfSignedCARootCertRotatorConfig{
			CheckInterval:      rootCertCheckInverval,
			caCertTTL:          caCertTTL,
			retryInterval:      cmd.ReadSigningCertRetryInterval,
			retryMax:           cmd.ReadSigningCertRetryMax,
			certInspector:      certutil.NewCertUtil(rootCertGracePeriodPercentile),
			caStorageNamespace: namespace,
			dualUse:            dualUse,
			org:                org,
			rootCertFile:       rootCertFile,
			enableJitter:       enableJitter,
			client:             client,
		},
	}
	if scrtErr != nil {
		pkiCaLog.Infof("Failed to get secret (error: %s), will create one", scrtErr)

		options := util.CertOptions{
			TTL:          caCertTTL,
			Org:          org,
			IsCA:         true,
			IsSelfSigned: true,
			RSAKeySize:   caRSAKeySize,
			IsDualUse:    dualUse,
		}
		pemCert, pemKey, ckErr := util.GenCertKeyFromOptions(options)
		if ckErr != nil {
			return nil, fmt.Errorf("unable to generate CA cert and key for self-signed CA (%v)", ckErr)
		}

		rootCerts, err := util.AppendRootCerts(pemCert, rootCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to append root certificates (%v)", err)
		}

		if caOpts.KeyCertBundle, err = util.NewVerifiedKeyCertBundleFromPem(pemCert, pemKey, nil, rootCerts); err != nil {
			return nil, fmt.Errorf("failed to create CA KeyCertBundle (%v)", err)
		}

		// Write the key/cert back to secret so they will be persistent when CA restarts.
		secret := k8ssecret.BuildSecret("", CASecret, namespace, nil, nil, nil, pemCert, pemKey, istioCASecretType)
		if _, err = client.Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{}); err != nil {
			pkiCaLog.Errorf("Failed to write secret to CA (error: %s). Abort.", err)
			return nil, fmt.Errorf("failed to create CA due to secret write error")
		}
		pkiCaLog.Infof("Using self-generated public key: %v", string(rootCerts))
	} else {
		pkiCaLog.Infof("Load signing key and cert from existing secret %s:%s", caSecret.Namespace, caSecret.Name)
		rootCerts, err := util.AppendRootCerts(caSecret.Data[CACertFile], rootCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to append root certificates (%v)", err)
		}
		if caOpts.KeyCertBundle, err = util.NewVerifiedKeyCertBundleFromPem(caSecret.Data[CACertFile],
			caSecret.Data[CAPrivateKeyFile], nil, rootCerts); err != nil {
			return nil, fmt.Errorf("failed to create CA KeyCertBundle (%v)", err)
		}
		pkiCaLog.Infof("Using existing public key: %v", string(rootCerts))
	}
	return caOpts, nil
}

// NewSelfSignedDebugIstioCAOptions returns a new IstioCAOptions instance using self-signed certificate produced by in-memory CA,
// which runs without K8s, and no local ca key file presented.
func NewSelfSignedDebugIstioCAOptions(rootCertFile string, caCertTTL, defaultCertTTL, maxCertTTL time.Duration,
	org string, caRSAKeySize int) (caOpts *IstioCAOptions, err error) {
	caOpts = &IstioCAOptions{
		CAType:         selfSignedCA,
		DefaultCertTTL: defaultCertTTL,
		MaxCertTTL:     maxCertTTL,
		CARSAKeySize:   caRSAKeySize,
	}

	options := util.CertOptions{
		TTL:          caCertTTL,
		Org:          org,
		IsCA:         true,
		IsSelfSigned: true,
		RSAKeySize:   caRSAKeySize,
		IsDualUse:    true, // hardcoded to true for K8S as well
	}
	pemCert, pemKey, ckErr := util.GenCertKeyFromOptions(options)
	if ckErr != nil {
		return nil, fmt.Errorf("unable to generate CA cert and key for self-signed CA (%v)", ckErr)
	}

	rootCerts, err := util.AppendRootCerts(pemCert, rootCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to append root certificates (%v)", err)
	}

	if caOpts.KeyCertBundle, err = util.NewVerifiedKeyCertBundleFromPem(pemCert, pemKey, nil, rootCerts); err != nil {
		return nil, fmt.Errorf("failed to create CA KeyCertBundle (%v)", err)
	}

	return caOpts, nil
}

// NewPluggedCertIstioCAOptions returns a new IstioCAOptions instance using given certificate.
func NewPluggedCertIstioCAOptions(fileBundle SigningCAFileBundle,
	defaultCertTTL, maxCertTTL time.Duration, caRSAKeySize int) (caOpts *IstioCAOptions, err error) {
	caOpts = &IstioCAOptions{
		CAType:         pluggedCertCA,
		DefaultCertTTL: defaultCertTTL,
		MaxCertTTL:     maxCertTTL,
		CARSAKeySize:   caRSAKeySize,
	}

	if caOpts.KeyCertBundle, err = util.NewVerifiedKeyCertBundleFromFile(
		fileBundle.SigningCertFile, fileBundle.SigningKeyFile, fileBundle.CertChainFiles, fileBundle.RootCertFile); err != nil {
		return nil, fmt.Errorf("failed to create CA KeyCertBundle (%v)", err)
	}

	// Validate that the passed in signing cert can be used as CA.
	// The check can't be done inside `KeyCertBundle`, since bundle could also be used to
	// validate workload certificates (i.e., where the leaf certificate is not a CA).
	b, err := os.ReadFile(fileBundle.SigningCertFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(b)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM encoded certificate")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse X.509 certificate")
	}
	if !cert.IsCA {
		return nil, fmt.Errorf("certificate is not authorized to sign other certificates")
	}

	return caOpts, nil
}

// IstioCA generates keys and certificates for Istio identities.
type IstioCA struct {
	defaultCertTTL time.Duration
	maxCertTTL     time.Duration
	caRSAKeySize   int

	keyCertBundle *util.KeyCertBundle

	// rootCertRotator periodically rotates self-signed root cert for CA. It is nil
	// if CA is not self-signed CA.
	rootCertRotator *SelfSignedCARootCertRotator
}

// NewIstioCA returns a new IstioCA instance.
func NewIstioCA(opts *IstioCAOptions) (*IstioCA, error) {
	ca := &IstioCA{
		maxCertTTL:    opts.MaxCertTTL,
		keyCertBundle: opts.KeyCertBundle,
		caRSAKeySize:  opts.CARSAKeySize,
	}

	if opts.CAType == selfSignedCA && opts.RotatorConfig != nil && opts.RotatorConfig.CheckInterval > time.Duration(0) {
		ca.rootCertRotator = NewSelfSignedCARootCertRotator(opts.RotatorConfig, ca)
	}

	// if CA cert becomes invalid before workload cert it's going to cause workload cert to be invalid too,
	// however citatel won't rotate if that happens, this function will prevent that using cert chain TTL as
	// the workload TTL
	defaultCertTTL, err := ca.minTTL(opts.DefaultCertTTL)
	if err != nil {
		return ca, fmt.Errorf("failed to get default cert TTL %s", err.Error())
	}
	ca.defaultCertTTL = defaultCertTTL

	return ca, nil
}

func (ca *IstioCA) Run(stopChan chan struct{}) {
	if ca.rootCertRotator != nil {
		// Start root cert rotator in a separate goroutine.
		go ca.rootCertRotator.Run(stopChan)
	}
}

// Sign takes a PEM-encoded CSR and cert opts, and returns a signed certificate.
func (ca *IstioCA) Sign(csrPEM []byte, certOpts CertOpts) (
	[]byte, error) {
	return ca.sign(csrPEM, certOpts.SubjectIDs, certOpts.TTL, true, certOpts.ForCA)
}

// SignWithCertChain is similar to Sign but returns the leaf cert and the entire cert chain.
func (ca *IstioCA) SignWithCertChain(csrPEM []byte, certOpts CertOpts) (
	[]string, error) {
	cert, err := ca.signWithCertChain(csrPEM, certOpts.SubjectIDs, certOpts.TTL, true, certOpts.ForCA)
	if err != nil {
		return nil, err
	}
	return []string{string(cert)}, nil
}

// GetCAKeyCertBundle returns the KeyCertBundle for the CA.
func (ca *IstioCA) GetCAKeyCertBundle() *util.KeyCertBundle {
	return ca.keyCertBundle
}

// GenKeyCert generates a certificate signed by the CA,
// returns the certificate chain and the private key.
func (ca *IstioCA) GenKeyCert(hostnames []string, certTTL time.Duration, checkLifetime bool) ([]byte, []byte, error) {
	opts := util.CertOptions{
		RSAKeySize: rsaKeySize,
	}

	// use the type of private key the CA uses to generate an intermediate CA of that type (e.g. CA cert using RSA will
	// cause intermediate CAs using RSA to be generated)
	_, signingKey, _, _ := ca.keyCertBundle.GetAll()
	if util.IsSupportedECPrivateKey(signingKey) {
		opts.ECSigAlg = util.EcdsaSigAlg
	}

	csrPEM, privPEM, err := util.GenCSR(opts)
	if err != nil {
		return nil, nil, err
	}

	certPEM, err := ca.signWithCertChain(csrPEM, hostnames, certTTL, checkLifetime, false)
	if err != nil {
		return nil, nil, err
	}

	return certPEM, privPEM, nil
}

func (ca *IstioCA) minTTL(defaultCertTTL time.Duration) (time.Duration, error) {
	certChainPem := ca.keyCertBundle.GetCertChainPem()
	if len(certChainPem) == 0 {
		return defaultCertTTL, nil
	}

	certChainExpiration, err := util.TimeBeforeCertExpires(certChainPem, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to get cert chain TTL %s", err.Error())
	}

	if certChainExpiration.Seconds() <= 0 {
		return 0, fmt.Errorf("cert chain has expired")
	}

	if defaultCertTTL.Seconds() > certChainExpiration.Seconds() {
		return certChainExpiration, nil
	}

	return defaultCertTTL, nil
}

func (ca *IstioCA) sign(csrPEM []byte, subjectIDs []string, requestedLifetime time.Duration, checkLifetime, forCA bool) ([]byte, error) {
	signingCert, signingKey, _, _ := ca.keyCertBundle.GetAll()
	if signingCert == nil {
		return nil, caerror.NewError(caerror.CANotReady, fmt.Errorf("Istio CA is not ready")) // nolint
	}

	csr, err := util.ParsePemEncodedCSR(csrPEM)
	if err != nil {
		return nil, caerror.NewError(caerror.CSRError, err)
	}

	lifetime := requestedLifetime
	// If the requested requestedLifetime is non-positive, apply the default TTL.
	if requestedLifetime.Seconds() <= 0 {
		lifetime = ca.defaultCertTTL
	}
	// If checkLifetime is set and the requested TTL is greater than maxCertTTL, return an error
	if checkLifetime && requestedLifetime.Seconds() > ca.maxCertTTL.Seconds() {
		return nil, caerror.NewError(caerror.TTLError, fmt.Errorf(
			"requested TTL %s is greater than the max allowed TTL %s", requestedLifetime, ca.maxCertTTL))
	}

	certBytes, err := util.GenCertFromCSR(csr, signingCert, csr.PublicKey, *signingKey, subjectIDs, lifetime, forCA)
	if err != nil {
		return nil, caerror.NewError(caerror.CertGenError, err)
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	cert := pem.EncodeToMemory(block)

	return cert, nil
}

func (ca *IstioCA) signWithCertChain(csrPEM []byte, subjectIDs []string, requestedLifetime time.Duration, lifetimeCheck,
	forCA bool) ([]byte, error) {
	cert, err := ca.sign(csrPEM, subjectIDs, requestedLifetime, lifetimeCheck, forCA)
	if err != nil {
		return nil, err
	}

	chainPem := ca.GetCAKeyCertBundle().GetCertChainPem()
	if len(chainPem) > 0 {
		cert = append(cert, chainPem...)
	}
	return cert, nil
}
