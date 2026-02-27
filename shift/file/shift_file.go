package file

import (
	"errors"
	"io"
	"rotate/shift"
)

var (
	CHUNK_SIZE = int64(1024)
)

// initialSeek shall jump to the start or end of the files (direction dependent) and return
// the current offset from the start of the file
func initialSeek(fileIn, fileOut io.Seeker, direction shift.DIRECTION) (int64, error) {
	seek := io.SeekStart
	if direction == shift.LEFT {
		seek = io.SeekEnd
	}
	location, err := fileIn.Seek(0, seek)
	if err != nil {
		return 0, err
	}
	_, err = fileOut.Seek(0, seek)
	if err != nil {
		return 0, err
	}

	return location, nil
}

// getFirstOverflowByte gets the overflow byte for the first chunk (shifting left) or last chunk (shifting right)
// in order to apply that overflow byte to the start or end of the file.
func getFirstOverflowByte(fileIn io.ReadSeeker, direction shift.DIRECTION) (uint8, error) {
	res := uint8(0)
	scratch := []byte{0}
	if direction == shift.LEFT {
		if _, err := fileIn.Seek(0, io.SeekStart); err != nil {
			return 0, err
		}
		if _, err := fileIn.Read(scratch); err != nil {
			return 0, err
		}
		res = scratch[0] >> 7
	} else {
		if _, err := fileIn.Seek(-1, io.SeekEnd); err != nil {
			return 0, err
		}
		if _, err := fileIn.Read(scratch); err != nil {
			return 0, err
		}
		res = scratch[0] << 7
	}
	return res, nil
}

// moveToNextLocation handles shifting the file position for the left moving operation
// We must move left by CHUNK_SIZE bytes, unless that pushes us past the start of the file, in which case
// we'll just go to the start and truncate chunk to the number of bytes left to read
func moveToNextLocation(curPos int64, chunk []uint8, fileIn io.Seeker, fileOut io.Seeker) (int64, []uint8, error) {
	curPos = curPos - CHUNK_SIZE
	if curPos < 0 {
		toRead := CHUNK_SIZE + curPos
		chunk = chunk[:toRead]
		curPos = 0
	}
	if _, err := fileOut.Seek(curPos, io.SeekStart); err != nil {
		return 0, nil, err
	}
	if _, err := fileIn.Seek(curPos, io.SeekStart); err != nil {
		return 0, nil, err
	}
	return curPos, chunk, nil
}

// ShiftFile will perform the shifting operation iteratively shifting the fileIn and writing to the
// fileOut, shifting in the appropriate direction
func ShiftFile(fileIn io.ReadSeeker, fileOut io.ReadWriteSeeker, direction shift.DIRECTION) error {
	chunk := make([]uint8, CHUNK_SIZE)
	overflowIn, err := getFirstOverflowByte(fileIn, direction)
	if err != nil {
		return err
	}
	var curPos int64
	if curPos, err = initialSeek(fileIn, fileOut, direction); err != nil {
		return err
	}
	for {
		if direction == shift.LEFT {
			if curPos, chunk, err = moveToNextLocation(curPos, chunk, fileIn, fileOut); err != nil {
				return err
			}
		}
		read, err := fileIn.Read(chunk)
		if err != nil && !errors.Is(io.EOF, err) {
			return err
		}
		if read != int(CHUNK_SIZE) {
			chunk = chunk[:read] //truncate the chunk
		}
		result := shift.ShiftChunk(chunk, overflowIn, direction)
		overflowIn = result.Overflow
		if _, err := fileOut.Write(result.Result); err != nil {
			return err
		}
		if read != int(CHUNK_SIZE) {
			break
		}
	}

	return nil
}
