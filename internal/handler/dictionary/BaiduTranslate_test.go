package dictionary

import (
	"context"
	"fmt"
	"github.com/tenhan/tip/internal/handler"
	"testing"
)

func TestBaiduTranslate_Handle(t *testing.T) {
	type args struct {
		ctx     context.Context
		keyword string
	}
	tests := []struct {
		name        string
		args        args
		wantResults []handler.Result
		wantErr     bool
	}{
		{
			name:        "hello",
			wantErr:     false,
		},
		{
			name:        "avenger",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BaiduTranslate{}
			gotResults, err := s.Handle(tt.args.ctx, tt.args.keyword)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(gotResults)
		})
	}
}
