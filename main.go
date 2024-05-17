package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

// header is 14 bytes
type FileHeader struct {
	Signature  uint16 // 2 bytes
	FileSize   uint32 // 4 bytes
	Reserved   uint32 //4 bytes
	DataOffset uint32 // 4 bytes
}

// infoheader is 40 bytes
type InfoHeader struct {
	Size            uint32 // 4 bytes
	Width           int32  // 4 bytes
	Height          int32  // 4 bytes
	Planes          uint16 // 2 bytes
	BitCount        uint16 // 2 bytes
	Compression     uint32 // 4 bytes
	ImageSize       uint32 // 4 bytes
	XpixelsPerM     uint32 // 4 bytes
	YpixelsPerM     uint32 // 4 bytes
	ColorsUsed      uint32 // 4 bytes
	ImportantColors uint32 // 4 bytes
}

// $ * NumColors bytes
// 0036h offset!
type ColorTable struct {
	Red      uint8
	Green    uint8
	Blue     uint8
	Reserved uint8
}

func main() {
	file, err := os.Open("snail.bmp")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var fileHeader FileHeader
	err = binary.Read(file, binary.LittleEndian, &fileHeader)
	if err != nil {
		fmt.Println(err)
		return
	}

	var infoHeader InfoHeader
	err = binary.Read(file, binary.LittleEndian, &infoHeader)
	if err != nil {
		fmt.Println(err)
		return
	}

	rowSize := ((int(infoHeader.BitCount)*int(infoHeader.Width) + 31) / 32) * 4

	_, err = file.Seek(int64(fileHeader.DataOffset), 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	pixelData := make([]byte, rowSize*int(infoHeader.Height))
	_, err = file.Read(pixelData)
	if err != nil {
		fmt.Println(err)
		return
	}

	asciiChars := "MNHQ$OC?7>!:-;. "
	width := int(infoHeader.Width)
	height := int(infoHeader.Height)

	// scaling factors
	scaleX, scaleY := 2, 4

	for y := height - 1; y >= 0; y -= scaleY {
		for x := 0; x < width; x += scaleX {
			var sumR, sumG, sumB float64
			count := 0

			for subY := 0; subY < scaleY && (y-subY) >= 0; subY++ {
				for subX := 0; subX < scaleX && (x+subX) < width; subX++ {
					offset := (y-subY)*rowSize + (x+subX)*3
					blue := pixelData[offset]
					green := pixelData[offset+1]
					red := pixelData[offset+2]

					sumR += float64(red)
					sumG += float64(green)
					sumB += float64(blue)
					count++
				}
			}

			avgR := sumR / float64(count)
			avgG := sumG / float64(count)
			avgB := sumB / float64(count)

			gray := 0.3*avgR + 0.59*avgG + 0.11*avgB
			charIndex := int((gray / 255.0) * float64(len(asciiChars)-1))
			fmt.Print(string(asciiChars[charIndex]))
		}
		fmt.Println()
	}
}
