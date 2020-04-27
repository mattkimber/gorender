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
	result, err := getChunkHeader(handle)
	return err == nil && result == magic
}

func getChunkHeader(handle io.Reader) (string, error) {
	limitedReader := io.LimitReader(handle, 4)
	result, err := ioutil.ReadAll(limitedReader)
	return string(result), err
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

func getRawVoxelsFromPointData(size voxelobject.Point, data []voxelobject.PointWithColour) voxelobject.RawVoxelObject {
	result := make([][][]byte, size.X)

	for x := range result {
		result[x] = make([][]byte, size.Y)
		for y := range result[x] {
			result[x][y] = make([]byte, size.Z)
		}
	}

	for _, p := range data {
		if p.Point.X < size.X && p.Point.Y < size.Y && p.Point.Z < size.Z {
			result[p.Point.X][p.Point.Y][p.Point.Z] = p.Colour - 2
		}
	}

	return result
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

func GetRawVoxels(handle io.Reader) (voxelobject.RawVoxelObject, error) {
	if !isHeaderValid(handle) {
		return nil, fmt.Errorf("header not valid")
	}

	size := voxelobject.Point{}
	pointData := make([]voxelobject.PointWithColour, 0)

	for {
		chunkType, err := getChunkHeader(handle)

		if err != nil {
			return nil, fmt.Errorf("error reading chunk header: %v", err)
		}

		if chunkType == "" {
			break
		}

		switch chunkType {
		case "SIZE":
			data, err := getSizeFromChunk(handle)
			if err != nil {
				return nil, fmt.Errorf("error reading size chunk: %v", err)
			}

			// We only expect one SIZE chunk, but use the last value
			size = data
		case "XYZI":
			data, err := getPointDataFromChunk(handle)
			if err != nil {
				return nil, fmt.Errorf("error reading size chunk: %v", err)
			}

			pointData = append(pointData, data...)
		default:
			skipUnhandledChunk(handle)
		}
	}

	if size.X == 0 || size.Y == 0 || size.Z == 0 {
		return nil, fmt.Errorf("invalid size %v", size)
	}

	rawVoxels := getRawVoxelsFromPointData(size, pointData)
	return rawVoxels, nil
}
