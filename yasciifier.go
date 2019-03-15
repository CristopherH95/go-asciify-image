package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()	// need to be provided an argument for the file path
		return
	}
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)		// setup for png or jpeg
	image.RegisterFormat("jpeg", "jpg", jpeg.Decode, jpeg.DecodeConfig)
	path := os.Args[1]	// grab the path
	if _, err := os.Stat(path); err == nil {	// if the file exists begin processing
		fmt.Printf("File: %s\n", path)
		fmt.Println("Processing image of size: ")
		fmt.Println(getImageSize(path)) // print the image size
		file := ConvertImageToAscii(path)	// convert
		fmt.Printf("Output saved to file: %s\n", file)	// tell user where the output can be found
	} else {
		log.Print(err)
		fmt.Println("Could not find image file to read in.")
	}
}

func ConvertImageToAscii(path string) string {
	pix, err := getImagePixels(path)	// get image pixels
	checkErr(err)	//
	pixVal := pixelMatrixToBrightness(pix)
	ascii := brightnessMatrixToAscii(pixVal)
	file := writeMatrixToFile(path, ascii)
	return file
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	msg := `usage:
				yasciifier <filename>\n`
	fmt.Print(msg)
}

func writeMatrixToFile(path string, ascii [][]byte) string {
	file := filepath.Base(path) + ".txt"
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkErr(err)
	for _, row := range ascii {
		_, err := f.Write(row)
		if err != nil {
			delErr := os.Remove(file)
			checkErr(delErr)
			log.Fatal(err)
		}
	}

	return file
}

func getImageSize(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Failed to open file with error %v", err)
		return -1, -1
	}
	img, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Printf("Failed to decode image with error %v", err)
		return -1, -1
	}
	return img.Width, img.Height
}

func getImagePixels(path string) ([][]Pixel, error) {
	log.Println("Reading in image pixel values")
	file, err := os.Open(path)
	checkErr(err)
	img, _, err := image.Decode(file)
	checkErr(err)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToRGB(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

func pixelMatrixToBrightness(pixels [][]Pixel) [][]uint32 {
	var brightness [][]uint32
	for y := 0; y < len(pixels); y++ {
		var row []uint32
		for x := 0; x < len(pixels[y]); x++ {
			row = append(row, uint32((pixels[y][x].R + pixels[y][x].G + pixels[y][x].B) / 3))
		}
		brightness = append(brightness, row)
	}

	return brightness
}

func brightnessMatrixToAscii(brightness [][]uint32) [][]byte {
	chars := "\"`^\\\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	log.Println("Converting brightness matrix to ascii")
	var ascii [][]byte
	for _, values := range brightness {
		var row []byte
		for _, val := range values {
			charVal := chars[getStringRelativeIndex(val, 255, chars)]
			log.Print(string(charVal))
			row = append(row, charVal)
		}
		row = append(row, '\n')
		log.Print("\n")
		ascii = append(ascii, row)
	}

	return ascii
}

func getStringRelativeIndex(val uint32, maxVal uint32, chars string) int {
	valPercent := float32(val) / float32(maxVal)
	idx := int(float32(len(chars)) * float32(valPercent))
	return clampInt(idx, 0, len(chars))
}

func clampInt(val int, minVal int, maxVal int) int {
	if val < minVal {
		return minVal
	} else if val > maxVal {
		return maxVal
	}
	return val
}

func rgbaToRGB(rgba ...uint32) Pixel {
	return Pixel{int(rgba[0] / 257), int(rgba[1] / 257), int(rgba[2] / 257)}
}

type Pixel struct {
	R int
	G int
	B int
}
