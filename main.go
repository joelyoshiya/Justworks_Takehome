// Written by: Joel Yoshiya Foster
// Email: joel.foster@gmail.com
// Date: 2022-12-29
// Description: This program reads a csv file containing transactions for a number of users, and outputs a csv file containing the minimum, maximum, and ending balance for each month for each user.

// Input CSV Format:
// `CustomerID, Date, Amount`

// Output CSV Format:
// `CustomerID, MM/YYYY, Min Balance, Max Balance, Ending Balance`

// APPROACH
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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
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
	Month         string
	MinBalance    int
	MaxBalance    int
	EndingBalance int
}

// Define a balances struct - a fixed size of array for each month
type Balances [12]Balance

// Define a user struct
type User struct {
	CustomerID   string
	Transactions []Transaction    // each item will be an individual transaction - multiple allowed per day, month, year
	YearBalances map[int]Balances // map where key is the year. Each year will hold slice of balance structs for each month
}

// Define a users struct

type Users struct {
	sync.RWMutex // for thread-safe access to the map
	UserMap      map[string]User
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

// our local storage solution - defined package-wide
var users = NewUsers()

// FUNCTIONS

// to validate the date format
func validateDate(date string) bool {
	date_arr := strings.Split(date, "/")
	if len(date_arr) != 3 {
		// fmt.Println("Error: Invalid date format.")
		return false
	}
	// check for valid month
	month, err := strconv.Atoi(date_arr[0])
	if err != nil {
		// fmt.Println("Error: Invalid month.")
		return false
	}
	if month < 1 || month > 12 {
		// fmt.Println("Error: Invalid month.")
		return false
	}
	// check for valid day
	day, err := strconv.Atoi(date_arr[1])
	if err != nil {
		// fmt.Println("Error: Invalid day.")
		return false
	}
	if day < 1 || day > 31 {
		// fmt.Println("Error: Invalid day.")
		return false
	}
	// check for valid year
	year, err := strconv.Atoi(date_arr[2])
	if err != nil {
		// fmt.Println("Error: Invalid year.")
		return false
	}
	if year < 1900 || year > 2050 {
		// fmt.Println("Error: Invalid year.")
		return false
	}
	return true
}

// opens a file and reads it into a list of transactions
func readCSV(filePath string) *[]Transaction {
	// check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		// if file does not exist, exit program
		os.Exit(1)
	}
	// initialize list of transactions
	transactions := make([]Transaction, 0)
	// initialize a csv reader
	csvReader := csv.NewReader(file)
	csvReader.TrimLeadingSpace = true // trim leading spaces
	csvReader.FieldsPerRecord = 3     // set number of fields per record
	// iterate through file
	for {
		// read line
		line, err := csvReader.Read()
		// if reached EOF or other error, break
		if err != nil {
			// fmt.Println("Error: ", err)
			break
		}
		// TODO: add error handling for invalid input - consider REGEX for cleaning
		// if invalid number of arguments, skip to next line
		if len(line) != 3 {
			// fmt.Println("Error: Invalid number of arguments.")
			continue
		}
		if line[0] == "" || line[1] == "" || line[2] == "" {
			// fmt.Println("Error: Invalid input.")
			continue
		}
		customerID := strings.TrimSpace(line[0])
		// parse date
		date := strings.TrimSpace(line[1])
		// check for valid date format
		if !validateDate(date) {
			// fmt.Println("Error: Invalid date format.")
			continue
		}
		// parse amount
		amount, err := strconv.Atoi(strings.TrimSpace(line[2]))
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

		// read lock
		users.RLock()
		user, ok := users.UserMap[custemerID]
		users.RUnlock()

		if ok { // if user exists, append transaction to user
			user.Transactions = append(user.Transactions, transaction) // update copy of user transactions
			// write lock
			users.Lock()
			users.UserMap[custemerID] = user // update user in local storage
			users.Unlock()
		} else { // if user does not exist, create user and append transaction to user
			newUser := User{
				CustomerID:   custemerID,
				Transactions: []Transaction{transaction},
				YearBalances: make(map[int]Balances),
			}
			// write lock
			users.Lock()
			users.UserMap[custemerID] = newUser // update user in local storage with new user object
			users.Unlock()
		}
	}

}

func calculateBalances() {
	// get transactions for each user
	for customerID, user := range users.UserMap {
		for _, transaction := range user.Transactions {
			// get month and year from date
			date := transaction.Date
			date_arr := strings.Split(date, "/")
			month, err := strconv.Atoi(date_arr[0])
			year, err := strconv.Atoi(date_arr[2])
			// check if year exists in user's yearBalances map
			yearBalances, ok := user.YearBalances[year]
			if ok { // if year exists, check if month exists in yearBalances
				monthBalance := yearBalances[month]
			}

			// add transaction amount to balance for that month
			// if balance is less than min balance, update min balance
			// if balance is greater than max balance, update max balance
			continue
		}

		// iterate through transactions
		// for each transaction, check the date
		// use a map to store the balances for each month
		// keep track of the min and max balances as updates are made
	}

}

func storeBalances() {
	// run calculateBalances()
	// store balances in local storage

}

func writeCSV() {
	// retrieve balances from local storage, sorted by customerID
	// write balances to CSV file

}

func main() {

	// TODO - should I initialize the users struct here?

	// Read CSV file
	transactions := readCSV(os.Args[1])

	// Store transactions in local storage
	storeTransactions(transactions)

	// print users and their transactions
	for customerID, user := range users.UserMap {
		// print customerID
		fmt.Printf("CustomerID: %v\n", customerID)
		// print transactions
		for _, transaction := range user.Transactions {
			fmt.Printf("\t%v\n", transaction)
		}
	}
	// print number of users
	println("Number of users: " + strconv.FormatInt(int64(len(users.UserMap)), 10))

	// Calculate balances for each month, for each user

	// Write list of strings to CSV file

}
