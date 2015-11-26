package main
import (
	"flag"
	"io/ioutil"
	"bufio"
	"bytes"
	"io"
	"strings"
	"os"
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"time"
	"encoding/hex"
)

var (
	source_folder = flag.String("path", ".", "file path to work on")

	config_file = flag.String("config", "config.txt", "config file path")
)

func main() {
	flag.Parse()

	fileModInfoMap, err := readConfig()
	if err != nil {
		panic(err)
	}

	newFileModInfoMap := make(map[string]modInfo)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		h := sha1.New()

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			return err
		}

		_, err = io.Copy(h, file)
		if err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return err
		}
		err = file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
			return err
		}

		sum := hex.EncodeToString(h.Sum(nil))

		old, has := fileModInfoMap[path]
		if !has || !strings.EqualFold(old.sum, sum) {
			// not same

			fmt.Printf("file change: %s\n", path)
			newFileModInfoMap[path] = modInfo{sum: sum, modTime:info.ModTime().Format(time.RFC3339Nano)}
		} else {
			newFileModInfoMap[path] = old
			// same, set file mod time
			modTime, err := time.Parse(time.RFC3339Nano, old.modTime)
			if err != nil {
				fmt.Printf("Error parsing time in config: %v\n", err)
				return err
			}

			err = os.Chtimes(path, modTime, modTime)
			if err != nil {
				fmt.Printf("Error changing file mod time: %v\n", err)
				return err
			}

		}

		return nil
	}


	err = filepath.Walk(*source_folder, walkFn)

	if err != nil {
		fmt.Printf("Walk folder error: %v\n", err)
		panic(err)
	}

	err = writeConfig(newFileModInfoMap)
	if err != nil{
		fmt.Printf("Error writing config file: %v\n", err)
	}
}

func writeConfig(info map[string]modInfo) error{
	// write new config
	file, err := os.Create(*config_file)
	if err != nil {
		return err
	}
	defer file.Close()
	for name, value := range info {
		fmt.Fprintf(file, "%v %v %v\r\n", value.sum, value.modTime, name)
	}

	return nil
}

func readConfig() (map[string]modInfo, error) {
	fileModInfoMap := make(map[string]modInfo)
	data, err := ioutil.ReadFile(*config_file)
	if err != nil {
		if os.IsNotExist(err) {
			return fileModInfoMap, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return fileModInfoMap, nil
	}

	bufReader := bufio.NewReader(bytes.NewReader(data))

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		parts := strings.SplitN(string(line), " ", 3)

		info := modInfo{sum: parts[0], modTime: parts[1]}
		fileModInfoMap[parts[2]] = info // fileName map to modInfo
	}
	return fileModInfoMap, nil
}

type modInfo struct {
	sum     string
	modTime string
}
