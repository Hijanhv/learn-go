// MODULES: go.mod and go.sum
//
// A MODULE is a collection of packages versioned together — one repo, usually.
// `go mod init learn-go` created go.mod at the repo root. That file is what makes
// this a module rather than a loose pile of .go files.
//
// go.mod  — what you DEPEND on. Hand-editable, committed to git.
// go.sum  — cryptographic CHECKSUMS of exactly what you downloaded.
//
//	Generated, never hand-edited, ALSO committed to git.
//
// This file imports a third-party package so you can watch both files change.
package main

import (
	"fmt"

	"rsc.io/quote" // NOT in the standard library — must be downloaded
)

func main() {
	fmt.Println("from the dependency:", quote.Hello())
	fmt.Println("and another:        ", quote.Go())
}
