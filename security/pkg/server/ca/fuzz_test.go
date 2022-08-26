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
	"fmt"
	"testing"

	"github.com/apache/dubbo-go-pixiu/pkg/fuzz"
	"github.com/apache/dubbo-go-pixiu/pkg/security"
	mockca "github.com/apache/dubbo-go-pixiu/security/pkg/pki/ca/mock"
	caerror "github.com/apache/dubbo-go-pixiu/security/pkg/pki/error"
	pb "istio.io/api/security/v1alpha1"
)

func FuzzCreateCertificate(f *testing.F) {
	fuzz.BaseCases(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		fg := fuzz.New(t, data)
		csr := fuzz.Struct[pb.IstioCertificateRequest](fg)
		ca := fuzz.Struct[mockca.FakeCA](fg)
		ca.SignErr = caerror.NewError(caerror.CSRError, fmt.Errorf("cannot sign"))
		server := &Server{
			ca:             &ca,
			Authenticators: []security.Authenticator{&mockAuthenticator{}},
			monitoring:     newMonitoringMetrics(),
		}
		_, _ = server.CreateCertificate(context.Background(), &csr)
	})
}