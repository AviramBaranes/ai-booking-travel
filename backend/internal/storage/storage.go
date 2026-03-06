package storage

import "encore.dev/storage/objects"

var ContentAttachments = objects.NewBucket("content-attachments", objects.BucketConfig{
	Versioned: false,
	Public:    true,
})

type ContentAttachment string

const (
	HertzLocationUpload ContentAttachment = "hertz-location-upload"
)
