package dto

import "github.com/rubengomes8/aircourts/internal/core/entity"

type SearchWithClubResponse struct {
	Results []Result
}

type Result struct {
	ClubID    string `json:"club_id"`
	ClubName  string `json:"club_name"`
	CourtID   string `json:"id"`
	CourtName string `json:"name"`
	Roof      string `json:"roof"`
	Slots     []Slot `json:"slots"`
}

type Slot struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Start   string `json:"start"`
	End     string `json:"end"`
	CourtID string `json:"court_id"`
	Locked  bool   `json:"locked"`
}

func (r Result) FreeSlots(date string) []entity.FreeSlot {

	var freeSlots []entity.FreeSlot

	for _, slot := range r.Slots {

		if !slot.Locked {

			freeSlot := entity.FreeSlot{
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
