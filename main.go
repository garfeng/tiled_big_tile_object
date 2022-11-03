package main

import (
	"flag"
	"fmt"
	"github.com/garfeng/tiled_big_tile_object/maker"
	"os"
	"path/filepath"
	"strings"
)

var (
	tileSize  = flag.Int("tileSize", 48, "Size of each tile.(RMVA:32, RMMV:48)")
	DstCols   = flag.Int("dstWidth", 480, "Width of dst image. will be tileSize x N")
	srcRoot   = flag.String("srcRoot", "", "src images dir")
	dstPrefix = flag.String("dstPrefix", "dst/tiled", "prefix of output image")
)

func main() {
	flag.Parse()
	if srcRoot == nil || *srcRoot == "" || (*tileSize) <= 0 || (*DstCols) < *tileSize {
		flag.PrintDefaults()
		return
	}

	sz := *tileSize

	m := &maker.Maker{
		TileSize: sz,
		DstCols:  (*DstCols) / sz * sz,
	}

	srcImages, err := scanPngs(*srcRoot)
	if err != nil {
		fmt.Println(err)
	}
	if len(srcImages) == 0 {
		fmt.Println("No images found in", *srcRoot)
		flag.PrintDefaults()
		return
	}

	err = m.Generate(srcImages, *dstPrefix)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func scanPngs(root string) ([]string, error) {
	fp, err := os.Open(root)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	names, err := fp.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, v := range names {
		ext := filepath.Ext(v)
		ext = strings.ToLower(ext)
		if ext == ".png" {
			res = append(res, filepath.Join(root, v))
		}
	}
	return res, nil
}
