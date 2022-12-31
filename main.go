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
	"log"
	"math"
	"os"
	"sort"
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
	MinBalance    int
	MaxBalance    int
	EndingBalance int
}

// Constructor to set default values for a balance struct
func NewBalance() Balance {
	return Balance{
		MinBalance:    math.MaxInt64,  // set to max int64 value
		MaxBalance:    -math.MaxInt64, // set to min int64 value
		EndingBalance: 0,
	}
}

// Define a balances struct - map of balances for each month, indexed by month in int format
type Balances map[int]Balance

// Constructor to set default values for a balances struct
func NewBalances() Balances {
	return make(Balances)
}

// Define a user struct
type User struct {
	CustomerID   string
	Transactions []Transaction    // each item will be an individual transaction - multiple allowed per day, month, year
	YearBalances map[int]Balances // map where key is the year. Each year will hold a map of balances for each month
}

// Define a users struct - map of users, indexed by customerID

type Users struct {
	sync.RWMutex // for thread-safe access to the map
	UserMap      map[string]User
}

// Define an output struct
type Output struct {
	CustomerID    string
	Month         int
	Year          int
	MinBalance    int
	MaxBalance    int
	EndingBalance int
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

// FUNCTIONS

func readCSV(filePath string) *csv.Reader {
	// check if file exists
	file, err := os.Open(filePath)
	if err != nil {
		// if file does not exist, exit program
		os.Exit(1)
	}
	// open csv reader and return the pointer
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = 3 // set number of fields per record

	return csvReader
}

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

func validateLine(line []string) bool {
	// check for valid customerID
	if line[0] == "" {
		// fmt.Println("Error: Invalid customerID.")
		return false
	}
	// check for valid date
	if !validateDate(line[1]) {
		// fmt.Println("Error: Invalid date.")
		return false
	}
	// check for valid amount
	amount, err := strconv.Atoi(line[2])
	if err != nil {
		// fmt.Println("Error: Invalid amount.")
		return false
	}
	if amount < -1000000 || amount > 1000000 {
		// fmt.Println("Error: Invalid amount.")
		return false
	}
	return true
}

// opens a file and reads it into a list of transactions
func processTransactions(csvReader *csv.Reader) *[]Transaction {
	// create a list of transactions
	transactions := make([]Transaction, 0)
	for {
		// read line
		line, err := csvReader.Read()
		// if reached EOF or other error, break
		if err != nil {
			// fmt.Println("Error: ", err)
			break
		}
		// validate line
		if !validateLine(line) {
			// fmt.Println("Error: Invalid line.")
			continue
		}
		// get customerID, date
		customerID, date := line[0], line[1]
		// get amount
		amount, err := strconv.Atoi(line[2])
		if err != nil {
			// fmt.Println("Error: Invalid amount.")
			continue
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

// stores transactions in local storage
func storeTransactions(users *Users, transactions *[]Transaction) {
	for _, transaction := range *transactions {

		// get customerID
		custemerID := transaction.CustomerID

		// check if user exists in local storage
		user, ok := users.UserMap[custemerID]

		if ok { // if user exists, append transaction to user
			user.Transactions = append(user.Transactions, transaction) // update copy of user transactions
			users.UserMap[custemerID] = user                           // update user in local storage
		} else { // if user does not exist, create user and append transaction to user
			newUser := User{
				CustomerID:   custemerID,
				Transactions: []Transaction{transaction},
				YearBalances: make(map[int]Balances),
			}
			users.UserMap[custemerID] = newUser // update user in local storage with new user object
		}
	}

}

func storeBalances(users *Users) {
	// get transactions for each user
	for _, user := range users.UserMap {
		for _, transaction := range user.Transactions {
			// get month and year from date
			date_arr := strings.Split(transaction.Date, "/")

			month, err := strconv.Atoi(date_arr[0])
			if err != nil { // skip to next transaction if error
				continue
			}
			year, err := strconv.Atoi(date_arr[2])
			if err != nil { // skip to next transaction if error
				continue
			}

			// check if year exists in user's yearBalances map
			_, ok := user.YearBalances[year]
			if !ok { // if year does not exist, create new year
				user.YearBalances[year] = NewBalances()
			}
			balances := user.YearBalances[year]

			// check if month exists in user's yearBalances map
			_, ok = balances[month]
			if !ok { // if month does not exist, create new month
				balances[month] = NewBalance()
			}
			balance := balances[month]

			// update balance
			balance.EndingBalance += transaction.Amount
			// check if current balance is max or min balance
			if balance.EndingBalance > balance.MaxBalance {
				balance.MaxBalance = balance.EndingBalance
			}
			if balance.EndingBalance < balance.MinBalance {
				balance.MinBalance = balance.EndingBalance
			}

			// update user's yearBalances map
			user.YearBalances[year][month] = balance

			// update user in local storage
			users.UserMap[user.CustomerID] = user
		}
	}

}

func createCSV(fileName string) *os.File {
	// open file writer
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// retrieve balances from local storage, sorted by customerID
// write balances to CSV file
// want sorted by customerID, then year, then month
func writeCSV(file *os.File, users *Users) {
	// grab userIDs from local storage, then sort userIDs
	sortedCustomerIDs := make([]string, 0)
	for customerID, _ := range users.UserMap {
		sortedCustomerIDs = append(sortedCustomerIDs, customerID)
	}
	sort.Strings(sortedCustomerIDs)

	// iterate through sorted userIDs
	for _, customerID := range sortedCustomerIDs {

		// get user from local storage
		user := users.UserMap[customerID]

		// grab years from user's yearBalances map, then sort years
		sortedYears := make([]int, 0)
		for year, _ := range user.YearBalances {
			sortedYears = append(sortedYears, year)
		}
		sort.Ints(sortedYears)

		// iterate through sorted years
		for _, year := range sortedYears {

			// grab months from user's yearBalances map, then sort months
			sortedMonths := make([]int, 0)
			for month, _ := range user.YearBalances[year] {
				sortedMonths = append(sortedMonths, month)
			}
			sort.Ints(sortedMonths)

			// iterate through sorted months
			for _, month := range sortedMonths {

				// get balance from user's yearBalances map
				balance := user.YearBalances[year][month]

				// write to file
				_, err := file.WriteString(fmt.Sprintf("%v,%v/%v,%v,%v,%v\n", customerID, month, year, balance.MinBalance, balance.MaxBalance, balance.EndingBalance))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	defer file.Close()
}

func main() {

	// Create local storage
	var users = NewUsers()

	// Read CSV file
	csvReader := readCSV(os.Args[1]) // filepath is first argument

	// Process transactions
	transactions := processTransactions(csvReader)

	// Store transactions in local storage
	storeTransactions(users, transactions)

	// Calculate and store balances for each month, for each user
	storeBalances(users)

	// Create CSV file
	file := createCSV(os.Args[2])

	// Write list of strings to CSV file
	writeCSV(file, users)

}
