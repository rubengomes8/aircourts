package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rubengomes8/aircourts/internal/domain"
)

/*
	CLUBES:
		311 - RACKETS PRO EUL
		48 - RACKETS PRO NACIONAL
		355 - RACKETS PRO SALDANHA
		96 - PADEL SPOT OLAIAS
		387 - VIVE PADEL
*/

func main() {

	clubIDs := []string{"311"}
	dates := []string{"2022-03-14", "2022-03-15", "2022-03-16", "2022-03-17", "2022-03-18"}
	startTime := "18:00"

	var clubs []domain.Club

	for _, clubID := range clubIDs {

		club := domain.Club{
			ClubID: clubID,
		}

		for _, date := range dates {

			club.Date = date

			url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=0&date=%s&start_time=%s", clubID, date, startTime)

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

			var clubCourts []domain.Court
			for _, result := range searchWithClubResponse.Results {

				if club.ClubName == "" {
					club.ClubName = result.ClubName
				}

				court := domain.Court{
					CourtID:   result.CourtID,
					CourtName: result.CourtName,
				}

				var courtFreeSlots []domain.FreeSlot
				for _, slot := range result.Slots {

					if !slot.Locked {

						freeSlot := domain.FreeSlot{
							Date:    date,
							Start:   slot.Start,
							End:     slot.End,
							CourtID: slot.CourtID,
						}

						courtFreeSlots = append(courtFreeSlots, freeSlot)
					}
				}

				court.FreeSlots = courtFreeSlots
				clubCourts = append(clubCourts, court)
			}

			//Convert the body to type string
			club.Courts = clubCourts

			clubs = append(clubs, club)
		}
	}

	fmt.Println(clubs)
}
