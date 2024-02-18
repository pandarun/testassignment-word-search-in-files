package text

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestCustomRuneReader_ReadRune(t *testing.T) {

	type fields struct {
		filter []rune
	}
	type args struct {
		input string
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantRune rune
		wantErr  bool
	}{
		{
			name: "test for OK",
			fields: fields{
				filter: []rune{rune(Dot)},
			},
			args: args{
				input: "4.",
			},
			wantRune: rune('4'),
			wantErr:  false,
		},
		{
			name: "test for SHY",
			fields: fields{
				filter: []rune{rune(SoftHyphen)},
			},
			args: args{
				input: "4­",
			},
			wantRune: rune('4'),
			wantErr:  false,
		},
		{
			name: "test for braces",
			fields: fields{
				filter: []rune{rune(RoundBraceLeft), rune(RoundBraceRight)},
			},
			args: args{
				input: "(4)",
			},
			wantRune: rune('4'),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := bufio.NewReader(strings.NewReader(tt.args.input))
			r := NewCustomRuneReader(in, tt.fields.filter...)
			gotRune, _, err := r.ReadRune()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadRune() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotRune) != string(tt.wantRune) {
				t.Errorf("test for OK failed - results not match\n %v \n %v", string(gotRune), string(tt.wantRune))
			}

		})
	}
}

func TestCustomRuneReader_Words(t *testing.T) {
	type fields struct {
		filter []rune
	}
	type args struct {
		input string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		wantWords []string
		wantErr   bool
	}{

		{
			name: "test for OK space",
			fields: fields{
				filter: GetDefaultFilteredRunes(),
			},
			args: args{
				input: "-  a, b.   (c)  d-e ",
			},
			wantWords: []string{"a", "b", "c", "d-e"},
			wantErr:   false,
		},
		{
			name: "test for OK",
			fields: fields{
				filter: []rune{rune(Dot)},
			},
			args: args{
				input: "4.",
			},
			wantWords: []string{"4"},
			wantErr:   false,
		},
		{
			name: "test for SHY",
			fields: fields{
				filter: []rune{rune(SoftHyphen)},
			},
			args: args{
				input: "4­",
			},
			wantWords: []string{"4"},
			wantErr:   false,
		},
		{
			name: "test for braces",
			fields: fields{
				filter: []rune{rune(RoundBraceLeft), rune(RoundBraceRight)},
			},
			args: args{
				input: "(4)",
			},
			wantWords: []string{"4"},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := bufio.NewReader(strings.NewReader(tt.args.input))
			gotWords, err := Words(in, tt.fields.filter...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(gotWords, tt.wantWords) != true {
				t.Errorf("Read() gotWords = %v, want %v", gotWords, tt.wantWords)
			}

		})
	}
}
