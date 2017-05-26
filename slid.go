package slid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// Atomically-incremented in the SLID constructors. A 32-bit value appended to a 64-bit nanosecond
// timestamp gives us 2.14e22 unique values per millsecond (assuming nanosecond resolution in the
// system's high-performance timer.
//
// Note that on Windows, the high-performance timer responds in 100ns increments, which reduces
// the possible number of unique values by two orders of magnitude.
var counter uint32

// SLID is a struct containing a 64-bit timestamp value and a 32-bit counter value. The timestamp
// is only valid until some time in the year 2262, so if you're for some reason using this around
// that timeframe, beware of overflow and probably also of giant xenoinsects.
type SLID struct {
	Timestamp uint64
	Counter   uint32
}

// New creates a SLID struct with the current time and the next integer in the sequence, incremented
// atomically.
func New() SLID {
	return SLID{
		Timestamp: uint64(time.Now().UnixNano()),
		Counter:   atomic.AddUint32(&counter, 1),
	}
}

// NewFrom created a SLID struct from the provided timestamp and the next integer in the sequence,
// incremented atomically.
func NewFrom(timestamp time.Time) SLID {
	return SLID{
		Timestamp: uint64(timestamp.UnixNano()),
		Counter:   atomic.AddUint32(&counter, 1),
	}
}

// Time returns the timestamp of a SLID as a Time object.
func (slid SLID) Time() time.Time {
	return time.Unix(0, int64(slid.Timestamp))
}

// Encode returns a byte slice of the encoded SLID. The first eight bytes comprise the reversed
// timestamp, and the following four bytes comprise the counter in natural byte order.
func (slid SLID) Encode() []byte {
	bs := make([]byte, 12) // We need 12 bytes to store uint64 + uint32.
	binary.LittleEndian.PutUint64(bs, uint64(slid.Timestamp))

	// Reverse the timestamp bytes so the SLID can be sorted lexicographically.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}

	binary.LittleEndian.PutUint32(bs[8:], counter)

	return bs
}

// Hex returns the hexadecimal string of a SLID.
func (slid SLID) Hex() string {
	return hex.EncodeToString(slid.Encode())
}

func (slid SLID) String() string {
	return hex.EncodeToString(slid.Encode())
}

// Decode accepts a 12-byte slice, decodes its value and returns a SLID.
func Decode(slid []byte) (SLID, error) {
	if len(slid) != 12 {
		return SLID{}, fmt.Errorf("Cannot create SLID with [%d]byte", len(slid))
	}

	// Reverse the first 8 bytes to get back the original byte representation of the timestamp.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		slid[i], slid[j] = slid[j], slid[i]
	}

	return SLID{
		Timestamp: binary.LittleEndian.Uint64(slid[0:8]),
		Counter:   binary.LittleEndian.Uint32(slid[8:]),
	}, nil
}
