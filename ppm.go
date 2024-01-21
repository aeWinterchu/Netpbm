package Netpbm

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           uint8
}

type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	var err error
	var magicNumber string = ""
	var width int
	var height int
	var maxval int
	var counter int
	var headersize int
	var splitfile []string

	// Read the entire file content.
	file, err := os.ReadFile(filename)
	if err != nil {
		// Handle file reading errors if any.
	}

	// Split the file content into lines based on the presence of carriage return.
	if strings.Contains(string(file), "\r") {
		splitfile = strings.SplitN(string(file), "\r\n", -1)
	} else {
		splitfile = strings.SplitN(string(file), "\n", -1)
	}

	// Parse the header information to determine image properties.
	for i, _ := range splitfile {
		if strings.Contains(splitfile[i], "P3") {
			magicNumber = "P3"
		} else if strings.Contains(splitfile[i], "P6") {
			magicNumber = "P6"
		} else {
			// Handle unrecognized magic numbers if any.
		}

		// Check for comments and update headersize accordingly.
		if strings.HasPrefix(splitfile[i], "#") && maxval != 0 {
			headersize = counter
		}

		// Parse width and height information.
		splitl := strings.SplitN(splitfile[i], " ", -1)
		if width == 0 && height == 0 && len(splitl) >= 2 {
			width, err = strconv.Atoi(splitl[0])
			height, err = strconv.Atoi(splitl[1])
			headersize = counter
		}

		// Parse maximum pixel value.
		if maxval == 0 && width != 0 {
			maxval, err = strconv.Atoi(splitfile[i])
			headersize = counter
		}

		counter++
	}

	// Initialize a 2D slice to store pixel data.
	data := make([][]Pixel, height)
	for j := 0; j < height; j++ {
		data[j] = make([]Pixel, width)
	}
	var splitdata []string

	// Parse and populate the pixel data.
	if counter > headersize {
		for i := 0; i < height; i++ {
			splitdata = strings.SplitN(splitfile[headersize+1+i], " ", -1)
			for j := 0; j < width*3; j += 3 {
				r, _ := strconv.Atoi(splitdata[j])
				if r > maxval {
					r = maxval
				}
				g, _ := strconv.Atoi(splitdata[j+1])
				if g > maxval {
					g = maxval
				}
				b, _ := strconv.Atoi(splitdata[j+2])
				if b > maxval {
					b = maxval
				}
				data[i][j/3] = Pixel{R: uint8(r), G: uint8(g), B: uint8(b)}
			}
		}
	}

	// Return a PPM struct with the parsed data.
	return &PPM{data: data, width: width, height: height, magicNumber: magicNumber, max: uint8(maxval)}, err
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = value
	}
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.magicNumber)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.width, ppm.height)
	if err != nil {
		file.Close()
		return err
	}
	_, err = fmt.Fprintln(file, ppm.max)
	if err != nil {
		file.Close()
		return err
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			if ppm.data[y][x].R > ppm.max || ppm.data[y][x].G > ppm.max || ppm.data[y][x].B > ppm.max {
				errors.New("data value is too high")
			} else {
				fmt.Fprint(file, ppm.data[y][x].R, ppm.data[y][x].G, ppm.data[y][x].B, " ")
			}
		}
		fmt.Fprintln(file)
	}
	return err
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for i := 0; i < len(ppm.data); i++ {
		for j := 0; j < len(ppm.data[0]); j++ {
			ppm.data[i][j].R = ppm.max - ppm.data[i][j].R
			ppm.data[i][j].G = ppm.max - ppm.data[i][j].G
			ppm.data[i][j].B = ppm.max - ppm.data[i][j].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	NumRows := ppm.width
	NumColums := ppm.height
	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColums/2; j++ {
			ppm.data[i][j], ppm.data[i][NumColums-j-1] = ppm.data[i][NumColums-j-1], ppm.data[i][j]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	numRows := len(ppm.data)
	if numRows == 0 {
		return
	}
	for i := 0; i < numRows/2; i++ {
		ppm.data[i], ppm.data[numRows-i-1] = ppm.data[numRows-i-1], ppm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	oldmax := ppm.max
	ppm.max = maxValue
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			// Convert each color component individually
			ppm.data[i][j].R = uint8(float64(ppm.data[i][j].R) * float64(ppm.max) / float64(oldmax))
			ppm.data[i][j].G = uint8(float64(ppm.data[i][j].G) * float64(ppm.max) / float64(oldmax))
			ppm.data[i][j].B = uint8(float64(ppm.data[i][j].B) * float64(ppm.max) / float64(oldmax))
		}
	}
}

// Rotate90CW rotates the PPM image 90Â° clockwise.
func (ppm *PPM) Rotate90CW() {
	NumRows := ppm.width
	NumColumns := ppm.height
	var newData [][]Pixel
	for i := 0; i < NumRows; i++ {
		newData = append(newData, make([]Pixel, NumColumns))
	}

	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColumns; j++ {
			newData[i][j] = ppm.data[NumColumns-j-1][i]
		}
	}
	ppm.data = newData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	var newNumber string
	if ppm.magicNumber == "P3" {
		newNumber = "P2"
	} else if ppm.magicNumber == "P6" {
		newNumber = "P5"
	}
	NumRows := ppm.width
	NumColumns := ppm.height
	var newData = make([][]uint8, NumColumns)

	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]uint8, NumRows)
		for j := 0; j < NumRows; j++ {
			newData[i][j] = uint8((int(ppm.data[i][j].R) + int(ppm.data[i][j].G) + int(ppm.data[i][j].B)) / 3)
		}
	}
	return &PGM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber, max: ppm.max}
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Determine the new magic number based on the original PPM magic number
	var newNumber string
	if ppm.magicNumber == "P3" {
		newNumber = "P1"
	} else if ppm.magicNumber == "P6" {
		newNumber = "P4"
	}
	// Set the dimensions for the new PBM image
	NumRows := ppm.width
	NumColumns := ppm.height
	// Create a 2D slice to store boolean data for the PBM image
	var newData = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]bool, NumRows)
		for j := 0; j < NumRows; j++ {
			// Convert RGB values to grayscale and check if it is less than half of the max value
			newData[i][j] = (uint8((int(ppm.data[i][j].R)+int(ppm.data[i][j].G)+int(ppm.data[i][j].B))/3) < ppm.max/2)
		}
	}
	// Return a pointer to a PBM struct containing the converted data
	return &PBM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber}
}

// Point represents a 2D point with X and Y coordinates.
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points in the PPM image using the specified color.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Calculate the differences in X and Y coordinates
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	// Determine the number of steps based on the maximum difference in coordinates
	steps := int(math.Max(math.Abs(dx), math.Abs(dy)))
	// Calculate the incremental changes in X and Y coordinates
	xIncrement := dx / float64(steps)
	yIncrement := dy / float64(steps)
	// Initialize the starting coordinates
	x, y := float64(p1.X), float64(p1.Y)
	// Draw the line by setting the specified color at each point along the line
	for i := 0; i <= steps; i++ {
		ppm.Set(int(x), int(y), color)
		x += xIncrement
		y += yIncrement
	}
}

func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	ppm.DrawLine(p1, p2, color)

	p3 := Point{p2.X, p2.Y + height}
	ppm.DrawLine(p2, p3, color)

	p4 := Point{p1.X, p1.Y + height}
	ppm.DrawLine(p3, p4, color)

	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	ppm.DrawRectangle(p1, width, height, color)
	for j := p1.Y + 1; j < p1.Y+height; j++ {
		for i := p1.X + 1; i < p1.X+width; i++ {
			ppm.Set(i, j, color)
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawFilledCircle dessine un cercle rempli.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// ...
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// ...
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// ...
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// ...
}
