# Justworks Takehome Assignment

For the Justworks take-home project

## Problem

Justworks wants to generate insight from a list of banking transactions occurring in customer accounts. We want to generate **minimum , maximum and ending balances** by month for all customers. You can assume starting balance at begining of month is 0. You should read transaction data from csv files and produce output in the format mentioned below. You can assume negative numbers as debit and positive numbers as credit

Please apply credit transactions first to calculate balance on a given day.  Please write clear instructions on how to run your program on a local machine. Please use dataset in Data Tab to test your program. You do not need to add Column Headers in the output. Please assume the input file does not have header row.

Input CSV Format:
`CustomerID, Date, Amount`

Output CSV Format:
`CustomerID, MM/YYYY, Min Balance, Max Balance, Ending Balance`

## Approach

Have a filereader that reads the csv file and parses the data into a list of **transactions**. Then, have a function that takes in a list of transactions and returns a list of **balances**. The function will iterate through the list of transactions and calculate the balance for each month, for each user. The function will return a list of balances, pertaining to each month, for each user.

*In the case that we are returning multiple months of balances for each user, we will return the balance items first in order of customer, then in order of month, by ascending order of both `CustomerID` and `MM/YYYY`.*

Then, have a function that takes in a list of balances and returns a list of strings that can be written to a csv file. The function will iterate through the list of balances and create a string for each balance. The function will return a list of strings.


Finally, have a filewriter that takes in a list of strings and writes them to a csv file. Output the

## How to run

## How to test

## Conclusion