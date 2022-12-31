package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

// note: each test is run in a separate goroutine, and memory is shared between tests

// expected balances
// C231,11/2022,616,43444,40000
// C865,11/2022,-2441,40013,21337
// C409,11/2022,-47654,-1878,-46118

// hold the expected balances
var expectedBalances = map[string]map[int]map[int]map[string]int{
	"C231": {
		2022: {
			11: {
				"MinBalance":    616,
				"MaxBalance":    43444,
				"EndingBalance": 40000,
			},
		},
	},
	"C865": {
		2022: {
			11: {
				"MinBalance":    -2441,
				"MaxBalance":    40013,
				"EndingBalance": 21337,
			},
		},
	},
	"C409": {
		2022: {
			11: {
				"MinBalance":    -47654,
				"MaxBalance":    -1878,
				"EndingBalance": -46118,
			},
		},
	},
}

// GOLDEN PATH TESTS
func Test_ReadCSV(t *testing.T) {
	transactions := processTransactions(readCSV("data_raw_1.csv"))
	if len((*transactions)) != 90 {
		t.Errorf("Expected 90 transactions, got %v", len((*transactions)))
	}
}

func Test_StoreTransactions(t *testing.T) {
	users := NewUsers()
	transactions := processTransactions(readCSV("data_raw_1.csv"))
	storeTransactions(users, transactions) // run only once (memory is shared between tests)
	if len(users.UserMap) != 3 {
		t.Errorf("Expected 3 users, got %v", len(users.UserMap))
	}
	for customerID, user := range users.UserMap {
		// check that length of transactions is greater than 0 for each user
		if len(user.Transactions) == 0 {
			t.Errorf("Expected transactions for user %v, got %v", customerID, len(user.Transactions))
		}
	}
}

func Test_CalculateBalances(t *testing.T) {
	users := NewUsers()
	transactions := processTransactions(readCSV("data_raw_1.csv"))
	storeTransactions(users, transactions)
	storeBalances(users)

	for customerID, user := range users.UserMap {
		// iterate through each year of balances
		for year, balances := range user.YearBalances {
			// iterate through each present month of balances
			for month, balance := range balances {
				if balance.MinBalance != expectedBalances[customerID][year][month]["MinBalance"] {
					t.Errorf("Expected MinBalance %v, got %v", expectedBalances[customerID][year][month]["MinBalance"], balance.MinBalance)
				}
				if balance.MaxBalance != expectedBalances[customerID][year][month]["MaxBalance"] {
					t.Errorf("Expected MaxBalance %v, got %v", expectedBalances[customerID][year][month]["MaxBalance"], balance.MaxBalance)
				}
				if balance.EndingBalance != expectedBalances[customerID][year][month]["EndingBalance"] {
					t.Errorf("Expected EndingBalance %v, got %v", expectedBalances[customerID][year][month]["EndingBalance"], balance.EndingBalance)
				}
			}
		}
	}

}

func Test_WriteCSV(t *testing.T) {
	users := NewUsers()
	transactions := processTransactions(readCSV("data_raw_1.csv"))
	storeTransactions(users, transactions)
	storeBalances(users)
	writeCSV(users)
	var expectedFile = "Output_Data.csv"
	var actualFile = "balances.csv"

	// check that file contents are identical
	expected, err := ioutil.ReadFile(expectedFile)
	if err != nil {
		t.Errorf("Error reading expected file: %v", err)
	}
	actual, err := ioutil.ReadFile(actualFile)
	if err != nil {
		t.Errorf("Error reading actual file: %v", err)
	}
	if !bytes.Equal(expected, actual) {
		t.Errorf("Expected file contents to be identical, got %v", string(actual))
	}

}

// BAD INPUT TESTS
