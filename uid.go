package uid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

// Atomically-incremented in the UID constructors. A 32-bit value appended to a 64-bit nanosecond
// timestamp gives us 2.14e22 unique values per millsecond (assuming nanosecond resolution in the
// system's high-performance timer.
//
// Note that on Windows, the high-performance timer responds in 100ns increments, which reduces
// the possible number of unique values by two orders of magnitude.
var counter uint32

// UID is a struct containing a 64-bit timestamp value and a 32-bit counter value. The timestamp
// is only valid until some time in the year 2262, so if you're for some reason using this around
// that timeframe, beware of overflow and probably also of giant xenoinsects.
type UID struct {
	Timestamp uint64
	Counter   uint32
}

// New creates a UID struct with the current time and the next integer in the sequence, incremented
// atomically.
func New() UID {
	return UID{
		Timestamp: uint64(time.Now().UnixNano()),
		Counter:   atomic.AddUint32(&counter, 1),
	}
}

// NewFrom created a UID struct from the provided timestamp and the next integer in the sequence,
// incremented atomically.
func NewFrom(timestamp time.Time) UID {
	return UID{
		Timestamp: uint64(timestamp.UnixNano()),
		Counter:   atomic.AddUint32(&counter, 1),
	}
}

// Time returns the timestamp of a UID as a Time object.
func (uid UID) Time() time.Time {
	return time.Unix(0, int64(uid.Timestamp))
}

// Encode returns a byte slice of the encoded UID. The first eight bytes comprise the reversed
// timestamp, and the following four bytes comprise the counter in natural byte order.
func (uid UID) Encode() []byte {
	bs := make([]byte, 12) // We need 12 bytes to store uint64 + uint32.
	binary.LittleEndian.PutUint64(bs, uint64(uid.Timestamp))

	// Reverse the timestamp bytes so the UID can be sorted lexicographically.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}

	binary.LittleEndian.PutUint32(bs[8:], counter)

	return bs
}

// Hex returns the hexadecimal string of a UID.
func (uid UID) Hex() string {
	return hex.EncodeToString(uid.Encode())
}

// Decode accepts a 12-byte slice, decodes its value and returns a UID.
func Decode(uid []byte) (UID, error) {
	if len(uid) != 12 {
		return UID{}, fmt.Errorf("Cannot create UID with [%d]byte", len(uid))
	}

	// Reverse the first 8 bytes to get back the original byte representation of the timestamp.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		uid[i], uid[j] = uid[j], uid[i]
	}

	return UID{
		Timestamp: binary.LittleEndian.Uint64(uid[0:8]),
		Counter:   binary.LittleEndian.Uint32(uid[8:]),
	}, nil
}
