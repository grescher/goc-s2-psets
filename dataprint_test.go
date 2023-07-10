package main

import (
	"reflect"
	"testing"
)

func TestNewUsersToPrint(t *testing.T) {
	type args struct {
		users   []User
		headers []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes UsersToPrint
	}{
		{
			name: "Test",
			args: args{
				users: []User{
					{"John Doe", 30, true, 80.0, []string{"Harry Potter", "1984"}},
					{"Jake Doe", 20, false, 60.0, []string{}},
					{" Jane Doe ", 150, true, .75, []string{"Harry Potter", "Game of Thrones"}},
				},
				headers: []string{"Name", "Age", "Active", "Mass", "Books"},
			},
			wantRes: UsersToPrint{
				Header: []string{"Name", "Age", "Active", "Mass", "Books"},
				Rows: []RowField{
					{
						"Name":   "John Doe",
						"Age":    "30",
						"Active": "yes",
						"Mass":   "80.0 kg",
						"Books":  "\"Harry Potter\"\n\"1984\"",
					},
					{
						"Name":   "Jake Doe",
						"Age":    "20",
						"Active": "-",
						"Mass":   "60.0 kg",
						"Books":  "",
					},
					{
						"Name":   " Jane Doe ",
						"Age":    "150",
						"Active": "yes",
						"Mass":   "0.75 qq",
						"Books":  "\"Harry Potter\"\n\"Game of Thrones\"",
					},
				},
				ColumnWidth: map[string]int{"Name": 10, "Age": 3, "Active": 6, "Mass": 7, "Books": 17},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := NewUsersToPrint(tt.args.users, tt.args.headers); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("NewUsersToPrint() = %#v, want %#v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_chunk(t *testing.T) {
	type args struct {
		str        string
		widthLimit int
	}
	tests := []struct {
		name       string
		args       args
		wantCh     string
		wantNewStr string
		wantOk     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCh, gotNewStr, gotOk := chunk(tt.args.str, tt.args.widthLimit)
			if gotCh != tt.wantCh {
				t.Errorf("chunk() gotCh = %v, want %v", gotCh, tt.wantCh)
			}
			if gotNewStr != tt.wantNewStr {
				t.Errorf("chunk() gotNewStr = %v, want %v", gotNewStr, tt.wantNewStr)
			}
			if gotOk != tt.wantOk {
				t.Errorf("chunk() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
