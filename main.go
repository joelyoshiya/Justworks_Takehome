// Written by: Joel Yoshiya Foster
// Email: joel.foster@gmail.com
// Date: 2022-12-29
// Description: This program reads a csv file containing transactions for a number of users, and outputs a csv file containing the minimum, maximum, and ending balance for each month for each user.

// Input CSV Format:
// `CustomerID, Date, Amount`

// Output CSV Format:
// `CustomerID, MM/YYYY, Min Balance, Max Balance, Ending Balance`

// APPRAOCH
// Have a filereader that reads the csv file and parses the data into a list of **transactions**.
// Customer IDs will be determined by the `CustomerID` column of the input csv file.
// Then, have a function that takes in a list of transactions and returns a list of **balances**.
// The function will iterate through the list of transactions and calculate the balance for each month, for each user.
// Balance will include the minimum balance, maximum balance, and ending balance for each month.
// The function will return a list of balances, pertaining to each month, for each user.
// *In the case that we are returning multiple months of balances for each user, we will return the balance items first in order of customer, then in order of month, by ascending order of both `CustomerID` followed by `MM/YYYY`.*
// Then, have a function that takes in a list of balances and returns a list of strings that can be written to a csv file. The function will iterate through the list of balances and create a string for each balance. The function will return a list of strings.
// Finally, have a filewriter that takes in a list of strings and writes them to a csv file. Output the

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// STRUCTS AND TYPES
// Input csv is mapped to a list of transactions
// transactions are each tied a user - custemerID is the unique identifier
// balances are each tied to a user - customerID is the unique identifier

// Define a transaction struct
type Transaction struct {
	CustomerID string
	Date       string
	Amount     int
}

// Define a balance struct
type Balance struct {
	CustomerID    string
	MonthYear     string
	MinBalance    int
	MaxBalance    int
	EndingBalance int
}

// Define a user struct
type User struct {
	CustomerID   string
	Transactions []Transaction // each item will be an individual transaction - multiple allowed per day, month, year
	Balances     []Balance     // each item will be a balance for a month
}

// Define a users struct

type Users struct {
	UserMap map[string]User
}

// Define a local storage for users
// In a production environment, this would most likely be a persistent storage solution, such as a relational database.
// Since we are dealing with transactions, my choice would be a relational database, such as MySQL or PostgreSQL.
// However, for the sake of this exercise, we will use a local storage solution.
// This will be a map of users, where the key is the CustomerID, and the value is the user struct.

// define a constructor for a pointer to a users struct
// allows passing of Users struct to other components, if needed down the line - See Referral 1
func NewUsers() *Users {
	return &Users{
		UserMap: make(map[string]User),
	}
}

// our local storage solution
var users = NewUsers()

// FUNCTIONS

// opens a file and reads it into a list of transactions
func readCSV() *[]Transaction {
	// grab csv file path from command line
	filePath := os.Args[1]
	// check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		// if file does not exist, exit program
		os.Exit(1)
	}
	// instantiate list of transactions
	var transactions []Transaction
	// initialize list of transactions
	transactions = make([]Transaction, 0)
	// use a buffered reader
	input := bufio.NewScanner(file)
	for input.Scan() {
		// parse line into transaction
		// split line into list of strings
		line := input.Text()
		lineList := strings.Split(line, ",")
		// parse list of strings into transaction
		// TODO: add error handling for invalid input
		// parse customerID
		customerID := lineList[0]
		// parse date
		date := lineList[1]
		// parse amount
		amount, err := strconv.Atoi(lineList[2])
		if err != nil {
			// if amount is not an integer, log error to stdout
			fmt.Printf("Error: Amount is not an integer. Error: %v", err)
		}
		// create transaction
		transaction := Transaction{
			CustomerID: customerID,
			Date:       date,
			Amount:     amount,
		}
		// append transaction to list of transactions
		transactions = append(transactions, transaction)
	}
	// return list of transactions
	return &transactions
}

func storeTransactions(transactions *[]Transaction) {
	// takes a pointer to a list of transactions
	// iterates through list of transactions
	for _, transaction := range *transactions {
		// check if user exists in local storage
		custemerID := transaction.CustomerID
		if user, ok := users.UserMap[custemerID]; ok {
			// if user exists, append transaction to user
			user.Transactions = append(user.Transactions, transaction) // update copy of user transactions
			users.UserMap[custemerID] = user                           // update user in local storage
		} else {
			// if user does not exist, create user and append transaction to user
			user := User{
				CustomerID:   custemerID,
				Transactions: []Transaction{transaction},
			}
			users.UserMap[custemerID] = user
		}
	}

}

func calculateBalances() {

}

func storeBalances() {

}

func writeCSV() {

}

func main() {

	// TODO - should I initialize the users struct here?
	// Read CSV file
	transactions := readCSV()
	// print transactions
	for _, transaction := range *transactions {
		fmt.Printf("%v\n", transaction)
	}
	println(strconv.FormatInt(int64(len(*transactions)), 10))
	// Store transactions in local storage
	storeTransactions(transactions)
	// Calculate balances for each month, for each user
	// Create list of strings to write to CSV file
	// Write list of strings to CSV file

}
