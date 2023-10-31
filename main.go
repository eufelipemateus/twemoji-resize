package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

// Create folder if do not exists
func createDir(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

const PATH_OUTPUT = "./output"
const PATH_SVG = "./svg"

func main() {
	sizeArg := flag.String("size", "", "Return Size resize image.")
	flag.Parse()

	if *sizeArg == "" {
		println("Size arg don't define")
		os.Exit(0)
	}

	size, _ := strconv.ParseInt(*sizeArg, 10, 32)
	pathPNG := fmt.Sprintf("%s/%dx%d", PATH_OUTPUT, size, size)

	createDir(pathPNG)

	entries, err := os.ReadDir(PATH_SVG)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		data, _ := os.Open(filepath.Join(PATH_SVG, e.Name()))

		if err != nil {
			panic(err)
		}
		defer data.Close()

		icon, _ := oksvg.ReadIconStream(data)
		icon.SetTarget(0, 0, float64(size), float64(size))

		rgba := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
		icon.Draw(rasterx.NewDasher(int(size), int(size), rasterx.NewScannerGV(int(size), int(size), rgba, rgba.Bounds())), 1)

		var extension = filepath.Ext(e.Name())
		var name = e.Name()[0 : len(e.Name())-len(extension)]
		newFile := fmt.Sprintf("%s/%s.png", pathPNG, name)
		println(newFile)

		file, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic("Cannot open file")
		}
		defer file.Close()

		err = png.Encode(file, rgba)
		if err != nil {
			panic(err)
		}
	}

}
