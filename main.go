package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/gographics/imagick.v2/imagick"
)

const (
	backgroundImagePath = "perk-background.png"
	rawFolderPath       = "raw"
	outputFolderPath    = "out"
	resizeWidth         = 560
	resizeHeight        = 560
	chunkSize           = 5
)

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	if err := ensureOutputDirectory(); err != nil {
		panic(fmt.Errorf("failed to ensure output directory: %v", err))
	}

	files, err := os.ReadDir(rawFolderPath)
	if err != nil {
		panic(fmt.Errorf("failed to read raw folder: %v", err))
	}

	chunks := chunkFiles(files, chunkSize)

	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(files []fs.DirEntry) {
			defer wg.Done()
			processImageChunk(files)
		}(chunk)
	}

	wg.Wait()
}

func processImageChunk(files []fs.DirEntry) {
	for _, file := range files {
		fmt.Println(fmt.Sprintf("Processing image at path: %s", file.Name()))

		rawFilePath := filepath.Join(rawFolderPath, file.Name())
		outputFilePath := filepath.Join(outputFolderPath, file.Name())

		if err := processImage(rawFilePath, outputFilePath); err != nil {
			fmt.Printf("Error processing image %s: %v\n", file.Name(), err)
		}
	}
}


func processImage(inputFile, outputFile string) error {
	background := imagick.NewMagickWand()
	if err := background.ReadImage(backgroundImagePath); err != nil {
		return fmt.Errorf("failed to read background image: %v", err)
	}

	perk := imagick.NewMagickWand()
	if err := perk.ReadImage(inputFile); err != nil {
		return fmt.Errorf("failed to read image %s: %v", filepath.Base(inputFile), err)
	}

	perk.ResizeImage(resizeWidth, resizeHeight, imagick.FILTER_BOX, 0)

	dx := background.GetImageHeight() - perk.GetImageHeight()
	dy := background.GetImageWidth() - perk.GetImageWidth()
	background.CompositeImage(perk, imagick.COMPOSITE_OP_ATOP, int(dx/2), int(dy/2))

	if err := background.WriteImage(outputFile); err != nil {
		return fmt.Errorf("failed to write output image %s: %v", filepath.Base(outputFile), err)
	}

	return nil
}

func chunkFiles(files []fs.DirEntry, chunkSize int) [][]fs.DirEntry {
	var chunks [][]fs.DirEntry

	for i := 0; i < len(files); i += chunkSize {
		end := i + chunkSize
		if end > len(files) {
			end = len(files)
		}
		chunks = append(chunks, files[i:end])
	}

	return chunks
}

func ensureOutputDirectory() error {
	_, err := os.Stat(outputFolderPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(outputFolderPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create output directory: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check output directory: %v", err)
	}

	return nil
}