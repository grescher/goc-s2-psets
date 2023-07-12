package main

import (
	"reflect"
	"testing"
)

func TestUserSlice_NewTable(t *testing.T) {
	type args struct {
		headers []string
	}
	tests := []struct {
		name    string
		users   UserSlice
		args    args
		wantRes Table
	}{
		{
			name: "Test",
			users: UserSlice{
				{"John Doe", 30, true, 80.0, []string{"Harry Potter", "1984"}},
				{"Jake Doe", 20, false, 60.0, []string{}},
				{" Jane Doe ", 150, true, .75, []string{"Harry Potter", "Game of Thrones"}},
			},
			args: args{
				headers: []string{"Name", "Age", "Active", "Mass", "Books"},
			},
			wantRes: Table{
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
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := tt.users.NewTable(tt.args.headers); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("Users.NewTable() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
