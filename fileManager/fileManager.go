package filemanager

import (
	"archive/zip"
	"bytes"
	"digitalDistribution/configuration"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Version struct {
	MajorVersion int
	MinorVersion int
}

type File struct {
	BuildDate       string
	FileData        []byte
	Version         Version
	DownloadCounter int
}

type Storage struct {
	Files []File
}

func (s Storage) Contains(version Version) bool {
	for _, element := range s.Files {
		element.Version.Equals(version)
		return true
	}
	return false
}

func (s Storage) Get(version Version) File {
	for _, element := range s.Files {
		if element.Version.Equals(version) {
			return element
		}
	}
	return File{}
}

func (v Version) Equals(version Version) bool {
	return v.MajorVersion == version.MajorVersion && v.MinorVersion == version.MinorVersion
}

var storage Storage

func Setup() {
	os.Mkdir("digitalFiles", 0777)
	f, err := os.OpenFile("./digitalFiles/files.json", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &storage)
	enforceUploadLimit()
	writeToFile()
}

func SaveFile(file multipart.File) {
	jsonFile := File{}
	versionInfo := getVersionInfo(file)
	versionStrings := strings.Split(versionInfo[1], ".")
	version := Version{}
	majorVersion, err := strconv.ParseInt(strings.TrimLeft(versionStrings[0], "v"), 0, 64)
	if err != nil {
		panic(err)
	}
	minorVersion, err := strconv.ParseInt(versionStrings[1], 0, 64)
	if err != nil {
		panic(err)
	}
	version.MajorVersion = int(majorVersion)
	version.MinorVersion = int(minorVersion)
	if !storage.Contains(version) {
		jsonFile.Version = version
		jsonFile.BuildDate = versionInfo[2]
		jsonFile.DownloadCounter = 0
		data, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}
		jsonFile.FileData = data
	}
	enforceUploadLimit()
	writeToFile()
}

func GetFile(version Version) *bytes.Reader {
	file := storage.Get(version)
	return bytes.NewReader(file.FileData)
}

func GetStorage() Storage {
	return storage
}

func IncreaseDownloadCounter(version Version) {
	for _, element := range storage.Files {
		if element.Version.Equals(version) {
			element.DownloadCounter++
		}
	}
	writeToFile()
}

func writeToFile() {
	file, _ := json.MarshalIndent(storage, "", " ")
	f, err := os.OpenFile("./digitalFiles/files.json", os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	io.Copy(f, bytes.NewBuffer(file))
}

func enforceUploadLimit() {
	limit := configuration.ReadConfig().SavesCount
	sort.Slice(storage.Files, func(i, j int) bool {
		if storage.Files[i].Version.MajorVersion == storage.Files[j].Version.MajorVersion {
			return storage.Files[i].Version.MinorVersion < storage.Files[j].Version.MinorVersion
		}
		return storage.Files[i].Version.MajorVersion < storage.Files[j].Version.MajorVersion
	})
	for limit < len(storage.Files) {
		_, storage.Files = storage.Files[0], storage.Files[1:]
	}
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
