package voxelobject

import (
	"colour"
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
	Palette  *colour.Palette
}

const normalRadius = 3
const normalAverageDistance = 1

func (r RawVoxelObject) GetProcessedVoxelObject(pal *colour.Palette) (p ProcessedVoxelObject) {
	p.Size = r.Size()
	p.Palette = pal

	p.setElements(r)
	p.calculateFirstPassData()
	p.calculateSecondPassData()

	return
}

func (p *ProcessedVoxelObject) calculateFirstPassData() {
	for x := 0; x < p.Size.X; x++ {
		for y := 0; y < p.Size.Y; y++ {
			for z := 0; z < p.Size.Z; z++ {
				p.Elements[x][y][z].IsSurface = p.isSurface(x, y, z)
				p.Elements[x][y][z].Normal = p.calculateNormal(x, y, z)
			}
		}
	}
}

func (p *ProcessedVoxelObject) calculateSecondPassData() {
	for x := 0; x < p.Size.X; x++ {
		for y := 0; y < p.Size.Y; y++ {
			for z := 0; z < p.Size.Z; z++ {
				p.Elements[x][y][z].AveragedNormal = p.getAverageNormal(x, y, z)
			}
		}
	}
}

func (p *ProcessedVoxelObject) getNormalRadius(index byte) (radius int) {
	radius = normalRadius + (p.Palette.GetSmoothness(index) * 2)
	if radius < 1 {
		return 1
	}
	return
}

func (p *ProcessedVoxelObject) calculateNormal(x, y, z int) (normal geometry.Vector3) {
	if !p.SafeGetData(x, y, z).IsSurface {
		return
	}

	radius := p.getNormalRadius(p.SafeGetData(x, y, z).Index)
	for i := -radius; i <= radius; i++ {
		for j := -radius; j <= radius; j++ {
			for k := -radius; k <= radius; k++ {
				if (i*i)+(j*j)+(k*k) <= (radius*radius) && p.SafeGetData(x+i, y+j, z+k).Index == 0 {
					normal = normal.Subtract(geometry.Vector3{X: float64(i), Y: float64(j), Z: float64(k)})
				}
			}
		}
	}

	if normal.Length() > 0.01 {
		return normal.Normalise()
	}

	return normal
}

func (p *ProcessedVoxelObject) getNormalAverageDistance(index byte) (distance int) {
	distance = normalAverageDistance + (p.Palette.GetSmoothness(index))
	if distance < 0 {
		return 0
	}
	return
}

func (p *ProcessedVoxelObject) getAverageNormal(x, y, z int) (normal geometry.Vector3) {
	if !p.SafeGetData(x, y, z).IsSurface {
		return
	}

	smoothness := p.Palette.GetSmoothness(p.SafeGetData(x, y, z).Index)

	distance := p.getNormalAverageDistance(p.SafeGetData(x, y, z).Index)
	for i := -distance; i <= distance; i++ {
		for j := -distance; j <= distance; j++ {
			for k := -distance; k <= distance; k++ {
				if p.SafeGetData(x+i, y+j, z+k).Index != 0 {
					if p.Palette.GetSmoothness(p.SafeGetData(x+i, y+i, z+i).Index) == smoothness {
						normal = normal.Add(p.SafeGetData(x+i, y+j, z+k).Normal)
					}
				}
			}
		}
	}

	if normal.Length() < 0.01 {
		return p.SafeGetData(x, y, z).Normal
	}

	return normal.Normalise()
}

func (p *ProcessedVoxelObject) isSurface(x, y, z int) bool {
	// A voxel is a surface voxel if any of the adjacent directions is zero
	// The edges of the voxel object are trivially surface voxels
	return p.Elements[x][y][z].Index != 0 && (x == 0 || y == 0 || z == 0 ||
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

func (pv *ProcessedVoxelObject) Invalid() bool {
	return pv.Size.X == 0 || pv.Size.Y == 0 || pv.Size.Z == 0
}
