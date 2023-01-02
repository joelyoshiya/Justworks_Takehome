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

## Assumptions

- Credits applied before debits for any given day (confirmed in FAQ)
- Maximum and Minimum balances are calculated based on the most recent balance as it pertains to one transaction. This is instead of the max/min being calculated based on the ending balance on a given day.
  - Confirmed via line 99 in FAQ
- No need to check for parsible yet invalid dates, such as 2021-02-31
  - Confirmed via line 136 in FAQ
- No need to check for exceeding amount, as long as representable as an integer (int64)
- EndingBalance relates to the balance at the end of the month, not the balance at the end of the day
  - in other words, balance starting from 0 after all debits and credits have been applied for the month

## Approach

Have a filereader that reads the csv file and parses the data into a list of **transactions**. Then, have a function that takes in a list of transactions and returns a list of **balances**. The function will iterate through the list of transactions and calculate the balance for each month, for each user. The function will return a list of balances, pertaining to each month, for each user.

*In the case that we are returning multiple months of balances for each user, we will return the balance items grouped by customer, and then in order of month, by ascending order of month. We don't guarantee any particular ordering of customers.*

Then, have a function that takes in a list of balances and returns a list of strings that can be written to a csv file. The function will iterate through the list of balances and create a string for each balance. The function will return a list of strings.

Finally, have a filewriter that takes in a list of strings and writes them to a csv file. Output the file.

### Handling edge cases

I've used the following approach to handle edge cases:

- clean up the input data
  - remove any empty lines
  - remove any lines that don't have 3 columns
  - remove any lines that have invalid dates
  - remove any lines that have invalid amounts
- handle any invalid data
  - if there are no valid lines, return an empty list of balances
  - if there are no valid lines, return an empty list of strings
  - if there are no valid lines, return an empty file

## Technology

- Go : go1.19.4 darwin/arm64
- Docker: 20.10.21

## How to run

If you have the specified version of Go installed, you can run the program locally. Otherwise, you can run the program in a docker container. If you wish to use the default input file (`data_raw_1.csv`, a replica of data in the Excel file), you can omit the input and output file names from either the `go run` or `docker run` commands.

### Locally

1. Clone the repo
2. Insert your input csv file into the `input` folder
3. Run `go run main.go [input_file_name] [output_file_name]` in the root directory of the project.
4. The output file will be in the `output` folder

### Docker

1. Clone the repo
2. Insert your input csv file into the `input` folder
3. Run `docker build -t justworks .` in the root directory of the project.
4. Run `docker run -it justworks [input_file_name] [output_file_name]` in the root directory of the project.
5. The output file will be in the `output` folder

## How to test

Run `go test` in the root directory of the project.

## Discussion

If I had time, here are the refactors I would do:

- Make transaction/balance processing per transaction/day instead of all at once from input
  - Make the code more modular
- Imitate a live scenario where new transactions instantly update balance live
- Use a Time.time object for the date field

## Conclusion

I enjoyed the exercise and how it challenged me to consider edge cases, and how to design objects and methods in an intuitive way.

## Submission

Tentatively, I plan to submit a zip file containing the following:

- `README.md` - this file
- `go.mod` - the go module file
- `main.go` - the main program
- `main_test.go` - the test file
- `input` - a folder containing the input csv file
- `output` - a folder containing the output csv file
- `testdata` - a folder containing the test data (used in `main_test.go`)
  - `testdata/input.csv` - the input csv file
  - `testdata/output.csv` - the output csv file
- Dockerfile - the dockerfile used to build the docker image

I also plan to leave a link to the github repo in the notes section of the submission form.

## References

- [The Go Init Function - TutorialEdge.net](https://tutorialedge.net/golang/the-go-init-function/)
- Go stdlib documentation
