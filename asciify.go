package asciify

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Invalid arguments")
		printUsage() // need to be provided an argument for the file path
		return
	}
	path := os.Args[1]                       // grab the path
	if _, err := os.Stat(path); err == nil { // if the file exists begin processing
		fmt.Printf("Processing file: %s\n", path)
		file, err := ConvertImageToAscii(path) // convert image to ascii
		if err != nil {
			fmt.Println("Failed to process and convert image")
			log.Fatal(err)
		}
		fmt.Printf("Output saved to file: %s\n", file) // tell user where the output can be found
	} else {
		log.Print(err)
		fmt.Println("Could not find image file to read in.")
	}
}

// Converts an image at the given file path to ascii and then writes the results to a file
// The new file will have the same file name as the image, but with .txt appended to it
func ConvertImageToAscii(path string) (string, error) {
	pix, err := getImagePixels(path) // get image pixels
	if err != nil {
		return "", err
	}
	pixVal := convertRGBToBrightness(pix)       // get brightness values from pixels
	ascii := convertBrightnessToAscii(pixVal)   // convert brightness to corresponding ascii
	file, err := writeMatrixToFile(path, ascii) // write results to a file
	return file, err
}

// Prints a very simple usage message
func printUsage() {
	var name string
	if len(os.Args) > 0 {
		name = os.Args[0]
	} else {
		name = "asciify"
	}
	msg := "usage:\n    %s <file-name>\nThis will convert file-name (a .png or .jpg) to ascii. " +
		"Then save the results to a new file with the same name as the original, but with .txt appended to it."
	fmt.Printf(msg, name)
}

// Writes the given ascii matrix to a new file with the same name as the given file path
// but with .txt appended to it
func writeMatrixToFile(path string, ascii [][]byte) (string, error) {
	file := path + ".txt"                                      // add .txt to original file name for new file
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644) // create file (or overwrite)
	if err != nil {
		return "", err
	}
	for _, row := range ascii {
		_, err := f.Write(row) // write data, if error encountered try to cleanup the file
		if err != nil {
			return "", err
		}
	}
	err = f.Close()

	return file, err
}

// Returns an image at the given path along with its width and height
// The image is re-sized if it's width or height is over 200 pixels each
// This is to prevent the ascii version from being too large
func getImageData(path string) (image.Image, int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, 0, 0, err
	}
	err = file.Close()
	if err != nil {
		return nil, 0, 0, err
	}
	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	if width > 200 || height > 200 {
		// image may be big, resize and preserve aspect ratio
		if width > height {
			img = imaging.Resize(img, 200, 0, imaging.Lanczos) // resize with width being larger
		} else {
			img = imaging.Resize(img, 0, 200, imaging.Lanczos) // resize with height being larger
		}
		bounds = img.Bounds()
		width = bounds.Max.X
		height = bounds.Max.Y
	}

	return img, width, height, nil
}

// Reads the image at the given file path and retrieves the pixels of the given width and height in the image
func getImagePixels(path string) ([][]Pixel, error) {
	log.Println("Reading in image pixel values")
	img, width, height, err := getImageData(path)
	if err != nil {
		return nil, err
	}
	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, convertRGBAToRGB(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// Converts all the given Pixels from their RGB format to a single brightness value
func convertRGBToBrightness(pixels [][]Pixel) [][]uint32 {
	var brightness [][]uint32
	for y := 0; y < len(pixels); y++ {
		var row []uint32
		for x := 0; x < len(pixels[y]); x++ {
			row = append(row, uint32((pixels[y][x].R+pixels[y][x].G+pixels[y][x].B)/3))
		}
		brightness = append(brightness, row)
	}

	return brightness
}

// Converts a matrix of brightness values into ascii values which approximate that brightness
func convertBrightnessToAscii(brightness [][]uint32) [][]byte {
	chars := "\"`^\\\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	log.Println("Converting brightness matrix to ascii")
	var ascii [][]byte
	for _, values := range brightness {
		var row []byte
		for _, val := range values {
			charVal := chars[getStringRelativeIndex(val, 255, chars)]
			row = append(row, charVal, charVal) // add item more than once to prevent final result from looking stretched
		}
		row = append(row, '\n')
		ascii = append(ascii, row)
	}

	return ascii
}

// Returns an index into a given string based on the percentage given by a value and its corresponding max value
func getStringRelativeIndex(val uint32, maxVal uint32, chars string) int {
	valPercent := float32(val) / float32(maxVal)
	idx := int(float32(len(chars)) * valPercent)
	return clampInt(idx, 0, len(chars)-1)
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
func convertRGBAToRGB(rgba ...uint32) Pixel {
	return Pixel{int(rgba[0] / 257), int(rgba[1] / 257), int(rgba[2] / 257)}
}

// Pixel struct that is RGB based
type Pixel struct {
	R int
	G int
	B int
}
