package main

import (
	"archive/zip"
	"disrupt/github_reader/internal/server"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const ZIP_TEMP_DIR = "C:\\disrupt\\zip"
const DIR_TEMP = "C:\\disrupt\\unzip"

func main() {
	server := server.NewServer()
	server.Run()

	//profile := "nktrr"
	//repository := "disrupt_old"
	//baseUrl := "https://api.github.com/repos"
	//baseUrl += "/" + profile
	//baseUrl += "/" + repository
	//url := baseUrl + "/zipball/master"
	//resp, err := http.Get(url)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	return
	//}
	//zipName := profile + "-" + repository + ".zip"
	//
	//out, err := os.Create(filepath.Join(ZIP_TEMP_DIR, zipName))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer out.Close()
	//_, err = io.Copy(out, resp.Body)
	//err = Unzip(filepath.Join(ZIP_TEMP_DIR, zipName), filepath.Join(DIR_TEMP, "nktrr"))
	//if err != nil {
	//	println(err.Error())
	//}
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
