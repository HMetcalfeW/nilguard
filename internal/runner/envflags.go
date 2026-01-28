package runner

import (
	"flag"
	"os"
	"path/filepath"
)

const toolName = "nilguard"

// SetupEnvDefaults sets safe defaults for Go cache locations when unset.
func SetupEnvDefaults() {
	base := filepath.Join(os.TempDir(), toolName)
	setIfEmpty("GOCACHE", filepath.Join(base, "go-build"))
	setIfEmpty("GOMODCACHE", filepath.Join(base, "go-mod"))
	setIfEmpty("GOTOOLCHAINDIR", filepath.Join(base, "toolchain"))
}

// RegisterEnvFlags adds CLI flags that override cache locations.
func RegisterEnvFlags() {
	flag.Func("cache-dir", "override GOCACHE directory", func(v string) error {
		return setEnvDir("GOCACHE", v)
	})
	flag.Func("mod-cache-dir", "override GOMODCACHE directory", func(v string) error {
		return setEnvDir("GOMODCACHE", v)
	})
	flag.Func("toolchain-dir", "override GOTOOLCHAINDIR directory", func(v string) error {
		return setEnvDir("GOTOOLCHAINDIR", v)
	})
}

func setIfEmpty(key, value string) {
	if _, ok := os.LookupEnv(key); ok {
		return
	}
	_ = os.Setenv(key, value)
}

func setEnvDir(key, value string) error {
	if value == "" {
		return nil
	}
	return os.Setenv(key, value)
}
