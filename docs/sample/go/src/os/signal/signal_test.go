// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package signal

import (
	"bytes"
	"flag"
	"fmt"
	"internal/testenv"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"testing"
	"time"
)

var testDeadline time.Time

func TestMain(m *testing.M) {
	flag.Parse()

	// TODO(golang.org/issue/28135): Remove this setup and use t.Deadline instead.
	timeoutFlag := flag.Lookup("test.timeout")
	if timeoutFlag != nil {
		if d := timeoutFlag.Value.(flag.Getter).Get().(time.Duration); d != 0 {
			testDeadline = time.Now().Add(d)
		}
	}

	os.Exit(m.Run())
}

func waitSig(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	waitSig1(t, c, sig, false)
}
func waitSigAll(t *testing.T, c <-chan os.Signal, sig os.Signal) {
	waitSig1(t, c, sig, true)
}

func waitSig1(t *testing.T, c <-chan os.Signal, sig os.Signal, all bool) {
	// Sleep multiple times to give the kernel more tries to
	// deliver the signal.
	for i := 0; i < 10; i++ {
		select {
		case s := <-c:
			// If the caller notified for all signals on
			// c, filter out SIGURG, which is used for
			// runtime preemption and can come at
			// unpredictable times.
			if all && s == syscall.SIGURG {
				continue
			}
			if s != sig {
				t.Fatalf("signal was %v, want %v", s, sig)
			}
			return

		case <-time.After(100 * time.Millisecond):
		}
	}
	t.Fatalf("timeout waiting for %v", sig)
}

// Test that basic signal handling works.
func TestSignal(t *testing.T) {
	// Ask for SIGHUP
	c := make(chan os.Signal, 1)
	Notify(c, syscall.SIGHUP)
	defer Stop(c)

	// Send this process a SIGHUP
	t.Logf("sighup...")
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSig(t, c, syscall.SIGHUP)

	// Ask for everything we can get. The buffer size has to be
	// more than 1, since the runtime might send SIGURG signals.
	// Using 10 is arbitrary.
	c1 := make(chan os.Signal, 10)
	Notify(c1)

	// Send this process a SIGWINCH
	t.Logf("sigwinch...")
	syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
	waitSigAll(t, c1, syscall.SIGWINCH)

	// Send two more SIGHUPs, to make sure that
	// they get delivered on c1 and that not reading
	// from c does not block everything.
	t.Logf("sighup...")
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSigAll(t, c1, syscall.SIGHUP)
	t.Logf("sighup...")
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSigAll(t, c1, syscall.SIGHUP)

	// The first SIGHUP should be waiting for us on c.
	waitSig(t, c, syscall.SIGHUP)
}

func TestStress(t *testing.T) {
	dur := 3 * time.Second
	if testing.Short() {
		dur = 100 * time.Millisecond
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	done := make(chan bool)
	finished := make(chan bool)
	go func() {
		sig := make(chan os.Signal, 1)
		Notify(sig, syscall.SIGUSR1)
		defer Stop(sig)
	Loop:
		for {
			select {
			case <-sig:
			case <-done:
				break Loop
			}
		}
		finished <- true
	}()
	go func() {
	Loop:
		for {
			select {
			case <-done:
				break Loop
			default:
				syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
				runtime.Gosched()
			}
		}
		finished <- true
	}()
	time.Sleep(dur)
	close(done)
	<-finished
	<-finished
	// When run with 'go test -cpu=1,2,4' SIGUSR1 from this test can slip
	// into subsequent TestSignal() causing failure.
	// Sleep for a while to reduce the possibility of the failure.
	time.Sleep(10 * time.Millisecond)
}

func testCancel(t *testing.T, ignore bool) {
	// Send SIGWINCH. By default this signal should be ignored.
	syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
	time.Sleep(100 * time.Millisecond)

	// Ask to be notified on c1 when a SIGWINCH is received.
	c1 := make(chan os.Signal, 1)
	Notify(c1, syscall.SIGWINCH)
	defer Stop(c1)

	// Ask to be notified on c2 when a SIGHUP is received.
	c2 := make(chan os.Signal, 1)
	Notify(c2, syscall.SIGHUP)
	defer Stop(c2)

	// Send this process a SIGWINCH and wait for notification on c1.
	syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)
	waitSig(t, c1, syscall.SIGWINCH)

	// Send this process a SIGHUP and wait for notification on c2.
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	waitSig(t, c2, syscall.SIGHUP)

	// Ignore, or reset the signal handlers for, SIGWINCH and SIGHUP.
	if ignore {
		Ignore(syscall.SIGWINCH, syscall.SIGHUP)
	} else {
		Reset(syscall.SIGWINCH, syscall.SIGHUP)
	}

	// At this point we do not expect any further signals on c1.
	// However, it is just barely possible that the initial SIGWINCH
	// at the start of this function was delivered after we called
	// Notify on c1. In that case the waitSig for SIGWINCH may have
	// picked up that initial SIGWINCH, and the second SIGWINCH may
	// then have been delivered on the channel. This sequence of events
	// may have caused issue 15661.
	// So, read any possible signal from the channel now.
	select {
	case <-c1:
	default:
	}

	// Send this process a SIGWINCH. It should be ignored.
	syscall.Kill(syscall.Getpid(), syscall.SIGWINCH)

	// If ignoring, Send this process a SIGHUP. It should be ignored.
	if ignore {
		syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
	}

	select {
	case s := <-c1:
		t.Fatalf("unexpected signal %v", s)
	case <-time.After(100 * time.Millisecond):
		// nothing to read - good
	}

	select {
	case s := <-c2:
		t.Fatalf("unexpected signal %v", s)
	case <-time.After(100 * time.Millisecond):
		// nothing to read - good
	}

	// Reset the signal handlers for all signals.
	Reset()
}

// Test that Reset cancels registration for listed signals on all channels.
func TestReset(t *testing.T) {
	testCancel(t, false)
}

// Test that Ignore cancels registration for listed signals on all channels.
func TestIgnore(t *testing.T) {
	testCancel(t, true)
}

// Test that Ignored correctly detects changes to the ignored status of a signal.
func TestIgnored(t *testing.T) {
	// Ask to be notified on SIGWINCH.
	c := make(chan os.Signal, 1)
	Notify(c, syscall.SIGWINCH)

	// If we're being notified, then the signal should not be ignored.
	if Ignored(syscall.SIGWINCH) {
		t.Errorf("expected SIGWINCH to not be ignored.")
	}
	Stop(c)
	Ignore(syscall.SIGWINCH)

	// We're no longer paying attention to this signal.
	if !Ignored(syscall.SIGWINCH) {
		t.Errorf("expected SIGWINCH to be ignored when explicitly ignoring it.")
	}

	Reset()
}

var checkSighupIgnored = flag.Bool("check_sighup_ignored", false, "if true, TestDetectNohup will fail if SIGHUP is not ignored.")

// Test that Ignored(SIGHUP) correctly detects whether it is being run under nohup.
func TestDetectNohup(t *testing.T) {
	if *checkSighupIgnored {
		if !Ignored(syscall.SIGHUP) {
			t.Fatal("SIGHUP is not ignored.")
		} else {
			t.Log("SIGHUP is ignored.")
		}
	} else {
		defer Reset()
		// Ugly: ask for SIGHUP so that child will not have no-hup set
		// even if test is running under nohup environment.
		// We have no intention of reading from c.
		c := make(chan os.Signal, 1)
		Notify(c, syscall.SIGHUP)
		if out, err := exec.Command(os.Args[0], "-test.run=TestDetectNohup", "-check_sighup_ignored").CombinedOutput(); err == nil {
			t.Errorf("ran test with -check_sighup_ignored and it succeeded: expected failure.\nOutput:\n%s", out)
		}
		Stop(c)
		// Again, this time with nohup, assuming we can find it.
		_, err := os.Stat("/usr/bin/nohup")
		if err != nil {
			t.Skip("cannot find nohup; skipping second half of test")
		}
		Ignore(syscall.SIGHUP)
		os.Remove("nohup.out")
		out, err := exec.Command("/usr/bin/nohup", os.Args[0], "-test.run=TestDetectNohup", "-check_sighup_ignored").CombinedOutput()

		data, _ := ioutil.ReadFile("nohup.out")
		os.Remove("nohup.out")
		if err != nil {
			t.Errorf("ran test with -check_sighup_ignored under nohup and it failed: expected success.\nError: %v\nOutput:\n%s%s", err, out, data)
		}
	}
}

var sendUncaughtSighup = flag.Int("send_uncaught_sighup", 0, "send uncaught SIGHUP during TestStop")

// Test that Stop cancels the channel's registrations.
func TestStop(t *testing.T) {
	sigs := []syscall.Signal{
		syscall.SIGWINCH,
		syscall.SIGHUP,
		syscall.SIGUSR1,
	}

	for _, sig := range sigs {
		// Send the signal.
		// If it's SIGWINCH, we should not see it.
		// If it's SIGHUP, maybe we'll die. Let the flag tell us what to do.
		if sig == syscall.SIGWINCH || (sig == syscall.SIGHUP && *sendUncaughtSighup == 1) {
			syscall.Kill(syscall.Getpid(), sig)
		}

		// The kernel will deliver a signal as a thread returns
		// from a syscall. If the only active thread is sleeping,
		// and the system is busy, the kernel may not get around
		// to waking up a thread to catch the signal.
		// We try splitting up the sleep to give the kernel
		// another chance to deliver the signal.
		time.Sleep(50 * time.Millisecond)
		time.Sleep(50 * time.Millisecond)

		// Ask for signal
		c := make(chan os.Signal, 1)
		Notify(c, sig)
		defer Stop(c)

		// Send this process that signal
		syscall.Kill(syscall.Getpid(), sig)
		waitSig(t, c, sig)

		Stop(c)
		time.Sleep(50 * time.Millisecond)
		select {
		case s := <-c:
			t.Fatalf("unexpected signal %v", s)
		case <-time.After(50 * time.Millisecond):
			// nothing to read - good
		}

		// Send the signal.
		// If it's SIGWINCH, we should not see it.
		// If it's SIGHUP, maybe we'll die. Let the flag tell us what to do.
		if sig != syscall.SIGHUP || *sendUncaughtSighup == 2 {
			syscall.Kill(syscall.Getpid(), sig)
		}

		time.Sleep(50 * time.Millisecond)
		select {
		case s := <-c:
			t.Fatalf("unexpected signal %v", s)
		case <-time.After(50 * time.Millisecond):
			// nothing to read - good
		}
	}
}

// Test that when run under nohup, an uncaught SIGHUP does not kill the program,
// but a
func TestNohup(t *testing.T) {
	// Ugly: ask for SIGHUP so that child will not have no-hup set
	// even if test is running under nohup environment.
	// We have no intention of reading from c.
	c := make(chan os.Signal, 1)
	Notify(c, syscall.SIGHUP)

	// When run without nohup, the test should crash on an uncaught SIGHUP.
	// When run under nohup, the test should ignore uncaught SIGHUPs,
	// because the runtime is not supposed to be listening for them.
	// Either way, TestStop should still be able to catch them when it wants them
	// and then when it stops wanting them, the original behavior should resume.
	//
	// send_uncaught_sighup=1 sends the SIGHUP before starting to listen for SIGHUPs.
	// send_uncaught_sighup=2 sends the SIGHUP after no longer listening for SIGHUPs.
	//
	// Both should fail without nohup and succeed with nohup.

	for i := 1; i <= 2; i++ {
		out, err := exec.Command(os.Args[0], "-test.run=TestStop", "-send_uncaught_sighup="+strconv.Itoa(i)).CombinedOutput()
		if err == nil {
			t.Fatalf("ran test with -send_uncaught_sighup=%d and it succeeded: expected failure.\nOutput:\n%s", i, out)
		}
	}

	Stop(c)

	// Skip the nohup test below when running in tmux on darwin, since nohup
	// doesn't work correctly there. See issue #5135.
	if runtime.GOOS == "darwin" && os.Getenv("TMUX") != "" {
		t.Skip("Skipping nohup test due to running in tmux on darwin")
	}

	// Again, this time with nohup, assuming we can find it.
	_, err := os.Stat("/usr/bin/nohup")
	if err != nil {
		t.Skip("cannot find nohup; skipping second half of test")
	}

	for i := 1; i <= 2; i++ {
		os.Remove("nohup.out")
		out, err := exec.Command("/usr/bin/nohup", os.Args[0], "-test.run=TestStop", "-send_uncaught_sighup="+strconv.Itoa(i)).CombinedOutput()

		data, _ := ioutil.ReadFile("nohup.out")
		os.Remove("nohup.out")
		if err != nil {
			t.Fatalf("ran test with -send_uncaught_sighup=%d under nohup and it failed: expected success.\nError: %v\nOutput:\n%s%s", i, err, out, data)
		}
	}
}

// Test that SIGCONT works (issue 8953).
func TestSIGCONT(t *testing.T) {
	c := make(chan os.Signal, 1)
	Notify(c, syscall.SIGCONT)
	defer Stop(c)
	syscall.Kill(syscall.Getpid(), syscall.SIGCONT)
	waitSig(t, c, syscall.SIGCONT)
}

// Test race between stopping and receiving a signal (issue 14571).
func TestAtomicStop(t *testing.T) {
	if os.Getenv("GO_TEST_ATOMIC_STOP") != "" {
		atomicStopTestProgram()
		t.Fatal("atomicStopTestProgram returned")
	}

	testenv.MustHaveExec(t)

	// Call Notify for SIGINT before starting the child process.
	// That ensures that SIGINT is not ignored for the child.
	// This is necessary because if SIGINT is ignored when a
	// Go program starts, then it remains ignored, and closing
	// the last notification channel for SIGINT will switch it
	// back to being ignored. In that case the assumption of
	// atomicStopTestProgram, that it will either die from SIGINT
	// or have it be reported, breaks down, as there is a third
	// option: SIGINT might be ignored.
	cs := make(chan os.Signal, 1)
	Notify(cs, syscall.SIGINT)
	defer Stop(cs)

	const execs = 10
	for i := 0; i < execs; i++ {
		timeout := "0"
		if !testDeadline.IsZero() {
			timeout = testDeadline.Sub(time.Now()).String()
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestAtomicStop", "-test.timeout="+timeout)
		cmd.Env = append(os.Environ(), "GO_TEST_ATOMIC_STOP=1")
		out, err := cmd.CombinedOutput()
		if err == nil {
			if len(out) > 0 {
				t.Logf("iteration %d: output %s", i, out)
			}
		} else {
			t.Logf("iteration %d: exit status %q: output: %s", i, err, out)
		}

		lost := bytes.Contains(out, []byte("lost signal"))
		if lost {
			t.Errorf("iteration %d: lost signal", i)
		}

		// The program should either die due to SIGINT,
		// or exit with success without printing "lost signal".
		if err == nil {
			if len(out) > 0 && !lost {
				t.Errorf("iteration %d: unexpected output", i)
			}
		} else {
			if ee, ok := err.(*exec.ExitError); !ok {
				t.Errorf("iteration %d: error (%v) has type %T; expected exec.ExitError", i, err, err)
			} else if ws, ok := ee.Sys().(syscall.WaitStatus); !ok {
				t.Errorf("iteration %d: error.Sys (%v) has type %T; expected syscall.WaitStatus", i, ee.Sys(), ee.Sys())
			} else if !ws.Signaled() || ws.Signal() != syscall.SIGINT {
				t.Errorf("iteration %d: got exit status %v; expected SIGINT", i, ee)
			}
		}
	}
}

// atomicStopTestProgram is run in a subprocess by TestAtomicStop.
// It tries to trigger a signal delivery race. This function should
// either catch a signal or die from it.
func atomicStopTestProgram() {
	// This test won't work if SIGINT is ignored here.
	if Ignored(syscall.SIGINT) {
		fmt.Println("SIGINT is ignored")
		os.Exit(1)
	}

	const tries = 10

	timeout := 2 * time.Second
	if !testDeadline.IsZero() {
		// Give each try an equal slice of the deadline, with one slice to spare for
		// cleanup.
		timeout = testDeadline.Sub(time.Now()) / (tries + 1)
	}

	pid := syscall.Getpid()
	printed := false
	for i := 0; i < tries; i++ {
		cs := make(chan os.Signal, 1)
		Notify(cs, syscall.SIGINT)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			Stop(cs)
		}()

		syscall.Kill(pid, syscall.SIGINT)

		// At this point we should either die from SIGINT or
		// get a notification on cs. If neither happens, we
		// dropped the signal. It is given 2 seconds to
		// deliver, as needed for gccgo on some loaded test systems.

		select {
		case <-cs:
		case <-time.After(timeout):
			if !printed {
				fmt.Print("lost signal on tries:")
				printed = true
			}
			fmt.Printf(" %d", i)
		}

		wg.Wait()
	}
	if printed {
		fmt.Print("\n")
	}

	os.Exit(0)
}

func TestTime(t *testing.T) {
	// Test that signal works fine when we are in a call to get time,
	// which on some platforms is using VDSO. See issue #34391.
	dur := 3 * time.Second
	if testing.Short() {
		dur = 100 * time.Millisecond
	}
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	done := make(chan bool)
	finished := make(chan bool)
	go func() {
		sig := make(chan os.Signal, 1)
		Notify(sig, syscall.SIGUSR1)
		defer Stop(sig)
	Loop:
		for {
			select {
			case <-sig:
			case <-done:
				break Loop
			}
		}
		finished <- true
	}()
	go func() {
	Loop:
		for {
			select {
			case <-done:
				break Loop
			default:
				syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
				runtime.Gosched()
			}
		}
		finished <- true
	}()
	t0 := time.Now()
	for t1 := t0; t1.Sub(t0) < dur; t1 = time.Now() {
	} // hammering on getting time
	close(done)
	<-finished
	<-finished
	// When run with 'go test -cpu=1,2,4' SIGUSR1 from this test can slip
	// into subsequent TestSignal() causing failure.
	// Sleep for a while to reduce the possibility of the failure.
	time.Sleep(10 * time.Millisecond)
}
