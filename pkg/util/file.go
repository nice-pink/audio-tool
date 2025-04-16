package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/nice-pink/goutil/pkg/filesystem"
)

func GetFilePath(baseFilePath string) string {
	if baseFilePath == "" {
		return ""
	}
	now := time.Now()
	return baseFilePath + "_" + strconv.FormatInt(now.Unix(), 10)
}

func CleanUp(folder string, sec int64, ignoreHiddenFiles bool, delete bool) {
	// get files
	files := filesystem.ListFiles(folder, sec, ignoreHiddenFiles)

	// logs
	ago := -time.Duration(sec) * time.Second
	dateThreshold := time.Now().Add(ago)
	fmt.Println("Files older than:", dateThreshold)
	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println(strconv.Itoa(len(files)), "files")
	fmt.Println()

	// delete
	if delete {
		fmt.Println("Delete files...")
		filesystem.DeleteFiles(files)
		fmt.Println("Done!")
	}
}

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
