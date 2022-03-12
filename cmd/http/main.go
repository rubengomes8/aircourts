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

	clubIDs := []string{"311"}
	sports := []string{"0"}
	dates := []string{"2022-03-12"}
	startTime := "11:30"

	url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=%s&date=%s&start_time=%s", clubIDs[0], sports[0], dates[0], startTime)

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

	//Convert the body to type string
	log.Println(searchWithClubResponse)
}
