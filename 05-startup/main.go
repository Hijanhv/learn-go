// HOW A GO PROGRAM STARTS
//
// Short answer to "enter, read, or compile?": Go COMPILES, ahead of time,
// all the way to native machine code. Nothing reads your .go files at runtime.
// By the time your program runs, the source is irrelevant — it has become a
// single self-contained binary full of ARM64/x86 instructions.
//
// The full journey, source -> running program:
//
//  1. go build   parses your .go files, type-checks them, optimizes, and emits
//     machine code for each package (.a archive files)
//  2. link       stitches those together WITH the Go runtime (garbage collector,
//     scheduler, memory allocator) into one executable file
//  3. OS exec    you run the binary; the OS loads it into memory and jumps to
//     the entry point — which is the RUNTIME's startup code, not your main
//  4. runtime    sets up the heap, the GC, and the first goroutine
//  5. package    every imported package is initialized, dependencies first:
//     init         package-level vars first, then init() functions
//  6. main.main  finally, the runtime calls your main() function
//  7. exit       when main() RETURNS, the process exits immediately —
//     it does not wait for anything still running
//
// Run this file and the printed numbers show steps 5 and 6 happening in order.
package main

import (
	"fmt"

	"learn-go/05-startup/vault" // import path = module name + folder path
)

// Package-level vars are computed before main() — and before init().
var limit = announce()

func announce() int {
	fmt.Println("  4. main: package-level var limit is being computed")
	return vault.Capacity / 2
}

// main's init() runs after main's vars, and after ALL imported packages are done.
func init() {
	fmt.Println("  5. main: init() runs")
}

// main is the entry point of YOUR code. Two hard rules:
//   - it must be in package main
//   - it takes no arguments and returns nothing (use os.Args / os.Exit instead)
func main() {
	fmt.Println("  6. main: main() finally runs")
	fmt.Printf("     vault.Capacity=%d, limit=%d\n", vault.Capacity, limit)
	fmt.Println("\nNotice the order: the imported package finished ALL of its setup")
	fmt.Println("before main's own variables were even computed.")
}

func init() {
	// A second init() in the same file is legal. They run top to bottom.
	fmt.Println("  5b. main: a second init() — they run in source order")
}
