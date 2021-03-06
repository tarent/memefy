package persistence

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const basePath = "files/"

func init() {
	err := os.MkdirAll(basePath, 0777)
	if err != nil {
		log.Fatal("Storage '%s' not available", basePath, err)
	}
}

func ListMemes() (memes []string, err error) {
	os.MkdirAll(basePath, 0755)
	contents, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	for i := range contents {
		if contents[i].IsDir() {
			memes = append(memes, contents[i].Name())
		}
	}
	return
}

func SaveMeme(name string, r io.Reader) {
	if err := os.MkdirAll(basePath+name, 0755); err != nil {
		log.Printf("Could not create dir: %s", err.Error())
		os.RemoveAll(basePath + name)
	}
	file, err := os.Create(basePath + name + "/out.mp4")
	if err != nil {
		log.Printf("Could not create file: %s", err.Error())
		os.RemoveAll(basePath + name)
	}
	_, err = io.Copy(file, r)
	if err != nil {
		log.Printf("Could not write file: %s", err.Error())
		os.RemoveAll(basePath + name)
	}
}
