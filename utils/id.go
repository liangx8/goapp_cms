package utils

import (
	"math/rand"
	"time"
)

var arr = []uint64{
	62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62 * 62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62 * 62,
	62 * 62 * 62 * 62,
	62 * 62 * 62,
	62 * 62,
	62,
}

func ser(v uint64) []uint8 {
	mx := make([]uint8, 0)
	for v >= arr[0]*62 {
		v = v - arr[0]*62
	}
	for _, a := range arr {
		d := uint8(0)
		for v > a {
			v -= a
			d++
		}
		mx = append(mx, d)
	}
	mx = append(mx, uint8(v))
	return mx
}
func char(idx int) rune {
	if idx >= 62 {
		panic("greater than 62")
	}
	if idx < 26 {
		return rune(int('a') + idx)
	}
	idx -= 26
	if idx < 10 {
		return rune(int('0' + idx))
	}
	idx -= 10
	return rune(int('A') + idx)

}
func init() {
	rand.Seed(time.Now().Unix())
}
func MakeID() string {
	id := make([]rune, 0)
	d8 := ser(rand.Uint64())

	for _, d := range d8 {
		id = append(id, char(int(d)))
	}
	l := len(d8)
	for i := l; i < 8; i++ {
		id = append(id, '=')
	}
	return string(id)
}
func RandomSalt(salt []byte) {
	if _, err := rand.Read(salt); err != nil {
		panic(err)
	}

}
