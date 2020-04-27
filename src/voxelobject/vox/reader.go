package vox

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"voxelobject"
)

const magic = "VOX "

func isHeaderValid(handle io.Reader) bool {
	limitedReader := io.LimitReader(handle, 4)
	result, err := ioutil.ReadAll(limitedReader)
	return err == nil && string(result) == magic
}

func getSizeFromChunk(handle io.Reader) (voxelobject.Point, error) {
	data, err := getChunkData(handle, 12)

	if err != nil {
		return voxelobject.Point{}, err
	}

	return voxelobject.Point{
		X: byte(binary.LittleEndian.Uint32(data[0:4])),
		Y: byte(binary.LittleEndian.Uint32(data[4:8])),
		Z: byte(binary.LittleEndian.Uint32(data[8:12])),
	}, nil
}

func getPointDataFromChunk(handle io.Reader) ([]voxelobject.PointWithColour, error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return getNilValueForPointDataFromChunk(), err
	}

	result := make([]voxelobject.PointWithColour, len(data)/4)

	for i := 0; i < len(data); i += 4 {
		point := voxelobject.PointWithColour{
			Point: voxelobject.Point{X: data[i], Y: data[i+1], Z: data[i+2]}, Colour: data[i+3],
		}

		result[i/4] = point
	}

	return result, nil
}

func skipUnhandledChunk(handle io.Reader) {
	_, _ = getChunkData(handle, 0)
}

func getChunkData(handle io.Reader, minSize int64) ([]byte, error) {
	parsedSize := getSize(handle)

	// Still need to read to the end even if the size
	// is invalid
	limitedReader := io.LimitReader(handle, parsedSize)
	data, err := ioutil.ReadAll(limitedReader)

	if parsedSize < minSize || parsedSize%4 != 0 {
		return nil, fmt.Errorf("invalid chunk size for xyzi")
	}

	if int64(len(data)) < parsedSize {
		return nil, fmt.Errorf("chunk size declared %d but was %d", parsedSize, len(data))
	}

	return data, err
}

func getSize(handle io.Reader) int64 {
	limitedReader := io.LimitReader(handle, 8)
	size, err := ioutil.ReadAll(limitedReader)

	if err != nil {
		return 0
	}

	parsedSize := int64(binary.LittleEndian.Uint64(size))
	return parsedSize
}

func getNilValueForPointDataFromChunk() []voxelobject.PointWithColour {
	return []voxelobject.PointWithColour{}
}
