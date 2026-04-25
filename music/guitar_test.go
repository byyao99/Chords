package music

import "testing"

func TestFindChordVoicing_Basic(t *testing.T) {
	tests := []string{"C", "D", "E", "F", "G", "A", "B", "Am", "Em", "Dm", "G7", "Fmaj7"}
	for _, name := range tests {
		v, ok := FindChordVoicing(name)
		if !ok {
			t.Errorf("FindChordVoicing(%q) not found", name)
			continue
		}
		if v.Name != name {
			t.Errorf("FindChordVoicing(%q).Name = %q", name, v.Name)
		}
	}
}

func TestFindChordVoicing_Alias(t *testing.T) {
	v, ok := FindChordVoicing("Db")
	if !ok {
		t.Fatal("FindChordVoicing(Db) not found")
	}
	if v.Name != "Db" {
		t.Errorf("alias should display as Db, got %q", v.Name)
	}
	// Should have same frets as C#
	c, _ := FindChordVoicing("C#")
	if v.Frets != c.Frets {
		t.Error("Db and C# should have same frets")
	}
}

func TestFindChordVoicing_NotFound(t *testing.T) {
	_, ok := FindChordVoicing("Xm7b5#11")
	if ok {
		t.Error("should not find unknown chord")
	}
}

func TestFindChordVoicing_Frets(t *testing.T) {
	v, _ := FindChordVoicing("Am")
	// Am = x02210
	want := [6]int{-1, 0, 2, 2, 1, 0}
	if v.Frets != want {
		t.Errorf("Am frets = %v, want %v", v.Frets, want)
	}
}

func TestAllChordNames(t *testing.T) {
	names := AllChordNames()
	if len(names) < 40 {
		t.Errorf("expected at least 40 chords, got %d", len(names))
	}
}
