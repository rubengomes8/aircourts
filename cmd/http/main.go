package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
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
	dateLayout = "2006-01-02"
)

func main() {

	startExecutionTime := time.Now()

	startDate := flag.String("start", time.Now().Format(dateLayout), "start date")
	endDate := flag.String("end", time.Now().AddDate(0, 0, 7).Format(dateLayout), "end date")
	startTime := flag.String("startTime", "18:30", "starting time")
	maxStart := flag.String("maxStart", "22:00", "maximum starting time")
	minSlots := flag.Int("slots", 3, "minimum slots in a row")
	allowFridays := flag.Bool("fridays", false, "allow fridays")
	allowWeekends := flag.Bool("weekends", false, "allow weekends")
	onlyIndoor := flag.Bool("indoor", true, "only indoor")
	includeStart := flag.Bool("includeStart", true, "include starting time")
	includeEnd := flag.Bool("includeEnd", true, "include ending time")
	sendEmail := flag.Bool("email", false, "send email with results")

	flag.Parse()

	validDates, err := utils.ValidDates(*startDate, *endDate, dateLayout)
	if err != nil {
		log.Fatalln(err)
	}
	if !validDates {
		log.Fatal("StartDate after EndDate.")
	}

	config, err := toml.LoadFile("./configuration/config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	senderEmail := config.Get("sender.email").(string)
	senderPwd := config.Get("sender.password").(string)

	sender := smtp.NewSender(senderEmail, senderPwd)

	clubIDs := []string{"355", "48", "311", "96", "441"}
	dates, err := utils.DatesBetween(*startDate, *endDate, dateLayout, *includeStart, *includeEnd, *allowFridays, *allowWeekends)
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

			url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=0&date=%s&start_time=%s", clubID, date, *startTime)

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

				if domain.DiscardCourt(result.CourtName, result.Roof, *onlyIndoor) {
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

			clubReport = domain.WantedSlots(club, *minSlots, *maxStart)

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

	if *sendEmail {

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

	if emailBody != "" {
		fmt.Println(emailBody)
	}

	fmt.Printf("### %v - Execution time: %v sec\n", time.Now(), time.Since(startExecutionTime).Seconds())
}
