// ERRORS ARE VALUES
//
// Go has no exceptions, no try/catch, no throw. A function that can fail simply
// RETURNS the failure as its last return value, and the caller checks it.
//
// `error` is not magic. It is an ordinary interface, defined in the standard
// library as exactly this:
//
//	type error interface {
//	    Error() string
//	}
//
// Anything with an Error() string method IS an error. That's the whole contract.
// A nil error means "no failure" — that's why `if err != nil` is everywhere.
package main

import (
	"errors"
	"fmt"
)

// --- 1. SENTINEL ERRORS ---
// Fixed error values you can compare against. Convention: name them Err...
// They're package-level vars because callers need to reference the same value.
var (
	ErrInvalidAmount = errors.New("amount must be positive")
	ErrWalletFrozen  = errors.New("wallet is frozen")
)

// --- 2. CUSTOM ERROR TYPES ---
// When the caller needs DATA about the failure, not just its identity, make a
// struct and give it an Error() string method. Now it satisfies the error interface.
type InsufficientFundsError struct {
	Requested float64
	Available float64
}

// This method is what makes the type an error.
func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf("insufficient funds: requested %.2f, available %.2f",
		e.Requested, e.Available)
}

// Because it's a real type, it can carry extra behaviour the caller can use.
func (e *InsufficientFundsError) Shortfall() float64 {
	return e.Requested - e.Available
}

type Wallet struct {
	ID      int
	Owner   string
	Balance float64
	Frozen  bool
}

// --- 3. CONSTRUCTORS RETURN (value, error) ---
// The standard Go shape. Note the pointer: callers get the real wallet, not a copy.
func NewWallet(id int, owner string, initial float64) (*Wallet, error) {
	if owner == "" {
		return nil, errors.New("owner name cannot be empty")
	}
	if initial < 0 {
		// %w WRAPS an error: the message is extended but the original stays
		// retrievable underneath. Use %w to wrap, %v to merely print.
		return nil, fmt.Errorf("initial balance %.2f: %w", initial, ErrInvalidAmount)
	}
	return &Wallet{ID: id, Owner: owner, Balance: initial}, nil
}

// Deposit now RETURNS an error instead of printing one. The caller decides what
// to do about it — that's the entire point. A library that prints is a bad library.
func (w *Wallet) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit %.2f: %w", amount, ErrInvalidAmount)
	}
	if w.Frozen {
		return fmt.Errorf("deposit into wallet %d: %w", w.ID, ErrWalletFrozen)
	}
	w.Balance += amount
	return nil // nil means success
}

func (w *Wallet) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw %.2f: %w", amount, ErrInvalidAmount)
	}
	if w.Frozen {
		return fmt.Errorf("withdraw from wallet %d: %w", w.ID, ErrWalletFrozen)
	}
	if amount > w.Balance {
		// Wrapping a custom error type works exactly the same way.
		return fmt.Errorf("withdraw from wallet %d: %w", w.ID,
			&InsufficientFundsError{Requested: amount, Available: w.Balance})
	}
	w.Balance -= amount
	return nil
}

// Display only reads, so it keeps a value receiver and returns nothing.
func (w Wallet) Display() {
	fmt.Printf("  wallet #%d | %s | %.2f\n", w.ID, w.Owner, w.Balance)
}

func main() {
	// --- THE CORE IDIOM ---
	// Call, check, handle. You will type `if err != nil` thousands of times.
	wallet, err := NewWallet(101, "Janhavi", 5000)
	if err != nil {
		fmt.Println("could not create wallet:", err)
		return
	}
	wallet.Display()

	fmt.Println("\n--- happy path ---")
	if err := wallet.Deposit(1500.50); err != nil {
		fmt.Println("deposit failed:", err)
	} else {
		fmt.Println("deposit ok")
	}
	wallet.Display()

	// --- errors.Is: "is this THAT specific error?" ---
	// It looks through the whole wrap chain, so wrapping never breaks the check.
	fmt.Println("\n--- errors.Is ---")
	err = wallet.Deposit(-50)
	fmt.Println("raw message:", err)
	if errors.Is(err, ErrInvalidAmount) {
		fmt.Println("errors.Is found ErrInvalidAmount underneath the wrapping")
	}
	// Never compare with == on a wrapped error; it fails:
	fmt.Println("err == ErrInvalidAmount is", err == ErrInvalidAmount, "<- why errors.Is exists")

	// --- errors.As: "is there an error of this TYPE, and give it to me?" ---
	fmt.Println("\n--- errors.As ---")
	err = wallet.Withdraw(99999)
	fmt.Println("raw message:", err)

	var insufficient *InsufficientFundsError
	if errors.As(err, &insufficient) {
		// We now hold the real struct and can use its data and methods.
		fmt.Printf("short by %.2f — offer the user a top-up of that amount\n",
			insufficient.Shortfall())
	}

	// --- errors.Unwrap: peel one layer ---
	fmt.Println("\n--- unwrapping ---")
	fmt.Println("outer:", err)
	fmt.Println("inner:", errors.Unwrap(err))

	// --- sentinel on a frozen wallet ---
	fmt.Println("\n--- frozen wallet ---")
	wallet.Frozen = true
	if err := wallet.Withdraw(10); errors.Is(err, ErrWalletFrozen) {
		fmt.Println("blocked:", err)
	}
	wallet.Frozen = false

	fmt.Println("\n--- final state ---")
	wallet.Display()

	// --- WHEN TO PANIC INSTEAD ---
	// Almost never. panic() is for genuinely unrecoverable bugs — an impossible
	// state, not a user typing a bad number. "Insufficient funds" is an expected
	// outcome of normal use, so it is an error value, not a panic.
}
