package main

import (
	"context"
	"fmt"
	"github.com/garfeng/tiled_big_tile_object/maker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) SelectImages() ([]string, error) {
	return runtime.OpenMultipleFilesDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select images. Press CTRL to select multiple",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "png",
				Pattern:     "*.png",
			},
		},
		ShowHiddenFiles:      false,
		CanCreateDirectories: false,
	})
}

func (a *App) SelectDstRoot() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                      "Select a directory to save dst image",
		CanCreateDirectories:       true,
		ResolvesAliases:            false,
		TreatPackagesAsDirectories: false,
	})
}

func (a *App) Generate(param Param) error {
	m := maker.Maker{
		TileSize:  param.TileSize,
		DstWidth:  param.DstWidth,
		DstHeight: param.DstHeight,
	}
	return m.Generate(param.SrcImages, filepath.Join(param.DstRoot, param.DstPrefix))
}

type Param struct {
	TileSize  int      `json:"tileSize"`
	DstWidth  int      `json:"dstWidth"`
	DstHeight int      `json:"dstHeight'"`
	SrcImages []string `json:"srcImages"`
	DstRoot   string   `json:"dstRoot"`
	DstPrefix string   `json:"dstPrefix"`
}
