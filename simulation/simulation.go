package simulation

import (
	"encoding/csv"
	"errors"
	"os"
	"sync"
)

type Storable interface {
	// first return -> Value names ; second return -> Values
	GetValues() ([]string, []string)
}

var file_access sync.Mutex

// StoreResults writes the results into the file <filename>.
// The data is drawn from the list of storable interfaces.
// During the simulation, the order of <objList> should be
// respected. Store Results is a thread safe function.
func StoreResults(dirName, fileName string, objList ...Storable) {
	file_access.Lock()
	defer file_access.Unlock()
	valueNames, values := make([]string, 0), make([]string, 0)
	for _, sto := range objList {
		nam, val := sto.GetValues()
		valueNames = append(valueNames, nam...)
		values = append(values, val...)
	}
	// Checking/Creating directory
	if err := CheckCreateDir(dirName); err != nil {
		panic("Error encountered: Directory could not be created.")
	}

	// Checking/Creating csv file
	if exists, err := FileExists(dirName, fileName); err == nil {
		if !exists {
			//Create file and write the first line
			f, errr := os.OpenFile(dirName+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if errr != nil {
				panic("Error encountered: CSV File could not be created.")
			}
			w := csv.NewWriter(f)
			w.Write(valueNames)
			w.Flush()
			if errrr := w.Error(); errrr != nil {
				panic("Error encountered: CSV File writing produced an error")
			}
		}
	} else {
		panic("Error encountered: CSV File could not be created/read.")
	}

	// Writing data to the file

	f, errr := os.OpenFile(dirName+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errr != nil {
		panic("Error encountered: CSV File could not be created.")
	}
	w := csv.NewWriter(f)
	w.Write(values)
	w.Flush()
	if err := w.Error(); err != nil {
		panic("Error encountered: CSV File writing produced an error")
	}
}

// Checks if given directory exists. If it doesn't it is created
func CheckCreateDir(dirName string) error {
	if _, err := os.Stat("dirName"); err != nil {
		if os.IsNotExist(err) {
			// file does not exist; Create file
			os.Mkdir(dirName, os.ModePerm)
		} else {
			return errors.New("Permission error; Probably")
		}
	}
	return nil
}

// Check if file exists, if it doesn't exist it creates it and returns false
// otherwise it resturns true
func FileExists(dirName, fileName string) (bool, error) {
	if _, err := os.Stat(dirName + "/" + fileName); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return true, nil
}
