# learn-go

Learning Go from scratch — notes written as runnable, heavily commented programs.

Each folder is a standalone lesson you can run on its own. Every concept is
demonstrated by code that actually prints the result, rather than described in prose.

## Requirements

Go 1.21+ (built with go1.26.5). No external dependencies — standard library only.

```bash
go version
```

## Running the lessons

From the repo root:

```bash
go run ./01-basics    # variables, int, float64, bool, strings
go run ./02-structs   # structs, custom types, initialization, value semantics
go run ./03-methods   # methods, value vs pointer receivers
go run ./04-wallet    # the Wallet exercise
```

## Lessons

| Folder | Covers |
| --- | --- |
| `01-basics` | `var` vs `:=`, zero values, `int`, `float64`, `bool`, strings, `Printf` verbs, float precision |
| `02-structs` | defining a struct, custom types, four ways to initialize, structs are copied on assignment |
| `03-methods` | methods and receivers, value vs pointer receivers, memory addresses proving the difference |
| `04-wallet` | `Wallet` struct with `Deposit`, `Withdraw`, and `Display` |

## The Wallet exercise

A `Wallet` struct with `ID int`, `Owner string`, `Balance float64`, and three methods:

- `Deposit(amount float64)` — adds to the balance
- `Withdraw(amount float64)` — deducts only if funds are sufficient, otherwise prints a message
- `Display()` — prints the wallet details

```
Initial wallet:
---------------------------
Wallet ID : 101
Owner     : Janhavi
Balance   : 5000.00
---------------------------

Transactions:
Deposited 1500.50. New balance: 6500.50
Withdrew 2000.25. New balance: 4500.25
Withdrawal failed: cannot withdraw 50000.00, balance is only 4500.25

Updated wallet:
---------------------------
Wallet ID : 101
Owner     : Janhavi
Balance   : 4500.25
---------------------------
```

## Why pointer receivers modify the original

The core idea the Wallet exercise exists to teach.

A Go struct is a **value**. Assigning it copies it:

```go
copyOfW1 := w1
copyOfW1.Balance = 999999
// w1.Balance is untouched
```

Methods follow the same rule, and the receiver decides which one you get:

```go
func (w Wallet) Display()   // gets a COPY   — writes are discarded → use for reading
func (w *Wallet) Deposit()  // gets the ADDRESS — writes stick     → use for writing
```

`03-methods` prints the memory address of each receiver to prove it: the pointer
receiver reports the same address as the original wallet, while every value
receiver reports a different one.

Two pieces of syntax sugar worth knowing:

1. `wallet.Deposit(500)` compiles to `(&wallet).Deposit(500)` — Go takes the address
   for you, which is why pointer methods look like ordinary calls.
2. `w.Balance` inside a pointer method means `(*w).Balance` — Go dereferences
   automatically, so field access never needs a `*`.

Rule of thumb: reads get value receivers, writes get pointer receivers. If any method
on a type needs a pointer receiver, give them all pointer receivers for consistency.

## Note on `float64` for money

`float64` is used here because the exercise calls for it, but it is the wrong choice
in production — `01-basics` shows `0.1 + 0.2` evaluating to `0.30000000000000004`.
Real financial code stores integer cents/paise and formats for display.
