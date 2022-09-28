package filemanager

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func SaveFile(file multipart.File, fileName string) {
	versionInfo := getVersionInfo(file)
	fileNameSplit := strings.Split(fileName, ".")
	filePath := "./digitalFiles/" + fileNameSplit[0] + "_" + versionInfo[1] + "." + fileNameSplit[1]
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		writeCSV(versionInfo)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		io.Copy(f, file)
	}
}

func writeCSV(versionInfo []string) {
	f, err := os.OpenFile("./digitalFiles/digitalReleses.csv", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	writer := csv.NewWriter(f)
	writer.Write(versionInfo[:])
	writer.Flush()
	f.Close()
}

func getVersionInfo(file multipart.File) []string {
	body, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		panic(err)
	}
	for _, zipFile := range zipReader.File {
		if zipFile.Name == "Digital/Version.txt" {
			reader, err := zipFile.Open()
			if err != nil {
				panic(err)
			}
			bytes, err := io.ReadAll(reader)
			if err != nil {
				panic(err)
			}
			text := string(bytes)
			lines := strings.Split(text, "\n")

			var valueArray [3]string
			for i := 0; i < 3; i++ {
				valueArray[i] = strings.TrimSpace(strings.Split(lines[i+2], ":")[1])
			}

			return valueArray[:]
		}
	}
	return nil
}

func Setup() {
	os.Mkdir("digitalFiles", 0777)
}
