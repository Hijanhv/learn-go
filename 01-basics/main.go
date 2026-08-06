// Every Go file starts by declaring which package it belongs to.
// "main" is special: it means "this builds into a runnable program".
package main

// fmt = "format". The standard library package for printing.
import "fmt"

func main() {
	// --- VARIABLES ---

	// Long form: var <name> <type> = <value>
	var owner string = "Janhavi"

	// Go can infer the type, so you can drop it:
	var country = "India"

	// Short form (only usable INSIDE a function). This is what you'll write 90% of the time.
	// := means "declare and assign, infer the type".
	city := "Mumbai"

	// Declare without a value -> Go gives it the ZERO VALUE for that type.
	// There is no "undefined"/"null" for these. Numbers start at 0, strings at "", bools at false.
	var uninitialized string
	fmt.Println(owner, country, city, "<-empty string here:", uninitialized, "|")

	// --- INT ---
	// int is a whole number. Its size follows your CPU (64-bit here).
	var id int = 1001
	age := 21 // inferred as int

	// Integer division TRUNCATES — it does not round.
	fmt.Println("7/2 as ints =", 7/2) // 3, not 3.5

	// --- FLOAT64 ---
	// float64 is a 64-bit decimal number. Untyped decimals default to float64.
	balance := 1500.75 // float64
	var rate float64 = 0.075

	// Go does NOT auto-convert between types. This would not compile:
	//     total := balance + age    // mismatched types float64 and int
	// You must convert explicitly:
	total := balance + float64(age)
	fmt.Println("id:", id, "balance:", balance, "rate:", rate, "total:", total)

	// Floats are approximate. Note the difference:
	fmt.Println("constant 0.1 + 0.2 =", 0.1+0.2) // 0.3  <- computed exactly at COMPILE time
	x, y := 0.1, 0.2
	fmt.Println("variable x + y     =", x+y) // 0.30000000000000004 <- real float64 math
	// Lesson: never compare floats with ==, and never store money as float64 in a real
	// production system (use integer cents). float64 is fine for learning.

	// --- BOOL ---
	// Only true / false. There is no "truthy" — 0 and "" are NOT false in Go.
	isActive := true
	hasEnough := balance >= 100.0
	fmt.Println("isActive:", isActive, "hasEnough:", hasEnough, "both:", isActive && hasEnough)

	// --- STRINGS ---
	// Strings are immutable and UTF-8. Double quotes only (single quotes = a rune/char).
	first := "Janhavi"
	last := "Chavada"
	full := first + " " + last // + concatenates
	fmt.Println(full, "length in bytes:", len(full))

	// Sprintf builds a formatted STRING (Printf prints it instead).
	//   %s string   %d integer   %.2f float with 2 decimals   %t bool   %v any value
	line := fmt.Sprintf("%s | id=%d | balance=%.2f | active=%t", full, id, balance, isActive)
	fmt.Println(line)
}
