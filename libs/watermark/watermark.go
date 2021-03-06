package watermark

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 水印的位置
const (
	TopLeft Pos = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

// 错误类型
var (
	ErrUnsupportedWatermarkType = errors.New("不支持的水印类型")
	ErrInvalidPos               = errors.New("无效的 pos 值")
)

// 允许做水印的图片类型
var allowExts = []string{
	".gif", ".jpg", ".jpeg", ".png",
}

// Pos 表示水印的位置
type Pos int

// watermark 用于给图片添加水印功能
// 目前支持gif、jpeg 和 png 三种图片格式。
// 若是 gif 图片，则取图片的第一帧；png 支持透明背景
type Watermark struct {
	image   image.Image // 水印图片
	padding int         // 水印留的边白
	pos     Pos         // 水印的位置
}

// New 声明一个Watermark对象
//
// path 为水印文件的路径
// padding 为水印在目标图像上的留白大小；
// pos 水印的位置。
func New(path string, padding int, pos Pos) (*Watermark, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var img image.Image

	switch strings.ToLower(filepath.Ext(path)) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(f)
	case ".png":
		img, err = png.Decode(f)
	case ".gif":
		img, err = gif.Decode(f)
	default:
		return nil, ErrUnsupportedWatermarkType
	}

	if err != nil {
		return nil, err
	}

	if pos < TopLeft || pos > Center {
		return nil, ErrInvalidPos
	}

	return &Watermark{
		image:   img,
		padding: padding,
		pos:     pos,
	}, nil
}

// IsAllowExt 该扩展名的图片是否允许使用水印
//
// ext 必须带上 . 符号
func IsAllowExt(ext string) bool {
	ext = strings.ToLower(ext)

	for _, e := range allowExts {
		if e == ext {
			return true
		}
	}
	return false
}

// MarkFile 给指定的文件打上水印
func (w *Watermark) MarkFile(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	return w.Mark(file, strings.ToLower(filepath.Ext(path)))
}

// Mark 将水印写入 src 中，由 ext 确定当前图片的类型。
func (w *Watermark) Mark(src io.ReadWriteSeeker, ext string) (err error) {
	var srcImg image.Image

	ext = strings.ToLower(ext)
	switch ext {
	case ".gif":
		return w.markGIF(src)
	case ".jpg", ".jpeg":
		srcImg, err = jpeg.Decode(src)
	case ".png":
		srcImg, err = png.Decode(src)
	default:
		return ErrUnsupportedWatermarkType
	}

	if err != nil {
		return err
	}

	point := w.getPoing(srcImg.Bounds().Dx(), srcImg.Bounds().Dy())
	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.ZP, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)

	if _, err = src.Seek(0, 0); err != nil {
		return err
	}

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(src, dstImg, nil)
	case ".png":
		return png.Encode(src, dstImg)
	default:
		return ErrUnsupportedWatermarkType
	}
}

func (w *Watermark) markGIF(src io.ReadWriteSeeker) error {
	srcGIF, err := gif.DecodeAll(src)
	if err != nil {
		return err
	}
	bound := srcGIF.Image[0].Bounds()
	point := w.getPoing(bound.Dx(), bound.Dy())

	for index, img := range srcGIF.Image {
		dstImg := image.NewPaletted(img.Bounds(), img.Palette)
		draw.Draw(dstImg, dstImg.Bounds(), img, image.ZP, draw.Src)
		draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)
		srcGIF.Image[index] = dstImg
	}

	if _, err = src.Seek(0, 0); err != nil {
		return err
	}
	return gif.EncodeAll(src, srcGIF)
}

func (w *Watermark) getPoing(width, height int) image.Point {
	var point image.Point

	switch w.pos {
	case TopLeft:
		point = image.Point{X: -w.padding, Y: -w.padding}
	case TopRight:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -w.padding,
		}
	case BottomLeft:
		point = image.Point{
			X: -w.padding,
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case BottomRight:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case Center:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()) / 2,
			Y: -(height - w.padding - w.image.Bounds().Dy()) / 2,
		}
	default:
		panic(ErrInvalidPos)
	}

	return point
}
