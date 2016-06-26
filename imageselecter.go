// imageselecter.go
package main

import (
	"math/rand"
	"time"
)

type ImageSelecter struct {
	Images  []string
	Counter int
	Repeat  bool
}

func NewImageSelecter(img []string, repeat bool) *ImageSelecter {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 0
	limg := len(img)
	IS := &ImageSelecter{
		Images:  make([]string, limg),
		Counter: 0,
		Repeat:  repeat,
	}
	for i := range img {
		n = r.Intn(limg)
		IS.Images[i], IS.Images[n] = img[n], img[i]
	}
	return IS
}

func (is *ImageSelecter) Mesh() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	limg := len(is.Images)
	n := 0
	for i := range is.Images {
		n = r.Intn(limg)
		is.Images[i], is.Images[n] = is.Images[n], is.Images[i]
	}
}

func (is *ImageSelecter) Next() (string, bool) {
	if is.Counter == len(is.Images)-1 {
		if !is.Repeat {
			return "", false
		} else {
			is.Counter = 0
			is.Mesh()
		}
	}
	is.Counter++
	return is.Images[is.Counter], true
}
