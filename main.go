package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

type NexusRequest struct {
	Items []Items `json:"items"`
}
type Checksum struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
	Sha256 string `json:"sha256"`
	Md5    string `json:"md5"`
}
type Maven2 struct {
	Extension  string `json:"extension"`
	GroupID    string `json:"groupId"`
	ArtifactID string `json:"artifactId"`
	Version    string `json:"version"`
}
type Items struct {
	DownloadURL    string    `json:"downloadUrl"`
	Path           string    `json:"path"`
	ID             string    `json:"id"`
	Repository     string    `json:"repository"`
	Format         string    `json:"format"`
	Checksum       Checksum  `json:"checksum"`
	ContentType    string    `json:"contentType"`
	LastModified   time.Time `json:"lastModified"`
	LastDownloaded time.Time `json:"lastDownloaded"`
	Uploader       string    `json:"uploader"`
	UploaderIP     string    `json:"uploaderIp"`
	FileSize       int       `json:"fileSize"`
	BlobCreated    time.Time `json:"blobCreated"`
	Maven2         Maven2    `json:"maven2"`
}

func main() {

	url := "http://172.30.86.46:18081/service/rest/v1/assets?repository=maven-public"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	//body, err := ioutil.ReadAll(res.Body)
	var t NexusRequest
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range t.Items {
		fmt.Println(v)
	}

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS nexus (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT,
			email TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	//if err != nil {

	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(string(body))
}
