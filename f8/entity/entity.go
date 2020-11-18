package entity

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/phanirithvij/fate/f8/xfs"
	"gorm.io/gorm"
)

// BaseEntity a base entity model
type BaseEntity struct {
	// ID entity id autogenerated
	ID string `gorm:"primaryKey;not null"`
	// https://gorm.io/docs/models.html#Field-Level-Permission
	// would be great if this was a map but unfortunately sql no maps
	Buckets []*xfs.Bucket `gorm:"polymorphic:Entity"`
	// Buckets []*xfs.Bucket `gorm:"polymorphic:Entity;<-:false"`
	// EntityType to pass it down
	EntityType string `gorm:"-"`
}

// From this vararg approach
// https://github.com/faiface/gui/commit/ee3366ded862f02a1a5ee4ea856a06e46bb889ee#diff-f9c3d4c5cce2eabfcf19a1c38214e739a6619bdd22b15373e34be2e3e4589247R89

// Option is a functional option to the entity constructor New.
type Option func(*options)
type options struct {
	id                string
	numBuckets        int
	defaultBucketName string
	bucketNames       []string
	tableName         string
}

// ID option sets the ID of the entity.
//
// Don't set it for an auto uuid
//
// If it's username better set it or you'll get uuids as usernames
func ID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

// BucketCount option sets the num of buckets initially
func BucketCount(numBuckets int) Option {
	return func(o *options) {
		o.numBuckets = numBuckets
	}
}

// TableName option sets the name of the entity
// REQUIRED
func TableName(tableName string) Option {
	return func(o *options) {
		o.tableName = tableName
	}
}

// BucketName option sets the default bucket Name for the entity
//
// Must specify this or BucketNames if using BucketCount() as buckets will be created
// as name-0, name-1, name-2, ..., name-{count-1}
func BucketName(bucketName string) Option {
	return func(o *options) {
		o.defaultBucketName = bucketName
	}
}

// BucketNames option sets the bucketNames for the entity
//
// Must specify these or BucketName if using BucketCount() as buckets will be created
// as names[0], names[1], names[2], ..., names[count-1]
func BucketNames(bucketNames []string) Option {
	return func(o *options) {
		o.bucketNames = bucketNames
	}
}

// NewBase a new base
func NewBase(opts ...Option) (*BaseEntity, error) {
	o := options{
		defaultBucketName: "default",
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.id == "" {
		o.id = uuid.New().String()
	}
	if o.tableName == "" {
		return nil, errors.New("Must specify the table name")
	}
	// whether we should use bucketNames[]
	usebNames := false
	if len(o.bucketNames) > 0 {
		if o.defaultBucketName != "default" {
			// both bucketNames and a name was specifed for an incremental bucket name
			return nil, errors.New("Use only one of BucketNames, BucketName")
		}
		if len(o.bucketNames) != o.numBuckets {
			return nil, errors.New("Number of bucket names must match the numBuckets")
		}
		usebNames = true
	}

	etype := o.tableName

	ent := &BaseEntity{
		ID:         o.id,
		Buckets:    []*xfs.Bucket{},
		EntityType: etype,
	}
	if o.numBuckets == 0 {
		// number of buckets was not specified
		// create a default bucket
		o.numBuckets = 1
	}

	// create initial buckets
	for i := 0; i < o.numBuckets; i++ {
		bID := o.defaultBucketName
		if usebNames {
			bID = o.bucketNames[i]
		} else {
			if i != 0 {
				// for i = 0 it will be default but not default-0
				bID = o.defaultBucketName + "-" + strconv.Itoa(i)
			}
		}
		ent.CreateBucket(bID)
	}
	return ent, nil
}

// CreateBucket creates a new bucket for the entity
// and appends it to the entity owned bucket list
func (e *BaseEntity) CreateBucket(bID string) (buck *xfs.Bucket) {
	buck = xfs.NewBucket(e.ID, e.EntityType, bID)
	e.Buckets = append(e.Buckets, buck)
	return buck
}

// AutoMigrate auto migrations required for the database
//
// Note: BaseEntity will not auto migrate because it's the parent's responsibility
func AutoMigrate(db *gorm.DB) error {
	return xfs.AutoMigrate(db)
}
