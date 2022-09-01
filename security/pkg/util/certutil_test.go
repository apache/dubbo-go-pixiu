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

package util

import (
	"os"
	"testing"
	"time"
)

const (
	// This cert has:
	//   NotBefore = 2017-08-23 19:00:40 +0000 UTC
	//   NotAfter  = 2017-08-24 19:00:40 +0000 UTC
	testCertFile = "testdata/cert-util.pem"
)

func TestGetWaitTime(t *testing.T) {
	testCert, err := os.ReadFile(testCertFile)
	if err != nil {
		t.Errorf("cannot read testing cert file")
		return
	}
	testCases := map[string]struct {
		cert             []byte
		now              time.Time
		expectedWaitTime int
		expectedErr      string
		minGracePeriod   time.Duration
	}{
		"Success": {
			// Now = 2017-08-23 21:00:40 +0000 UTC
			// Cert TTL is 24h, and grace period is 50% of TTL, that is 12h.
			// The cert expires at 2017-08-24 19:00:40 +0000 UTC, so the grace period starts at 2017-08-24 07:00:40 +0000 UTC
			// The wait time is the duration from fake now to the grace period start time, which is 10h39s (36039s).
			cert:             testCert,
			now:              time.Date(2017, time.August, 23, 21, 0, 0, 40, time.UTC),
			expectedWaitTime: 36039,
			minGracePeriod:   time.Duration(0),
		},
		"Success with min grace period": {
			// Now = 2017-08-23 21:00:40 +0000 UTC
			// Cert TTL is 24h, and grace period is 50% of TTL, that is 12h.
			// The cert expires at 2017-08-24 19:00:40 +0000 UTC, so the grace period starts at 2017-08-24 07:00:40 +0000 UTC
			// The wait time is the duration from fake now to the grace period start time, which is 10h39s (36039s).
			cert:             testCert,
			now:              time.Date(2017, time.August, 23, 21, 0, 0, 40, time.UTC),
			expectedWaitTime: 32439,
			minGracePeriod:   13 * time.Hour,
		},
		"Cert expired": {
			// Now = 2017-08-25 21:00:40 +0000 UTC.
			// Now is later than cert's NotAfter 2017-08-24 19:00:40 +0000 UTC.
			cert: testCert,
			now:  time.Date(2017, time.August, 25, 21, 0, 0, 40, time.UTC),
			expectedErr: "certificate already expired at 2017-08-24 19:00:40 +0000" +
				" UTC, but now is 2017-08-25 21:00:00.00000004 +0000 UTC",
			minGracePeriod: time.Duration(0),
		},
		"Renew now": {
			// Now = 2017-08-24 16:00:40 +0000 UTC
			// Now is later than the start of grace period 2017-08-24 07:00:40 +0000 UTC, but earlier than
			// cert expiration 2017-08-24 19:00:40 +0000 UTC.
			cert:           testCert,
			now:            time.Date(2017, time.August, 24, 16, 0, 0, 40, time.UTC),
			expectedErr:    "got a certificate that should be renewed now",
			minGracePeriod: time.Duration(0),
		},
		"Renew now with min grace period": {
			// Now = 2017-08-24 6:00:40 +0000 UTC
			// Now is earlier than the start of default grace period 2017-08-24 07:00:40 +0000 UTC, but with min grace period,
			// the cert is in grace period after 2017-08-24 5:00:40 +0000 UTC.
			cert:           testCert,
			now:            time.Date(2017, time.August, 24, 6, 0, 40, 0, time.UTC),
			expectedErr:    "got a certificate that should be renewed now",
			minGracePeriod: 14 * time.Hour,
		},
		"Invalid cert pem": {
			cert:           []byte(`INVALIDCERT`),
			now:            time.Date(2017, time.August, 23, 21, 0, 0, 40, time.UTC),
			expectedErr:    "invalid PEM encoded certificate",
			minGracePeriod: time.Duration(0),
		},
	}

	cu := NewCertUtil(50) // Grace period percentage is set to 50
	for id, c := range testCases {
		waitTime, err := cu.GetWaitTime(c.cert, c.now, c.minGracePeriod)
		if c.expectedErr != "" {
			if err == nil {
				t.Errorf("%s: no error is returned.", id)
			}
			if err.Error() != c.expectedErr {
				t.Errorf("%s: incorrect error message: %s VS %s", id, err.Error(), c.expectedErr)
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error: %v", id, err)
			}
			if int(waitTime.Seconds()) != c.expectedWaitTime {
				t.Errorf("%s: incorrect waittime. Expected %ds, but got %ds.", id, c.expectedWaitTime, int(waitTime.Seconds()))
			}
		}
	}
}
