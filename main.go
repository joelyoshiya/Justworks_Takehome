// Written by: Joel Yoshiya Foster
// Email: joel.foster@gmail.com
// Date: 2022-12-29
// Description: This program reads a csv file containing transactions for a number of users, and outputs a csv file containing the minimum, maximum, and ending balance for each month for each user.

// Input CSV Format:
// `CustomerID, Date, Amount`

// Output CSV Format:
// `CustomerID, MM/YYYY, Min Balance, Max Balance, Ending Balance`

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// DEFAULT PATHS
var defaultInputFP = "data_raw_1.csv" // holds default input
var defaultOutputFP = "output.csv"    // holds default

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
// Since we are dealing with transactions, my choice would be an ACID compliant relational database, such as MySQL or PostgreSQL.
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
		log.Fatal(err)
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

// to validate that the customerID, date, and amount are valid
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
	_, err := strconv.Atoi(line[2])
	return err == nil
}

// to clean the line for whitespace
func cleanLine(line []string) []string {
	// clean each line for whitespace
	for i, item := range line {
		line[i] = strings.TrimSpace(item)
	}
	return line
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
		// clean lines for whitespace
		line = cleanLine(line)
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

// stores transactions with the pertinent user
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

// calculates and stores balances based on transactions for a single user
func storeBalances(users *Users) {
	// get transactions for each user
	for _, user := range users.UserMap {
		sort.Slice(user.Transactions, func(i, j int) bool {
			// if dates are equal, sort by amount
			if user.Transactions[i].Date == user.Transactions[j].Date {
				// start by calculating credits before debits for any given day
				return user.Transactions[i].Amount > user.Transactions[j].Amount
			}
			// transactions processed chronologically
			return user.Transactions[i].Date < user.Transactions[j].Date
		})
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

// generate a brand new CSV file with user-created name
func createCSV(fileName string) *os.File {
	// open file writer
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// retrieves balances from local storage, binds to Output struct, and writes to CSV file
func writeCSV(file *os.File, users *Users) {
	outputs := make([]Output, 0)
	// grab all balance data for each user, map fields to an Output struct
	for _, user := range users.UserMap {
		for year, balances := range user.YearBalances {
			for month, balance := range balances {
				output := Output{
					CustomerID:    user.CustomerID,
					Year:          year,
					Month:         month,
					MinBalance:    balance.MinBalance,
					MaxBalance:    balance.MaxBalance,
					EndingBalance: balance.EndingBalance,
				}
				outputs = append(outputs, output)
			}
		}
	}
	// sorting based on customerID, then Year, and then Month, sort Output structs
	sort.Slice(outputs, func(i, j int) bool {
		if outputs[i].CustomerID == outputs[j].CustomerID {
			if outputs[i].Year == outputs[j].Year {
				return outputs[i].Month < outputs[j].Month
			}
			return outputs[i].Year < outputs[j].Year
		}
		return outputs[i].CustomerID < outputs[j].CustomerID
	})
	// write sorted Output structs to file
	for _, output := range outputs {
		_, err := file.WriteString(fmt.Sprintf("%v,%v/%v,%v,%v,%v\n", output.CustomerID, output.Month, output.Year, output.MinBalance, output.MaxBalance, output.EndingBalance))
		if err != nil {
			log.Fatal(err)
		}
	}
	defer file.Close()
}

func main() {
	var input, output string
	// Check for correct number of arguments
	if len(os.Args) != 3 {
		input = defaultInputFP
		output = defaultOutputFP
	} else {
		input = os.Args[1]
		output = os.Args[2]
	}

	// Create local storage
	var users = NewUsers()

	// Read CSV file
	csvReader := readCSV("input/" + input) // input filepath is first argument

	// Process transactions
	transactions := processTransactions(csvReader)

	// Store transactions in local storage
	storeTransactions(users, transactions)

	// Calculate and store balances for each month, for each user
	storeBalances(users)

	// Create CSV file
	file := createCSV("output/" + output) // output filepath is second argument

	// Write list of strings to CSV file
	writeCSV(file, users)

	// do unix cat on output file
	cmd := exec.Command("cat", "output/"+output)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
