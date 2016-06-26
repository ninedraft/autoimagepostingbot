// imageselecter.go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type ImageSelecter struct {
	Config BotConfig
	Name   string
	Images []os.FileInfo
	Cursor int
	Size   int
}

func NewSelecter(Config BotConfig, dir os.FileInfo) (*ImageSelecter, error) {
	files, err := ioutil.ReadDir(Config.Root + "/" + dir.Name())
	if err != nil {
		return nil, err
	}
	images := []os.FileInfo{}
	for _, f := range files {
		if !f.IsDir() {
			images = append(images, f)
		}
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("can't find any image")
	}
	ims := &ImageSelecter{
		Config: Config,
		Name:   dir.Name(),
		Images: make([]os.FileInfo, len(images)),
		Cursor: 0,
		Size:   len(images),
	}
	copy(ims.Images, images)
	return ims, nil
}

func (ims *ImageSelecter) Next() (string, bool) {
	n := Rand.Intn(ims.Size-ims.Cursor) + ims.Cursor
	next := ims.Images[n]
	ims.Images[n], ims.Images[ims.Cursor] = ims.Images[ims.Cursor], ims.Images[n]
	ims.Cursor++
	ok := true
	if ims.Cursor == ims.Size {
		if ims.Config.Repeat {
			ims.Cursor = 0
		} else {
			ok = false
		}
	}
	image := ims.Config.Root + "/" + ims.Name + "/" + next.Name()
	return image, ok
}
