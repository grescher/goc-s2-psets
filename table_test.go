package main

import (
	"reflect"
	"testing"
)

func TestTable_setColumnWidth(t *testing.T) {
	type fields struct {
		Header      []string
		Rows        []RowField
		ColumnWidth map[string]int
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]int
	}{
		{
			name: "General",
			fields: fields{
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
			want: map[string]int{"Active": 6, "Age": 3, "Books": 17, "Mass": 7, "Name": 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Table{
				Header:      tt.fields.Header,
				Rows:        tt.fields.Rows,
				ColumnWidth: tt.fields.ColumnWidth,
			}
			tr.setColumnWidth()
			if !reflect.DeepEqual(tr.ColumnWidth, tt.want) {
				t.Errorf("table.setColumnWidth() = %#v, want: %#v\n", tr.ColumnWidth, tt.want)
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
		{
			name:       "Empty string",
			args:       args{str: "", widthLimit: 0},
			wantCh:     "  ",
			wantNewStr: "",
			wantOk:     true,
		},
		{
			name:       "Empty string and limit=10",
			args:       args{str: "", widthLimit: 10},
			wantCh:     "            ",
			wantNewStr: "",
			wantOk:     true,
		},
		{
			name:       "Exceeding Width Limit",
			args:       args{str: "The width of this string exceeds the limit", widthLimit: 10},
			wantCh:     " The width  ",
			wantNewStr: "of this string exceeds the limit",
			wantOk:     false,
		},
		{
			name:       "Exceeding Width Limit Locale",
			args:       args{str: "Ширина цієї стрічки перевищує ліміт", widthLimit: 10},
			wantCh:     " Ширина ціє ",
			wantNewStr: "ї стрічки перевищує ліміт",
			wantOk:     false,
		},
		{
			name:       "Short string",
			args:       args{str: "Hello", widthLimit: 10},
			wantCh:     " Hello      ",
			wantNewStr: "",
			wantOk:     true,
		},
		{
			name:       "Books: Short name",
			args:       args{str: "\"1984\"\n\"Harry Potter\"", widthLimit: 10},
			wantCh:     " \"1984\"     ",
			wantNewStr: "\"Harry Potter\"",
			wantOk:     false,
		},
		{
			name:       "Books: Long name",
			args:       args{str: "\"Harry Potter\"\n\"1984\"", widthLimit: 10},
			wantCh:     " \"Harry Pot ",
			wantNewStr: "ter\"\n\"1984\"",
			wantOk:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCh, gotNewStr, gotOk := chunk(tt.args.str, tt.args.widthLimit)
			if gotCh != tt.wantCh {
				t.Errorf("chunk() gotCh = %#v, want %#v", gotCh, tt.wantCh)
			}
			if gotNewStr != tt.wantNewStr {
				t.Errorf("chunk() gotNewStr = %#v, want %#v", gotNewStr, tt.wantNewStr)
			}
			if gotOk != tt.wantOk {
				t.Errorf("chunk() gotOk = %#v, want %#v", gotOk, tt.wantOk)
			}
		})
	}
}
