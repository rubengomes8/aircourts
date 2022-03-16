package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pelletier/go-toml"
	"github.com/rubengomes8/aircourts/internal/communication/http"
	"github.com/rubengomes8/aircourts/internal/communication/smtp"
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

/*
	roof: 0 - Descoberto
	roof: 1 - Coberto
	roof: 2 - Indoor
*/

// brew install dlv
// export GOPATH="$HOME/go"
// export PATH="$GOPATH/bin:$PATH"
// dlv debug cmd/http/main.go --headless --listen=:2345 --log

const (
	startDate = "2022-03-21"
	endDate   = "2022-04-01"
	startTime = "18:30"
	minSlots  = 3
	maxStart  = "21:00"

	allowFridays  = false
	allowWeekends = false
	onlyIndoor    = true
	includeStart  = true
	includeEnd    = true
	dateLayout    = "2006-01-02"

	sendEmail = false
)

func main() {

	startExecutionTime := time.Now()

	config, err := toml.LoadFile("./internal/configuration/config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	senderEmail := config.Get("sender.email").(string)
	senderPwd := config.Get("sender.password").(string)

	sender := smtp.NewSender(senderEmail, senderPwd)

	clubIDs := []string{"355", "48", "311", "96", "387", "441"}
	dates, err := utils.DatesBetween(startDate, endDate, dateLayout, includeStart, includeEnd, allowFridays, allowWeekends)
	if err != nil {
		log.Fatalln(err)
	}

	var buffer bytes.Buffer
	w := bufio.NewWriter(&buffer)

	var clubReport *domain.ClubReport

	for _, date := range dates {

		var club domain.Club
		for _, clubID := range clubIDs {

			club = domain.Club{
				ClubID: clubID,
			}

			club.Date = date

			url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=0&date=%s&start_time=%s", clubID, date, startTime)

			body, err := http.HTTPGet(url)
			if err != nil {
				log.Fatalln(err)
			}

			var searchWithClubResponse domain.SearchWithClubResponse
			err = json.Unmarshal(body, &searchWithClubResponse)
			if err != nil {
				log.Fatalln(err)
			}

			for _, result := range searchWithClubResponse.Results {

				if domain.DiscardCourt(result.CourtName, result.Roof, onlyIndoor) {
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
				club.Courts = append(club.Courts, court)
			}

			clubReport = domain.WantedSlots(club, minSlots, maxStart)

			if clubReport == nil {
				continue
			}

			err = domain.ReportWantedSlots(w, clubReport)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	w.Flush()
	emailBody := buffer.String()

	if sendEmail {

		if emailBody != "" {
			subject := config.Get("email.subject").(string)
			from := config.Get("email.from").(string)
			receivers := config.Get("email.to").([]interface{})

			for _, to := range receivers {

				email := smtp.Email{
					To:      to.(string),
					From:    from,
					Subject: subject,
					Body:    emailBody,
				}

				fmt.Println("Sending Email...")
				err = sender.SendEmail(email)
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	}

	fmt.Printf("\nExecution time: %v sec\n", time.Since(startExecutionTime).Seconds())
}
