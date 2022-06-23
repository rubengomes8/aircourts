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
	"github.com/rubengomes8/aircourts/internal/core/dto"
	"github.com/rubengomes8/aircourts/internal/core/entity"
	"github.com/rubengomes8/aircourts/internal/utils"
)

/*
	CLUBES:
	355 - RACKETS PRO SALDANHA
	311 - RACKETS PRO EUL
	48 - RACKETS PRO NACIONAL
	96 - PADEL SPOT OLAIAS
	441 - W Padel
	387 - VIVE PADEL
	110 - Indoor Padel Center
	106 - Padel CIF
	390 - Padel EXPO
	316 - Padel Benfica
	56 - Ténis e Padel Boa hora
	89 - Lambert
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

	// clubIDs := []string{"355", "48", "311", "96", "441", "387", "110", "106", "390", "316", "56", "89"}
	clubIDs := []string{"355", "48", "311", "96", "441", "387", "110", "106"}

	dates, err := utils.DatesBetween(*startDate, *endDate, dateLayout, *includeStart, *includeEnd, *allowFridays, *allowWeekends)
	if err != nil {
		log.Fatalln(err)
	}

	var buffer bytes.Buffer
	w := bufio.NewWriter(&buffer)

	var clubReport *entity.ClubReport

	for _, date := range dates {

		var club entity.Club
		for _, clubID := range clubIDs {

			club = entity.Club{
				ClubID: clubID,
			}

			club.Date = date

			url := fmt.Sprintf("https://www.aircourts.com/index.php/api/search_with_club/%s?sport=0&date=%s&start_time=%s", clubID, date, *startTime)

			body, err := http.HTTPGet(url)
			if err != nil {
				log.Fatalln(err)
			}

			var searchWithClubResponse dto.SearchWithClubResponse
			err = json.Unmarshal(body, &searchWithClubResponse)
			if err != nil {
				log.Fatalln(err)
			}

			for _, result := range searchWithClubResponse.Results {

				if entity.DiscardCourt(result.CourtName, result.Roof, *onlyIndoor) {
					continue
				}

				if club.ClubName == "" {
					club.ClubName = result.ClubName
				}

				court := entity.Court{
					CourtID:   result.CourtID,
					CourtName: result.CourtName,
				}

				freeSlots := result.FreeSlots(date)

				court.FreeSlots = freeSlots
				club.Courts = append(club.Courts, court)
			}

			clubReport = club.WantedSlots(*minSlots, *maxStart)

			if clubReport == nil {
				continue
			}

			err = entity.ReportWantedSlots(w, clubReport)
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

			err := sender.SendEmails(subject, emailBody, from, receivers)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if emailBody != "" {
		fmt.Println(emailBody)
	}

	fmt.Printf("### %v - Execution time: %v sec\n", time.Now(), time.Since(startExecutionTime).Seconds())
}
