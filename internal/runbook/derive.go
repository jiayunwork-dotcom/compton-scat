package runbook

import (
	"sort"
)

func (b *Book) AverageRecoil() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.RecoilKEV
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func (b *Book) MaxLoss() (float64, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	best := 0.0
	id := ""
	for _, e := range b.items {
		if e.RecoilKEV > best {
			best = e.RecoilKEV
			id = e.ID
		}
	}
	return best, id
}

func (b *Book) ForwardScatterCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.AngleDeg == 0 {
			n++
		}
	}
	return n
}

func (b *Book) BackscatterCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.AngleDeg == 180 {
			n++
		}
	}
	return n
}

func (b *Book) Similar(target Entry, angleTol float64) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range b.items {
		if e.ID == target.ID {
			continue
		}
		dAngle := e.AngleDeg - target.AngleDeg
		if dAngle < 0 {
			dAngle = -dAngle
		}
		dEnergy := relDiff(e.EnergyKEV, target.EnergyKEV)
		if dAngle <= angleTol && dEnergy <= angleTol {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (b *Book) MeanScattered() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.ScatteredKEV
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func relDiff(a, b float64) float64 {
	if b == 0 {
		return 1
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d / b
}
