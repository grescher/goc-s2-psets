// Implementing stringer methods for the User's fields
package main

import (
	"reflect"
	"testing"
)

func TestName_String(t *testing.T) {
	tests := []struct {
		name string
		n    Name
		want string
	}{
		{
			name: "Escape sequences",
			n:    "\tJohn\nAppleseed \b",
			want: "\\tJohn\\nAppleseed \\b",
		},
		{
			name: "Bytes",
			n:    "\x00\x10\x20\x30\x40\x50\x60\x70",
			want: "\\x00\\x10 0@P`p",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("Name.String() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestAge_String(t *testing.T) {
	tests := []struct {
		name string
		a    Age
		want string
	}{
		{name: "Positive", a: 42, want: "42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Age.String() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestActive_String(t *testing.T) {
	tests := []struct {
		name string
		a    ActiveIndex
		want string
	}{
		{name: "True 0", a: 0b0001, want: "yes"},
		{name: "True 1", a: 0b0010, want: "yes"},
		{name: "True 2", a: 0b0100, want: "yes"},
		{name: "True 3", a: 0b1000, want: "yes"},
		{name: "True 7", a: 0b10000000, want: "yes"},
		{name: "False", a: 0b0, want: "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Active.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMass_String(t *testing.T) {
	tests := []struct {
		name string
		m    Mass
		want string
	}{
		{name: "0.001 qq", m: 0.001, want: "0.001 qq"},
		{name: "0.999 qq", m: 0.999, want: "0.999 qq"},
		{name: "90 kg", m: 90, want: "90.0 kg"},
		{name: "621 oz", m: 621, want: "621.0 oz"},
		{name: "8080 oz", m: 8080, want: "8080.0 oz"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.String(); got != tt.want {
				t.Errorf("Mass.String() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestBooks_String(t *testing.T) {
	tests := []struct {
		name string
		b    Books
		want string
	}{
		{
			name: "Normal",
			b:    Books([]string{"1984", "Harry Potter", "Cat's Cradle"}),
			want: "\"1984\"\n\"Harry Potter\"\n\"Cat's Cradle\"",
		},
		{
			name: "EscSequencies",
			b:    Books([]string{"1984\n", "Harry O\"Potter\t", "Cat's Cradle"}),
			want: "\"1984\\n\"\n\"Harry O\\\"Potter\\t\"\n\"Cat's Cradle\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Books.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
				{"John Doe", 30, 0b00000001, 80.0, []string{"Harry Potter", "1984"}},
				{"Jake Doe", 20, 0b0, 60.0, []string{}},
				{" Jane Doe 1", 150, 0b00000100, .75, []string{"Harry Potter", "Game of Thrones"}},
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
