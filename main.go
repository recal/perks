package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()
	files, err := ioutil.ReadDir("raw")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		fmt.Println(fmt.Sprintf("---> %s", file.Name()))
		background := imagick.NewMagickWand()
		if err := background.ReadImage("perk-background.png"); err != nil {
			panic(err)
		}

		path := fmt.Sprintf("raw/%s", file.Name())
		out := fmt.Sprintf("with-background/%s", file.Name())
		perk := imagick.NewMagickWand()

		if err := perk.ReadImage(path); err != nil {
			panic(err)
		}

		perk.ResizeImage(560, 560, imagick.FILTER_BOX, 0)

		dx := background.GetImageHeight() - perk.GetImageHeight()
		dy := background.GetImageWidth() - perk.GetImageWidth()
		background.CompositeImage(perk, imagick.COMPOSITE_OP_ATOP, int(dx/2), int(dy/2))

		if err := background.WriteImage(out); err != nil {
			panic(err)
		}
	}
}
