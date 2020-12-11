// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program can be used as go_darwin_arm_exec by the Go tool.
// It executes binaries on an iOS device using the XCode toolchain
// and the ios-deploy program: https://github.com/phonegap/ios-deploy
//
// This script supports an extra flag, -lldb, that pauses execution
// just before the main program begins and allows the user to control
// the remote lldb session. This flag is appended to the end of the
// script's arguments and is not passed through to the underlying
// binary.
//
// This script requires that three environment variables be set:
// 	GOIOS_DEV_ID: The codesigning developer id or certificate identifier
// 	GOIOS_APP_ID: The provisioning app id prefix. Must support wildcard app ids.
// 	GOIOS_TEAM_ID: The team id that owns the app id prefix.
// $GOROOT/misc/ios contains a script, detect.go, that attempts to autodetect these.
package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const debug = false

var tmpdir string

var (
	devID    string
	appID    string
	teamID   string
	bundleID string
	deviceID string
)

// lock is a file lock to serialize iOS runs. It is global to avoid the
// garbage collector finalizing it, closing the file and releasing the
// lock prematurely.
var lock *os.File

func main() {
	log.SetFlags(0)
	log.SetPrefix("go_darwin_arm_exec: ")
	if debug {
		log.Println(strings.Join(os.Args, " "))
	}
	if len(os.Args) < 2 {
		log.Fatal("usage: go_darwin_arm_exec a.out")
	}

	// e.g. B393DDEB490947F5A463FD074299B6C0AXXXXXXX
	devID = getenv("GOIOS_DEV_ID")

	// e.g. Z8B3JBXXXX.org.golang.sample, Z8B3JBXXXX prefix is available at
	// https://developer.apple.com/membercenter/index.action#accountSummary as Team ID.
	appID = getenv("GOIOS_APP_ID")

	// e.g. Z8B3JBXXXX, available at
	// https://developer.apple.com/membercenter/index.action#accountSummary as Team ID.
	teamID = getenv("GOIOS_TEAM_ID")

	// Device IDs as listed with ios-deploy -c.
	deviceID = os.Getenv("GOIOS_DEVICE_ID")

	parts := strings.SplitN(appID, ".", 2)
	// For compatibility with the old builders, use a fallback bundle ID
	bundleID = "golang.gotest"
	if len(parts) == 2 {
		bundleID = parts[1]
	}

	exitCode, err := runMain()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	os.Exit(exitCode)
}

func runMain() (int, error) {
	var err error
	tmpdir, err = ioutil.TempDir("", "go_darwin_arm_exec_")
	if err != nil {
		return 1, err
	}
	if !debug {
		defer os.RemoveAll(tmpdir)
	}

	appdir := filepath.Join(tmpdir, "gotest.app")
	os.RemoveAll(appdir)

	if err := assembleApp(appdir, os.Args[1]); err != nil {
		return 1, err
	}

	// This wrapper uses complicated machinery to run iOS binaries. It
	// works, but only when running one binary at a time.
	// Use a file lock to make sure only one wrapper is running at a time.
	//
	// The lock file is never deleted, to avoid concurrent locks on distinct
	// files with the same path.
	lockName := filepath.Join(os.TempDir(), "go_darwin_arm_exec-"+deviceID+".lock")
	lock, err = os.OpenFile(lockName, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return 1, err
	}
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX); err != nil {
		return 1, err
	}

	if err := uninstall(bundleID); err != nil {
		return 1, err
	}

	if err := install(appdir); err != nil {
		return 1, err
	}

	if err := mountDevImage(); err != nil {
		return 1, err
	}

	// Kill any hanging debug bridges that might take up port 3222.
	exec.Command("killall", "idevicedebugserverproxy").Run()

	closer, err := startDebugBridge()
	if err != nil {
		return 1, err
	}
	defer closer()

	if err := run(appdir, bundleID, os.Args[2:]); err != nil {
		// If the lldb driver completed with an exit code, use that.
		if err, ok := err.(*exec.ExitError); ok {
			if ws, ok := err.Sys().(interface{ ExitStatus() int }); ok {
				return ws.ExitStatus(), nil
			}
		}
		return 1, err
	}
	return 0, nil
}

func getenv(envvar string) string {
	s := os.Getenv(envvar)
	if s == "" {
		log.Fatalf("%s not set\nrun $GOROOT/misc/ios/detect.go to attempt to autodetect", envvar)
	}
	return s
}

func assembleApp(appdir, bin string) error {
	if err := os.MkdirAll(appdir, 0755); err != nil {
		return err
	}

	if err := cp(filepath.Join(appdir, "gotest"), bin); err != nil {
		return err
	}

	pkgpath, err := copyLocalData(appdir)
	if err != nil {
		return err
	}

	entitlementsPath := filepath.Join(tmpdir, "Entitlements.plist")
	if err := ioutil.WriteFile(entitlementsPath, []byte(entitlementsPlist()), 0744); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(appdir, "Info.plist"), []byte(infoPlist(pkgpath)), 0744); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(appdir, "ResourceRules.plist"), []byte(resourceRules), 0744); err != nil {
		return err
	}

	cmd := exec.Command(
		"codesign",
		"-f",
		"-s", devID,
		"--entitlements", entitlementsPath,
		appdir,
	)
	if debug {
		log.Println(strings.Join(cmd.Args, " "))
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("codesign: %v", err)
	}
	return nil
}

// mountDevImage ensures a developer image is mounted on the device.
// The image contains the device lldb server for idevicedebugserverproxy
// to connect to.
func mountDevImage() error {
	// Check for existing mount.
	cmd := idevCmd(exec.Command("ideviceimagemounter", "-l", "-x"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("ideviceimagemounter: %v", err)
	}
	var info struct {
		Dict struct {
			Data []byte `xml:",innerxml"`
		} `xml:"dict"`
	}
	if err := xml.Unmarshal(out, &info); err != nil {
		return fmt.Errorf("mountDevImage: failed to decode mount information: %v", err)
	}
	dict, err := parsePlistDict(info.Dict.Data)
	if err != nil {
		return fmt.Errorf("mountDevImage: failed to parse mount information: %v", err)
	}
	if dict["ImagePresent"] == "true" && dict["Status"] == "Complete" {
		return nil
	}
	// Some devices only give us an ImageSignature key.
	if _, exists := dict["ImageSignature"]; exists {
		return nil
	}
	// No image is mounted. Find a suitable image.
	imgPath, err := findDevImage()
	if err != nil {
		return err
	}
	sigPath := imgPath + ".signature"
	cmd = idevCmd(exec.Command("ideviceimagemounter", imgPath, sigPath))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("ideviceimagemounter: %v", err)
	}
	return nil
}

// findDevImage use the device iOS version and build to locate a suitable
// developer image.
func findDevImage() (string, error) {
	cmd := idevCmd(exec.Command("ideviceinfo"))
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ideviceinfo: %v", err)
	}
	var iosVer, buildVer string
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		spl := bytes.SplitN(line, []byte(": "), 2)
		if len(spl) != 2 {
			continue
		}
		key, val := string(spl[0]), string(spl[1])
		switch key {
		case "ProductVersion":
			iosVer = val
		case "BuildVersion":
			buildVer = val
		}
	}
	if iosVer == "" || buildVer == "" {
		return "", errors.New("failed to parse ideviceinfo output")
	}
	verSplit := strings.Split(iosVer, ".")
	if len(verSplit) > 2 {
		// Developer images are specific to major.minor ios version.
		// Cut off the patch version.
		iosVer = strings.Join(verSplit[:2], ".")
	}
	sdkBase := "/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/DeviceSupport"
	patterns := []string{fmt.Sprintf("%s (%s)", iosVer, buildVer), fmt.Sprintf("%s (*)", iosVer), fmt.Sprintf("%s*", iosVer)}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(sdkBase, pattern, "DeveloperDiskImage.dmg"))
		if err != nil {
			return "", fmt.Errorf("findDevImage: %v", err)
		}
		if len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", fmt.Errorf("failed to find matching developer image for iOS version %s build %s", iosVer, buildVer)
}

// startDebugBridge ensures that the idevicedebugserverproxy runs on
// port 3222.
func startDebugBridge() (func(), error) {
	errChan := make(chan error, 1)
	cmd := idevCmd(exec.Command("idevicedebugserverproxy", "3222"))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("idevicedebugserverproxy: %v", err)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				errChan <- fmt.Errorf("idevicedebugserverproxy: %s", stderr.Bytes())
			} else {
				errChan <- fmt.Errorf("idevicedebugserverproxy: %v", err)
			}
		}
		errChan <- nil
	}()
	closer := func() {
		cmd.Process.Kill()
		<-errChan
	}
	// Dial localhost:3222 to ensure the proxy is ready.
	delay := time.Second / 4
	for attempt := 0; attempt < 5; attempt++ {
		conn, err := net.DialTimeout("tcp", "localhost:3222", 5*time.Second)
		if err == nil {
			conn.Close()
			return closer, nil
		}
		select {
		case <-time.After(delay):
			delay *= 2
		case err := <-errChan:
			return nil, err
		}
	}
	closer()
	return nil, errors.New("failed to set up idevicedebugserverproxy")
}

// findDeviceAppPath returns the device path to the app with the
// given bundle ID. It parses the output of ideviceinstaller -l -o xml,
// looking for the bundle ID and the corresponding Path value.
func findDeviceAppPath(bundleID string) (string, error) {
	cmd := idevCmd(exec.Command("ideviceinstaller", "-l", "-o", "xml"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return "", fmt.Errorf("ideviceinstaller: -l -o xml %v", err)
	}
	var list struct {
		Apps []struct {
			Data []byte `xml:",innerxml"`
		} `xml:"array>dict"`
	}
	if err := xml.Unmarshal(out, &list); err != nil {
		return "", fmt.Errorf("failed to parse ideviceinstaller output: %v", err)
	}
	for _, app := range list.Apps {
		values, err := parsePlistDict(app.Data)
		if err != nil {
			return "", fmt.Errorf("findDeviceAppPath: failed to parse app dict: %v", err)
		}
		if values["CFBundleIdentifier"] == bundleID {
			if path, ok := values["Path"]; ok {
				return path, nil
			}
		}
	}
	return "", fmt.Errorf("failed to find device path for bundle: %s", bundleID)
}

// Parse an xml encoded plist. Plist values are mapped to string.
func parsePlistDict(dict []byte) (map[string]string, error) {
	d := xml.NewDecoder(bytes.NewReader(dict))
	values := make(map[string]string)
	var key string
	var hasKey bool
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if tok, ok := tok.(xml.StartElement); ok {
			if tok.Name.Local == "key" {
				if err := d.DecodeElement(&key, &tok); err != nil {
					return nil, err
				}
				hasKey = true
			} else if hasKey {
				var val string
				var err error
				switch n := tok.Name.Local; n {
				case "true", "false":
					// Bools are represented as <true/> and <false/>.
					val = n
					err = d.Skip()
				default:
					err = d.DecodeElement(&val, &tok)
				}
				if err != nil {
					return nil, err
				}
				values[key] = val
				hasKey = false
			} else {
				if err := d.Skip(); err != nil {
					return nil, err
				}
			}
		}
	}
	return values, nil
}

func uninstall(bundleID string) error {
	cmd := idevCmd(exec.Command(
		"ideviceinstaller",
		"-U", bundleID,
	))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("ideviceinstaller -U %q: %s", bundleID, err)
	}
	return nil
}

func install(appdir string) error {
	attempt := 0
	for {
		cmd := idevCmd(exec.Command(
			"ideviceinstaller",
			"-i", appdir,
		))
		if out, err := cmd.CombinedOutput(); err != nil {
			// Sometimes, installing the app fails for some reason.
			// Give the device a few seconds and try again.
			if attempt < 5 {
				time.Sleep(5 * time.Second)
				attempt++
				continue
			}
			os.Stderr.Write(out)
			return fmt.Errorf("ideviceinstaller -i %q: %v (%d attempts)", appdir, err, attempt)
		}
		return nil
	}
}

func idevCmd(cmd *exec.Cmd) *exec.Cmd {
	if deviceID != "" {
		// Inject -u device_id after the executable, but before the arguments.
		args := []string{cmd.Args[0], "-u", deviceID}
		cmd.Args = append(args, cmd.Args[1:]...)
	}
	return cmd
}

func run(appdir, bundleID string, args []string) error {
	var env []string
	for _, e := range os.Environ() {
		// Don't override TMPDIR, HOME, GOCACHE on the device.
		if strings.HasPrefix(e, "TMPDIR=") || strings.HasPrefix(e, "HOME=") || strings.HasPrefix(e, "GOCACHE=") {
			continue
		}
		env = append(env, e)
	}
	attempt := 0
	for {
		// The device app path reported by the device might be stale, so retry
		// the lookup of the device path along with the lldb launching below.
		deviceapp, err := findDeviceAppPath(bundleID)
		if err != nil {
			// The device app path might not yet exist for a newly installed app.
			if attempt == 5 {
				return err
			}
			attempt++
			time.Sleep(5 * time.Second)
			continue
		}
		lldb := exec.Command(
			"python",
			"-", // Read script from stdin.
			appdir,
			deviceapp,
		)
		lldb.Args = append(lldb.Args, args...)
		lldb.Env = env
		lldb.Stdin = strings.NewReader(lldbDriver)
		lldb.Stdout = os.Stdout
		var out bytes.Buffer
		lldb.Stderr = io.MultiWriter(&out, os.Stderr)
		err = lldb.Start()
		if err == nil {
			// Forward SIGQUIT to the lldb driver which in turn will forward
			// to the running program.
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGQUIT)
			proc := lldb.Process
			go func() {
				for sig := range sigs {
					proc.Signal(sig)
				}
			}()
			err = lldb.Wait()
			signal.Stop(sigs)
			close(sigs)
		}
		// If the program was not started it can be retried without papering over
		// real test failures.
		started := bytes.HasPrefix(out.Bytes(), []byte("lldb: running program"))
		if started || err == nil || attempt == 5 {
			return err
		}
		// Sometimes, the app was not yet ready to launch or the device path was
		// stale. Retry.
		attempt++
		time.Sleep(5 * time.Second)
	}
}

func copyLocalDir(dst, src string) error {
	if err := os.Mkdir(dst, 0755); err != nil {
		return err
	}

	d, err := os.Open(src)
	if err != nil {
		return err
	}
	defer d.Close()
	fi, err := d.Readdir(-1)
	if err != nil {
		return err
	}

	for _, f := range fi {
		if f.IsDir() {
			if f.Name() == "testdata" {
				if err := cp(dst, filepath.Join(src, f.Name())); err != nil {
					return err
				}
			}
			continue
		}
		if err := cp(dst, filepath.Join(src, f.Name())); err != nil {
			return err
		}
	}
	return nil
}

func cp(dst, src string) error {
	out, err := exec.Command("cp", "-a", src, dst).CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
	}
	return err
}

func copyLocalData(dstbase string) (pkgpath string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	finalPkgpath, underGoRoot, err := subdir()
	if err != nil {
		return "", err
	}
	cwd = strings.TrimSuffix(cwd, finalPkgpath)

	// Copy all immediate files and testdata directories between
	// the package being tested and the source root.
	pkgpath = ""
	for _, element := range strings.Split(finalPkgpath, string(filepath.Separator)) {
		if debug {
			log.Printf("copying %s", pkgpath)
		}
		pkgpath = filepath.Join(pkgpath, element)
		dst := filepath.Join(dstbase, pkgpath)
		src := filepath.Join(cwd, pkgpath)
		if err := copyLocalDir(dst, src); err != nil {
			return "", err
		}
	}

	if underGoRoot {
		// Copy timezone file.
		//
		// Typical apps have the zoneinfo.zip in the root of their app bundle,
		// read by the time package as the working directory at initialization.
		// As we move the working directory to the GOROOT pkg directory, we
		// install the zoneinfo.zip file in the pkgpath.
		err := cp(
			filepath.Join(dstbase, pkgpath),
			filepath.Join(cwd, "lib", "time", "zoneinfo.zip"),
		)
		if err != nil {
			return "", err
		}
		// Copy src/runtime/textflag.h for (at least) Test386EndToEnd in
		// cmd/asm/internal/asm.
		runtimePath := filepath.Join(dstbase, "src", "runtime")
		if err := os.MkdirAll(runtimePath, 0755); err != nil {
			return "", err
		}
		err = cp(
			filepath.Join(runtimePath, "textflag.h"),
			filepath.Join(cwd, "src", "runtime", "textflag.h"),
		)
		if err != nil {
			return "", err
		}
	}

	return finalPkgpath, nil
}

// subdir determines the package based on the current working directory,
// and returns the path to the package source relative to $GOROOT (or $GOPATH).
func subdir() (pkgpath string, underGoRoot bool, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false, err
	}
	cwd, err = filepath.EvalSymlinks(cwd)
	if err != nil {
		log.Fatal(err)
	}
	goroot, err := filepath.EvalSymlinks(runtime.GOROOT())
	if err != nil {
		return "", false, err
	}
	if strings.HasPrefix(cwd, goroot) {
		subdir, err := filepath.Rel(goroot, cwd)
		if err != nil {
			return "", false, err
		}
		return subdir, true, nil
	}

	for _, p := range filepath.SplitList(build.Default.GOPATH) {
		pabs, err := filepath.EvalSymlinks(p)
		if err != nil {
			return "", false, err
		}
		if !strings.HasPrefix(cwd, pabs) {
			continue
		}
		subdir, err := filepath.Rel(pabs, cwd)
		if err == nil {
			return subdir, false, nil
		}
	}
	return "", false, fmt.Errorf(
		"working directory %q is not in either GOROOT(%q) or GOPATH(%q)",
		cwd,
		runtime.GOROOT(),
		build.Default.GOPATH,
	)
}

func infoPlist(pkgpath string) string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>CFBundleName</key><string>golang.gotest</string>
<key>CFBundleSupportedPlatforms</key><array><string>iPhoneOS</string></array>
<key>CFBundleExecutable</key><string>gotest</string>
<key>CFBundleVersion</key><string>1.0</string>
<key>CFBundleIdentifier</key><string>` + bundleID + `</string>
<key>CFBundleResourceSpecification</key><string>ResourceRules.plist</string>
<key>LSRequiresIPhoneOS</key><true/>
<key>CFBundleDisplayName</key><string>gotest</string>
<key>GoExecWrapperWorkingDirectory</key><string>` + pkgpath + `</string>
</dict>
</plist>
`
}

func entitlementsPlist() string {
	return `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>keychain-access-groups</key>
	<array><string>` + appID + `</string></array>
	<key>get-task-allow</key>
	<true/>
	<key>application-identifier</key>
	<string>` + appID + `</string>
	<key>com.apple.developer.team-identifier</key>
	<string>` + teamID + `</string>
</dict>
</plist>
`
}

const resourceRules = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>rules</key>
	<dict>
		<key>.*</key>
		<true/>
		<key>Info.plist</key>
		<dict>
			<key>omit</key>
			<true/>
			<key>weight</key>
			<integer>10</integer>
		</dict>
		<key>ResourceRules.plist</key>
		<dict>
			<key>omit</key>
			<true/>
			<key>weight</key>
			<integer>100</integer>
		</dict>
	</dict>
</dict>
</plist>
`

const lldbDriver = `
import sys
import os
import signal

exe, device_exe, args = sys.argv[1], sys.argv[2], sys.argv[3:]

env = []
for k, v in os.environ.items():
	env.append(k + "=" + v)

sys.path.append('/Applications/Xcode.app/Contents/SharedFrameworks/LLDB.framework/Resources/Python')

import lldb

debugger = lldb.SBDebugger.Create()
debugger.SetAsync(True)
debugger.SkipLLDBInitFiles(True)

err = lldb.SBError()
target = debugger.CreateTarget(exe, None, 'remote-ios', True, err)
if not target.IsValid() or not err.Success():
	sys.stderr.write("lldb: failed to setup up target: %s\n" % (err))
	sys.exit(1)

target.modules[0].SetPlatformFileSpec(lldb.SBFileSpec(device_exe))

listener = debugger.GetListener()
process = target.ConnectRemote(listener, 'connect://localhost:3222', None, err)
if not err.Success():
	sys.stderr.write("lldb: failed to connect to remote target: %s\n" % (err))
	sys.exit(1)

# Don't stop on signals.
sigs = process.GetUnixSignals()
for i in range(0, sigs.GetNumSignals()):
	sig = sigs.GetSignalAtIndex(i)
	sigs.SetShouldStop(sig, False)
	sigs.SetShouldNotify(sig, False)

event = lldb.SBEvent()
running = False
prev_handler = None
while True:
	if not listener.WaitForEvent(1, event):
		continue
	if not lldb.SBProcess.EventIsProcessEvent(event):
		continue
	if running:
		# Pass through stdout and stderr.
		while True:
			out = process.GetSTDOUT(8192)
			if not out:
				break
			sys.stdout.write(out)
		while True:
			out = process.GetSTDERR(8192)
			if not out:
				break
			sys.stderr.write(out)
	state = process.GetStateFromEvent(event)
	if state in [lldb.eStateCrashed, lldb.eStateDetached, lldb.eStateUnloaded, lldb.eStateExited]:
		if running:
			signal.signal(signal.SIGQUIT, prev_handler)
		break
	elif state == lldb.eStateConnected:
		process.RemoteLaunch(args, env, None, None, None, None, 0, False, err)
		if not err.Success():
			sys.stderr.write("lldb: failed to launch remote process: %s\n" % (err))
			process.Kill()
			debugger.Terminate()
			sys.exit(1)
		# Forward SIGQUIT to the program.
		def signal_handler(signal, frame):
			process.Signal(signal)
		prev_handler = signal.signal(signal.SIGQUIT, signal_handler)
		# Tell the Go driver that the program is running and should not be retried.
		sys.stderr.write("lldb: running program\n")
		running = True
		# Process stops once at the beginning. Continue.
		process.Continue()

exitStatus = process.GetExitStatus()
process.Kill()
debugger.Terminate()
sys.exit(exitStatus)
`
