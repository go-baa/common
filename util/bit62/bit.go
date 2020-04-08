package bit62

import (
	"math"
)

var bit uint64 = 62
var chars = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

// Bit10to62 10进制到62进制
func Bit10to62(i uint64) string {
	var t []byte
	var r uint64
	for i >= bit {
		r = i % bit
		t = append(t, chars[r])
		i = i / bit
	}
	t = append(t, chars[i])

	tr := make([]byte, len(t))
	for i, b := range t {
		tr[len(t)-i-1] = b
	}
	return string(tr)
}

// Bit62to10 62进制到10进制
func Bit62to10(i string) uint64 {
	var t uint64
	for r := 0; r < len(i); r++ {
		t += toUint64(i[r]) * uint64(math.Pow(float64(bit), float64(len(i)-r-1)))
	}
	return uint64(t)
}

func toUint64(b byte) uint64 {
	i := charsRevert[b]
	return uint64(i)
}

var charsRevert map[byte]int

func init() {
	charsRevert = make(map[byte]int)
	for i, b := range chars {
		charsRevert[b] = i
	}
}
