package main

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

// define a constructor for the users struct
func NewUsers() *Users {
	return &Users{
		UserMap: make(map[string]User),
	}
}

// our local storage solution
var users = NewUsers()

// FUNCTIONS

func readCSV() {
	// Read CSV file
	// Parse CSV file into list of transactions
	// Return list of transactions
}

func storeTransactions() {
	// takes a list of transactions, and ties them to internal user structs
}

func storeBalances() {

}

func main() {
	// Read CSV file
	// Parse CSV file into list of transactions
	// Calculate balances for each month, for each user
	// Create list of strings to write to CSV file
	// Write list of strings to CSV file

}
