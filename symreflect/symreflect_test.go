package symreflect_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cretz/interflect/symreflect/internal/symreflecttest/testpkg"
	"github.com/stretchr/testify/require"
)

func TestLoadFuncsInTest(t *testing.T) {
	testpkg.TestFuncs(require.New(t))
}

func TestLoadFuncsInProg(t *testing.T) {
	buildAndRunProgTest(t, "funcs")
}

func TestLoadSymbolsInTest(t *testing.T) {
	// Symbols are elided in "go test" builds
	testpkg.TestSymbols(require.New(t), false)
}

func TestLoadSymbolsInProg(t *testing.T) {
	buildAndRunProgTest(t, "symbols")
}

func TestLoadTypesInTest(t *testing.T) {
	testpkg.TestTypes(require.New(t))
}

func TestLoadTypesInProg(t *testing.T) {
	buildAndRunProgTest(t, "types")
}

func buildAndRunProgTest(t *testing.T, arg string) {
	// Temp dir
	tmp, err := os.MkdirTemp("", "symreflect-test-")
	require.NoError(t, err)
	t.Logf("Building in %v", tmp)
	defer os.RemoveAll(tmp)

	// Build into that dir
	_, currFile, _, _ := runtime.Caller(0)
	exe := filepath.Join(tmp, "main")
	if runtime.GOOS == "windows" {
		exe += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", exe, ".")
	cmd.Dir = filepath.Join(currFile, "../internal/symreflecttest/prog")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "failed building: %s", out)
	t.Logf("Build output: %s", out)

	// Run it
	out, err = exec.Command(exe, arg).CombinedOutput()
	require.NoError(t, err, "failed running: %s", out)
	t.Logf("Run output: %s", out)
}
