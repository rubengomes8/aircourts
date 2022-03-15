package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rubengomes8/aircourts/internal/consuming"
	"github.com/rubengomes8/aircourts/internal/domain"
	"github.com/rubengomes8/aircourts/internal/utils"
)

/*
	CLUBES:
		311 - RACKETS PRO EUL
		48 - RACKETS PRO NACIONAL
		355 - RACKETS PRO SALDANHA
		96 - PADEL SPOT OLAIAS
		387 - VIVE PADEL
		441 - W Padel
*/

// brew install dlv
// export GOPATH="$HOME/go"
// export PATH="$GOPATH/bin:$PATH"
// dlv debug cmd/http/main.go --headless --listen=:2345 --log

const (
	startDate = "2022-03-21"
	endDate   = "2022-03-25"
	startTime = "18:30"
	minSlots  = 3
	maxStart  = "21:00"

	onlyIndoor = true
	dateLayout = "2006-01-02"
)

func main() {

	clubIDs := []string{"355", "48", "311", "96", "387", "441"}
	dates, err := utils.DatesBetween(startDate, endDate, dateLayout, true, true)
	if err != nil {
		log.Fatalln(err)
	}

	for _, clubID := range clubIDs {

		club := domain.Club{
			ClubID: clubID,
		}

		for _, date := range dates {

			club.Date = date

			url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=0&date=%s&start_time=%s", clubID, date, startTime)

			body, err := consuming.HTTPGet(url)
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

				if domain.DiscardCourt(result.CourtName, onlyIndoor) {
					continue
				}

				if club.ClubName == "" {
					club.ClubName = result.ClubName
				}

				court := domain.Court{
					CourtID:   result.CourtID,
					CourtName: result.CourtName,
				}

				freeSlots := domain.FreeSlots(result, date)

				court.FreeSlots = freeSlots
				clubCourts = append(clubCourts, court)
			}

			club.Courts = clubCourts
		}

		clubReport := domain.WantedSlots(club, minSlots, maxStart)

		if clubReport == nil {
			continue
		}

		err := domain.ReportWantedSlots(os.Stdout, clubReport)
		if err != nil {
			log.Fatalln(err)
		}

		// clubJsonIndent, err := json.MarshalIndent(clubReport, "", "     ")
		// if err != nil {
		// 	log.Fatalln(err)
		// }

		// // fmt.Printf("%s\n", clubJsonIndent)
	}

}
