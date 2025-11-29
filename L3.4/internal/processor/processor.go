package processor

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	. "image-processor/internal/models"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
)

type Processor struct{}

func NewProcessor() *Processor {
	return &Processor{}
}

func (p *Processor) Process(task Task, data []byte) ([]byte, error) {
	// декодируем исходное изображение
	src, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var result image.Image

	// выбираем действие
	switch task.Type {
	case "resize":
		result = imaging.Resize(src, task.Params.Width, task.Params.Height, imaging.Lanczos)

	case "thumbnail":
		result = imaging.Thumbnail(src, 300, 300, imaging.Lanczos)

	case "watermark":
		wm, err := makeWatermark(task.Params.Watermark)
		if err != nil {
			return nil, fmt.Errorf("failed to create watermark: %w", err)
		}
		if wm == nil {
			return nil, errors.New("watermark image is nil")
		}
		result = imaging.Overlay(
			src,
			wm,
			image.Pt(30, 30),
			0.7,
		)

	default:
		return nil, fmt.Errorf("unknown task type: %s", task.Type)
	}

	// кодируем обратно
	buf := new(bytes.Buffer)
	err = imaging.Encode(buf, result, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func makeWatermark(text string) (image.Image, error) {
	const (
		w = 500
		h = 120
	)

	rgba := image.NewRGBA(image.Rect(0, 0, w, h))

	// полностью прозрачный фон
	draw.Draw(rgba, rgba.Bounds(), &image.Uniform{color.NRGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	// загружаем шрифт
	fontBytes, err := os.ReadFile("assets/Roboto-Bold.ttf")
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(42)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.NewUniform(color.NRGBA{255, 255, 255, 160})) // белый + прозрачность

	pt := freetype.Pt(30, 70)
	_, err = c.DrawString(text, pt)
	if err != nil {
		return nil, err
	}

	return rgba, nil
}
