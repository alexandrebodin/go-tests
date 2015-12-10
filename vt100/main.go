package main

import "os"

const (
	escapeCode = 033
	vertChar   = '_'
	horzChar   = '|'
	angleChar  = '+'
)

func main() {

	var output []byte
	var width = 10
	var height = 10

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if (i == 0 || i == height-1) && (j == 0 || j == width-1) {
				output = append(output, escapeCode, '[', 'K', angleChar)
			} else if i == 0 || i == height-1 {
				output = append(output, escapeCode, '[', 'K', '-')
			} else if j == 0 || j == width-1 {
				output = append(output, escapeCode, '[', 'K', '|')
			} else {
				output = append(output, ' ')
			}
		}
		output = append(output, '\n')
	}
	os.Stdout.Write(output)
}
