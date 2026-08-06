# learn-go

Learning Go from scratch — notes written as runnable, heavily commented programs.

Each folder is a standalone lesson you can run on its own. Every concept is
demonstrated by code that actually prints the result, rather than described in prose.

## Requirements

Go 1.21+ (built with go1.26.5). One third-party dependency (`rsc.io/quote`), used
only by `06-modules` to demonstrate how `go.mod` and `go.sum` work.

```bash
go version
go mod download   # fetch dependencies
```

## Running the lessons

From the repo root:

```bash
go run ./01-basics    # variables, int, float64, bool, strings
go run ./02-structs   # structs, custom types, initialization, value semantics
go run ./03-methods   # methods, value vs pointer receivers
go run ./04-wallet    # the Wallet exercise
go run ./05-startup   # what runs before main(), and in what order
go run ./06-modules   # go.mod, go.sum, and a real dependency
go run ./07-errors    # errors as values: the wallet, done idiomatically
```

## Lessons

| Folder | Covers |
| --- | --- |
| `01-basics` | `var` vs `:=`, zero values, `int`, `float64`, `bool`, strings, `Printf` verbs, float precision |
| `02-structs` | defining a struct, custom types, four ways to initialize, structs are copied on assignment |
| `03-methods` | methods and receivers, value vs pointer receivers, memory addresses proving the difference |
| `04-wallet` | `Wallet` struct with `Deposit`, `Withdraw`, and `Display` |
| `05-startup` | the compile-then-run pipeline, package init order, `init()`, `main()` |
| `06-modules` | modules, `go.mod` vs `go.sum`, direct vs indirect dependencies |
| `07-errors` | errors as values, sentinel errors, custom error types, wrapping, `errors.Is` / `errors.As` |

## Errors are values

Go has no exceptions. `error` is an ordinary interface — anything with an
`Error() string` method satisfies it — and failures are returned, not thrown:

```go
func (w *Wallet) Withdraw(amount float64) error {
    if amount > w.Balance {
        return fmt.Errorf("withdraw from wallet %d: %w", w.ID,
            &InsufficientFundsError{Requested: amount, Available: w.Balance})
    }
    w.Balance -= amount
    return nil // nil means success
}
```

`07-errors` rewrites the wallet this way. The pieces:

- **Sentinel errors** (`var ErrInvalidAmount = errors.New(...)`) — fixed values for
  failures the caller identifies by identity.
- **Custom error types** — a struct with an `Error() string` method, for when the
  caller needs *data* about the failure (how much they were short by).
- **`%w` wrapping** — `fmt.Errorf("...: %w", err)` adds context while keeping the
  original retrievable. `%v` would flatten it to text and lose it.
- **`errors.Is(err, target)`** — "is that specific error anywhere in the chain?"
  Use it instead of `==`, which fails on wrapped errors.
- **`errors.As(err, &target)`** — "is there an error of this type in the chain?"
  and hands you the struct so you can read its fields.

Libraries return errors; they never print them. Printing is the caller's decision.
`panic` is for unrecoverable bugs, not for a user entering a bad number.

## How a Go program starts

Go is **compiled ahead of time to native machine code**. Nothing reads the `.go`
files at runtime — by the time the program runs, the source is irrelevant.

```
your .go files
   │  go build: parse → type-check → optimize → machine code
   ▼
compiled packages
   │  link: bundle with the Go runtime (GC, scheduler, allocator)
   ▼
one self-contained binary
   │  OS exec: loads it, jumps to the RUNTIME entry point (not your main)
   ▼
runtime setup → package init → main.main() → exit when main returns
```

Package initialization, which `05-startup` prints in order:

1. imported packages are fully initialized first, dependencies deepest-first
2. within a package: package-level `var`s, then `init()` functions in source order
3. only then does `main()` run
4. when `main()` **returns, the process exits immediately** — it does not wait for
   anything still running in the background

Because it compiles first, a type error stops the whole program from building —
no part of it runs, not even lines above the mistake.

## go.mod and go.sum

| File | Purpose |
| --- | --- |
| `go.mod` | the module's name, its Go version, and what it depends on. Hand-editable. Commit it. |
| `go.sum` | cryptographic checksums of the exact dependency content that was downloaded. Generated, never hand-edited. Commit it too. |

`go.sum` is not a lockfile — `go.mod` already pins exact versions. `go.sum` is a
tamper check: if a published version's bytes ever change, the build fails instead of
silently using different code. Dependencies marked `// indirect` are things your
dependencies need, not things you imported yourself.

Useful commands:

```bash
go mod init <name>   # create go.mod
go mod tidy          # add what's imported, drop what isn't, refresh go.sum
go get <pkg>@<ver>   # add or upgrade one dependency
go build ./...       # compile everything
go vet ./...         # report suspicious code the compiler still accepts
gofmt -l .           # list files that aren't formatted canonically
```

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
