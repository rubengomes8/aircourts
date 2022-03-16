package domain

import (
	"strings"
)

type Club struct {
	ClubID   string  `json:"club_id"`
	ClubName string  `json:"club_name"`
	Date     string  `json:"date"`
	Courts   []Court `json:"courts"`
}

type Court struct {
	CourtID   string     `json:"court_id"`
	CourtName string     `json:"court_name"`
	FreeSlots []FreeSlot `json:"free_slots"`
}

type FreeSlot struct {
	Date    string `json:"date"`
	Start   string `json:"start"`
	End     string `json:"end"`
	CourtID string `json:"court_id"`
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

func FreeSlots(courtResult Result, date string) []FreeSlot {

	var freeSlots []FreeSlot

	for _, slot := range courtResult.Slots {

		if !slot.Locked {

			freeSlot := FreeSlot{
				Date:    date,
				Start:   slot.Start,
				End:     slot.End,
				CourtID: slot.CourtID,
			}

			freeSlots = append(freeSlots, freeSlot)
		}
	}

	return freeSlots
}

func WantedSlots(club Club, minSlots int, maxStart string) *ClubReport {

	var clubReport ClubReport

	var ws WantedSlot
	isFirst := true
	for _, court := range club.Courts {

		numSlotsInSeq := 0
		prevEnd := ""
		var slotInit string
		var freeSlot FreeSlot
		for _, freeSlot = range court.FreeSlots {

			if numSlotsInSeq == 0 || freeSlot.Start == prevEnd {

				if freeSlot.Start > maxStart {
					continue
				}

				if numSlotsInSeq == 0 {
					slotInit = freeSlot.Start
				}

				numSlotsInSeq++
				prevEnd = freeSlot.End
			} else {
				if numSlotsInSeq >= minSlots {

					if isFirst {
						clubReport.ClubName = club.ClubName
						isFirst = false
					}

					ws = WantedSlot{
						CourtName: court.CourtName,
						Date:      freeSlot.Date,
						Start:     slotInit,
						End:       prevEnd,
					}

					clubReport.WantedSlots = append(clubReport.WantedSlots, ws)
				}
				numSlotsInSeq = 0
				prevEnd = ""
			}
		}

		if numSlotsInSeq >= minSlots {

			if isFirst {
				clubReport.ClubName = club.ClubName
				isFirst = false
			}

			ws = WantedSlot{
				CourtName: court.CourtName,
				Date:      freeSlot.Date,
				Start:     slotInit,
				End:       prevEnd,
			}

			clubReport.WantedSlots = append(clubReport.WantedSlots, ws)
		}
	}

	if isFirst {
		return nil
	}

	return &clubReport
}
