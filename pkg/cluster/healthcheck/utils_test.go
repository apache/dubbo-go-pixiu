package healthcheck

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func TestTcpConn(t *testing.T) {

	go func() {
		listener, err := net.Listen("tcp", ":12347")
		if err != nil {
			log.Fatalf("listen failed: %v", err)
		}

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("accept new client failed: %v", err)
				continue
			}

			client.Close()
		}
	}()

	time.Sleep(100 * time.Millisecond) // Give the server some time to start

	type TestCase struct {
		name         string
		addr         string
		port         int
		timeout      time.Duration
		shouldListen bool
		expected     bool
	}

	testCases := []TestCase{
		{
			name:     "12347",
			addr:     "127.0.0.1",
			port:     12347,
			timeout:  100 * time.Millisecond,
			expected: true,
		},
		{
			name:     "12348",
			addr:     "127.0.0.1",
			port:     12348,
			timeout:  100 * time.Millisecond,
			expected: false,
		},
	}

	var wg sync.WaitGroup

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := TcpConn(tc.addr, tc.port, tc.timeout)
			if actual != tc.expected {
				t.Fail()
			}
		})
	}

	wg.Wait()
}
