package random

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"math/rand"
)

// String generate a random English string containing both lowercase and uppercase letters (a-z, A-Z)
func String(length uint) string {
	if length == 0 {
		return ""
	}
	b := make([]byte, length)
	for i := uint(0); i < length; i++ {
		base := 65
		if Bool() {
			base = 97
		}
		b[i] = byte(rand.Intn(26) + base)
	}
	return string(b)
}

// UInt Generate a random non-negative integer within the range [min, max).
func UInt(min, max uint) uint {
	if min > max {
		min, max = max, min
	}

	if min == max {
		return min
	} else {
		return uint(rand.Intn(int(max-min)) + int(min))
	}
}

// Bool Generate a random true or false.
func Bool() bool {
	return 0 == rand.Intn(2)
}

// Bin Generate a random binary.
func Bin(size int) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	for i := 0; i < size/4; i++ {
		var b uint32
		b = uint32(UInt(0, 4294967295))
		_ = binary.Write(bytesBuffer, binary.BigEndian, &b)
	}
	return bytesBuffer.Bytes()
}

// Json Generate a random JSON file.
func Json(sizeMin, sizeMax uint) []byte {
	j := make(map[string]string)
	targetSize := UInt(sizeMin, sizeMax)
	var curSize uint = 0
	for curSize == 0 || curSize < targetSize {
		j[String(UInt(1, 10))] = String(UInt(1, 10))
		data, err := json.Marshal(j)
		if err != nil {
			panic(err)
		}
		curSize = uint(len(data))
	}

	data, _ := json.Marshal(j)
	return data
}

// Png Generate a random PNG file.
func Png() []byte {
	width := UInt(255, 1024)
	height := UInt(255, 1024)
	var rgbList [10][3]uint
	for i := 0; i < 10; i++ {
		rgbList[i] = [3]uint{UInt(1, 255), UInt(1, 255), UInt(1, 255)}
	}
	img := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8((x + y) & 255),
				G: uint8((x + y) << 1 & 255),
				B: uint8((x + y) << 2 & 255),
				A: 255,
			})
		}
	}

	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := png.Encode(bytesBuffer, img); err != nil {
		panic(err)
	}

	return bytesBuffer.Bytes()
}
