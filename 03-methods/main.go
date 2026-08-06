package main

import "fmt"

type Wallet struct {
	ID      int
	Owner   string
	Balance float64
}

// --- WHAT IS A METHOD? ---
// A method is just a function with an extra parameter in front, called the RECEIVER.
//
//	func (w Wallet) Name(args) returnType { ... }
//	     ^^^^^^^^^^ the receiver — "this method belongs to Wallet"
//
// Go has no classes. You attach methods to your own types instead.
// `w` is like `this`/`self`, except YOU name it (convention: first letter of the type).

// --- VALUE RECEIVER ---  (w Wallet)
// The method gets a COPY of the wallet. Perfect for read-only things like Display.
func (w Wallet) Display() {
	fmt.Printf("[value receiver] wallet #%d owned by %s holds %.2f (copy lives at %p)\n",
		w.ID, w.Owner, w.Balance, &w)
}

// A value receiver CANNOT change the original — it only edits its own copy,
// which is thrown away when the method returns.
func (w Wallet) BrokenDeposit(amount float64) {
	w.Balance += amount // modifies the COPY
}

// --- POINTER RECEIVER ---  (w *Wallet)
// The method gets the ADDRESS of the wallet, so it edits the real thing.
func (w *Wallet) Deposit(amount float64) {
	w.Balance += amount // Go auto-derefs: this means (*w).Balance += amount
}

func (w *Wallet) Withdraw(amount float64) {
	if amount > w.Balance {
		fmt.Println("insufficient funds")
		return
	}
	w.Balance -= amount
}

// Prints the wallet's real memory address so you can compare it with Display's copy.
func (w *Wallet) WhereAmI() {
	fmt.Printf("[pointer receiver] I am editing the ORIGINAL at %p\n", w)
}

func main() {
	w := Wallet{ID: 1, Owner: "Janhavi", Balance: 1000}
	fmt.Printf("the real wallet lives at %p\n\n", &w)

	// Call a method with a dot, like a field:
	w.Display()
	w.WhereAmI()
	// ^ Notice: Display's address is DIFFERENT (it's a copy).
	//   WhereAmI's address MATCHES the real wallet.

	fmt.Println("\n--- value receiver tries to deposit 500 ---")
	w.BrokenDeposit(500)
	fmt.Println("balance:", w.Balance, "<- unchanged. The copy got the 500 and died.")

	fmt.Println("\n--- pointer receiver deposits 500 ---")
	w.Deposit(500)
	fmt.Println("balance:", w.Balance, "<- the original actually changed")

	fmt.Println("\n--- withdrawals ---")
	w.Withdraw(200)
	fmt.Println("after withdrawing 200:", w.Balance)
	w.Withdraw(99999) // guarded — prints the message, changes nothing
	fmt.Println("after failed withdrawal:", w.Balance)

	// --- WHY DOES `w.Deposit(500)` WORK IF w IS A VALUE, NOT A POINTER? ---
	// Sugar. Go rewrites it to (&w).Deposit(500) for you, because w is addressable.
	// You can write it out yourself and it's identical:
	(&w).Deposit(1)
	fmt.Println("\nafter (&w).Deposit(1):", w.Balance)

	// The same sugar works the other way: a pointer can call a value method.
	p := &w
	p.Display() // Go rewrites to (*p).Display()

	// --- THE MENTAL MODEL ---
	// value receiver   -> Go copies the struct     -> edits are lost      -> use for READING (Display)
	// pointer receiver -> Go passes the address    -> edits stick         -> use for WRITING (Deposit/Withdraw)
	//
	// Rule of thumb: if ANY method on a type needs a pointer receiver,
	// give ALL its methods pointer receivers, for consistency.
}
