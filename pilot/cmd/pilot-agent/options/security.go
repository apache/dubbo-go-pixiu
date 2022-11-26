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

package options

import (
	"fmt"
	"strings"
)

import (
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/pkg/log"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	securityModel "github.com/apache/dubbo-go-pixiu/pilot/pkg/security/model"
	"github.com/apache/dubbo-go-pixiu/pkg/config/constants"
	"github.com/apache/dubbo-go-pixiu/pkg/jwt"
	"github.com/apache/dubbo-go-pixiu/pkg/security"
	"github.com/apache/dubbo-go-pixiu/security/pkg/credentialfetcher"
	"github.com/apache/dubbo-go-pixiu/security/pkg/nodeagent/cafile"
	"github.com/apache/dubbo-go-pixiu/security/pkg/nodeagent/plugin/providers/google/stsclient"
	"github.com/apache/dubbo-go-pixiu/security/pkg/stsservice/tokenmanager"
)

func NewSecurityOptions(proxyConfig *meshconfig.ProxyConfig, stsPort int, tokenManagerPlugin string) (*security.Options, error) {
	o := &security.Options{
		CAEndpoint:                     caEndpointEnv,
		CAProviderName:                 caProviderEnv,
		PilotCertProvider:              features.PilotCertProvider,
		OutputKeyCertToDir:             outputKeyCertToDir,
		ProvCert:                       provCert,
		ClusterID:                      clusterIDVar.Get(),
		FileMountedCerts:               fileMountedCertsEnv,
		WorkloadNamespace:              PodNamespaceVar.Get(),
		ServiceAccount:                 serviceAccountVar.Get(),
		XdsAuthProvider:                xdsAuthProvider.Get(),
		TrustDomain:                    trustDomainEnv,
		Pkcs8Keys:                      pkcs8KeysEnv,
		ECCSigAlg:                      eccSigAlgEnv,
		SecretTTL:                      secretTTLEnv,
		FileDebounceDuration:           fileDebounceDuration,
		SecretRotationGracePeriodRatio: secretRotationGracePeriodRatioEnv,
		STSPort:                        stsPort,
		CertSigner:                     certSigner.Get(),
		CARootPath:                     cafile.CACertFilePath,
		CertChainFilePath:              security.DefaultCertChainFilePath,
		KeyFilePath:                    security.DefaultKeyFilePath,
		RootCertFilePath:               security.DefaultRootCertFilePath,
	}

	o, err := SetupSecurityOptions(proxyConfig, o, jwtPolicy.Get(),
		credFetcherTypeEnv, credIdentityProvider)
	if err != nil {
		return o, err
	}

	var tokenManager security.TokenManager
	if stsPort > 0 || xdsAuthProvider.Get() != "" {
		// tokenManager is gcp token manager when using the default token manager plugin.
		tokenManager, err = tokenmanager.CreateTokenManager(tokenManagerPlugin,
			tokenmanager.Config{CredFetcher: o.CredFetcher, TrustDomain: o.TrustDomain})
	}
	o.TokenManager = tokenManager

	return o, err
}

func SetupSecurityOptions(proxyConfig *meshconfig.ProxyConfig, secOpt *security.Options, jwtPolicy,
	credFetcherTypeEnv, credIdentityProvider string) (*security.Options, error) {
	var jwtPath string
	switch jwtPolicy {
	case jwt.PolicyThirdParty:
		log.Info("JWT policy is third-party-jwt")
		jwtPath = constants.TrustworthyJWTPath
	case jwt.PolicyFirstParty:
		log.Info("JWT policy is first-party-jwt")
		jwtPath = securityModel.K8sSAJwtFileName
	default:
		log.Info("Using existing certs")
	}

	o := secOpt

	// If not set explicitly, default to the discovery address.
	if o.CAEndpoint == "" {
		o.CAEndpoint = proxyConfig.DiscoveryAddress
		o.CAEndpointSAN = istiodSAN.Get()
	}

	o.CredIdentityProvider = credIdentityProvider
	credFetcher, err := credentialfetcher.NewCredFetcher(credFetcherTypeEnv, o.TrustDomain, jwtPath, o.CredIdentityProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential fetcher: %v", err)
	}
	log.Infof("using credential fetcher of %s type in %s trust domain", credFetcherTypeEnv, o.TrustDomain)
	o.CredFetcher = credFetcher

	if o.CAProviderName == security.GkeWorkloadCertificateProvider {
		if !security.CheckWorkloadCertificate(security.GkeWorkloadCertChainFilePath,
			security.GkeWorkloadKeyFilePath, security.GkeWorkloadRootCertFilePath) {
			return nil, fmt.Errorf("GKE workload certificate files (%v, %v, %v) not present",
				security.GkeWorkloadCertChainFilePath, security.GkeWorkloadKeyFilePath, security.GkeWorkloadRootCertFilePath)
		}
		if o.ProvCert != "" {
			return nil, fmt.Errorf(
				"invalid options: PROV_CERT and FILE_MOUNTED_CERTS of GKE workload cert are mutually exclusive")
		}
		o.FileMountedCerts = true
		o.CertChainFilePath = security.GkeWorkloadCertChainFilePath
		o.KeyFilePath = security.GkeWorkloadKeyFilePath
		o.RootCertFilePath = security.GkeWorkloadRootCertFilePath
		return o, nil
	}

	// Default the CA provider where possible
	if strings.Contains(o.CAEndpoint, "googleapis.com") {
		o.CAProviderName = security.GoogleCAProvider
	}
	// TODO extract this logic out to a plugin
	if o.CAProviderName == security.GoogleCAProvider || o.CAProviderName == security.GoogleCASProvider {
		var err error
		o.TokenExchanger, err = stsclient.NewSecureTokenServiceExchanger(o.CredFetcher, o.TrustDomain)
		if err != nil {
			return nil, err
		}
	}

	if o.ProvCert != "" && o.FileMountedCerts {
		return nil, fmt.Errorf("invalid options: PROV_CERT and FILE_MOUNTED_CERTS are mutually exclusive")
	}
	return o, nil
}
