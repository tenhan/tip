package dictionary

import (
	"context"
	"fmt"
	"testing"
)

func TestJargon_Handle(t *testing.T) {
	type args struct {
		ctx     context.Context
		keyword string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "RTFM", args: args{
				ctx:     context.TODO(),
				keyword: "RTFM",
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Jargon{}
			gotResults, err := s.Handle(tt.args.ctx, tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotResults)
		})
	}
}
