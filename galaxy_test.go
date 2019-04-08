package galaxy

import "testing"

var locCases = []struct {
	name   string
	dns    DNS
	sector int
	want   float64
}{
	{
		name: "zero",
		want: 0,
	},
	{
		name: "sector is zero",
		dns: DNS{
			X:   2.5,
			Y:   2.4,
			Z:   1.0,
			Vel: 3.2,
		},
		want: 3.2,
	},
	{
		name: "ok #1",
		dns: DNS{
			X:   2.5,
			Y:   2.4,
			Z:   1.0,
			Vel: 3.2,
		},
		sector: 4,
		want:   26.8,
	},
	{
		name: "ok #2",
		dns: DNS{
			X:   4.5,
			Y:   2.4,
			Z:   3.0,
			Vel: 2.4,
		},
		sector: 6,
		want:   61.8,
	},
}

func TestDNS(t *testing.T) {
	for _, tt := range locCases {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dns.Loc(tt.sector)

			if got != tt.want {
				t.Errorf("Expected DNS(%v).Loc(%v) = %v, got %v instead",
					tt.dns, tt.sector, tt.want, got)
			}
		})
	}
}
