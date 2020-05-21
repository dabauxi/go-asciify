package main

import (
	"flag"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

var asciiChars = [13]string{"@", "#", "S", "%", "?", "*", "+", ";", ":", ",", "."}

func main() {

	imagePath := flag.String("image", "", "provide a path to an image here")
	width := flag.Int("width", 0, "provide the desired width of the resulting image.")
	flag.Parse()

	if *imagePath == "" {
		fmt.Println("Provide a correct filepath.")
		os.Exit(1)
	}
	if *width < 0 {
		fmt.Println("Provide a valid width.")
		os.Exit(1)
	}
	_, err := os.Stat(*imagePath)
	if os.IsNotExist(err) {
		fmt.Println("File does not exist.")
		os.Exit(1)
	}

	img, err := getImageFromFilePath(*imagePath)
	if err != nil {
		fmt.Print("error:")
		fmt.Print(err)
	}
	resizedImg := resize.Resize(uint(*width), 0, img, resize.Lanczos3)
	resizedGrayScaleImg := convertToGrayScale(resizedImg)
	convertToASCII(resizedGrayScaleImg)
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	importedImage, _, err := image.Decode(f)
	return importedImage, err
}

func writeImageToFile(filePath string, img image.Image) {
	outputFile, err := os.Create(filePath)
	if err != nil {
		fmt.Print("error:")
		fmt.Print(err)
	}
	if strings.HasSuffix(filePath, ".jpg") {
		jpeg.Encode(outputFile, img, nil)
	} else if strings.HasSuffix(filePath, ".png") {
		png.Encode(outputFile, img)
	}
	outputFile.Close()
}

func convertToGrayScale(img image.Image) (image.Image) {
	b := img.Bounds()
	imgSet := image.NewRGBA(b)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			oldPixel := img.At(x, y)
			pixel := color.GrayModel.Convert(oldPixel)
			imgSet.Set(x, y, pixel)
		}
	}
	return imgSet
}

func convertToASCII(img image.Image) (image.Image, *ImageArray) {
	b := img.Bounds()
	array := ImageArrayFactory(b.Max.X, b.Max.Y)
	for y := 0; y < b.Max.Y; y++ {
		var tmp = ""
		for x := 0; x < b.Max.X; x++ {

			pixel := img.At(x,y)
			r, _, _, _ := pixel.RGBA()
			tmp += getASCIIChar(uint8(r))
			array.Put(x, y, tmp)
		}
		fmt.Println(tmp)
	}
	return img, array
}

func getASCIIChar(value uint8) (char string) {
	return asciiChars[value/25]
}


type ImageArray struct {
	data []string
	xSpan int
	ySpan int
}

func ImageArrayFactory(xspan, yspan int) *ImageArray {
	return &ImageArray{data: make([]string, xspan*yspan), xSpan: xspan, ySpan: yspan}
}

func (td *ImageArray) Put(x int, y int, value string) {
	td.data[x*td.ySpan+y] = value
}

func (td *ImageArray) Get(x, y int) string {
	return td.data[x*td.ySpan+y]
}

func writeToAsciiImage(asciiImage *ImageArray) {
	// TODO
}