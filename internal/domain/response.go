package domain

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
	// Schedules Schedules `json:"schedules"`
}

type Slot struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Start   string `json:"start"`
	End     string `json:"end"`
	CourtID string `json:"court_id"`
	Locked  bool   `json:"locked"`
	// Status  string `json:"status"`
	// LockReason string `json:"lock_reason"`
	// Durations  []string `json:"durations"`
}

// type Schedules struct {
// 	Schedule    Schedule `json:"schedule"`
// 	Opening     string   `json:"opening"`      // "2022-03-12 07:00:00"
// 	Closing     string   `json:"closing"`      // "2022-03-12 07:00:00"
// 	OpeningTime string   `json:"opening_time"` // "07:00:00"
// 	ClosingTime string   `json:"closing_time"` // "07:00:00"
// 	// Offline     bool     `json:"offline"`
// }

// type Schedule struct {
// 	CourtID        string `json:"court_id"`
// 	MondayOpen     string `json:"monday_open"` // "20:00:00"
// 	MondayClose    string `json:"monday_close"`
// 	TuesdayOpen    string `json:"tuesday_open"`
// 	TuesdayClose   string `json:"tuesday_close"`
// 	WednesdayOpen  string `json:"wednesday_open"`
// 	WednesdayClose string `json:"wednesday_close"`
// 	ThursdayOpen   string `json:"thursday_open"`
// 	ThursdayClose  string `json:"thursday_close"`
// 	FridayOpen     string `json:"friday_open"`
// 	FridayClose    string `json:"friday_close"`
// 	SaturdayOpen   string `json:"saturday_open"`
// 	SaturdayClose  string `json:"saturday_close"`
// 	SundayOpen     string `json:"sunday_open"`
// 	SundayClose    string `json:"sunday_close"`
// }
