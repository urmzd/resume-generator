package resume

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestPartialDateUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantYear  int
		wantMonth time.Month
		wantPrec  DatePrecision
	}{
		{"year only quoted", `"2017"`, 2017, time.January, PrecisionYear},
		{"year only unquoted", `2017`, 2017, time.January, PrecisionYear},
		{"month-year", `"2016-05"`, 2016, time.May, PrecisionMonth},
		{"full date", `"2016-05-01T00:00:00Z"`, 2016, time.May, PrecisionFull},
		{"date only", `2019-01-15`, 2019, time.January, PrecisionFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pd PartialDate
			if err := yaml.Unmarshal([]byte(tt.input), &pd); err != nil {
				t.Fatalf("UnmarshalYAML(%q) error: %v", tt.input, err)
			}
			if pd.Time.Year() != tt.wantYear {
				t.Errorf("year = %d, want %d", pd.Time.Year(), tt.wantYear)
			}
			if pd.Time.Month() != tt.wantMonth {
				t.Errorf("month = %v, want %v", pd.Time.Month(), tt.wantMonth)
			}
			if pd.Precision != tt.wantPrec {
				t.Errorf("precision = %d, want %d", pd.Precision, tt.wantPrec)
			}
		})
	}
}

func TestPartialDateUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantYear  int
		wantMonth time.Month
		wantPrec  DatePrecision
	}{
		{"year only", `"2017"`, 2017, time.January, PrecisionYear},
		{"month-year", `"2016-05"`, 2016, time.May, PrecisionMonth},
		{"full RFC3339", `"2016-05-01T00:00:00Z"`, 2016, time.May, PrecisionFull},
		{"date only", `"2019-01-15"`, 2019, time.January, PrecisionFull},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pd PartialDate
			if err := json.Unmarshal([]byte(tt.input), &pd); err != nil {
				t.Fatalf("UnmarshalJSON(%q) error: %v", tt.input, err)
			}
			if pd.Time.Year() != tt.wantYear {
				t.Errorf("year = %d, want %d", pd.Time.Year(), tt.wantYear)
			}
			if pd.Time.Month() != tt.wantMonth {
				t.Errorf("month = %v, want %v", pd.Time.Month(), tt.wantMonth)
			}
			if pd.Precision != tt.wantPrec {
				t.Errorf("precision = %d, want %d", pd.Precision, tt.wantPrec)
			}
		})
	}
}

func TestPartialDateRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		pd   PartialDate
	}{
		{"year", NewYearDate(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))},
		{"month", NewMonthDate(time.Date(2020, 6, 1, 0, 0, 0, 0, time.UTC))},
		{"full", NewPartialDate(time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC), PrecisionFull)},
	}

	for _, tt := range tests {
		t.Run(tt.name+" yaml", func(t *testing.T) {
			data, err := yaml.Marshal(tt.pd)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			var got PartialDate
			if err := yaml.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if got.Precision != tt.pd.Precision {
				t.Errorf("precision = %d, want %d", got.Precision, tt.pd.Precision)
			}
		})

		t.Run(tt.name+" json", func(t *testing.T) {
			data, err := json.Marshal(tt.pd)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			var got PartialDate
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if got.Precision != tt.pd.Precision {
				t.Errorf("precision = %d, want %d", got.Precision, tt.pd.Precision)
			}
		})
	}
}

func TestDateRangeYAML(t *testing.T) {
	input := `
start: "2017"
end: "2019-05"
`
	var dr DateRange
	if err := yaml.Unmarshal([]byte(input), &dr); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if dr.Start.Precision != PrecisionYear {
		t.Errorf("start precision = %d, want PrecisionYear", dr.Start.Precision)
	}
	if dr.End == nil {
		t.Fatal("end should not be nil")
	}
	if dr.End.Precision != PrecisionMonth {
		t.Errorf("end precision = %d, want PrecisionMonth", dr.End.Precision)
	}
}
