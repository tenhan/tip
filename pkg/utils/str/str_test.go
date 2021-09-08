package str

import "testing"

func TestTrim(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{str: ""}, ""},
		{"2", args{str: " "}, ""},
		{"3", args{str: " abc \n abc"}, "abc \n abc"},
		{"4", args{str: "abc \n abc\n\n"}, "abc \n abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Trim(tt.args.str); got != tt.want {
				t.Errorf("Trim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicatedWhiteSpace(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{str: ""}, ""},
		{"4", args{str: " "}, " "},
		{"4", args{str: "  "}, " "},
		{"2", args{str: "abc   dd   abc abc"}, "abc dd abc abc"},
		{"3", args{str: "    abc    "}, " abc "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicatedWhiteSpace(tt.args.str); got != tt.want {
				t.Errorf("ClearDuplicatedWhiteSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartWith(t *testing.T) {
	type args struct {
		s      string
		subStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s:      "Hello world",
				subStr: "Hello",
			},
			want: true,
		},
		{
			name: "1",
			args: args{
				s:      "Hello world",
				subStr: "Hello ",
			},
			want: true,
		},
		{
			name: "1",
			args: args{
				s:      "Hello world",
				subStr: "Hello world!",
			},
			want: false,
		},
		{
			name: "1",
			args: args{
				s:      "",
				subStr: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StartWith(tt.args.s, tt.args.subStr); got != tt.want {
				t.Errorf("StartWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndWith(t *testing.T) {
	type args struct {
		s      string
		subStr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				s:      "Hello world",
				subStr: "world",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				s:      "Hello world",
				subStr: " ld",
			},
			want: false,
		},
		{
			name: "3",
			args: args{
				s:      "Hello world",
				subStr: " wld",
			},
			want: false,
		},
		{
			name: "4",
			args: args{
				s:      "Hello world",
				subStr: "123Hello world",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndWith(tt.args.s, tt.args.subStr); got != tt.want {
				t.Errorf("EndWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
