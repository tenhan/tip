package str

import "testing"

func TestIsAlpha(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"abc", args{str: "abc"}, true},
		{"Abc", args{str: "abc"}, true},
		{"", args{str: ""}, false},
		{" ", args{str: " "}, false},
		{"123", args{str: "123"}, false},
		{"abc Ac", args{str: "abc Ac"}, false},
		{"ab123", args{str: "ab123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.args.str); got != tt.want {
				t.Errorf("IsAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEnglish(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"abc", args{str: "abc"}, true},
		{"Abc", args{str: "abc"}, true},
		{"", args{str: ""}, false},
		{" ", args{str: " "}, false},
		{"你好", args{str: "你好"}, false},
		{"?", args{str: "?"}, false},
		{"!hello", args{str: "!hello"}, false},
		{"123", args{str: "123"}, false},
		{"abc Ac", args{str: "abc Ac"}, true},
		{"ab123", args{str: "ab123"}, false},
		{"ab cd", args{str: "ab cd"}, true},
		{"ab123", args{str: "ab123"}, false},
		{"an anti-pattern", args{str: "anti-pattern"}, true},
		{"I'm fine.", args{str: "I'm fine."}, true},
		{"I'm fine!", args{str: "I'm fine!"}, true},
		{"Are you OK?", args{str: "Are you OK?"}, true},
		{"Yes, I am.", args{str: "Yes, I am."}, true},
		{"Hello world", args{str: "Hello world"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEnglish(tt.args.str); got != tt.want {
				t.Errorf("%s: IsEnglish() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
