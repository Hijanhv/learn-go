package main

import "fmt"

// Wallet is our custom type: a bundle of the three fields the task asked for.
type Wallet struct {
	ID      int     // reference number for the wallet
	Owner   string  // who it belongs to
	Balance float64 // how much is in it
}

// Deposit adds amount to the balance.
// POINTER receiver (*Wallet) because it must change the real wallet.
func (w *Wallet) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Printf("Deposit failed: %.2f is not a valid amount.\n", amount)
		return
	}
	w.Balance += amount
	fmt.Printf("Deposited %.2f. New balance: %.2f\n", amount, w.Balance)
}

// Withdraw removes amount from the balance, but only if there's enough.
// POINTER receiver for the same reason.
func (w *Wallet) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Printf("Withdrawal failed: %.2f is not a valid amount.\n", amount)
		return
	}
	if amount > w.Balance {
		fmt.Printf("Withdrawal failed: cannot withdraw %.2f, balance is only %.2f\n",
			amount, w.Balance)
		return // early return — nothing is changed
	}
	w.Balance -= amount
	fmt.Printf("Withdrew %.2f. New balance: %.2f\n", amount, w.Balance)
}

// Display prints the wallet details.
// VALUE receiver (w Wallet) because it only reads — it never modifies anything.
func (w Wallet) Display() {
	fmt.Println("---------------------------")
	fmt.Printf("Wallet ID : %d\n", w.ID)
	fmt.Printf("Owner     : %s\n", w.Owner)
	fmt.Printf("Balance   : %.2f\n", w.Balance)
	fmt.Println("---------------------------")
}

func main() {
	// Variables for the owner's name and the starting balance.
	owner := "Janhavi"
	initialBalance := 5000.00

	// Build the wallet from those variables.
	wallet := Wallet{
		ID:      101,
		Owner:   owner,
		Balance: initialBalance,
	}

	fmt.Println("Initial wallet:")
	wallet.Display()

	fmt.Println("\nTransactions:")
	wallet.Deposit(1500.50)
	wallet.Withdraw(2000.25)
	wallet.Withdraw(50000) // rejected — not enough money

	fmt.Println("\nUpdated wallet:")
	wallet.Display()
}
