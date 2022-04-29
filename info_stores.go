package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type InfoStoreResponse struct {
	Name              string              `json:"name"`
	CompanyName       string              `json:"companyName"`
	Coordinates       Coordinates         `json:"coordinates"`
	RestaurantAddress []RestaurantAddress `json:"restaurantAddress"`
}

type InfoStoreFormated struct {
	Name              string
	CompanyName       string
	Coordinates       string
	RestaurantAddress string
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RestaurantAddress struct {
	City     string `json:"city"`
	ZipCode  string `json:"zipCode"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
	Address4 string `json:"address4"`
}

var (
	urlBaseInfoStore = "https://ws.mcdonalds.fr/api/restaurant/%s/?responseGroups=RG.RESTAURANT.FACILITIES"
)

func getStoresInfo(storeID string) (*InfoStoreFormated, error) {
	//slow down
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf(urlBaseInfoStore, storeID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	setHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var infoStoreResponse InfoStoreResponse
	err = json.Unmarshal(b, &infoStoreResponse)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	formated := formateInfoStore(infoStoreResponse)
	return &formated, nil
}

func formateInfoStore(i InfoStoreResponse) InfoStoreFormated {
	coordinates := fmt.Sprintf("%f;%f", i.Coordinates.Latitude, i.Coordinates.Longitude)
	address1 := i.RestaurantAddress[0]
	restaurantAddress := fmt.Sprintf("%s %s %s", address1.Address1, address1.City, address1.ZipCode)

	removeEurlSarl := func(n string) string{
		indexFirstEscape := strings.Index(n," ")
		if strings.Index(n," ")>0 {
			return n[:indexFirstEscape]
		}
		return n
	}

	companyName := removeEurlSarl(i.CompanyName)
	return InfoStoreFormated{
		Name:              i.Name,
		CompanyName:       companyName,
		Coordinates:       coordinates,
		RestaurantAddress: restaurantAddress,
	}
}
