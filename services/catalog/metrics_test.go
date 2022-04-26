package catalog

import (
	"testing"
)

func Test_metrics_IncreaseFind(t *testing.T) {
	tests := []struct {
		name                          string
		m                             *metrics
		wantFind, wantDone, wantError int
	}{
		{
			name: "First",
			m: &metrics{
				Find:  1,
				Done:  0,
				Error: 0,
			},
			wantFind:  2,
			wantDone:  0,
			wantError: 0,
		},
		{
			name: "Second",
			m: &metrics{
				Find:  2,
				Done:  1,
				Error: 1,
			},
			wantFind:  3,
			wantDone:  1,
			wantError: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.IncreaseFind()
		})

		if tt.m.Find != tt.wantFind {
			t.Errorf("metrics.GetStatistic() = %v, want Find %v", tt.m.Find, tt.wantFind)
		}

		if tt.m.Done != tt.wantDone {
			t.Errorf("metrics.GetStatistic() = %v, want Done %v", tt.m.Done, tt.wantDone)
		}

		if tt.m.Error != tt.wantError {
			t.Errorf("metrics.GetStatistic() = %v, want Done %v", tt.m.Error, tt.wantError)
		}

	}
}
