package collector

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestApiHealth(t *testing.T) {

	tcs := map[string]string{
		"1.7.6": ` {"share":{"total":1000,"successful":998,"failed":2},"api_stats":[{"api_name":"Api_1","api_requests":10,"api_requests_latency":1},{"api_name":"Api_12","api_requests":20,"api_requests_latency":2}]}`,
		"2.4.5": ` {"share":{"total":2000,"successful":1998,"failed":2},"api_stats":[{"api_name":"Api_1","api_requests":20,"api_requests_latency":2},{"api_name":"Api_12","api_requests":20,"api_requests_latency":3}]}`,
	}
	for ver, out := range tcs {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, out)
		}))
		defer ts.Close()

		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("Failed to parse URL: %s", err)
		}

		buf := &bytes.Buffer{}
		logger := log.New(buf, "", log.Lshortfile|log.LstdFlags)

		c := NewApiHealth(*logger, http.DefaultClient, u)
		chr, err := c.fetchAndDecodeApiStats()
		if err != nil {
			t.Fatalf("Failed to fetch or decode api health: %s", err)
		}
		t.Logf("[%s] Api Health Response: %+v", ver, chr)
		if chr.ApiStats[0].ApiName == "Api_1" {
			fmt.Println("ok")
		}

	}
}
