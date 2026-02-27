package shift

const (
	LEFT = iota
	RIGHT
)

type DIRECTION int

type ByteShiftResult struct {
	Result   uint8
	Overflow uint8
}

type ChunkShiftResult struct {
	Result   []uint8
	Overflow uint8
}

// shiftByte will shift a single byte, b, in a single byte in the specified direction. An overflow byte
// may be passed in which shall be ORed with the byte depending on the direction specified
func shiftByte(b uint8, overflowIn uint8, direction DIRECTION) ByteShiftResult {
	overflow := uint8(0)
	result := uint8(0)
	if direction == LEFT {
		result = (b << 1) | overflowIn
		overflow = b >> 7
	} else {
		result = (b >> 1) | overflowIn
		overflow = b << 7
	}
	return ByteShiftResult{
		Result:   result,
		Overflow: overflow,
	}
}

// ShiftChunk shall take an array of bytes, chunk, and shift them all one bit in the specified direction.
// An oveflow bit may be passed in. The output shall be the shifted chunk and an output bit.
func ShiftChunk(chunk []uint8, overflowIn uint8, direction DIRECTION) ChunkShiftResult {
	var output = make([]uint8, len(chunk))
	var overflow = overflowIn
	if direction == LEFT {
		for idx := len(chunk) - 1; idx >= 0; idx-- {
			byteResult := shiftByte(chunk[idx], overflow, direction)
			output[idx] = byteResult.Result
			overflow = byteResult.Overflow
		}
	} else {
		for idx := 0; idx < len(chunk); idx++ {
			byteResult := shiftByte(chunk[idx], overflow, direction)
			output[idx] = byteResult.Result
			overflow = byteResult.Overflow
		}
	}
	return ChunkShiftResult{
		Result:   output,
		Overflow: overflow,
	}
}
