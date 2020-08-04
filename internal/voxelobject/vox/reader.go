package vox

import (
	"encoding/binary"
	"fmt"
	"github.com/mattkimber/gorender/internal/geometry"
	"github.com/mattkimber/gorender/internal/utils/byteutils"
	"io"
	"io/ioutil"
)

const magic = "VOX "

type MagicaVoxelObject [][][]byte

func isHeaderValid(handle io.Reader) bool {
	result, err := getChunkHeader(handle)
	return err == nil && result == magic
}

func getChunkHeader(handle io.Reader) (string, error) {
	limitedReader := io.LimitReader(handle, 4)
	result, err := ioutil.ReadAll(limitedReader)
	return string(result), err
}

func getSizeFromChunk(handle io.Reader) (geometry.Point, error) {
	data, err := getChunkData(handle, 12)

	if err != nil {
		return geometry.Point{}, err
	}

	return geometry.Point{
		X: int(binary.LittleEndian.Uint32(data[0:4])),
		Y: int(binary.LittleEndian.Uint32(data[4:8])),
		Z: int(binary.LittleEndian.Uint32(data[8:12])),
	}, nil
}

func getPointDataFromChunk(handle io.Reader) ([]geometry.PointWithColour, error) {
	data, err := getChunkData(handle, 4)

	if err != nil {
		return getNilValueForPointDataFromChunk(), err
	}

	result := make([]geometry.PointWithColour, len(data)/4)

	for i := 0; i < len(data); i += 4 {
		point := geometry.PointWithColour{
			Point: geometry.Point{X: int(data[i]), Y: int(data[i+1]), Z: int(data[i+2])}, Colour: data[i+3],
		}

		result[i/4] = point
	}

	return result, nil
}

func getVoxelObjectFromPointData(size geometry.Point, data []geometry.PointWithColour) MagicaVoxelObject {
	result := byteutils.Make3DByteSlice(size)

	for _, p := range data {
		if p.Point.X < size.X && p.Point.Y < size.Y && p.Point.Z < size.Z && p.Colour != 0 {
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

	parsedSize := int64(binary.LittleEndian.Uint32(size[0:4]))
	return parsedSize
}

func getNilValueForPointDataFromChunk() []geometry.PointWithColour {
	return []geometry.PointWithColour{}
}

func GetMagicaVoxelObject(handle io.Reader) (MagicaVoxelObject, error) {
	if !isHeaderValid(handle) {
		return nil, fmt.Errorf("header not valid")
	}
	getChunkHeader(handle)

	size := geometry.Point{}
	pointData := make([]geometry.PointWithColour, 0)

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

	rawVoxels := getVoxelObjectFromPointData(size, pointData)
	return rawVoxels, nil
}

func (v *MagicaVoxelObject) GetFromReader(handle io.Reader) (err error) {
	*v, err = GetMagicaVoxelObject(handle)
	return
}
