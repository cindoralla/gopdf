package gopdf

import (
	"os"
	"path/filepath"
	"fmt"
	"time"

	"github.com/cindoralla/gopdf/core"
)

type Image struct {
	pdf           *core.Report
	autobreak     bool
	path          string
	width, height float64
	margin        core.Scope
	temppath      []string
}

func NewImage(path string, pdf *core.Report) *Image {
	return NewImageWithWidthAndHeight(path, 0, 0, pdf)
}

func NewImageWithWidthAndHeight(path string, width, height float64, pdf *core.Report) *Image {
	contentWidth, contentHeight := pdf.GetContentWidthAndHeight()
	if width > contentWidth {
		width = contentWidth
	}
	if height > contentHeight {
		height = contentHeight
	}

	var temppath []string
	var tempDir string = os.TempDir()
	if _, err := os.Stat(path); err != nil {
		path = fmt.Sprintf("%s/%v.png", tempDir, time.Now().Unix())
		temppath = append(temppath, path)
		DrawPNG(path)
	}

	srcPath, _ := filepath.Abs(path)

	// if the file is png format, do not change to jpeg, otherwise, change it to jpeg
	var dstPath string
	fd, _ := os.Open(srcPath)
	_, pictureType, _ := image.DecodeConfig(fd)
	switch pictureType {
	case PNG:
		dstPath = fmt.Sprintf("%s/%v.png", tempDir, time.Now().UnixNano())

		input, _ := ioutil.ReadFile(srcPath)
		ioutil.WriteFile(dstPath, input, 0644)

		//destination, _ := os.Create(dstPath)
		//io.Copy(destination, fd)
		//destination.Seek(0, 0)

	default:
		dstPath = fmt.Sprintf("%s/%v.jpeg", tempDir, time.Now().UnixNano())
		Convert2JPEG(srcPath, dstPath)
	}

	temppath = append(temppath, dstPath)

	image := &Image{
		pdf:      pdf,
		path:     dstPath,
		temppath: temppath,
	}

	image.imageWidthHeight(image.path, width, height)

	if dstPath != "" {
		pdf.AddCallBack(image.delTempImage)
	}

	return image
}

func (image *Image) imageWidthHeight(path string, width, height float64) *Image {
	w, h := GetImageWidthAndHeight(path)
	if width > 0 && height > 0 {
		if float64(h)*width/float64(w) > height {
			width = float64(w) * height / float64(h)
		} else {
			height = float64(h) * width / float64(w)
		}
	} else if width > 0 {
		height = float64(h) * width / float64(w)
	} else if height > 0 {
		width = float64(w) * height / float64(h)
	} else {
		width, height = float64(w), float64(h)
	}

	image.width = width
	image.height = height

	return image
}

func (image *Image) SetMargin(margin core.Scope) *Image {
	margin.ReplaceMarign()
	image.margin = margin
	return image
}

func (image *Image) GetHeight() float64 {
	return image.height
}
func (image *Image) GetWidth() float64 {
	return image.width
}

func (image *Image) SetAutoBreak() {
	image.autobreak = true
}

// 自动换行
func (image *Image) GenerateAtomicCell() (pagebreak, over bool, err error) {
	var (
		sx, sy = image.pdf.GetXY()
	)

	x, y := sx+image.margin.Left, sy+image.margin.Top
	pageEndX, pageEndY := image.pdf.GetPageEndXY()
	if y < pageEndY && y+float64(image.height) > pageEndY {
		if image.autobreak {
			image.pdf.AddNewPage(false)
			goto draw
		}

		return true, false, nil
	}

draw:
	image.pdf.Image(image.path, x, y, x+float64(image.width), y+float64(image.height))
	if x+float64(image.width) >= pageEndX {
		sx, _ = image.pdf.GetPageStartXY()
		image.pdf.SetXY(sx, y+float64(image.height)+image.margin.Bottom)
	} else {
		image.pdf.SetXY(x+float64(image.width), y)
	}

	if image.autobreak {
		sx, _ = image.pdf.GetPageStartXY()
		_, sy := image.pdf.GetXY()
		image.pdf.SetXY(sx, sy+image.height)
	}

	return false, true, nil
}

func (image *Image) delTempImage(report *core.Report) {
	if image.temppath == nil {
		return
	}

	for _, path := range image.temppath {
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			os.Remove(path)
		}
	}
}
