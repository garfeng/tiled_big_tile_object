package maker

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"sort"
)

type Maker struct {
	TileSize int
	DstCols  int
}

func (m *Maker) Extend(src gocv.Mat, dst *gocv.Mat) {
	tmp := gocv.NewMatWithSize(src.Rows()+m.TileSize*2, src.Cols()+m.TileSize*2, src.Type())
	defer tmp.Close()
	dstRect := image.Rect(m.TileSize, m.TileSize, tmp.Cols()-m.TileSize, tmp.Rows()-m.TileSize)
	dstRegion := tmp.Region(dstRect)
	defer dstRegion.Close()
	src.CopyTo(&dstRegion)
	dstRegion.CopyTo(dst)
}

func (m *Maker) SplitSrc(src gocv.Mat) []Object {
	chns := gocv.Split(src)
	defer func() {
		for _, v := range chns {
			v.Close()
		}
	}()

	alpha := chns[3]
	thres := gocv.NewMat()
	defer thres.Close()

	gocv.Threshold(alpha, &thres, 10, 255, gocv.ThresholdBinary)
	contours := gocv.FindContours(thres, gocv.RetrievalExternal, gocv.ChainApproxNone)
	defer contours.Close()

	contoursNumber := contours.Size()

	tileSize := TileSize(m.TileSize)

	objects := []Object{}

	for i := 0; i < contoursNumber; i++ {
		c := contours.At(i)
		rect := gocv.BoundingRect(c)
		rect.Min = tileSize.PointToTilePoint(rect.Min)
		maxPoint := tileSize.PointToTilePoint(rect.Max)
		if maxPoint.X < rect.Max.X {
			maxPoint.X += m.TileSize
		}
		if maxPoint.Y < rect.Max.Y {
			maxPoint.Y += m.TileSize
		}
		rect.Max = maxPoint

		cols := rect.Dx() / m.TileSize
		rows := rect.Dy() / m.TileSize

		existed := false
		for j, v := range objects {
			interset := v.Rect.Intersect(rect)
			if interset.Dx() > 10 || interset.Dy() > 10 {
				existed = true
				combine := v.Rect.Union(rect)
				objects[j].Rect = combine
				objects[j].Cols = combine.Dx() / m.TileSize
				objects[j].Rows = combine.Dy() / m.TileSize

				break
			}
		}

		if !existed {
			objects = append(objects, Object{
				Rect: rect,
				Cols: cols,
				Rows: rows,
			})
		}
	}

	return objects
}

func (m *Maker) Classify(objects []Object) []*ObjectGroup {
	groups := []*ObjectGroup{}

	for _, v := range objects {
		find := false
		for _, g := range groups {
			if g.Cols == v.Cols && g.Rows == v.Rows {
				g.Objects = append(g.Objects, v)
				find = true
				break
			}
		}
		if !find {
			newGroup := &ObjectGroup{
				Cols:    v.Cols,
				Rows:    v.Rows,
				Objects: []Object{v},
			}
			groups = append(groups, newGroup)
		}
	}
	return groups
}

type Object struct {
	Rect image.Rectangle
	Cols int
	Rows int
}

type ObjectGroup struct {
	Cols int
	Rows int

	Objects []Object
}

func (o *ObjectGroup) Less(i, j int) bool {
	y1 := o.Objects[i].Rect.Min.Y
	y2 := o.Objects[j].Rect.Min.Y

	if y1 != y2 {
		return y1 < y2
	}

	x1 := o.Objects[i].Rect.Min.X
	x2 := o.Objects[j].Rect.Min.X

	return x1 < x2
}

func (o *ObjectGroup) Len() int {
	return len(o.Objects)
}

func (o *ObjectGroup) Swap(i, j int) {
	tmp := o.Objects[i]
	o.Objects[i] = o.Objects[j]
	o.Objects[j] = tmp
}

func (o *ObjectGroup) Sort() {
	sort.Sort(o)
}

func (o *ObjectGroup) GenerateImage(src gocv.Mat, tileSize, dstWidth int, prefix string) {
	objectCols := dstWidth / o.Cols / tileSize
	objectRows := o.Len()/objectCols + 1
	dstHeight := objectRows * o.Rows * tileSize
	dst := gocv.NewMatWithSize(dstHeight, dstWidth, src.Type())

	rect0 := image.Rect(0, 0, o.Cols*tileSize, o.Rows*tileSize)

	srcroi := image.Rect(0, 0, src.Cols(), src.Rows())

	for i, v := range o.Objects {

		xId := i % objectCols
		yId := i / objectCols

		x := xId * tileSize * o.Cols
		y := yId * tileSize * o.Rows

		dstRect := rect0.Add(image.Pt(x, y))

		if !v.Rect.Intersect(srcroi).Eq(v.Rect) {
			fmt.Println("Err", i, "/", o.Len(), src.Cols(), src.Rows(), v.Rect.Max, "|", dst.Cols(), dst.Rows(), dstRect.Max)
			continue
		}

		dstRegion := dst.Region(dstRect)

		srcRegion := src.Region(v.Rect)
		srcRegion.CopyTo(&dstRegion)

		//srcRegion.Close()
		//dstRegion.Close()

	}

	name := fmt.Sprintf("%s_%dx%d.png", prefix, o.Cols, o.Rows)

	gocv.IMWrite(name, dst)
}

type TileSize int

func (t TileSize) PixToId(pix int) int {
	return pix / int(t)
}

func (t TileSize) IdToPix(id int) int {
	return id * int(t)
}

func (t TileSize) PixToTilePix(pix int) int {
	return pix / int(t) * int(t)
}

func (t TileSize) PointToTilePoint(pt image.Point) image.Point {
	return image.Point{
		X: t.PixToTilePix(pt.X),
		Y: t.PixToTilePix(pt.Y),
	}
}
