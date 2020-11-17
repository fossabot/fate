package xfs

import (
	"os"
	"time"

	"gorm.io/gorm"
)

// Bucket is equivalient to a filesystem with a name
type Bucket struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	XFileSystem
	/*
		Bucket ID will not be unique as evey entity has a `default` bucket
		Entity ID and bucket ID together will be unique
	*/
	// https://stackoverflow.com/a/63409572/8608146
	// https://gorm.io/docs/composite_primary_key.html
	ID       string `gorm:"uniqueIndex:compositeindex;primaryKey"`
	EntityID string `gorm:"uniqueIndex:compositeindex"`
}

// NewBucket returns a new bucket, if id is empty ID is default
func NewBucket(id string) *Bucket {
	if id == "" {
		id = "default"
	}
	// Provision a bucket with an empty file system
	return &Bucket{
		ID:          id,
		XFileSystem: XFileSystem{FileDirs: []FileDir{}},
	}
}

// Exists checks if the bucket already exists
func (b *Bucket) Exists(db *gorm.DB) bool {
	// cache bucket results after querying once
	return false
}

func (b *Bucket) String() string {
	return b.ID
}

// XFileSystem a simple filesystem implementation
type XFileSystem struct {
	FileDirs []FileDir `gorm:"foreignKey:BucketID"`
}

// copied from os.FileInfo interface{}

// FileDir a file or directory
type FileDir struct {
	gorm.Model
	Name string `gorm:"primarykey"` // base name of the file
	// TODO: Once a file or dir is created it is our job to populate these fields
	Size     int64       // length in bytes for regular files; system-dependent for others
	Mode     os.FileMode // file mode bits
	ModTime  time.Time   // modification time
	IsDir    bool        // abbreviation for Mode.IsDir
	BucketID string      `gorm:"primarykey"`
	info     os.FileInfo `gorm:"-"`
}

// NewFile retuns a new file
func NewFile(name string) *FileDir {
	if name == "" {
		return nil
	}
	return &FileDir{
		Name:  name,
		IsDir: false,
	}
}

// NewDir returns a new directory
func NewDir(name string) *FileDir {
	if name == "" {
		return nil
	}
	return &FileDir{
		Name:  name,
		IsDir: true,
	}
}

// AutoMigrate for xfs
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Bucket{}, &FileDir{})
}
