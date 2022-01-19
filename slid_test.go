package slid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Overwrite time.Now() to always return a fixed time value.
	now = func() time.Time {
		t, _ := time.Parse(time.RFC3339, time.RFC3339)
		return t
	}
}

func TestNew(t *testing.T) {
	got := New()
	want := []byte{161, 178, 3, 235, 61, 26, 0, 0, 1, 0, 0, 0}

	assert.Equal(t, want, got.Encode())
}

func TestDecode(t *testing.T) {
	assert := assert.New(t)

	slid := New()
	encoded := slid.Encode()

	got, err := Decode(encoded)
	if !assert.NoError(err) {
		t.FailNow()
	}

	assert.Equal(slid, got)
}

func TestTime(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	slid := NewFrom(now)
	got := slid.Time()

	assert.True(now.Equal(got))
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}

func BenchmarkEncode(b *testing.B) {
	slid := New()

	for n := 0; n < b.N; n++ {
		slid.Encode()
	}
}

func BenchmarkNewEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New().Encode()
	}
}

func TestHex(t *testing.T) {
	got := New().Hex()
	want := "a1b203eb3d1a000004000000"

	assert.Equal(t, want, got)
}

func BenchmarkHex(b *testing.B) {
	slid := New()

	for n := 0; n < b.N; n++ {
		slid.Hex()
	}
}
