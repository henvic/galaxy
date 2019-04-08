package galaxy

// DNS of a sector.
type DNS struct {
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
	Z   float64 `json:"z"`
	Vel float64 `json:"vel"`
}

// Loc of a sector.
func (d *DNS) Loc(sectorID int) float64 {
	sid := float64(sectorID)
	return d.X*sid + d.Y*sid + d.Z*sid + d.Vel
}
