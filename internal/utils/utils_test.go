package utils

import (
	"reflect"
	"testing"
)

func TestDatesBetween(t *testing.T) {
	type args struct {
		startDate     string
		endDate       string
		layout        string
		allowFridays  bool
		allowWeekends bool
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "march_tt",
			args: args{
				startDate:     "2022-03-08",
				endDate:       "2022-03-13",
				layout:        "2006-01-02",
				allowFridays:  true,
				allowWeekends: true,
			},
			want:    []string{"2022-03-08", "2022-03-09", "2022-03-10", "2022-03-11", "2022-03-12", "2022-03-13"},
			wantErr: false,
		},
		{
			name: "jan_feb_tt",
			args: args{
				startDate:     "2022-01-30",
				endDate:       "2022-02-02",
				layout:        "2006-01-02",
				allowFridays:  true,
				allowWeekends: true,
			},
			want:    []string{"2022-01-30", "2022-01-31", "2022-02-01", "2022-02-02"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DatesBetween(tt.args.startDate, tt.args.endDate, tt.args.layout, tt.args.allowFridays, tt.args.allowWeekends)
			if (err != nil) != tt.wantErr {
				t.Errorf("DatesBetween() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatesBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
