package main

import "fmt"

// --- CREATING A CUSTOM TYPE ---
// A struct is a TYPE you define yourself: a bundle of named fields.
// Think "a row in a table" or "a class with only data, no methods (yet)".
//
// Field names starting with a CAPITAL letter are exported (visible to other
// packages). Lowercase = private to this package. That's Go's entire access-control system.
type Wallet struct {
	ID      int
	Owner   string
	Balance float64
}

// You can also define a custom type on top of an existing one.
// Currency is its OWN type — you can't accidentally mix it with a plain string.
type Currency string

// And on numbers, which is how you'd get type-safe units:
type Rupees float64

func main() {
	// --- 4 WAYS TO INITIALIZE A STRUCT ---

	// 1. Field names (BEST — order doesn't matter, missing fields get zero values).
	w1 := Wallet{ID: 1, Owner: "Janhavi", Balance: 5000.50}

	// 2. Positional — must match declaration order exactly. Brittle: if someone
	//    reorders the struct, your code silently breaks. Avoid.
	w2 := Wallet{2, "Aryan", 250.00}

	// 3. Zero value — every field gets its type's zero value (0, "", false).
	var w3 Wallet
	w3.ID = 3 // assign fields afterwards with dot notation
	w3.Owner = "Riya"
	// w3.Balance stays 0

	// 4. Partial — unlisted fields are zero. Handy for "new empty account".
	w4 := Wallet{ID: 4, Owner: "Dev"}

	fmt.Println(w1)
	fmt.Println(w2)
	fmt.Println(w3)
	fmt.Println(w4)

	// %+v prints field NAMES too — your best friend for debugging.
	fmt.Printf("%v\n", w1)  // {1 Janhavi 5000.5}
	fmt.Printf("%+v\n", w1) // {ID:1 Owner:Janhavi Balance:5000.5}
	fmt.Printf("%#v\n", w1) // main.Wallet{ID:1, Owner:"Janhavi", Balance:5000.5}

	// Reading and writing a field:
	fmt.Println("Owner of w1:", w1.Owner)
	w1.Balance = w1.Balance + 100
	fmt.Println("w1 balance after +100:", w1.Balance)

	// --- THE BIG ONE: STRUCTS ARE VALUES, NOT REFERENCES ---
	// Assigning a struct COPIES it. copy is a completely separate wallet.
	copyOfW1 := w1
	copyOfW1.Balance = 999999

	fmt.Println("original w1.Balance:", w1.Balance)      // unchanged
	fmt.Println("copyOfW1.Balance:  ", copyOfW1.Balance) // changed
	// ^ Remember this. It's the whole reason pointer receivers exist (next lesson).

	// Custom types in action:
	var c Currency = "INR"
	var amount Rupees = 500.25
	fmt.Printf("%s %.2f\n", c, amount)

	// This would NOT compile — Rupees and float64 are different types:
	//     var plain float64 = amount
	// You must convert on purpose:
	var plain float64 = float64(amount)
	fmt.Println("converted:", plain)
}
