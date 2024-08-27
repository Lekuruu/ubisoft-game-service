package common

import (
	"math"
)

// GS XOR encryption algorithm
func GSXOREncrypt(input []byte) []byte {
	size := len(input)
	result := make([]byte, size)
	copy(result, input)
	for i := 0; i < size; i++ {
		result[i] ^= byte((i - 119) & 0xff)
	}

	sizeRoot := int(math.Sqrt(float64(size)))
	if sizeRoot*sizeRoot < size {
		sizeRoot++
	}

	newSize := 2 * sizeRoot * sizeRoot
	buf := make([]byte, newSize)
	for i := range buf {
		buf[i] = 0xff
	}

	a, b := 0, 0
	for i := 0; i < size; i++ {
		if a < sizeRoot {
			if b < 0 {
				b = a
				a = 0
			}
		} else {
			a = b + 2
			b = sizeRoot - 1
		}
		buf[a+sizeRoot*b] = result[i]
		a++
		b--
	}

	idx := 0
	for j := 0; j < sizeRoot; j++ {
		for k := 0; k < sizeRoot; k++ {
			if buf[k+sizeRoot*j] != 0xff {
				result[idx] = buf[k+sizeRoot*j]
				idx++
			}
		}
	}

	return result
}

// GS XOR decryption algorithm
func GSXORDecrypt(input []byte) []byte {
	size := len(input)
	result := make([]byte, size)
	copy(result, input)

	root := math.Sqrt(float64(size))
	sizeRoot := int(root)
	if float64(sizeRoot) < root {
		sizeRoot++
	}
	newSize := sizeRoot * sizeRoot
	buf := make([]byte, newSize)

	a, b := 0, 0
	if size > 0 {
		sizeCpy := size
		for sizeCpy > 0 {
			if b < sizeRoot {
				if a < 0 {
					a = b
					b = 0
				}
			} else {
				b = a + 2
				a = sizeRoot - 1
			}
			buf[b+sizeRoot*a] = 1
			a--
			b++
			sizeCpy--
		}
	}

	c, d := 0, 0
	if size > 0 {
		count := 0
		for d < size {
			if c >= sizeRoot {
				count += sizeRoot
				c = 0
			}
			if buf[count+c] > 0 {
				buf[count+c] = input[d]
				d++
			}
			c++
		}
	}

	e, f := 0, 0
	for idx := 0; idx < size; idx++ {
		if f < sizeRoot {
			if e < 0 {
				e = f
				f = 0
			}
		} else {
			f = e + 2
			e = sizeRoot - 1
		}
		result[idx] = buf[f+sizeRoot*e]
		e--
		f++
	}

	for i := 0; i < size; i++ {
		result[i] ^= byte((i - 119) & 0xff)
	}

	return result
}
