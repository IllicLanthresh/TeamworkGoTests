// Package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

// read chunks of fixed size in goroutines
// insert as they read
// pseudocode of goroutine:
// for each line, try to get the email address(regex), detect any csv errors and any email format errors.
// after getting email, extract the domain of it.
// REMEMBER TO USE MUTEX
// if the list is empty just insert it and forget
// check if domain is already in the list (maybe use an array to store visited?)
// if its in the list just add to the counter
// if it's not in the list walk in order to find the best position and insert a new struct
// struct should have: domain, count and mutex (maybe?)
// list should be array? (maybe there's best approaches with array like data structures best suited for indexing like sets?)
// Repeated elements have no special treatment, they count as one more

// Another approach is to use radix buckets, to investigate later
