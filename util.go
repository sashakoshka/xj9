package main

import "os"
import "image"
import "os/exec"
import "github.com/faiface/pixel"

func loadPicture(path string) (picture pixel.Picture) {
	path = "images/" + path

	file, err := os.Open(path)
	defer file.Close()
	if err != nil { panic(err) }

	img, _, err := image.Decode(file)
	if err != nil { panic(err) }

	return pixel.PictureDataFromImage(img)
}

func speak(message string) () {
	exec.Command("espeak", "-v", "annie", "-p", "75", message)
}
