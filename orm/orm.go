package orm

import "gorm.io/gorm"

type NexusRequest struct {
	Items             []Item `json:"items"`
	ContinuationToken string `json:"continuationToken"`
}

type Item struct {
	ID         string  `json:"id"`
	Repository string  `json:"repository"`
	Format     string  `json:"format"`
	Group      string  `json:"group"`
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	Assets     []Asset `json:"assets"`
}

type Asset struct {
	DownloadURL    string      `json:"downloadUrl"`
	Path           string      `json:"path"`
	ID             string      `json:"id"`
	Repository     string      `json:"repository"`
	Format         string      `json:"format"`
	Checksum       Checksum    `json:"checksum"`
	ContentType    string      `json:"contentType"`
	LastModified   string      `json:"lastModified"`
	LastDownloaded string      `json:"lastDownloaded"`
	Uploader       string      `json:"uploader"`
	UploaderIP     string      `json:"uploaderIp"`
	FileSize       int64       `json:"fileSize"`
	BlobCreated    string      `json:"blobCreated"`
	Maven2         Maven2Class `json:"maven2,omitempty"`
	Npm            Npm         `json:"npm,omitempty"`
}

type Checksum struct {
	Sha1   string `json:"sha1"`
	Sha512 string `json:"sha512"`
	Sha256 string `json:"sha256"`
	Md5    string `json:"md5"`
}

type Maven2Class struct {
	Extension  string  `json:"extension"`
	GroupID    string  `json:"groupId"`
	ArtifactID string  `json:"artifactId"`
	Version    string  `json:"version"`
	Classifier *string `json:"classifier,omitempty"`
}

type Npm struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type MavenRepository struct {
	gorm.Model
	DownloadURL    string `gorm:"index"`
	GroupID        string
	ArtifactID     string
	Version        string
	Path           string
	LocalFilePath  string
	Extension      string
	DownLoadStatus bool `gorm:"default:false"`
	UpLoadStatus   bool `gorm:"default:false"`
}

type NpmRepository struct {
	gorm.Model
	DownloadURL    string `gorm:"index"`
	Name           string
	Path           string
	LocalFilePath  string
	Version        string
	DownLoadStatus bool `gorm:"default:false"`
	UpLoadStatus   bool `gorm:"default:false"`
}
