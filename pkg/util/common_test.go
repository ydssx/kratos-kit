package util

import (
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestJsonToMap(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", args{`{"Name": "Alice", "Age": 30, "Address": "123 Main St."}`}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToMap(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("JsonToMap() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestCompareRequests(t *testing.T) {
	type Request struct {
		Field1 string
		Field2 int
	}

	type args struct {
		requests []interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"", args{requests: []interface{}{&Request{Field1: "value", Field2: 10}, &Request{Field1: "value", Field2: 10}}}, true},
		{"", args{requests: []interface{}{&Request{Field1: "value", Field2: 10}, &Request{Field1: "value", Field2: 11}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareRequests(tt.args.requests...); got != tt.want {
				t.Errorf("CompareRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	type args struct {
		data interface{}
	}
	type MyStruct struct {
		Name    string `default:"John"`
		Age     int    `default:"30"`
		Enabled bool   `default:"true"`
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"", args{&MyStruct{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefaults(tt.args.data)
			t.Logf("%#+v", tt.args.data)
		})
	}
}

func TestGenerateOrderNumber(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{"", "1234567890", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateOrderNumber()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOrderNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateOrderNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLastDayOfMonth(t *testing.T) {
	type args struct {
		year  int
		month int
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		// TODO: Add test cases.
		{"", args{2020, 1}, time.Date(2020, 1, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 2}, time.Date(2020, 2, 29, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 3}, time.Date(2020, 3, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 4}, time.Date(2020, 4, 30, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 5}, time.Date(2020, 5, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 6}, time.Date(2020, 6, 30, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 7}, time.Date(2020, 7, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 8}, time.Date(2020, 8, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 9}, time.Date(2020, 9, 30, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 10}, time.Date(2020, 10, 31, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 11}, time.Date(2020, 11, 30, 0, 0, 0, 0, time.Local)},
		{"", args{2020, 12}, time.Date(2020, 12, 31, 0, 0, 0, 0, time.Local)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLastDayOfMonth(tt.args.year, tt.args.month); !reflect.DeepEqual(got, tt.want) {
				log.Printf("GetLastDayOfMonth() = %v, want %v", got.Format("2006-01-02"), tt.want)
				t.Errorf("GetLastDayOfMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateEmailAddress(t *testing.T) {
	tests := []struct {
		name string
		want *regexp.Regexp
	}{
		{
			name: "Valid email format",
			want: regexp.MustCompile(`^[a-zA-Z0-9]{5,15}@(gmail\.com|yahoo\.com|hotmail\.com|outlook\.com)$`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateEmailAddress()
			if !tt.want.MatchString(got) {
				t.Errorf("GenerateEmailAddress() = %v, want to match regex %v", got, tt.want)
			}
		})
	}
}

func TestGenerateEmailAddressUniqueness(t *testing.T) {
	emailSet := make(map[string]bool)
	iterations := 50

	for i := 0; i < iterations; i++ {
		email := GenerateEmailAddress()
		if emailSet[email] {
			t.Errorf("GenerateEmailAddress() produced duplicate email: %s", email)
		}
		emailSet[email] = true
	}
}

func TestGenerateEmailAddressDomainDistribution(t *testing.T) {
	domains := map[string]int{
		"gmail.com":   0,
		"yahoo.com":   0,
		"hotmail.com": 0,
		"outlook.com": 0,
	}
	iterations := 1000

	for i := 0; i < iterations; i++ {
		email := GenerateEmailAddress()
		for domain := range domains {
			if regexp.MustCompile("@" + domain + "$").MatchString(email) {
				domains[domain]++
				break
			}
		}
	}

	for domain, count := range domains {
		if count == 0 {
			t.Errorf("Domain %s was never used in %d iterations", domain, iterations)
		}
	}
}


func TestTimer(t *testing.T) {
	testFunc := func(x int) int {
		time.Sleep(100 * time.Millisecond)
		return x * 2
	}

	timedFunc := Timer(testFunc)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"Positive input", 5, 10},
		{"Zero input", 0, 0},
		{"Negative input", -3, -6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timedFunc(tt.input)
			if result != tt.expected {
				t.Errorf("Timer(%d) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTimerWithString(t *testing.T) {
	testFunc := func(s string) int {
		time.Sleep(50 * time.Millisecond)
		return len(s)
	}

	timedFunc := Timer(testFunc)

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Non-empty string", "hello", 5},
		{"Long string", "This is a long string for testing", 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timedFunc(tt.input)
			if result != tt.expected {
				t.Errorf("Timer(%q) = %d; want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTimerExecutionTime(t *testing.T) {
	testFunc := func(duration time.Duration) time.Duration {
		time.Sleep(duration)
		return duration
	}

	timedFunc := Timer(testFunc)

	input := 200 * time.Millisecond
	start := time.Now()
	result := timedFunc(input)
	executionTime := time.Since(start)

	if result != input {
		t.Errorf("Timer returned incorrect result: got %v, want %v", result, input)
	}

	if executionTime < input {
		t.Errorf("Timer execution time too short: got %v, want at least %v", executionTime, input)
	}

	tolerance := 50 * time.Millisecond
	if executionTime > input+tolerance {
		t.Errorf("Timer execution time too long: got %v, want at most %v", executionTime, input+tolerance)
	}
}
