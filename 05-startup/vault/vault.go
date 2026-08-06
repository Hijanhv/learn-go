// Package vault exists only to prove ONE thing: an imported package is fully
// initialized BEFORE the package that imports it runs any of its own code.
package vault

import "fmt"

// Package-level variables are initialized before any init() in this file.
var Capacity = computeCapacity()

func computeCapacity() int {
	fmt.Println("  2. vault: package-level var Capacity is being computed")
	return 500
}

// init() is a special function: no arguments, no return value, you never call it.
// The runtime calls it for you, exactly once, after this package's vars are set.
// A package can have several init() functions, even across multiple files.
func init() {
	fmt.Println("  3. vault: init() runs")
}
