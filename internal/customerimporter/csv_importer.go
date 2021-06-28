// Package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each emailDomain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package customerimporter

import (
	"encoding/csv"
	"fmt"
	"github.com/IllicLanthresh/TeamworkGoTests/pkg/helperTypes"
	"github.com/IllicLanthresh/TeamworkGoTests/pkg/radixSorter"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// customerimporter just a demo of how this module could hold different implementations for different email sources
// like: databases, APIs...
type customerImporter interface {
	// emailAddressesGenerator is a generator like closure that feeds from customer records
	// and returns a channel of emailAddress objects,
	emailAddressesGenerator() (emailAddresses chan emailAddress, err error)
	// CustomerCountByDomain feeds from emailAddressesGenerator and outputs the count of customers for each email domain
	CustomerCountByDomain() (customerCountByDomain map[string]int)
}

// csvCustomerImporter is the CSV implementation of a customerImporter, it recieves a csv file path and an email key
// representing the header name for the email column
type csvCustomerImporter struct {
	csvPath  string
	emailKey string
}

// NewCsvCustomerImporter constructor for csvCustomerImporter
func NewCsvCustomerImporter(csvPath string, emailKey string) (*csvCustomerImporter, error) {
	if csvPath == "" || len(csvPath) < 4 || !strings.HasSuffix(csvPath, ".csv") {
		return nil, CsvPathInvalidError{path: csvPath}
	}
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("path %s does not exist", csvPath)
	}
	if emailKey == "" {
		return nil, MissingEmailKey{}
	}
	importer := csvCustomerImporter{
		csvPath:  csvPath,
		emailKey: emailKey,
	}
	return &importer, nil
}

type emailAddress struct {
	Address string
	Err     error
}

// emailAddressesGenerator takes a path `csvPath` to a csv file with customer records,
// the csv file is supposed to have a headers row and `emailKey` is the name of the row holding email addresses.
// Any email addresses not complying with RFC 5322 will be ignored
func (imp *csvCustomerImporter) emailAddressesGenerator() (emailAddresses chan emailAddress, err error) {
	regex := regexp.MustCompile(emailRegex)

	fileReader, err := os.Open(imp.csvPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %s: %w", imp.csvPath, err)
	}
	csvReader := csv.NewReader(fileReader)
	headers, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("couldn't read headers row on CSV file: %w", err)
	}

	emailIndex := helperTypes.StringSlice(headers).IndexOf(imp.emailKey)
	if emailIndex == -1 {
		return nil, KeyNotFoundError{
			key:   imp.emailKey,
			slice: headers,
		}
	}

	emailAddresses = make(chan emailAddress)

	go func() {
		defer func() {
			err = fileReader.Close()
			if err != nil {
				log.Printf("error trying to close file: %s", err)
			}
		}()
		defer close(emailAddresses)

		for {
			row, err := csvReader.Read()
			if err != nil {
				switch err {
				case io.EOF:
					return
				// We could explore adding a case here for ErrFieldCount, maybe there's a missing field but the email is still there?
				default:
					emailAddresses <- emailAddress{Address: "", Err: err}
					continue
				}
			}
			if len(row) == 0 {
				// csv.Reader does not error out with ErrFieldCount when the row is empty
				continue
			}
			if !regex.MatchString(row[emailIndex]) {
				log.Printf("ignoring row %v, email address not compliant with RFC 5322", row)
				continue
			}
			emailAddresses <- emailAddress{Address: row[emailIndex], Err: nil}
		}
	}()
	return emailAddresses, nil
}

type emailDomain struct {
	Domain        string
	CustomerCount int
}

//CustomerCountByDomain outputs the count of customers for each email domain in the csv file you introduced in the constructor
func (imp *csvCustomerImporter) CustomerCountByDomain() (sortedDomains []emailDomain, err error) {
	emailAddresses, err := imp.emailAddressesGenerator()
	if err != nil {
		return nil, fmt.Errorf("couldn't create address generator: %w", err)
	}

	customerCountByDomain := make(map[string]int)
	sorter := radixSorter.NewRadixSorter()

	for address := range emailAddresses {
		if address.Err != nil {
			log.Printf("couldn't process row: %s", address.Err)
			continue
		}
		splitAddress := strings.Split(address.Address, "@")
		domain := strings.ToLower(splitAddress[len(splitAddress)-1])

		if _, exists := customerCountByDomain[domain]; !exists {
			sorter.Add(domain)
		}
		customerCountByDomain[domain] += 1
	}

	for _, domain := range sorter.Sort() {
		sortedDomains = append(sortedDomains, emailDomain{
			Domain:        domain,
			CustomerCount: customerCountByDomain[domain],
		})
	}
	return
}
