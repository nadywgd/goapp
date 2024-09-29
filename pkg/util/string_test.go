package util

import "testing"

func TestRandHexString(t *testing.T) {
	for i := 0; i < 100; i++ {
		length := 10
		result := RandHexString(length)
		if len(result) != length {
			t.Errorf("expected length %d but got %d", length, len(result))
		}

		for _, char := range result {
			if !isHexChar(char) {
				t.Errorf("unexpected character %c in result %s", char, result)
			}
		}
	}

}

func isHexChar(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')
}

func BenchmarkRandHexString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandHexString(10)
	}
}
