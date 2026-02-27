package shift

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	LEFT_SET  = 0b10000000
	RIGHT_SET = 0b00000001
)

// Some basic test cases to ensure that the byte shifting works as expected
func TestShiftByte(t *testing.T) {
	type TestCase struct {
		name       string
		input      uint8
		overflowIn uint8
		direction  DIRECTION
		want       ByteShiftResult
	}

	tt := []TestCase{
		{
			name:       "All zeroes, left",
			input:      0,
			overflowIn: 0,
			direction:  LEFT,
			want: ByteShiftResult{
				Result:   0,
				Overflow: 0,
			},
		}, {
			name:       "All zeroes, right",
			input:      0,
			overflowIn: 0,
			direction:  RIGHT,
			want: ByteShiftResult{
				Result:   0,
				Overflow: 0,
			},
		}, {
			name:       "Basic no overflow, left",
			input:      RIGHT_SET,
			overflowIn: 0,
			direction:  LEFT,
			want: ByteShiftResult{
				Result:   0b00000010,
				Overflow: 0,
			},
		}, {
			name:       "Basic no overflow, right",
			input:      LEFT_SET,
			overflowIn: 0,
			direction:  RIGHT,
			want: ByteShiftResult{
				Result:   0b01000000,
				Overflow: 0,
			},
		}, {
			name:       "Basic overflow out, right",
			input:      RIGHT_SET,
			overflowIn: 0,
			direction:  RIGHT,
			want: ByteShiftResult{
				Result:   0,
				Overflow: LEFT_SET,
			},
		}, {
			name:       "Basic overflow out, left",
			input:      LEFT_SET,
			overflowIn: 0,
			direction:  LEFT,
			want: ByteShiftResult{
				Result:   0,
				Overflow: RIGHT_SET,
			},
		}, {
			name:       "All zeroes, left, overflowIn",
			input:      0,
			overflowIn: RIGHT_SET,
			direction:  LEFT,
			want: ByteShiftResult{
				Result:   RIGHT_SET,
				Overflow: 0,
			},
		}, {
			name:       "All zeroes, right, overflowIn",
			input:      0,
			overflowIn: LEFT_SET,
			direction:  LEFT,
			want: ByteShiftResult{
				Result:   LEFT_SET,
				Overflow: 0,
			},
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			output := shiftByte(testCase.input, testCase.overflowIn, testCase.direction)
			assert.Equal(t, testCase.want.Result, output.Result)
			assert.Equal(t, testCase.want.Overflow, output.Overflow)
		})
	}
}

// Fuzz test, for any byte input, shift one direction and then back again should result in the original byte
func FuzzByteShiftThereAndBack(f *testing.F) {
	f.Fuzz(func(t *testing.T, input byte) {
		output1 := shiftByte(input, 0, LEFT)
		result := shiftByte(output1.Result, output1.Overflow, RIGHT)
		assert.Equal(t, 0, result.Overflow)
		assert.Equal(t, input, result.Result)
		output1 = shiftByte(input, 0, RIGHT)
		result = shiftByte(output1.Result, output1.Overflow, LEFT)
		assert.Equal(t, 0, result.Overflow)
		assert.Equal(t, input, result.Result)
	})
}

func TestShiftChunk(t *testing.T) {
	type TestCase struct {
		name       string
		input      []uint8
		overflowIn uint8
		direction  DIRECTION
		want       ChunkShiftResult
	}

	tt := []TestCase{
		{
			name:       "All zeros, left",
			input:      []uint8{0, 0, 0},
			overflowIn: 0,
			direction:  LEFT,
			want: ChunkShiftResult{
				Result:   []uint8{0, 0, 0},
				Overflow: 0,
			},
		}, {
			name:       "All zeros, right",
			input:      []uint8{0, 0, 0},
			overflowIn: 0,
			direction:  RIGHT,
			want: ChunkShiftResult{
				Result:   []uint8{0, 0, 0},
				Overflow: 0,
			},
		}, {
			name:       "All ones, left",
			input:      []uint8{255, 255, 255, 255, 255},
			overflowIn: 0,
			direction:  LEFT,
			want: ChunkShiftResult{
				Result:   []uint8{255, 255, 255, 255, 0b11111110},
				Overflow: RIGHT_SET,
			},
		}, {
			name:       "All ones, right",
			input:      []uint8{255, 255, 255, 255, 255},
			overflowIn: 0,
			direction:  RIGHT,
			want: ChunkShiftResult{
				Result:   []uint8{0b01111111, 255, 255, 255, 255},
				Overflow: LEFT_SET,
			},
		},
		{
			name:       "All ones, overflow in, left",
			input:      []uint8{255, 255, 255, 255, 255},
			overflowIn: RIGHT_SET,
			direction:  LEFT,
			want: ChunkShiftResult{
				Result:   []uint8{255, 255, 255, 255, 255},
				Overflow: RIGHT_SET,
			},
		}, {
			name:       "All ones, overflow in, right",
			input:      []uint8{255, 255, 255, 255, 255},
			overflowIn: LEFT_SET,
			direction:  RIGHT,
			want: ChunkShiftResult{
				Result:   []uint8{255, 255, 255, 255, 255},
				Overflow: LEFT_SET,
			},
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			output := ShiftChunk(testCase.input, testCase.overflowIn, testCase.direction)
			assert.Equal(t, testCase.want.Result, output.Result)
			assert.Equal(t, testCase.want.Overflow, output.Overflow)
		})
	}
}

func FuzzChunkShiftThereAndBack(f *testing.F) {
	f.Fuzz(func(t *testing.T, chunkLen uint8) {
		chunk := make([]uint8, chunkLen)
		for i := range chunk {
			chunk[i] = uint8(rand.Uint32())
		}
		output := ShiftChunk(chunk, 0, LEFT)
		result := ShiftChunk(output.Result, output.Overflow, RIGHT)
		assert.Equal(t, chunk, result.Result)

		output = ShiftChunk(chunk, 0, RIGHT)
		result = ShiftChunk(output.Result, output.Overflow, LEFT)
		assert.Equal(t, chunk, result.Result)
	})
}
