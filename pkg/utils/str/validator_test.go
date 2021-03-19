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
		{"abc Ac", args{str: "123"}, false},
		{"ab123", args{str: "123"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAlpha(tt.args.str); got != tt.want {
				t.Errorf("IsAlpha() = %v, want %v", got, tt.want)
			}
		})
	}
}
