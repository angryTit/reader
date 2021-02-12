package service

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/juliangruber/go-intersect"

	"github.com/angryTit/reader/types"
)

const (
	newLine = '\n'
)

func parse(source []string) map[string][]string {
	result := make(map[string][]string)
	for _, each := range source {
		arr := strings.Split(each, ",")
		if len(arr) != 3 {
			log.Printf("[WARN] wrong record format [%v]", each)
			continue
		}
		userId := strings.TrimSpace(arr[0])
		ip := strings.TrimSpace(arr[1])
		target, ok := result[userId]
		if !ok {
			slice := []string{ip}
			result[userId] = slice
			continue
		}

		target = append(target, ip)
		result[userId] = target
	}
	return result
}

func readFrom(file io.ReadSeeker, startPosition int64) ([]string, *int64, error) {
	_, err := file.Seek(startPosition, 0)
	if err != nil {
		log.Printf("[ERROR] fail to seek position : %v", err)
		return nil, nil, err
	}
	result := make([]string, 0)
	reader := bufio.NewReader(file)
	currentPosition := startPosition
	for {
		data, err := reader.ReadBytes(newLine)
		if err != nil && err != io.EOF {
			log.Printf("[ERROR] fail to read from file : %v", err)
			return nil, nil, err
		}

		currentPosition += int64(len(data))
		if len(data) > 0 {
			result = append(result, string(data))
		}

		if err == io.EOF {
			return result, &currentPosition, nil
		}
	}
}

func FillStorage(file io.ReadSeeker, startPosition int64, storage *types.Storage) (*int64, error) {
	arr, currentPosition, err := readFrom(file, startPosition)
	if err != nil {
		return nil, err
	}

	tmp := parse(arr)

	for eachUserId, eachIps := range tmp {
		storage.Set(eachUserId, eachIps)
	}

	return currentPosition, nil
}

func IsSame(sourceUserId, targetUserId string, storage *types.Storage) bool {
	if sourceUserId == targetUserId {
		return true
	}

	sourceIpsSlice := storage.Get(sourceUserId)
	targetIpsSlice := storage.Get(targetUserId)
	if sourceIpsSlice == nil || targetIpsSlice == nil {
		return false
	}

	sIps := *sourceIpsSlice.GetSlice()
	tIps := *targetIpsSlice.GetSlice()
	resP := intersect.Simple(sIps, tIps)
	result := resP.([]interface{})
	return len(result) >= 2
}

//call in separate thread only
func UpdateStorageInBackground(filePath string, startPosition int64, storage *types.Storage, duration time.Duration) {
	for {
		time.Sleep(duration)
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("[ERROR] fail to open file [%v] : %v", filePath, err)
		}

		position, err := FillStorage(file, startPosition, storage)
		if err != nil {
			os.Exit(1)
		}
		err = file.Close()
		if err != nil {
			log.Fatalf("[ERROR] fail to close file [%v] : %v", filePath, err)
		}

		startPosition = *position
		//log.Printf("[INFO] successfully update storage at [%v]", time.Now())
	}
}
