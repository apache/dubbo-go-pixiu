// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Build initialization (after flag parsing).

package work

import (
	"cmd/go/internal/base"
	"cmd/go/internal/cfg"
	"cmd/go/internal/load"
	"cmd/internal/sys"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func BuildInit() {
	load.ModInit()
	instrumentInit()
	buildModeInit()

	// Make sure -pkgdir is absolute, because we run commands
	// in different directories.
	if cfg.BuildPkgdir != "" && !filepath.IsAbs(cfg.BuildPkgdir) {
		p, err := filepath.Abs(cfg.BuildPkgdir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "go %s: evaluating -pkgdir: %v\n", flag.Args()[0], err)
			base.SetExitStatus(2)
			base.Exit()
		}
		cfg.BuildPkgdir = p
	}
}

func instrumentInit() {
	if !cfg.BuildRace && !cfg.BuildMSan {
		return
	}
	if cfg.BuildRace && cfg.BuildMSan {
		fmt.Fprintf(os.Stderr, "go %s: may not use -race and -msan simultaneously\n", flag.Args()[0])
		base.SetExitStatus(2)
		base.Exit()
	}
	if cfg.BuildMSan && !sys.MSanSupported(cfg.Goos, cfg.Goarch) {
		fmt.Fprintf(os.Stderr, "-msan is not supported on %s/%s\n", cfg.Goos, cfg.Goarch)
		base.SetExitStatus(2)
		base.Exit()
	}
	if cfg.BuildRace {
		if !sys.RaceDetectorSupported(cfg.Goos, cfg.Goarch) {
			fmt.Fprintf(os.Stderr, "go %s: -race is only supported on linux/amd64, linux/ppc64le, linux/arm64, freebsd/amd64, netbsd/amd64, darwin/amd64 and windows/amd64\n", flag.Args()[0])
			base.SetExitStatus(2)
			base.Exit()
		}
	}
	mode := "race"
	if cfg.BuildMSan {
		mode = "msan"
		// MSAN does not support non-PIE binaries on ARM64.
		// See issue #33712 for details.
		if cfg.Goos == "linux" && cfg.Goarch == "arm64" && cfg.BuildBuildmode == "default" {
			cfg.BuildBuildmode = "pie"
		}
	}
	modeFlag := "-" + mode

	if !cfg.BuildContext.CgoEnabled {
		fmt.Fprintf(os.Stderr, "go %s: %s requires cgo; enable cgo by setting CGO_ENABLED=1\n", flag.Args()[0], modeFlag)
		base.SetExitStatus(2)
		base.Exit()
	}
	forcedGcflags = append(forcedGcflags, modeFlag)
	forcedLdflags = append(forcedLdflags, modeFlag)

	if cfg.BuildContext.InstallSuffix != "" {
		cfg.BuildContext.InstallSuffix += "_"
	}
	cfg.BuildContext.InstallSuffix += mode
	cfg.BuildContext.BuildTags = append(cfg.BuildContext.BuildTags, mode)
}

func buildModeInit() {
	gccgo := cfg.BuildToolchainName == "gccgo"
	var codegenArg string

	// Configure the build mode first, then verify that it is supported.
	// That way, if the flag is completely bogus we will prefer to error out with
	// "-buildmode=%s not supported" instead of naming the specific platform.

	switch cfg.BuildBuildmode {
	case "archive":
		pkgsFilter = pkgsNotMain
	case "c-archive":
		pkgsFilter = oneMainPkg
		if gccgo {
			codegenArg = "-fPIC"
		} else {
			switch cfg.Goos {
			case "darwin":
				switch cfg.Goarch {
				case "arm", "arm64":
					codegenArg = "-shared"
				}

			case "dragonfly", "freebsd", "illumos", "linux", "netbsd", "openbsd", "solaris":
				// Use -shared so that the result is
				// suitable for inclusion in a PIE or
				// shared library.
				codegenArg = "-shared"
			}
		}
		cfg.ExeSuffix = ".a"
		ldBuildmode = "c-archive"
	case "c-shared":
		pkgsFilter = oneMainPkg
		if gccgo {
			codegenArg = "-fPIC"
		} else {
			switch cfg.Goos {
			case "linux", "android", "freebsd":
				codegenArg = "-shared"
			case "windows":
				// Do not add usual .exe suffix to the .dll file.
				cfg.ExeSuffix = ""
			}
		}
		ldBuildmode = "c-shared"
	case "default":
		switch cfg.Goos {
		case "android":
			codegenArg = "-shared"
			ldBuildmode = "pie"
		case "darwin":
			switch cfg.Goarch {
			case "arm", "arm64":
				codegenArg = "-shared"
			}
			fallthrough
		default:
			ldBuildmode = "exe"
		}
		if gccgo {
			codegenArg = ""
		}
	case "exe":
		pkgsFilter = pkgsMain
		ldBuildmode = "exe"
		// Set the pkgsFilter to oneMainPkg if the user passed a specific binary output
		// and is using buildmode=exe for a better error message.
		// See issue #20017.
		if cfg.BuildO != "" {
			pkgsFilter = oneMainPkg
		}
	case "pie":
		if cfg.BuildRace {
			base.Fatalf("-buildmode=pie not supported when -race is enabled")
		}
		if gccgo {
			codegenArg = "-fPIE"
		} else if cfg.Goos != "aix" {
			codegenArg = "-shared"
		}
		ldBuildmode = "pie"
	case "shared":
		pkgsFilter = pkgsNotMain
		if gccgo {
			codegenArg = "-fPIC"
		} else {
			codegenArg = "-dynlink"
		}
		if cfg.BuildO != "" {
			base.Fatalf("-buildmode=shared and -o not supported together")
		}
		ldBuildmode = "shared"
	case "plugin":
		pkgsFilter = oneMainPkg
		if gccgo {
			codegenArg = "-fPIC"
		} else {
			codegenArg = "-dynlink"
		}
		cfg.ExeSuffix = ".so"
		ldBuildmode = "plugin"
	default:
		base.Fatalf("buildmode=%s not supported", cfg.BuildBuildmode)
	}

	if !sys.BuildModeSupported(cfg.BuildToolchainName, cfg.BuildBuildmode, cfg.Goos, cfg.Goarch) {
		base.Fatalf("-buildmode=%s not supported on %s/%s\n", cfg.BuildBuildmode, cfg.Goos, cfg.Goarch)
	}

	if cfg.BuildLinkshared {
		if !sys.BuildModeSupported(cfg.BuildToolchainName, "shared", cfg.Goos, cfg.Goarch) {
			base.Fatalf("-linkshared not supported on %s/%s\n", cfg.Goos, cfg.Goarch)
		}
		if gccgo {
			codegenArg = "-fPIC"
		} else {
			forcedAsmflags = append(forcedAsmflags, "-D=GOBUILDMODE_shared=1")
			codegenArg = "-dynlink"
			forcedGcflags = append(forcedGcflags, "-linkshared")
			// TODO(mwhudson): remove -w when that gets fixed in linker.
			forcedLdflags = append(forcedLdflags, "-linkshared", "-w")
		}
	}
	if codegenArg != "" {
		if gccgo {
			forcedGccgoflags = append([]string{codegenArg}, forcedGccgoflags...)
		} else {
			forcedAsmflags = append([]string{codegenArg}, forcedAsmflags...)
			forcedGcflags = append([]string{codegenArg}, forcedGcflags...)
		}
		// Don't alter InstallSuffix when modifying default codegen args.
		if cfg.BuildBuildmode != "default" || cfg.BuildLinkshared {
			if cfg.BuildContext.InstallSuffix != "" {
				cfg.BuildContext.InstallSuffix += "_"
			}
			cfg.BuildContext.InstallSuffix += codegenArg[1:]
		}
	}

	switch cfg.BuildMod {
	case "":
		// ok
	case "readonly", "vendor", "mod":
		if !cfg.ModulesEnabled && !inGOFLAGS("-mod") {
			base.Fatalf("build flag -mod=%s only valid when using modules", cfg.BuildMod)
		}
	default:
		base.Fatalf("-mod=%s not supported (can be '', 'mod', 'readonly', or 'vendor')", cfg.BuildMod)
	}
	if !cfg.ModulesEnabled {
		if cfg.ModCacheRW && !inGOFLAGS("-modcacherw") {
			base.Fatalf("build flag -modcacherw only valid when using modules")
		}
		if cfg.ModFile != "" && !inGOFLAGS("-mod") {
			base.Fatalf("build flag -modfile only valid when using modules")
		}
	}
}

func inGOFLAGS(flag string) bool {
	for _, goflag := range base.GOFLAGS() {
		name := goflag
		if strings.HasPrefix(name, "--") {
			name = name[1:]
		}
		if i := strings.Index(name, "="); i >= 0 {
			name = name[:i]
		}
		if name == flag {
			return true
		}
	}
	return false
}
