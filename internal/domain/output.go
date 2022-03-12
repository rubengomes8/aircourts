package domain

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
