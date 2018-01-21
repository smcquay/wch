package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var all = flag.Bool("a", false, "list all instances of executables found (instead of just the first one of each).")
var silent = flag.Bool("s", false, "No output, just return 0 if all of the executables are found, or 1 if some were not found.")

func main() {
	flag.Parse()

	paths := filepath.SplitList(os.Getenv("PATH"))

	rc := 0
	for _, cmd := range flag.Args() {
		found := false
		for _, dir := range paths {
			path := filepath.Join(dir, cmd)
			if err := findExecutable(path); err == nil {
				found = true
				if !*silent {
					fmt.Println(path)
				}
				if !*all {
					break
				}
			}
		}
		if !found {
			if !*silent {
				fmt.Printf("%v not found\n", cmd)
			}
			rc++
		}
	}
	if rc > 0 {
		rc = 1
	}
	os.Exit(rc)
}

// findExecutable is from the stdlib: https://golang.org/src/os/exec/lp_unix.go?s=458:647#L19
func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}
