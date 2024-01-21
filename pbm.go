package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var magicNumber string
	// Function to read the next line without comment
	readNextLine := func() (string, error) {
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Ignore empty lines or lines beginning with "#".
			if line != "" && !strings.HasPrefix(line, "#") {
				return line, nil
			}
		}
		return "", scanner.Err()
	}
	// Read the first uncommented line to obtain the magic number
	if magicNumber, err = readNextLine(); err != nil {
		return nil, err
	}
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, errors.New("unsupported file type")
	}
	// Read dimensions
	dimensions, err := readNextLine()
	if err != nil {
		return nil, err
	}
	dimTokens := strings.Fields(dimensions)
	if len(dimTokens) != 2 {
		return nil, errors.New("invalid image dimensions")
	}
	width, _ := strconv.Atoi(dimTokens[0])
	height, _ := strconv.Atoi(dimTokens[1])
	var data [][]bool
	// If the image is empty, initialize data with an empty slice
	if width > 0 && height > 0 {
		data = make([][]bool, height)
		for i := range data {
			data[i] = make([]bool, width)
		}
		if magicNumber == "P1" {
			for i := 0; i < height; i++ {
				line, err := readNextLine()
				if err != nil {
					return nil, err
				}

				tokens := strings.Fields(line)
				for j, token := range tokens {
					pixel, err := strconv.Atoi(token)
					if err != nil {
						return nil, err
					}
					data[i][j] = pixel == 1
				}
			}
		} else if magicNumber == "P4" {
			// Calculate the number of padding bits
			paddingBits := (8 - width%8) % 8
			// Calculate the number of bytes per row, considering padding
			bytesPerRow := (width + paddingBits + 7) / 8
			// Create a buffer to read binary data
			buffer := make([]byte, bytesPerRow)
			for i := 0; i < height; i++ {
				_, err := file.Read(buffer)
				if err != nil {
					return nil, err
				}
				// Process the bits from the buffer
				for j := 0; j < width; j++ {
					// Get the byte containing the bit
					byteIndex := j / 8
					bitIndex := 7 - (j % 8)
					bit := (buffer[byteIndex] >> bitIndex) & 1
					data[i][j] = bit == 1
				}
			}
		}
	}
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	if x >= 0 && x < pbm.width || y >= 0 && y < pbm.height {
		return pbm.data[y][x]
	}
	return false
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	if x >= 0 && x < pbm.width && y >= 0 && y < pbm.height {
		pbm.data[y][x] = value
	}
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	// Write the magic number
	_, err = fmt.Fprintln(writer, pbm.magicNumber)
	if err != nil {
		return err
	}
	// Write the dimensions
	_, err = fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)
	if err != nil {
		return err
	}
	// Write pixel data based on the magic number
	if pbm.magicNumber == "P1" {
		for i := 0; i < pbm.height; i++ {
			for j := 0; j < pbm.width; j++ {
				if pbm.data[i][j] {
					_, err = fmt.Fprint(writer, "1 ")
				} else {
					_, err = fmt.Fprint(writer, "0 ")
				}
				if err != nil {
					return err
				}
			}
			_, err = fmt.Fprintln(writer)
			if err != nil {
				return err
			}
		}
	}
	return writer.Flush()
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j, k := 0, pbm.width-1; j < k; j, k = j+1, k-1 {
			pbm.data[i][j], pbm.data[i][k] = pbm.data[i][k], pbm.data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for i, j := 0, pbm.height-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
