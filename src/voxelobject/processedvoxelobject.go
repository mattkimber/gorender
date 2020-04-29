package voxelobject

import (
	"geometry"
)

type ProcessedElement struct {
	Normal         geometry.Vector3
	AveragedNormal geometry.Vector3
	Index          byte
	IsSurface      bool
}

type ProcessedVoxelObject struct {
	Elements [][][]ProcessedElement
	Size     geometry.Point
}

const normalRadius = 5

func (r RawVoxelObject) GetProcessedVoxelObject() (p ProcessedVoxelObject) {
	p.Size = r.Size()

	p.setElements(r)
	p.calculateNormals()

	return
}

func (p *ProcessedVoxelObject) calculateNormals() {
	for x := 0; x < p.Size.X; x++ {
		for y := 0; y < p.Size.Y; y++ {
			for z := 0; z < p.Size.Z; z++ {
				p.Elements[x][y][z].IsSurface = p.isSurface(x, y, z)
				p.Elements[x][y][z].Normal = p.calculateNormal(x, y, z)
			}
		}
	}
}

func (p *ProcessedVoxelObject) calculateNormal(x, y, z int) (normal geometry.Vector3) {
	if !p.SafeGetData(x,y,z).IsSurface {
		return
	}

	for i := -normalRadius; i <= normalRadius; i++ {
		for j := -normalRadius; j <= normalRadius; j++ {
			for k := -normalRadius; k <= normalRadius; k++ {
				if (i*i)+(j*j)+(k*k) <= (normalRadius*normalRadius) && p.SafeGetData(x+i, y+j, z+k).Index == 0 {
					normal = normal.Subtract(geometry.Vector3{X: float64(i), Y: float64(j), Z: float64(k)})
				}
			}
		}
	}

	return normal.Normalise()
}

func (p *ProcessedVoxelObject) isSurface(x, y, z int) bool {
	// A voxel is a surface voxel if any of the adjacent directions is zero
	// The edges of the voxel object are trivially surface voxels
	return  p.Elements[x][y][z].Index != 0 && (
			x == 0 || y == 0 || z == 0 ||
			x == p.Size.X-1 || y == p.Size.Y-1 || z == p.Size.Z-1 ||
			p.Elements[x+1][y][z].Index == 0 || p.Elements[x-1][y][z].Index == 0 ||
			p.Elements[x][y+1][z].Index == 0 || p.Elements[x][y-1][z].Index == 0 ||
			p.Elements[x][y][z-1].Index == 0 || p.Elements[x-1][y][z+1].Index == 0 ||
			p.Elements[x][y][z+1].Index == 0 || p.Elements[x+1][y][z+1].Index == 0 ||
			p.Elements[x][y-1][z+1].Index == 0 || p.Elements[x][y+1][z+1].Index == 0)
}

func (p *ProcessedVoxelObject) setElements(r RawVoxelObject) {
	p.Elements = make([][][]ProcessedElement, p.Size.X)
	for x := 0; x < p.Size.X; x++ {
		p.Elements[x] = make([][]ProcessedElement, p.Size.Y)
		for y := 0; y < p.Size.Y; y++ {
			p.Elements[x][y] = make([]ProcessedElement, p.Size.Z)
			for z := 0; z < p.Size.Z; z++ {
				p.Elements[x][y][z].Index = r[x][y][z]
			}
		}
	}
}

func (pv *ProcessedVoxelObject) SafeGetData(x, y, z int) (pe ProcessedElement) {
	if x >= 0 && y >= 0 && z >= 0 && x < pv.Size.X && y < pv.Size.Y && z < pv.Size.Z {
		pe = pv.Elements[x][y][z]
	}

	return
}
