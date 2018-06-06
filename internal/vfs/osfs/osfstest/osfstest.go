package osfstest

import (
	"io/ioutil"
	"os"
	"path"
)

type testing interface {
	Fatal(v ...interface{})
}

// SetupPath creates a temporary directory and a test file, based on the passed
// file name.  The path to the temporary test file and a teardown function for
// deleting the temporary directory are returned.
// On failure the Fatal method of t will be executed.
func SetupPath(t testing, file string) (fileName string, teardown func()) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	if file == "" {
		file = "test.dat"
	}
	return path.Join(dir, file), func() {
		os.RemoveAll(dir)
	}
}
