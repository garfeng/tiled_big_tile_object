package maker

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
)

type Maker struct {
	TileSize  int
	DstWidth  int
	DstHeight int
}

func (m *Maker) Extend(src gocv.Mat, dst *gocv.Mat) {
	tmp := gocv.NewMatWithSize(src.Rows()+m.TileSize*2, src.Cols()+m.TileSize*2, src.Type())
	defer tmp.Close()
	dstRect := image.Rect(m.TileSize, m.TileSize, tmp.Cols()-m.TileSize, tmp.Rows()-m.TileSize)
	dstRegion := tmp.Region(dstRect)
	defer dstRegion.Close()
	src.CopyTo(&dstRegion)
	tmp.CopyTo(dst)
}

func (m *Maker) SplitOneSrc(srcId int, src gocv.Mat) []Object {
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
		rect.Min = tileSize.PointToTilePoint(rect.Min, floor)
		rect.Max = tileSize.PointToTilePoint(rect.Max, ceil)

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
				Rect:  rect,
				Cols:  cols,
				Rows:  rows,
				SrcId: srcId,
			})
		}
	}

	return objects
}

func (m *Maker) SplitSrc(src []gocv.Mat) []Object {
	objects := []Object{}

	for i, one := range src {
		objects = append(objects, m.SplitOneSrc(i, one)...)
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

func (m *Maker) Generate(srcNames []string, dstPrefix string) error {
	srcList := []gocv.Mat{}
	labelList := []gocv.Mat{}
	for _, name := range srcNames {
		src, err := ImRead(name)
		if err != nil {
			return err
		}
		m.Extend(src, &src)
		srcList = append(srcList, src)
		labelList = append(labelList, src.Clone())
	}

	objects := m.SplitSrc(srcList)
	groups := m.Classify(objects)
	dstRoot, _ := filepath.Split(dstPrefix)
	os.MkdirAll(dstRoot, 0755)

	for _, v := range groups {
		v.Sort()
		v.GenerateImage(srcList, labelList, m.TileSize, m.DstWidth, m.DstHeight, dstPrefix)
	}
	for i, v := range labelList {
		_, name := filepath.Split(srcNames[i])
		err := ImWrite(filepath.Join(dstRoot, name), v)

		if err != nil {
			return err
		}

		v.Close()
		srcList[i].Close()
	}
	return nil
}

func ImRead(name string) (gocv.Mat, error) {
	buff, err := ioutil.ReadFile(name)
	if err != nil {
		return gocv.Mat{}, err
	}
	img, err := gocv.IMDecode(buff, gocv.IMReadUnchanged)
	if err != nil {
		return gocv.Mat{}, err
	}
	return img, nil
}
func ImWrite(name string, mat gocv.Mat) error {
	buff, err := gocv.IMEncode(".png", mat)
	if err != nil {
		return err
	}
	defer buff.Close()
	return ioutil.WriteFile(name, buff.GetBytes(), 0644)
}

type Object struct {
	Rect  image.Rectangle
	Cols  int
	Rows  int
	SrcId int
}

type ObjectGroup struct {
	Cols int
	Rows int

	Objects []Object
}

func (o *ObjectGroup) Less(i, j int) bool {
	id1 := o.Objects[i].SrcId
	id2 := o.Objects[j].SrcId
	if id1 != id2 {
		return id1 < id2
	}

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

func GetColor(w, h int) color.RGBA {
	c := color.RGBA{
		R: uint8((w * 100) % 256),
		G: uint8((h * 100) % 256),
		A: 255,
	}
	bi := c.A/10*3 + c.B/5*3
	b := (128 - int(bi)) * 10
	if b > 255 {
		b = 255
	}
	if b < 0 {
		b = 0
	}
	c.B = uint8(b)
	return c
}

func (o *ObjectGroup) GenerateImage(src []gocv.Mat, srcPutLabel []gocv.Mat, tileSize, dstWidth, dstHeight int, prefix string) {
	objectCols := dstWidth / o.Cols / tileSize
	dstWidth = objectCols * tileSize * o.Cols
	oneDstObjectRows := dstHeight / o.Rows / tileSize
	dstHeight = oneDstObjectRows * tileSize * o.Rows

	objectRows := ceil(float64(o.Len()) / float64(objectCols))
	dstImageNumber := ceil(float64(objectRows) / float64(oneDstObjectRows))

	dst := []gocv.Mat{}
	for i := 0; i < dstImageNumber; i++ {
		oneDst := gocv.NewMatWithSize(dstHeight, dstWidth, src[0].Type())
		defer oneDst.Close()
		dst = append(dst, oneDst)
	}

	rect0 := image.Rect(0, 0, o.Cols*tileSize, o.Rows*tileSize)

	drawPadding := image.Pt(tileSize/10+1, tileSize/10+1)
	label := fmt.Sprintf("%dx%d", o.Cols, o.Rows)
	drawColor := GetColor(o.Cols, o.Rows)

	for i, v := range o.Objects {

		xId := i % objectCols
		yId := i / objectCols

		objectImageId := yId / oneDstObjectRows
		yId = yId % oneDstObjectRows

		x := xId * tileSize * o.Cols
		y := yId * tileSize * o.Rows

		dstRect := rect0.Add(image.Pt(x, y))

		srcroi := image.Rect(0, 0, src[v.SrcId].Cols(), src[v.SrcId].Rows())

		if !v.Rect.Intersect(srcroi).Eq(v.Rect) {
			fmt.Println("Err", i, "/", o.Len(), srcroi, v.Rect.Max, "|", dst[objectImageId].Cols(), dst[objectImageId].Rows(), dstRect.Max)
			continue
		}

		dstRegion := dst[objectImageId].Region(dstRect)

		srcRegion := src[v.SrcId].Region(v.Rect)

		rectForDraw := v.Rect
		rectForDraw.Min = rectForDraw.Min.Add(drawPadding)
		rectForDraw.Max = rectForDraw.Max.Sub(drawPadding)
		//gocv.Rectangle(srcPutLabel, image.Rectangle{Min: rectForDraw.Min, Max: rectForDraw.Min.Add(image.Pt(40, 30))}, color.RGBA{A: 255}, -1)
		gocv.Rectangle(&srcPutLabel[v.SrcId], rectForDraw, drawColor, drawPadding.X/2+1)

		gocv.PutText(&srcPutLabel[v.SrcId], label, image.Pt(rectForDraw.Min.X+2, rectForDraw.Min.Y+20), gocv.FontHersheySimplex, 0.5, drawColor, 2)

		srcRegion.CopyTo(&dstRegion)

		srcRegion.Close()
		dstRegion.Close()

	}

	for i := 0; i < dstImageNumber; i++ {
		name := fmt.Sprintf("%s_%dx%d_%d.png", prefix, o.Cols, o.Rows, i+1)
		gocv.IMWrite(name, dst[i])
	}
}

var (
	Green = color.RGBA{0, 255, 0, 255}
)

type TileSize int

func (t TileSize) PixToTilePix(pix int, method roundMethod) int {
	return method(float64(pix)/float64(t)) * int(t)
}

func (t TileSize) PointToTilePoint(pt image.Point, method roundMethod) image.Point {
	return image.Point{
		X: t.PixToTilePix(pt.X, method),
		Y: t.PixToTilePix(pt.Y, method),
	}
}

type roundMethod func(x float64) int

func floor(x float64) int {
	return int(math.Floor(x))
}
func ceil(x float64) int {
	return int(math.Ceil(x))
}
