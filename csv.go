package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func initCsv() (*os.File, error) {
	header := []string{"storeID", "storeName"}

	f, err := os.Create("store.csv")
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write(header); err != nil {
		return nil, err
	}

	return f, err
}

func getFile(csvName string) (*os.File, error) {
	file, err := os.OpenFile(csvName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if errors.Is(err, os.ErrNotExist) {
		return initCsv()
	}
	return file, err
}

func writeStoresIDs(storeIds []string) error {
	fmt.Printf("writing %d storesIds\n",len(storeIds))
	f, err := getFile("store.csv")
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	for _, id := range storeIds {
		w.Write([]string{id})
	}
	w.Flush()
	return nil
}

func writeCsv(records [][]string, csvName string)error{
	f, err := getFile(csvName)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.WriteAll(records)
	w.Flush()
	return nil
}
func readFromCsv(csvName string) ([][]string,error) {
	csvFile, err := os.Open(csvName)
	if err != nil {
		return nil,err
	}
	
	defer csvFile.Close()
    
    return csv.NewReader(csvFile).ReadAll()
}
