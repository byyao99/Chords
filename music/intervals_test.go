package music

import "testing"

func TestQualityIntervals_KnownQualities(t *testing.T) {
	cases := []struct {
		quality   string
		intervals []int
	}{
		{"", []int{0, 4, 7}},
		{"m", []int{0, 3, 7}},
		{"7", []int{0, 4, 7, 10}},
		{"m7", []int{0, 3, 7, 10}},
		{"maj7", []int{0, 4, 7, 11}},
		{"aug", []int{0, 4, 8}},
		{"dim", []int{0, 3, 6}},
		{"m7b5", []int{0, 3, 6, 10}},
		{"sus4", []int{0, 5, 7}},
		{"add9", []int{0, 4, 7, 14}},
	}
	for _, tc := range cases {
		f, ok := QualityIntervals(tc.quality)
		if !ok {
			t.Errorf("quality %q not found", tc.quality)
			continue
		}
		if len(f.Intervals) != len(tc.intervals) {
			t.Errorf("quality %q: got %v, want %v", tc.quality, f.Intervals, tc.intervals)
			continue
		}
		for i, iv := range f.Intervals {
			if iv != tc.intervals[i] {
				t.Errorf("quality %q interval[%d]: got %d, want %d", tc.quality, i, iv, tc.intervals[i])
			}
		}
	}
}

func TestQualityIntervals_Unknown(t *testing.T) {
	_, ok := QualityIntervals("xyz")
	if ok {
		t.Error("expected unknown quality to return false")
	}
}

func TestQualityIntervals_RequiredIsSubset(t *testing.T) {
	for quality, formula := range qualityToFormula {
		if formula.Required == nil {
			continue
		}
		allSet := toMod12Set(formula.Intervals)
		for _, r := range formula.Required {
			if !allSet[r%12] {
				t.Errorf("quality %q: required interval %d not in Intervals", quality, r)
			}
		}
	}
}

func TestChordTones(t *testing.T) {
	cases := []struct {
		name string
		want []string
	}{
		{"Am7", []string{"A", "C", "E", "G"}},
		{"C", []string{"C", "E", "G"}},
		{"F#m", []string{"F#", "A", "C#"}},
		{"Bbmaj7", []string{"Bb", "D", "F", "A"}},
		{"Edim", []string{"E", "G", "A#"}},
	}
	for _, tc := range cases {
		got, ok := ChordTones(tc.name)
		if !ok {
			t.Errorf("ChordTones(%q) returned false", tc.name)
			continue
		}
		if len(got) != len(tc.want) {
			t.Errorf("ChordTones(%q) = %v, want %v", tc.name, got, tc.want)
			continue
		}
		for i, n := range got {
			if n != tc.want[i] {
				t.Errorf("ChordTones(%q)[%d] = %q, want %q", tc.name, i, n, tc.want[i])
			}
		}
	}
}

func TestQualityIntervals_StartsWithZero(t *testing.T) {
	for quality, formula := range qualityToFormula {
		if len(formula.Intervals) == 0 || formula.Intervals[0] != 0 {
			t.Errorf("quality %q: Intervals must start with 0, got %v", quality, formula.Intervals)
		}
	}
}
