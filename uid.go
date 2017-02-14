package uid

import (
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"
)

// Atomically-incremented in the UID constructors. A 32-bit value appended to a 64-bit nanosecond
// timestamp gives us 2.14e22 unique values per millsecond.
var counter uint32

// UID is a struct containing a 64-bit timestamp value and a 32-bit counter value. The timestamp
// is only valid until some time in the year 2262, so if you're for some reason using this around
// that timeframe, beware of overflow and probably also of giant xenoinsects.
type UID struct {
	timestamp uint64
	counter   uint32
}

// New creates a UID struct with the current time and the next integer in the sequence, incremented
// atomically.
func New() UID {
	return UID{
		timestamp: uint64(time.Now().UnixNano()),
		counter:   atomic.AddUint32(&counter, 1),
	}
}

// NewFrom created a UID struct from the provided timestamp and the next integer in the sequence,
// incremented atomically.
func NewFrom(timestamp time.Time) UID {
	return UID{
		timestamp: uint64(timestamp.UnixNano()),
		counter:   atomic.AddUint32(&counter, 1),
	}
}

// Time returns the timestamp of a UID as a Time object.
func (uid UID) Time() time.Time {
	return time.Unix(0, int64(uid.timestamp))
}

// Encode
func (uid UID) Encode() []byte {
	bs := make([]byte, 12) // We need 12 bytes to store uint64 + uint32.
	binary.LittleEndian.PutUint64(bs, uint64(uid.timestamp))

	// Reverse the timestamp bytes so the UID can be sorted lexicographically.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		bs[i], bs[j] = bs[j], bs[i]
	}

	binary.LittleEndian.PutUint32(bs[8:], counter)

	return bs
}

// Decode accepts a 12-byte slice, decodes its value and returns a UID.
func Decode(uid []byte) UID {
	if len(uid) != 12 {
		return UID{}
	}

	// Reverse the first 8 bytes to get back the original byte representation of the timestamp.
	for i, j := 0, 7; i < j; i, j = i+1, j-1 {
		uid[i], uid[j] = uid[j], uid[i]
	}

	return UID{
		timestamp: binary.LittleEndian.Uint64(uid[0:8]),
		counter:   binary.LittleEndian.Uint32(uid[8:]),
	}
}

// Hex returns the hexadecimal string of a UID.
func (uid UID) Hex() string {
	return hex.EncodeToString(uid.Encode())
}
