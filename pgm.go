package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// ReadPGM reads a PGM image from a file and returns a struct representing the image.
func ReadPGM(filename string) (*PGM, error) {
	// Open the file and handle any errors that may occur.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var width, height, max int
	var data [][]uint8
	scanner := bufio.NewScanner(file)
	// Read and validate the PGM magic number.
	scanner.Scan()
	magicNumber := scanner.Text()
	if magicNumber != "P2" && magicNumber != "P5" {
		return nil, errors.New("unsupported file type")
	}
	// Read width and height from the file.
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			_, err := fmt.Sscanf(line, "%d %d", &width, &height)
			if err == nil {
				break
			} else {
				fmt.Println("Invalid width or height:", err)
			}
		}
	}
	// Read the maximum pixel value.
	scanner.Scan()
	max, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, errors.New("invalid maximum pixel value")
	}
	// Read pixel data from the file.
	for scanner.Scan() {
		line := scanner.Text()
		if magicNumber == "P2" {
			row := make([]uint8, 0)
			for _, char := range strings.Fields(line) {
				pixel, err := strconv.Atoi(char)
				if err != nil {
					fmt.Println("Error converting to integer:", err)
				}
				if pixel >= 0 && pixel <= max {
					row = append(row, uint8(pixel))
				} else {
					fmt.Println("Invalid pixel value:", pixel)
				}
			}
			data = append(data, row)
		}
	}
	// Return a PGM struct with the parsed data.
	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[x][y]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[x][y] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing the header
	fmt.Fprintf(file, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Writing image data
	for _, row := range pgm.data {
		for _, pixel := range row {
			fmt.Fprintf(file, "%d ", pixel)
		}
		fmt.Fprintln(file)
	}

	return nil
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for i := 0; i < len(pgm.data); i++ {
		for j := 0; j < len(pgm.data[i]); j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	NumRows := pgm.width
	NumColumns := pgm.height
	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColumns/2; j++ {
			pgm.data[i][j], pgm.data[i][NumColumns-j-1] = pgm.data[i][NumColumns-j-1], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	numRows := len(pgm.data)
	if numRows == 0 {
		return
	}
	for i := 0; i < numRows/2; i++ {
		pgm.data[i], pgm.data[numRows-i-1] = pgm.data[numRows-i-1], pgm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	oldmax := pgm.max
	pgm.max = maxValue
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {

			pgm.data[i][j] = pgm.data[i][j] * uint8(5) / oldmax
		}
	}
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	NumRows := pgm.width
	NumColumns := pgm.height
	var newData [][]uint8
	for i := 0; i < NumRows; i++ {
		newData = append(newData, make([]uint8, NumColumns))
	}

	for i := 0; i < NumRows; i++ {
		for j := 0; j < NumColumns; j++ {
			newData[i][j] = pgm.data[NumRows-j-1][i]
		}
	}
	pgm.data = newData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	var newNumber string
	if pgm.magicNumber == "P2" {
		newNumber = "P1"
	} else if pgm.magicNumber == "P5" {
		newNumber = "P4"
	}
	NumRows := pgm.width
	NumColumns := pgm.height
	var newData = make([][]bool, NumColumns)
	for i := 0; i < NumColumns; i++ {
		newData[i] = make([]bool, NumRows)
		for j := 0; j < NumRows; j++ {
			newData[i][j] = (pgm.data[i][j] < pgm.max/2)
		}
	}
	return &PBM{data: newData, width: NumRows, height: NumColumns, magicNumber: newNumber}
}
