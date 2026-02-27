package file

import (
	"io"
	"os"
	"path"
	"rotate/shift"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftFile(t *testing.T) {
	// Shrink the chunk size for these simple tests so the input doesn't need to be crazy big
	oldChunkSize := CHUNK_SIZE
	CHUNK_SIZE = 2
	defer func() { CHUNK_SIZE = oldChunkSize }()

	type TestCase struct {
		name      string
		input     []uint8
		output    []uint8
		direction shift.DIRECTION
	}
	tt := []TestCase{
		{
			name:      "All ones, right",
			input:     []uint8{255, 255, 255},
			output:    []uint8{255, 255, 255},
			direction: shift.RIGHT,
		}, {
			name:      "All ones, left",
			input:     []uint8{255, 255, 255},
			output:    []uint8{255, 255, 255},
			direction: shift.LEFT,
		}, {
			name:      "All zeroes, right",
			input:     []uint8{0, 0, 0, 0, 0, 0},
			output:    []uint8{0, 0, 0, 0, 0, 0},
			direction: shift.RIGHT,
		}, {
			name:      "All zeroes, left",
			input:     []uint8{0, 0, 0, 0, 0},
			output:    []uint8{0, 0, 0, 0, 0},
			direction: shift.LEFT,
		}, {
			name:      "255 in the middle, right",
			input:     []uint8{0, 0, 255, 0, 0},
			output:    []uint8{0, 0, 0b01111111, 0b10000000, 0},
			direction: shift.RIGHT,
		}, {
			name:      "255 in the middle, left",
			input:     []uint8{0, 0, 255, 0, 0},
			output:    []uint8{0, 0b00000001, 0b11111110, 0, 0},
			direction: shift.LEFT,
		}, {
			name:      "two 255s start of the file, left",
			input:     []uint8{255, 255, 0},
			output:    []uint8{255, 0b11111110, 0b00000001},
			direction: shift.LEFT,
		}, {
			name:      "two 255s end of the file, right",
			input:     []uint8{0, 255, 255},
			output:    []uint8{0b10000000, 0b01111111, 255},
			direction: shift.RIGHT,
		}, {
			name:      "two 255s end of the file, left",
			input:     []uint8{0, 255, 255},
			output:    []uint8{0b00000001, 255, 0b11111110},
			direction: shift.LEFT,
		}, {
			name:      "two 255s start of the file, right",
			input:     []uint8{255, 255, 0},
			output:    []uint8{0b01111111, 255, 0b10000000},
			direction: shift.RIGHT,
		},
	}
	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			testDir, err := os.MkdirTemp("", "*")
			assert.NoError(t, err)
			defer os.RemoveAll(testDir)
			fIn, err := os.Create(path.Join(testDir, "fileIn"))
			assert.NoError(t, err)
			defer fIn.Close()
			_, err = fIn.Write(testCase.input)
			assert.NoError(t, err)
			fOut, err := os.Create(path.Join(testDir, "fileOut"))
			assert.NoError(t, err)
			defer fOut.Close()

			assert.NoError(t, ShiftFile(fIn, fOut, testCase.direction))
			assert.NoError(t, fOut.Sync())
			_, err = fOut.Seek(0, io.SeekStart)
			assert.NoError(t, err)
			result, err := io.ReadAll(fOut)
			assert.NoError(t, err)
			assert.Equal(t, testCase.output, result)
		})
	}
}
