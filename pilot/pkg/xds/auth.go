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

package xds

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"istio.io/pkg/env"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/spiffe"
)

var AuthPlaintext = env.RegisterBoolVar("XDS_AUTH_PLAINTEXT", false,
	"Authenticate plain text requests - used if Istiod is behind a gateway handling TLS").Get()

// authenticate authenticates the ADS request using the configured authenticators.
// Returns the validated principals or an error.
// If no authenticators are configured, or if the request is on a non-secure
// stream ( 15010 ) - returns an empty list of principals and no errors.
func (s *DiscoveryServer) authenticate(ctx context.Context) ([]string, error) {
	if !features.XDSAuth {
		return nil, nil
	}

	// Authenticate - currently just checks that request has a certificate signed with the our key.
	// Protected by flag to avoid breaking upgrades - should be enabled in multi-cluster/meshexpansion where
	// XDS is exposed.
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("invalid context")
	}
	// Not a TLS connection, we will not perform authentication
	// TODO: add a flag to prevent unauthenticated requests ( 15010 )
	// request not over TLS on the insecure port
	if _, ok := peerInfo.AuthInfo.(credentials.TLSInfo); !ok && !AuthPlaintext {
		return nil, nil
	}
	authFailMsgs := []string{}
	for _, authn := range s.Authenticators {
		u, err := authn.Authenticate(ctx)
		// If one authenticator passes, return
		if u != nil && u.Identities != nil && err == nil {
			return u.Identities, nil
		}
		authFailMsgs = append(authFailMsgs, fmt.Sprintf("Authenticator %s: %v", authn.AuthenticatorType(), err))
	}

	log.Errorf("Failed to authenticate client from %s: %s", peerInfo.Addr.String(), strings.Join(authFailMsgs, "; "))
	return nil, errors.New("authentication failure")
}

func (s *DiscoveryServer) authorize(con *Connection, identities []string) error {
	if con == nil || con.proxy == nil {
		return nil
	}

	if features.EnableXDSIdentityCheck && identities != nil {
		// TODO: allow locking down, rejecting unauthenticated requests.
		id, err := checkConnectionIdentity(con.proxy, identities)
		if err != nil {
			log.Warnf("Unauthorized XDS: %v with identity %v: %v", con.peerAddr, identities, err)
			return status.Newf(codes.PermissionDenied, "authorization failed: %v", err).Err()
		}
		con.proxy.VerifiedIdentity = id
	}
	return nil
}

func checkConnectionIdentity(proxy *model.Proxy, identities []string) (*spiffe.Identity, error) {
	for _, rawID := range identities {
		spiffeID, err := spiffe.ParseIdentity(rawID)
		if err != nil {
			continue
		}
		if proxy.ConfigNamespace != "" && spiffeID.Namespace != proxy.ConfigNamespace {
			continue
		}
		if proxy.Metadata.ServiceAccount != "" && spiffeID.ServiceAccount != proxy.Metadata.ServiceAccount {
			continue
		}
		return &spiffeID, nil
	}
	return nil, fmt.Errorf("no identities (%v) matched %v/%v", identities, proxy.ConfigNamespace, proxy.Metadata.ServiceAccount)
}
