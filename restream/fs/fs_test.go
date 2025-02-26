package fs

import (
	"strings"
	"testing"
	"time"

	"github.com/datarhei/core/io/fs"
	"github.com/stretchr/testify/require"
)

func TestMaxFiles(t *testing.T) {
	memfs := fs.NewMemFilesystem(fs.MemConfig{
		Base:  "/",
		Size:  1024,
		Purge: false,
	})

	cleanfs := New(Config{
		FS: memfs,
	})

	cleanfs.Start()

	cleanfs.SetCleanup("foobar", []Pattern{
		{
			Pattern:    "/*.ts",
			MaxFiles:   3,
			MaxFileAge: 0,
		},
	})

	cleanfs.Store("/chunk_0.ts", strings.NewReader("chunk_0"))
	cleanfs.Store("/chunk_1.ts", strings.NewReader("chunk_1"))
	cleanfs.Store("/chunk_2.ts", strings.NewReader("chunk_2"))

	require.Eventually(t, func() bool {
		return cleanfs.Files() == 3
	}, 3*time.Second, time.Second)

	cleanfs.Store("/chunk_3.ts", strings.NewReader("chunk_3"))

	require.Eventually(t, func() bool {
		if cleanfs.Files() != 3 {
			return false
		}

		names := []string{}

		for _, f := range cleanfs.List("/*.ts") {
			names = append(names, f.Name())
		}

		require.ElementsMatch(t, []string{"/chunk_1.ts", "/chunk_2.ts", "/chunk_3.ts"}, names)

		return true
	}, 3*time.Second, time.Second)

	cleanfs.Stop()
}

func TestMaxAge(t *testing.T) {
	memfs := fs.NewMemFilesystem(fs.MemConfig{
		Base:  "/",
		Size:  1024,
		Purge: false,
	})

	cleanfs := New(Config{
		FS: memfs,
	})

	cleanfs.Start()

	cleanfs.SetCleanup("foobar", []Pattern{
		{
			Pattern:    "/*.ts",
			MaxFiles:   0,
			MaxFileAge: 3 * time.Second,
		},
	})

	cleanfs.Store("/chunk_0.ts", strings.NewReader("chunk_0"))
	cleanfs.Store("/chunk_1.ts", strings.NewReader("chunk_1"))
	cleanfs.Store("/chunk_2.ts", strings.NewReader("chunk_2"))

	require.Eventually(t, func() bool {
		return cleanfs.Files() == 0
	}, 5*time.Second, time.Second)

	cleanfs.Store("/chunk_3.ts", strings.NewReader("chunk_3"))

	require.Eventually(t, func() bool {
		if cleanfs.Files() != 1 {
			return false
		}

		names := []string{}

		for _, f := range cleanfs.List("/*.ts") {
			names = append(names, f.Name())
		}

		require.ElementsMatch(t, []string{"/chunk_3.ts"}, names)

		return true
	}, 2*time.Second, time.Second)

	cleanfs.Stop()
}

func TestUnsetCleanup(t *testing.T) {
	memfs := fs.NewMemFilesystem(fs.MemConfig{
		Base:  "/",
		Size:  1024,
		Purge: false,
	})

	cleanfs := New(Config{
		FS: memfs,
	})

	cleanfs.Start()

	cleanfs.SetCleanup("foobar", []Pattern{
		{
			Pattern:    "/*.ts",
			MaxFiles:   3,
			MaxFileAge: 0,
		},
	})

	cleanfs.Store("/chunk_0.ts", strings.NewReader("chunk_0"))
	cleanfs.Store("/chunk_1.ts", strings.NewReader("chunk_1"))
	cleanfs.Store("/chunk_2.ts", strings.NewReader("chunk_2"))

	require.Eventually(t, func() bool {
		return cleanfs.Files() == 3
	}, 3*time.Second, time.Second)

	cleanfs.Store("/chunk_3.ts", strings.NewReader("chunk_3"))

	require.Eventually(t, func() bool {
		if cleanfs.Files() != 3 {
			return false
		}

		names := []string{}

		for _, f := range cleanfs.List("/*.ts") {
			names = append(names, f.Name())
		}

		require.ElementsMatch(t, []string{"/chunk_1.ts", "/chunk_2.ts", "/chunk_3.ts"}, names)

		return true
	}, 3*time.Second, time.Second)

	cleanfs.UnsetCleanup("foobar")

	cleanfs.Store("/chunk_4.ts", strings.NewReader("chunk_4"))

	require.Eventually(t, func() bool {
		if cleanfs.Files() != 4 {
			return false
		}

		names := []string{}

		for _, f := range cleanfs.List("/*.ts") {
			names = append(names, f.Name())
		}

		require.ElementsMatch(t, []string{"/chunk_1.ts", "/chunk_2.ts", "/chunk_3.ts", "/chunk_4.ts"}, names)

		return true
	}, 3*time.Second, time.Second)

	cleanfs.Stop()
}
