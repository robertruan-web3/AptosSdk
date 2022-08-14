package main

import (
	"AptosSdk/pkg/aptos"
	"fmt"
)

func main() {
	// two seeds for two accounts
	seedA := "abcd1234abcd1234abcd1234abcd1234"
	seedB := "abcd5678abcd5678abcd5678abcd5678"

	accountA, err := aptos.NewAccount(seedA)
	if err != nil {
		panic("generator new address by seed failed:%s" + err.Error())
	}

	accountB, err := aptos.NewAccount(seedB)
	if err != nil {
		panic("generator new address by seed failed:%s" + err.Error())
	}

	// print address of account A and B
	{
		fmt.Println("public address & private address.")
		fmt.Println(accountA.PublicKey(), accountA.PublicAddress())
		fmt.Println(accountB.PublicKey(), accountB.PublicAddress())
	}

	if false {
		aptos.FoundAccount(accountA, 100000)
		aptos.FoundAccount(accountB, 0)
	}

	// get account balance
	if true {
		balanceA, err := aptos.AccountGetBalance(accountA)
		if err != nil {
			panic("A AccountGetBalance failed:" + err.Error())
		}

		balanceB, err := aptos.AccountGetBalance(accountB)
		if err != nil {
			panic("A AccountGetBalance failed:" + err.Error())
		}

		fmt.Printf("before transaction, balance A: %d, balance B: %d\n", balanceA, balanceB)
	}

	// send coin from account A to B
	if true {
		aptos.Transfer(accountA, accountB.PublicAddress(), 10000)
	}

	// get account balance after transfer
	if true {
		balanceA, err := aptos.AccountGetBalance(accountA)
		if err != nil {
			panic("A AccountGetBalance failed:" + err.Error())
		}

		balanceB, err := aptos.AccountGetBalance(accountB)
		if err != nil {
			panic("A AccountGetBalance failed:" + err.Error())
		}

		fmt.Printf("after transaction, balance A: %d, balance B: %d\n", balanceA, balanceB)
	}

	fmt.Println("\n=== end of main ===")
}
