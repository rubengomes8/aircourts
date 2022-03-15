package domain

import (
	"io"
	"os"
	"text/template"
)

const (
	templ = `
Club: {{.ClubName}}:
{{range .WantedSlots}}----------------------------------------
Court:  {{.CourtName}}
Date:   {{.Date}} {{.Start}} - {{.End}}
{{end}}`
)

type ClubReport struct {
	ClubName    string       `json:"club_name"`
	WantedSlots []WantedSlot `json:"wanted_slots"`
}

type WantedSlot struct {
	CourtName string `json:"court_name"`
	Date      string `json:"date"`
	Start     string `json:"start"`
	End       string `json:"end"`
}

func ReportWantedSlots(wr io.Writer, clubReport interface{}) error {

	var report = template.Must(template.New("wantedSlots").Parse(templ))

	if err := report.Execute(os.Stdout, &clubReport); err != nil {
		return err
	}
	return nil
}
