package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Println("need at least one arg list-stores || info-stores || info-societes")
		os.Exit(0)
	}
	arg := os.Args[1]
	log.Printf("args are %s\n", os.Args)

	if arg == "list-stores" || arg == "ls" {
		getStoresIdsAndWrite()
	} else if arg == "info-stores" || arg == "is" {
		getStoresInfosAndWrite()
	} else if arg == "info-societes" || arg == "ic" {
		getInfoCompanyAndWrite()
	}
}

func getStoresIdsAndWrite() {
	storesIDs, err := getListStores()
	if err != nil {
		log.Fatal("can't get storesIds")
	}
	writeStoresIDs(storesIDs)
}

func getStoresInfosAndWrite() {
	csvLines, err := readFromCsv("store.csv")
	if err != nil {
		log.Fatalf("can't read csv=> %s", err)
	}

	for index, l := range csvLines {
		storeID := l[0]
		storeInfo, err := getStoresInfo(storeID)
		if err != nil {
			log.Fatal(err)
		}
		l = append(l, storeInfo.Name)
		l = append(l, storeInfo.CompanyName)
		l = append(l, storeInfo.Coordinates)
		l = append(l, storeInfo.RestaurantAddress)
		csvLines[index] = l
	}

	writeCsv(csvLines, "info_stores.csv")
}

func getInfoCompanyAndWrite() {
	csvLines, err := readFromCsv("info_stores.csv")
	if err != nil {
		log.Fatalf("can't read csv=> %s", err)
	}
	for index, l := range csvLines {
		companyName := l[2]
		companyInfos := getInfoCompany(companyName)
		if err != nil {
			log.Fatal(err)
		}
		l = append(l, companyInfos...)

		csvLines[index] = l
	}
	writeCsv(csvLines, "info_company.csv")
}
