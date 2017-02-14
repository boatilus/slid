package uid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	got := New()

	t.Log(got)
	t.Log(New())

	counter = 0
}

func TestEncode(t *testing.T) {
	t.Log(New().Encode())

	counter = 0
}

func TestDecode(t *testing.T) {
	assert := assert.New(t)

	uid := New()
	encoded := uid.Encode()

	got := Decode(encoded)

	assert.Equal(uid, got)
}

func TestTime(t *testing.T) {
	assert := assert.New(t)

	now := time.Now()
	uid := NewFrom(now)
	got := uid.Time()

	assert.Equal(now, got)
}

func BenchmarkNew(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New()
	}
}

func BenchmarkEncode(b *testing.B) {
	uid := New()

	for n := 0; n < b.N; n++ {
		uid.Encode()
	}
}

func BenchmarkNewEncode(b *testing.B) {
	for n := 0; n < b.N; n++ {
		New().Encode()
	}
}

func TestHex(t *testing.T) {
	got := New().Hex()

	t.Log(got)
}

func BenchmarkHex(b *testing.B) {
	uid := New()

	for n := 0; n < b.N; n++ {
		uid.Hex()
	}
}
