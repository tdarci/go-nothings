package filters

import (
	"image"
	"image/color"
)

func Redden(im image.Image) image.Image {
	outImage := image.NewRGBA(im.Bounds())
	colorModel := outImage.ColorModel()

	var blendingRed uint32 = 65535
	var redChangeIncrement uint32 = 65535 / uint32(im.Bounds().Max.X-im.Bounds().Min.X)

	for x := im.Bounds().Min.X; x <= im.Bounds().Max.X; x++ {

		// decrease our red-ness each time we move across horizontally
		blendingRed -= redChangeIncrement
		if blendingRed < 0 {
			blendingRed = 0
		}

		for y := im.Bounds().Min.Y; y <= im.Bounds().Max.Y; y++ {
			cRed, cGreen, cBlue, _ := im.At(x, y).RGBA()
			// calculate a new red value that averages the current red with our red
			newRed := (cRed + blendingRed) / 2

			// WHY divide by 0x101?
			// see here: https://jimsaunders.net/2015/05/22/manipulating-colors-in-go.html
			// and here: https://stackoverflow.com/questions/35374300/why-does-golang-rgba-rgba-method-use-and
			// and here: https://blog.golang.org/image
			newColor := color.RGBA{R: uint8(newRed / 0x101), G: uint8(cGreen / 0x101), B: uint8(cBlue / 0x101), A: 255}
			outImage.Set(x, y, colorModel.Convert(newColor))
		}
	}
	return outImage
}
