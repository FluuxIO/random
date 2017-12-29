package random // import "fluux.io/random"

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

const testsRandomNumber = 100000

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// TestRandomSize checks RandomUnsafe.Size returns valid values.
func TestRandomSize(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	r := NewRandomUnsafe()
	for i := 0; i < testsRandomNumber; i++ {
		n := r.Size()
		if n == nil {
			continue
		}
		if n.Value < 1 {
			t.Errorf("%v < 1)", n.Value)
		}
	}
}

// TestRandomString checks that length of string returned by RandString
// is falling within the specified bounds.
func TestRandomString(t *testing.T) {
	r := NewRandomUnsafe()

	for i := 0; i < testsRandomNumber; i++ {
		min := int(r.Int(100))
		max := min + int(r.Int(100))
		s := r.String(min, max)
		l := len(s)
		if l < min || l > max {
			t.Errorf("wrong length: %q, (%d, %d))", s, min, max)
		}
	}
}

func TestRandomFixedLen(t *testing.T) {
	r := NewRandomUnsafe()
	for i := 0; i < testsRandomNumber; i++ {
		l := int(r.Int(100))
		s := r.FixedLenString(l)
		if len(s) != l {
			t.Errorf("wrong length: %q, (%d))", l, len(s))
		}
	}
}

// TestRandomBool checks that TestRandBool is balanced.
func TestRandomBool(tt *testing.T) {
	r := NewRandomUnsafe()

	var t, f int
	for i := 0; i < testsRandomNumber; i++ {
		b := r.Bool()
		if b {
			t++
		} else {
			f++
		}
	}

	min := testsRandomNumber/2 - (testsRandomNumber * 10 / 100)
	max := testsRandomNumber/2 + (testsRandomNumber * 10 / 100)
	if t < min || t > max || f < min || f > max {
		tt.Errorf("RandBool is not balanced, (%d true, %d false))", t, f)
	}
}

//=============================================================================
// Benchmarks

func BenchmarkRandString(b *testing.B) {
	min := 15
	max := 25

	for i := 0; i < b.N; i++ {
		randString(min, max)
	}
}

func BenchmarkRandomString(b *testing.B) {
	r := NewRandomUnsafe()
	for i := 0; i < b.N; i++ {
		r.String(15, 25)
	}
}

func BenchmarkRandId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randId("test")
	}
}

func BenchmarkRandomId(b *testing.B) {
	r := NewRandomUnsafe()
	for i := 0; i < b.N; i++ {
		r.RandomId("test")
	}
}

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Intn(1000)
	}
}

func BenchmarkRandomInt(b *testing.B) {
	r := NewRandomUnsafe()
	for i := 0; i < b.N; i++ {
		r.Int(1000)
	}
}

//=============================================================================
// Benchmark helpers

func randString(min, max int) string {
	n := rand.Intn(max-min) + min
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randId(prefix string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	id := []string{"test", randString(10, 20), timestamp}
	return strings.Join(id, "_")
}
