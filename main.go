package main

import (
	"github.com/garfeng/tiled_big_tile_object/maker"
	"gocv.io/x/gocv"
	"os"
)

func main() {
	m := &maker.Maker{
		TileSize: 48,
		DstCols:  480,
	}

	src := gocv.IMRead("./examples/src/Map005.png", gocv.IMReadUnchanged)
	m.Extend(src, &src)
	objects := m.SplitSrc(src)
	groups := m.Classify(objects)
	os.RemoveAll("tmp")
	os.Mkdir("tmp", 0755)
	for _, v := range groups {
		v.Sort()
		v.GenerateImage(src, m.TileSize, m.DstCols, "tmp/Map005")
	}
}
