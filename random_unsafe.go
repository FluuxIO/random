/*
Fast Random Generator for use from single threaded data injection code. It is
especially useful if you need to load test a service while generating random
data. It that case, the random generator locks can become a bottleneck.
*/
package random // import "fluux.io/random"

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/ptypes/wrappers"
)

//=============================================================================
// Unsafe but faster random generator

// RandomUnsafe is a structure wrapping non-thread safe random generator
// to use from a single go routine.
// It is more efficient than the default generator as it avoid using the mutex
// locks used as default for thread safety.
// It is intended to be used in part of code that use random value heavily.
type RandomUnsafe struct {
	src *rand.Rand
	// preallocated random string
	prealloc []byte
	// Cache for generating boolean number more efficiently
	boolcache int64
	boolcount int
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numbers = "0123456789"
const stringSeedSize = 10000

// NewRandomUnsafe creates an initialized random generator to use from a
// single go routine.
func NewRandomUnsafe() RandomUnsafe {
	src := rand.New(rand.NewSource(time.Now().Unix()))
	prealloc := make([]byte, stringSeedSize)
	for i := range prealloc {
		prealloc[i] = letters[src.Int63()%int64(len(letters))]
	}
	return RandomUnsafe{src: src, prealloc: prealloc}
}

// NumString returns a random string containing numbers.
func (r *RandomUnsafe) NumString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = numbers[r.src.Int63()%int64(len(numbers))]
	}
	return ByteSliceToString(b)
}

// Length generates an integer between min and max.
func (r *RandomUnsafe) Length(min, max int) int {
	if min > max {
		return 0
	}
	if min == max {
		return min
	}
	return r.src.Intn(max-min) + min
}

// String returns a random string of random length between min and max.
func (r *RandomUnsafe) String(min, max int) string {
	length := r.Length(min, max)
	return r.FixedLenString(length)
}

// FixedLenString returns a random string of n chars.
func (r *RandomUnsafe) FixedLenString(n int) string {
	pos := r.src.Intn(stringSeedSize - n)
	return ByteSliceToString(r.prealloc[pos : pos+n])
}

// Bool returns a random boolean. This function uses a cache to only trigger call to random number
// generator every 63 calls. We generate a 63 bits number and then use each bits as one random boolean.
func (r *RandomUnsafe) Bool() bool {
	if r.boolcount == 0 {
		r.boolcache, r.boolcount = r.src.Int63(), 63
	}

	result := r.boolcache&0x01 == 1
	r.boolcache >>= 1
	r.boolcount--

	return result
}

// OptBool return an optional random boolean.
func (r *RandomUnsafe) OptBool() *wrappers.BoolValue {
	if !r.Bool() {
		return nil
	}
	return &wrappers.BoolValue{Value: r.Bool()}
}

// Int returns a random int32.
func (r *RandomUnsafe) Int(n int) int32 {
	return int32(r.src.Intn(n))
}

// OptInt32 returns a optional random int32.
func (r *RandomUnsafe) OptInt32(n int) *wrappers.Int32Value {
	if !r.Bool() {
		return nil
	}
	return &wrappers.Int32Value{Value: r.Int(n)}
}

// OptInt64 returns a optional random int64.
func (r *RandomUnsafe) OptInt64(n int) *wrappers.Int64Value {
	if !r.Bool() {
		return nil
	}
	return &wrappers.Int64Value{Value: int64(r.Int(n))}
}

// Date returns a random recent date formatted as string.
func (r *RandomUnsafe) Date() string {
	min := time.Now().AddDate(0, 0, -5).Unix() // 5 days ago
	max := time.Now().Unix()                   // Now
	delta := max - min

	sec := r.src.Int63n(delta) + min
	return time.Unix(sec, 0).Format(time.RFC3339)
}

// OptString returns an optional random string of random length between min and
// max.
func (r *RandomUnsafe) OptString(min, max int) *wrappers.StringValue {
	if !r.Bool() {
		return nil
	}
	return &wrappers.StringValue{Value: r.String(min, max)}
}

// Size returns a physical measure for an object using a normal distribution.
func (r *RandomUnsafe) Size() *wrappers.Int32Value {
	if !r.Bool() {
		return nil
	}
	var size int32
	for ; size <= 0; size = int32(r.src.NormFloat64()*2500 + 3000) {
	}
	return &wrappers.Int32Value{Value: size}
}

// RandomId returns a random string to use as id, starting with prefix.
func (r *RandomUnsafe) RandomId(prefix string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	id := []string{prefix, r.String(10, 20), timestamp}
	return strings.Join(id, "_")
}

func (r *RandomUnsafe) Code(prefix string, i int) string {
	code := []string{prefix, strconv.Itoa(i), r.String(10, 32)}
	return strings.Join(code, "_")
}

// ByteSliceToString is used when you really want to convert a slice // of bytes to a string without incurring overhead.
// It is only safe to use if you really know the byte slice is not going to change // in the lifetime of the string.
// The unsafe pointer operation will not be needed in Go 1.10, when this feature is added:
// https://github.com/golang/go/issues/18990
func ByteSliceToString(bs []byte) string {
	// This is copied from runtime. It relies on the string
	// header being a prefix of the slice header!
	return *(*string)(unsafe.Pointer(&bs))
}
