package runbook

import (
	"errors"
	"testing"
)

func TestBookAddGetRemove(t *testing.T) {
	b := NewBook(4)
	e := Entry{
		ID:           "run-1",
		EnergyKEV:    511,
		AngleDeg:     90,
		ScatteredKEV: 255.5,
		RecoilKEV:    255.5,
	}
	if err := b.Add(e); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 1 || b.NextSeq() != 1 {
		t.Fatalf("len=%d seq=%d", b.Len(), b.NextSeq())
	}
	got, ok := b.Get("run-1")
	if !ok || got.EnergyKEV != 511 {
		t.Fatalf("get failed: %+v", got)
	}
	if err := b.Add(e); !errors.Is(err, ErrExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if !b.Remove("run-1") {
		t.Fatal("remove failed")
	}
}

func TestBookRenameFreezeSetNote(t *testing.T) {
	b := NewBook(8)
	if err := b.Add(Entry{ID: "a", EnergyKEV: 100, AngleDeg: 60, ScatteredKEV: 80, RecoilKEV: 20}); err != nil {
		t.Fatal(err)
	}
	if err := b.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := b.Freeze("b"); err != nil {
		t.Fatal(err)
	}
	if err := b.SetNote("b", "changed"); !errors.Is(err, ErrFrozen) {
		t.Fatalf("frozen set note err=%v", err)
	}
}

func TestValidateRejectsBadEntries(t *testing.T) {
	bad := []Entry{
		{ID: "", EnergyKEV: 1, AngleDeg: 0, ScatteredKEV: 1, RecoilKEV: 0},
		{ID: "x", EnergyKEV: 0, AngleDeg: 0, ScatteredKEV: 1, RecoilKEV: 0},
		{ID: "x", EnergyKEV: 1, AngleDeg: -1, ScatteredKEV: 1, RecoilKEV: 0},
		{ID: "x", EnergyKEV: 1, AngleDeg: 181, ScatteredKEV: 1, RecoilKEV: 0},
		{ID: "x", EnergyKEV: 1, AngleDeg: 0, ScatteredKEV: -1, RecoilKEV: 0},
		{ID: "x", EnergyKEV: 1, AngleDeg: 0, ScatteredKEV: 1, RecoilKEV: -1},
	}
	for i, e := range bad {
		if err := e.Validate(); err == nil {
			t.Fatalf("entry %d should fail", i)
		}
	}
}

func TestDerivedStats(t *testing.T) {
	b := NewBook(16)
	for _, e := range []Entry{
		{ID: "a", EnergyKEV: 100, AngleDeg: 0, ScatteredKEV: 100, RecoilKEV: 0},
		{ID: "b", EnergyKEV: 200, AngleDeg: 180, ScatteredKEV: 50, RecoilKEV: 150},
		{ID: "c", EnergyKEV: 300, AngleDeg: 90, ScatteredKEV: 120, RecoilKEV: 180},
	} {
		if err := b.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	avg, n := b.AverageRecoil()
	if n != 3 || avg != 110 {
		t.Fatalf("avg=%v n=%d", avg, n)
	}
	maxLoss, id := b.MaxLoss()
	if id != "c" || maxLoss != 180 {
		t.Fatalf("maxLoss=%v id=%s", maxLoss, id)
	}
	if b.ForwardScatterCount() != 1 || b.BackscatterCount() != 1 {
		t.Fatalf("forward=%d back=%d", b.ForwardScatterCount(), b.BackscatterCount())
	}
	sim := b.Similar(Entry{ID: "probe", EnergyKEV: 105, AngleDeg: 0, ScatteredKEV: 105, RecoilKEV: 0}, 0.1)
	if len(sim) != 1 || sim[0].ID != "a" {
		t.Fatalf("similar=%+v", sim)
	}
	mean, mn := b.MeanScattered()
	if mn != 3 || mean != 90 {
		t.Fatalf("mean scattered=%v n=%d", mean, mn)
	}
}
