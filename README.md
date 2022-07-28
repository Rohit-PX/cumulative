# Calculate cumulative vesting schedule

### This is a command line program that reads in a file of vesting events and outputs a vesting schedule to stdout.

####  Requires go 1.17 or higher 

#### Running the program
 * Create a path anywhere on your system: mkdir -p /root/git/go/src/github.com/
 * Unzip the folder into the above path
 * You should now have the path /root/git/go/src/github.com/cumulative
 * Set GOPATH as `export GOPATH=/root/git/go/`
 * `cd /root/git/go/src/github.com/cumulative`
 * The binary file `cumulative` is already present in the folder
 * The program takes 3 params:
 * 1. Path to the CSV file
 * 2. Date in YYYY-MM-DD format
 * 3. Precision in integer format for Quantity of vested shares 
 * You can now run the program as follows:
 ```
 	./cumulative <path-to-your-csv-file> <YYYY-MM--DD> <int>
 ```
 * Example:
 ```
 	./cumulative <path-to-your-csv-file> 2020-04-01 2
 ```
 * Unit tests can be run as follows:
```
   go test pkg/testing/vesting_test.go -v
``` 

#### Code has been divided into vesting and utils packages 
#### Vesting package has vesting related methods/APIs
#### Utils package as utility methods such as date conversion and validation
#### Testing package has unit tests to test APIs from vesting package
#### Some tests csv files have been provided in pkg/testing folder
 

