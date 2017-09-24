package terminfo

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
)

var termNameCache = struct {
	names map[string]string
	sync.Mutex
}{}

var fileRE = regexp.MustCompile("^([0-9]+|[a-zA-Z])/")

func terms(t *testing.T) map[string]string {
	termNameCache.Lock()
	defer termNameCache.Unlock()

	if termNameCache.names == nil {
		termNameCache.names = make(map[string]string)
		for _, dir := range termDirs() {
			werr := filepath.Walk(dir, func(file string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if fi.IsDir() || !fileRE.MatchString(file[len(dir)+1:]) || fi.Mode()&os.ModeSymlink != 0 {
					return nil
				}

				term := filepath.Base(file)
				/*if term != "kterm" {
					continue
				}*/
				if term == "xterm-old" {
					return nil
				}

				termNameCache.names[term] = file
				return nil
			})
			if werr != nil {
				t.Fatalf("could not walk directory, got: %v", werr)
			}
		}
	}

	return termNameCache.names
}

func termDirs() []string {
	return []string{"/lib/terminfo", "/usr/share/terminfo"}
}
