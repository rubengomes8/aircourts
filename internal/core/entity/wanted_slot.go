package entity

import "strings"

type WantedSlot struct {
	CourtName string `json:"court_name"`
	Date      string `json:"date"`
	Start     string `json:"start"`
	End       string `json:"end"`
}

func DiscardCourt(courtName, roof string, onlyIndoor bool) bool {
	if strings.Contains(strings.ToLower(courtName), "ténis") {
		return true
	}

	if onlyIndoor && strings.Contains(strings.ToLower(courtName), "descoberto") {
		return true
	}

	if onlyIndoor && roof == "0" {
		return true
	}

	return false
}
