package model

import (
	"math"
	"testing"
)

func TestDiffPositions(t *testing.T) {
	type args struct {
		oldPositions []Position
		newPositions []Position
	}
	tests := []struct {
		name string
		args args
		want []PositionDiff
	}{
		{
			args: args{oldPositions: nil, newPositions: []Position{{Symbol: "A", Amount: 1.5}, {Symbol: "B", Amount: 3.2}}},
			want: []PositionDiff{
				{Symbol: "A", Amount: 1.5, Change: 1.5, Status: StatusNew},
				{Symbol: "B", Amount: 3.2, Change: 3.2, Status: StatusNew},
			},
		},
		{
			args: args{oldPositions: []Position{{Symbol: "A", Amount: 1.5}}, newPositions: []Position{{Symbol: "A", Amount: 0.85}, {Symbol: "B", Amount: 3.2}}},
			want: []PositionDiff{
				{Symbol: "A", Amount: 0.85, Change: -0.65, Status: StatusDecreased},
				{Symbol: "B", Amount: 3.2, Change: 3.2, Status: StatusNew},
			},
		},
		{
			args: args{oldPositions: []Position{{Symbol: "A", Amount: 1.5}, {Symbol: "B", Amount: 1.4}}, newPositions: []Position{{Symbol: "B", Amount: 3.0}}},
			want: []PositionDiff{
				{Symbol: "B", Amount: 3.0, Change: 1.6, Status: StatusIncreased},
				{Symbol: "A", Change: -1.5, Status: StatusClosed},
			},
		},
		{
			args: args{
				oldPositions: []Position{{Symbol: "A", Amount: 1.5}, {Symbol: "B", Amount: 1.4}},
				newPositions: []Position{{Symbol: "A", Amount: 1.5}, {Symbol: "B", Amount: 3.0}},
			},
			want: []PositionDiff{
				{Symbol: "B", Amount: 3.0, Change: 1.6, Status: StatusIncreased},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DiffPositions(tt.args.oldPositions, tt.args.newPositions)
			if len(got) != len(tt.want) {
				t.Errorf("DiffPositions() = %v, want %v", got, tt.want)
			}
			for i, g := range got {
				if !g.Equal(tt.want[i]) {
					t.Errorf("DiffPositions() at %v = %v, want %v", i, g, tt.want[i])
				}
			}
		})
	}
}

func fn() float64 {
	return 2.1
}

func TestMath(t *testing.T) {
	a := 4.5 / (fn() - 2.0)
	if math.IsInf(a, 0) {
		a = 3.0
	}
	t.Log(a)
}
