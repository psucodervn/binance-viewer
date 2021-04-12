package binance

import (
	"context"
	"testing"
)

func TestGetIdolPositions(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{args: args{uid: "9745A111F31F836D6D2E9F758DA3A07B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIdolPositions(context.Background(), tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIdolPositions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				t.Log(got.Positions)
			}
		})
	}
}
