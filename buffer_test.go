package neat

import "testing"

func TestBufferReadWrite(t *testing.T) {
	b := NewBuffer().Write(int16(-25), "kaka", uint64(216126), float64(1.23456), []byte{255, 162})
	readable := b.Readable()
	if v := readable.Read16(); v != -25 {
		t.Log(v, readable)
		t.Fail()
	}
	if v := readable.ReadStr(); v != "kaka" {
		t.Log(v, readable)
		t.Fail()
	}
	if v := readable.ReadU64(); v != 216126 {
		t.Log(v, readable)
		t.Fail()
	}
	if v := readable.ReadF64(); v != 1.23456 {
		t.Log(v, readable)
		t.Fail()
	}
	if v := string(readable.ReadBytes()); v != string([]byte{255, 162}) {
		t.Log(v, readable)
		t.Fail()
	}
}
