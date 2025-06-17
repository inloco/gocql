package varint

import (
	"math"
	"math/big"
	"math/rand"
	"testing"
)

func TestDec2BigInt(t *testing.T) {
	t.Parallel()

	genData := func(v int64) []byte {
		data := []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
		out := make([]byte, 0)
		add := false
		for i, b := range data {
			if !add {
				if v < 0 {
					if b != 255 || b == 255 && data[i+1] < 128 {
						add = true
					} else {
						continue
					}
				} else {
					if b != 0 || b == 0 && data[i+1] > 127 {
						add = true
					} else {
						continue
					}
				}
			}
			out = append(out, b)
		}

		return out
	}

	t.Run("positive", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(rand.Int63()))
		for i := int64(math.MaxInt16); i < 1<<23; i = i + int64(rnd.Int31n(300)) {
			data := genData(i)
			expected := big.NewInt(i)

			received := Dec2BigInt(data)
			if expected.Cmp(received) != 0 {
				t.Fatalf("%d\nexpected:%s\nreceived:%s", i, expected, received)
			}

			_ = DecBigInt(data, received)
			if expected.Cmp(received) != 0 {
				t.Fatalf("%d\nexpected:%s\nreceived:%s", i, expected, received)
			}
		}
	})

	t.Run("negative", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(rand.Int63()))
		for i := int64(math.MinInt16); i > -1<<23; i = i - int64(rnd.Int31n(300)) {
			data := genData(i)
			expected := big.NewInt(i)

			received := Dec2BigInt(data)
			if expected.Cmp(received) != 0 {
				t.Fatalf("%d\nexpected:%s\nreceived:%s", i, expected, received)
			}

			_ = DecBigInt(data, received)
			if expected.Cmp(received) != 0 {
				t.Fatalf("%d\nexpected:%s\nreceived:%s", i, expected, received)
			}
		}
	})
}
