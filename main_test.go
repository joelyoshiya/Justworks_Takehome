package main

import "testing"

func Test_ReadCSV(t *testing.T) {
	transactions := readCSV("data_raw_1.csv")
	if len((*transactions)) != 90 {
		t.Errorf("Expected 90 transactions, got %v", len((*transactions)))
	}
}

func Test_StoreTransactions(t *testing.T) {
	transactions := readCSV("data_raw_1.csv")
	storeTransactions(transactions)
	if len(users.UserMap) != 3 {
		t.Errorf("Expected 3 users, got %v", len(users.UserMap))
	}
}

// func Test_InputBalances(t *testing.T) {

// }
