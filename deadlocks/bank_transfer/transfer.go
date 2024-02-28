package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type BankAccount struct {
	id      string
	balance int
	mu      sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id, 100, sync.Mutex{},
	}
}

func (src *BankAccount) TransferDeadlock(to *BankAccount, amount int, exID int) {
	fmt.Printf("%d Locking %s's account\n", exID, src.id)
	src.mu.Lock()
	fmt.Printf("%d Locking %s's account\n", exID, to.id)
	to.mu.Lock()
	src.balance -= amount
	to.balance += amount
	to.mu.Unlock()
	src.mu.Unlock()
	fmt.Printf("%d Unlocked %s and %s\n", exID, src.id, to.id)
}

func (src *BankAccount) TransferWellDone(to *BankAccount, amount int, exID int) {
	accounts := []*BankAccount{src, to}
	sort.Slice(accounts, func(a, b int) bool {
		return accounts[a].id < accounts[b].id
	})
	fmt.Printf("%d Locking %s's account\n", exID, accounts[0].id)
	accounts[0].mu.Lock()
	fmt.Printf("%d Locking %s's account\n", exID, accounts[1].id)
	accounts[1].mu.Lock()
	src.balance -= amount
	to.balance += amount
	to.mu.Unlock()
	src.mu.Unlock()
	fmt.Printf("%d Unlocked %s and %s\n", exID, src.id, to.id)
}

func main() {
	accounts := []*BankAccount{
		NewBankAccount("cr7"),
		NewBankAccount("vini"),
		NewBankAccount("bale"),
		NewBankAccount("ramos"),
	}
	total := len(accounts)
	for i := 0; i < 4; i++ {
		go func(exID int) {
			for j := 1; j < 1000; j++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].TransferDeadlock(accounts[to], 10, exID)
				//accounts[from].TransferWellDone(accounts[to], 10, exID)
			}
			fmt.Println(exID, "COMPLETE")
		}(i)
	}
	time.Sleep(2 * time.Second)
}

/*
最后几行的输出，我们是看不到COMPLETE的
0 Locking cr7's account
1 Locking ramos's account
1 Locking vini's account
3 Locking vini's account
2 Locking cr7's account

*/
