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

// STRUCTS
// Define a transaction struct
// Define a balance struct
// Define a user struct

func main() {
	// Read CSV file
	// Parse CSV file into list of transactions
	// Calculate balances for each month, for each user
	// Create list of strings to write to CSV file
	// Write list of strings to CSV file

}
