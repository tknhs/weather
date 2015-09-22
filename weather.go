package main

import (
	"image"
	"image/color"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/url"
	"os"

	"github.com/soniakeys/quant/median"
)

type Weather struct {
	YahooToken string
	Filename   string
}

func (w *Weather) downloadImage(date, place string) (image.Image, error) {
	reqestUrl := "http://map.olp.yahooapis.jp/OpenLocalPlatform/V1/static"

	// base params
	values := url.Values{}
	values.Add("appid", w.YahooToken)
	values.Add("mode", "blankmap")
	values.Add("output", "png")
	values.Add("scalebar", "off")
	values.Add("overlay", "type:rainfall|datelabel:on|date:"+date)

	if place == "main" {
		values.Add("style", "bm")
		values.Add("bbox", "34.710834,137.726126,37.242060,137.700088")
		values.Add("pin1", "35.662713,139.709088,,blue")
		values.Add("pin2", "37.442060,138.819511,,red")
		values.Add("pin3", "35.0010405,135.7593765,,yellow")
		values.Add("z", "9")
		values.Add("width", "900")
		values.Add("height", "900")
	} else {
		values.Add("style", "bm.c.city:off|bm.p.13113:ccc|bm.p.15202:ccc|bm.p.26106:ccc")
		values.Add("height", "300")
		values.Add("width", "300")
		values.Add("z", "11")

		if place == "tokyo" {
			// tokyo
			values.Add("pin1", "35.662713,139.709088,,blue")
		} else if place == "nagaoka" {
			// nagaoka
			values.Add("pin2", "37.442060,138.819511,,red")
		} else if place == "kyoto" {
			// kyoto
			values.Add("pin3", "35.0010405,135.7593765,,yellow")
		}
	}

	resp, err := http.Get(reqestUrl + "?" + values.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	src, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return src, err
}

func (w *Weather) CreateGifImage(dateArray []string) error {
	g := &gif.GIF{
		Image:     []*image.Paletted{},
		Delay:     []int{},
		LoopCount: 0,
	}

	for _, date := range dateArray {
		srcMain, err := w.downloadImage(date, "main")
		if err != nil {
			return err
		}
		srcTokyo, err := w.downloadImage(date, "tokyo")
		if err != nil {
			return err
		}
		srcNagaoka, err := w.downloadImage(date, "nagaoka")
		if err != nil {
			return err
		}
		srcKyoto, err := w.downloadImage(date, "kyoto")
		if err != nil {
			return err
		}

		q := median.Quantizer(256)
		p := q.Quantize(make(color.Palette, 0, 256), srcMain)

		rMain := srcMain.Bounds()
		rTokyo := srcTokyo.Bounds()
		rNagaoka := srcNagaoka.Bounds()

		// 画像を合成
		minX := rMain.Min.X
		maxX := rMain.Max.X + rTokyo.Max.X
		minY := rMain.Min.Y
		maxY := rMain.Max.Y
		paletSize := image.Rect(0, 0, maxX, maxY)
		dst := image.NewPaletted(paletSize, p)

		for x := minX; x < maxX; x++ {
			divideX := x - rMain.Max.X
			for y := minY; y < maxY; y++ {
				if x < rMain.Max.X && y < rMain.Max.Y {
					// whole map
					dst.Set(x, y, srcMain.At(x, y))
				} else if x > rMain.Max.X {
					if y < rNagaoka.Max.Y {
						// nagaoka
						dst.Set(x, y, srcNagaoka.At(divideX, y))
					} else if y > rNagaoka.Max.Y && y < rNagaoka.Max.Y+rTokyo.Max.Y {
						// tokyo
						divideY := y - rNagaoka.Max.Y
						dst.Set(x, y, srcTokyo.At(divideX, divideY))
					} else if y > rNagaoka.Max.Y+rTokyo.Max.Y {
						// kyoto
						divideY := y - (rNagaoka.Max.Y + rTokyo.Max.Y)
						dst.Set(x, y, srcKyoto.At(divideX, divideY))
					}
				} else if x == rMain.Max.X /* || (x == rMain.Max.X && y == rNagaoka.Max.Y)*/ {
					// 境界線塗りつぶし
					dst.Set(x, y, color.RGBA{uint8(0), uint8(0), uint8(0), uint8(0)})
				}
			}
		}

		g.Image = append(g.Image, dst)
		g.Delay = append(g.Delay, 100)
	}

	out, err := os.Create(w.Filename + ".gif")
	if err != nil {
		return err
	}
	defer out.Close()
	err = gif.EncodeAll(out, g)
	if err != nil {
		return err
	}
	return nil
}

func NewWeather(config *Config) *Weather {
	w := &Weather{
		YahooToken: config.Yahoo.Token,
		Filename:   config.General.Filename,
	}
	return w
}
