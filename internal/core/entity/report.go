package entity

import (
	"html/template"
	"io"
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

func ReportWantedSlots(wr io.Writer, clubReport interface{}) error {

	if clubReport == nil {
		return nil
	}

	var report = template.Must(template.New("wantedSlots").Parse(templ))

	if err := report.Execute(wr, &clubReport); err != nil {
		return err
	}

	return nil
}
