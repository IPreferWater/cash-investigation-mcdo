package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
)

type ListStoreResponse struct {
	Data map[int]StoreID `json:"data"`
}

type StoreID struct {
	StoreID string `json:"store_id"`
}

var (
	urlBaseStoresIds = "https://api.woosmap.com/tiles/%s.grid.json?key=%s"
	key              = "woos-77bec2e5-8f40-35ba-b483-67df0d5401be&_=1650092957"
	zoom             = 8
	xMin, xMax       = 122, 134
	yMin, yMax       = 85, 96
)

func getListStores() ([]string, error) {
	log.Println("starting get list stores")
	//store the stores in map first, to avoid doublons
	mapStoresIds := make(map[string]int)

	for i := xMin; i <= xMax; i++ {
		for j := yMin; j <= yMax; j++ {

			var listStoreResponse ListStoreResponse

			//slow down
			time.Sleep(1 * time.Second)
			paramsUrl := fmt.Sprintf("%d-%d-%d", zoom, i, j)
			url := fmt.Sprintf(urlBaseStoresIds, paramsUrl, key)
			log.Printf("doing request for params %s\n", paramsUrl)
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return nil, err
			}

			setHeaders(req)

			res, err := client.Do(req)
			if err != nil {
				return nil, err
			}

			defer res.Body.Close()

			// do br decoding
			reader := brotli.NewReader(res.Body)
			b, err := ioutil.ReadAll(reader)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &listStoreResponse)
			if err != nil {
				return nil, err
			}

			log.Printf("request returned %d storeIds\n", len(listStoreResponse.Data))
			for _, storeID := range listStoreResponse.Data {
				mapStoresIds[storeID.StoreID]++
			}
			log.Printf("mapStoresIds is now %d size\n", len(mapStoresIds))
		}
	}

	//store the IDs in a slice
	arrStoreIDs := make([]string, len(mapStoresIds))
	index := 0
	for keyStoreID := range mapStoresIds {
		
		/* this part was to check if there is doublon in the slices
		if countStoreID>1 {
			log.Fatalf("some stores are not unique, check the map %v", m)
		}*/

		arrStoreIDs[index] = keyStoreID
		index++
	}

	return arrStoreIDs, nil
}
