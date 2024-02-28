package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// In this way, we avoid deadlocking, since the resources are
// locked only if they are all available.

type Arbitrator struct {
	accountsInUse map[string]bool
	cond          *sync.Cond
}

func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: make(map[string]bool),
		cond:          sync.NewCond(new(sync.Mutex)),
	}
}

/*
Next, we need to implement a function that allows us to block the accounts if they are
free or to suspend the execution of the goroutine if they’re not.

If any of the accounts are in use, the gorou-
tine calls Wait() on the condition variable. This suspends the execution of the gorou-
tine and unlocks the mutex. Once the execution is resumed, the goroutine reacquires
the mutex, and this check is repeated until all the accounts are free.
*/

func (a *Arbitrator) LockAccounts(ids ...string) {
	a.cond.L.Lock()
	for allAvailable := false; !allAvailable; {
		allAvailable = true

		for _, id := range ids {
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
			}
		}
	}
	for _, id := range ids {
		a.accountsInUse[id] = true
	}
	a.cond.L.Unlock()
}

/*
Once the goroutine is done with its transfer logic, it needs to mark the accounts as no
longer in use.

This has the effect of waking up
any suspended goroutines, which will then go ahead and check to see if their accounts
have become available.
*/

func (a *Arbitrator) UnlockAccounts(ids ...string) {
	a.cond.L.Lock()
	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	a.cond.Broadcast()
	a.cond.L.Unlock()
}

type BankAccount struct {
	id      string
	balance int
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id, 100,
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, tellerID int, arb *Arbitrator) {
	fmt.Printf("%d Locking %s and %s\n", tellerID, src.id, to.id)
	arb.LockAccounts(src.id, to.id)
	src.balance -= amount
	to.balance += amount
	arb.UnlockAccounts(to.id, src.id)
	fmt.Printf("%d Unlocked %s and %s\n", tellerID, src.id, to.id)
}

func main() {
	accounts := []*BankAccount{
		NewBankAccount("cr7"),
		NewBankAccount("vini"),
		NewBankAccount("bale"),
		NewBankAccount("ramos"),
	}
	total := len(accounts)
	arb := NewArbitrator()
	for i := 0; i < 4; i++ {
		go func(tellerID int) {
			for j := 1; j < 1000; j++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(accounts[to], 10, tellerID, arb)
			}
			fmt.Println(tellerID, "COMPLETE")
		}(i)
	}
	time.Sleep(1 * time.Second)
}
