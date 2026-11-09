package bundled

import (
	"fmt"
	"os"
	"sync"

	"github.com/microsoft/TypeScript/tsc/internal/osutil"
	"github.com/microsoft/TypeScript/tsc/internal/tspath"
	"github.com/microsoft/TypeScript/tsc/internal/vfs/osvfs"
)

var executableDir = sync.OnceValue(func() tspath.RootedDirectoryPath {
	exe, err := osutil.Executable()
	if err != nil {
		panic(fmt.Sprintf("bundled: failed to get executable path: %v", err))
	}
	realExe := osvfs.FS().Realpath(tspath.RootedFilePathFromAbsolute(exe).AsPath())
	return realExe.Directory()
})

var noembedLibPath = sync.OnceValue(func() string {
	dirStr := os.Getenv("TSC_BUNDLE_DIR")
	if dirStr == "" {
		return ""
	}
	dir := tspath.RootedDirectoryPathFromAbsolute(dirStr)

	libdts := dir.ResolveFile("lib.d.ts")
	if info := osvfs.FS().Stat(libdts.AsPath()); info == nil {
		panic(fmt.Sprintf("bundled: %v does not exist; this executable may be misplaced", libdts))
	}

	return dir.AsString()
})

