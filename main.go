package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/angryTit/reader/service"
	"github.com/angryTit/reader/types"
)

const (
	filePath = "log.txt"
	duration = 5 * time.Second
)

var storage *types.Storage

func init() {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("[ERROR] fail to open file [%v] : %v", filePath, err)
	}
	storage = types.NewStorage()
	position, err := service.FillStorage(f, 0, storage)
	if err != nil {
		os.Exit(1)
	}

	go service.UpdateStorageInBackground(filePath, *position, storage, duration)
	log.Println("[INFO] ----- ready -----")
}

func main() {
	http.HandleFunc("/", server)
	http.ListenAndServe(":8080", nil)
}

func server(w http.ResponseWriter, r *http.Request) {
	path := strings.Replace(r.URL.Path, "/", "", 1)
	arr := strings.Split(path, "/")
	fmt.Fprintf(w, `{"dupes":%v}`, service.IsSame(arr[0], arr[1], storage))
}
