// Package storage declares the app's object storage buckets.
package storage

import "encore.dev/storage/objects"

// PrivateFiles holds generated documents that must never be reachable without going through the
// API — monthly billing reports today, any other internal file tomorrow. It is deliberately not
// tied to a single kind of document: keep files apart by giving each kind its own key prefix
// (see monthlyReportKey in the billing service) rather than by adding another bucket.
var PrivateFiles = objects.NewBucket("private-files", objects.BucketConfig{
	Versioned: false,
})
