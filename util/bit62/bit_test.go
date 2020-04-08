package bit62

import "testing"

func TestTo62(t *testing.T) {
	var i1 uint64 = 10000000000
	var i2 uint64 = 1000000000
	var i3 uint64 = 101
	var s1 = "Aukyoa"
	var s2 = "15ftgG"
	var s3 = "1d"

	ii1 := Bit10to62(i1)
	ii2 := Bit10to62(i2)
	ii3 := Bit10to62(i3)

	if ii1 != s1 {
		t.Errorf("%d to bit62 should be: %s, but: %s", i1, s1, ii1)
	}
	if ii2 != s2 {
		t.Errorf("%d to bit62 should be: %s, but: %s", i2, s2, ii2)
	}
	if ii3 != s3 {
		t.Errorf("%d to bit62 should be: %s, but: %s", i3, s3, ii3)
	}
}

func TestTo10(t *testing.T) {
	var i1 uint64 = 10000000000
	var i2 uint64 = 1000000000
	var i3 uint64 = 102
	var s1 = "Aukyoa"
	var s2 = "15ftgG"
	var s3 = "1e"

	ii1 := Bit62to10(s1)
	ii2 := Bit62to10(s2)
	ii3 := Bit62to10(s3)

	if ii1 != i1 {
		t.Errorf("%s to bit10 should be: %d, but: %d", s1, i1, ii1)
	}
	if ii2 != i2 {
		t.Errorf("%s to bit10 should be: %d, but: %d", s2, i2, ii2)
	}
	if ii3 != i3 {
		t.Errorf("%s to bit10 should be: %d, but: %d", s3, i3, ii3)
	}
}
