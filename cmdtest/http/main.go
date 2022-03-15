package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rubengomes8/aircourts/internal/domain"
)

func main() {

	clubID := 311
	date := "2022-03-15"
	startTime := "18:00"

	url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%v?sport=0&date=%v&start_time=%v", clubID, date, startTime)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var searchWithClubResponse domain.SearchWithClubResponse
	err = json.Unmarshal(body, &searchWithClubResponse)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%+v\n", searchWithClubResponse)
}
