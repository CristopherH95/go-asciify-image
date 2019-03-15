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
		fmt.Printf("Processing file: %s\n", path)
		file := ConvertImageToAscii(path)	// convert image to ascii
		fmt.Printf("Output saved to file: %s\n", file)	// tell user where the output can be found
	} else {
		log.Print(err)
		fmt.Println("Could not find image file to read in.")
	}
}

// Converts an image at the given file path to ascii and then writes the results to a file
// The new file will have the same file name as the image, but with .txt appended to it
func ConvertImageToAscii(path string) string {
	width, height := getImageSize(path)		// get image bounds
	pix, err := getImagePixels(path, width, height)	// get image pixels
	checkErr(err)
	pixVal := pixelMatrixToBrightness(pix)	// get brightness values from pixels
	ascii := brightnessMatrixToAscii(pixVal)	// convert brightness to corresponding ascii
	file := writeMatrixToFile(path, ascii)	// write results to a file
	return file
}

// Checks if an error is nil, if it is not nil the error is logged and execution is terminated
func checkErr(err error) {
	if err != nil {
		fmt.Println("Failure encountered while attempting to process and convert image")
		log.Fatal(err)
	}
}

// Prints a very simple usage message
func printUsage() {
	msg := `usage:
				yasciifier <filename>\n
			This will convert the image file (.png or .jpg) to ascii and then save the results 
			to a new file with the same name as the image, but with .txt appended to it.\n`
	fmt.Print(msg)
}

// Writes the given ascii matrix to a new file with the same name as the given file path
// but with .txt appended to it
func writeMatrixToFile(path string, ascii [][]byte) string {
	file := filepath.Base(path) + ".txt"	// add .txt to original file name for new file
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)	// create or append
	checkErr(err)
	for _, row := range ascii {
		_, err := f.Write(row)	// write data, if error encountered try to cleanup the file
		if err != nil {
			delErr := os.Remove(file)
			checkErr(delErr)
			log.Fatal(err)
		}
	}

	return file
}

// Retrieves the image dimensions of the image file at the given path
func getImageSize(path string) (int, int) {
	file, err := os.Open(path)
	checkErr(err)
	img, _, err := image.DecodeConfig(file)
	checkErr(err)
	return img.Width, img.Height
}

// Reads the image at the given file path and retrieves the pixels of the given width and height in the image
func getImagePixels(path string, width int, height int) ([][]Pixel, error) {
	log.Println("Reading in image pixel values")
	file, err := os.Open(path)
	checkErr(err)
	img, _, err := image.Decode(file)
	checkErr(err)
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

// Converts all the given Pixels from their RGB format to a single brightness value
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

// Converts a matrix of brightness values into ascii values which approximate that brightness
func brightnessMatrixToAscii(brightness [][]uint32) [][]byte {
	chars := "\"`^\\\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	log.Println("Converting brightness matrix to ascii")
	var ascii [][]byte
	for _, values := range brightness {
		var row []byte
		for _, val := range values {
			charVal := chars[getStringRelativeIndex(val, 255, chars)]
			row = append(row, charVal)
		}
		row = append(row, '\n')
		ascii = append(ascii, row)
	}

	return ascii
}

// Returns an index into a given string based on the percentage given by a value and its corresponding max value
func getStringRelativeIndex(val uint32, maxVal uint32, chars string) int {
	valPercent := float32(val) / float32(maxVal)
	idx := int(float32(len(chars)) * float32(valPercent))
	return clampInt(idx, 0, len(chars))
}

// Clamps the given integer between min and max values
func clampInt(val int, minVal int, maxVal int) int {
	if val < minVal {
		return minVal
	} else if val > maxVal {
		return maxVal
	}
	return val
}

// Converts RGBA to an RGB based Pixel
func rgbaToRGB(rgba ...uint32) Pixel {
	return Pixel{int(rgba[0] / 257), int(rgba[1] / 257), int(rgba[2] / 257)}
}

// Pixel struct that is RGB based
type Pixel struct {
	R int
	G int
	B int
}
