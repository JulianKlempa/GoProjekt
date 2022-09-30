package filemanager

import (
	"archive/zip"
	"bytes"
	"digitalDistribution/configuration"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Setup() {
	os.Mkdir("digitalFiles", 0777)
	enforceUploadLimit()
}

func SaveFile(file multipart.File, fileName string) {
	versionInfo := getVersionInfo(file)
	filePath := "./digitalFiles/Digital_" + versionInfo[1] + ".zip"
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		writeCSV(versionInfo, filePath)
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		file.Seek(0, io.SeekStart)
		io.Copy(f, file)
		enforceUploadLimit()
	}
}

func GetFile(fileName string) *bytes.Reader {
	f, err := os.OpenFile("./digitalFiles/"+fileName, os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func GetCurrentData() [][]string {
	f, err := os.OpenFile("./digitalFiles/digitalReleases.csv", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}
	return data
}

func writeNewCSV(data [][]string) {
	if err := os.Truncate("./digitalFiles/digitalReleases.csv", 0); err != nil {
		panic(err)
	}
	f, err := os.OpenFile("./digitalFiles/digitalReleases.csv", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	writer.WriteAll(data)
	writer.Flush()
}

func writeCSV(versionInfo []string, filepath string) {
	f, err := os.OpenFile("./digitalFiles/digitalReleases.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	writer := csv.NewWriter(f)
	versionInfo = append(versionInfo, filepath)
	writer.Write(versionInfo[:])
	writer.Flush()
}

func enforceUploadLimit() {
	limit := configuration.ReadConfig().SavesCount
	data := GetCurrentData()
	sort.Slice(data, func(i, j int) bool {
		iStrings := strings.Split(strings.TrimLeft(data[i][1], "v"), ".")
		jStrings := strings.Split(strings.TrimLeft(data[j][1], "v"), ".")

		iMajor, err1 := strconv.ParseInt(iStrings[0], 0, 64)
		iMinor, err2 := strconv.ParseInt(iStrings[1], 0, 64)
		jMajor, err3 := strconv.ParseInt(jStrings[0], 0, 64)
		jMinor, err4 := strconv.ParseInt(jStrings[1], 0, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			panic(err1.Error() + err2.Error() + err3.Error() + err4.Error())
		}
		if iMajor == jMajor {
			return iMinor < jMinor
		}
		return iMajor < jMajor
	})
	for limit < len(data) {
		var element []string
		element, data = data[0], data[1:]
		err := os.Remove(element[3])
		if err != nil {
			fmt.Println(err)
		}
	}
	writeNewCSV(data)
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
