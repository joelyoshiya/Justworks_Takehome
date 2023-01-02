# Justworks Takehome Assignment

For the Justworks take-home project

## Problem

Justworks wants to generate insight from a list of banking transactions occurring in customer accounts. We want to generate **minimum , maximum and ending balances** by month for all customers. You can assume starting balance at begining of month is 0.

You should read transaction data from csv files and produce output in the format mentioned below.

You can assume negative numbers as debit and positive numbers as credit.

- Please apply credit transactions first to calculate balance on a given day.  
- Please write clear instructions on how to run your program on a local machine.
- Please use dataset in Data Tab to test your program.
- You do not need to add Column Headers in the output.
  - Please assume the input file does not have header row.

This is a command line program that takes in a csv file as input and outputs a csv file.

Input CSV Format:
`CustomerID, Date, Amount`

Output CSV Format:
`CustomerID, MM/YYYY, Min Balance, Max Balance, Ending Balance`

## Approach

Have a filereader that reads the csv file and parses the data into a list of **transactions**. Then, have a function that takes in a list of transactions and returns a list of **balances**. The function will iterate through the list of transactions and calculate the balance for each month, for each user. The function will return a list of balances, pertaining to each month, for each user.

*In the case that we are returning multiple months of balances for each user, we will return the balance items grouped by customer, and then in order of month, by ascending order of month. We don't guarantee any particular ordering of customers.*

Then, have a function that takes in a list of balances and returns a list of strings that can be written to a csv file. The function will iterate through the list of balances and create a string for each balance. The function will return a list of strings.

Finally, have a filewriter that takes in a list of strings and writes them to a csv file. Output the file.

### Handling edge cases

## Technology

- Go : go1.19.4 darwin/arm64
- Docker: 20.10.21

## How to run

blah

## How to test

blah

## Discussion

// TODO
// Reconsider separation of transactions and balance logic
// We want to imitate a live scenario -> as soon as we have a transaction, we want to calculate the balance for that month
// As new transactions role in for that monthh -> we want to update the balance for that month
// After all transactions have been processed for that month -> we will output balances for each month for which there is at least one transaction
// First step to doing this: make operations atomic

// TODO consider using a Time.time object for the date field
## Conclusion

blah.

## Submission

Tentatively, I plan to submit a zip file containing the following:

- `README.md` - this file
- `go.mod` - the go module file
- `main.go` - the main program
- `main_test.go` - the test file
- `testdata` - a folder containing the test data
  - `testdata/input.csv` - the input csv file
  - `testdata/output.csv` - the output csv file

I also plan to leave a link to the github repo in the notes section of the submission form.

## References

- [The Go Init Function - TutorialEdge.net](https://tutorialedge.net/golang/the-go-init-function/)
