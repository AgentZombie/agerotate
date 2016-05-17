package bucket

import (
	"time"

	"agerotate"
)

func Cleanup(sortedRanges []agerotate.Range, objects agerotate.Objects) error {
	now := time.Now()
	buckets := makeBuckets(sortedRanges)
	overflow, err := readObjects(objects, buckets, now)
	if err != nil {
		return err
	}

	if err = cleanupBuckets(buckets); err != nil {
		return err
	}

	if err = cleanupOverflow(overflow); err != nil {
		return err
	}

	return nil
}

func makeBuckets(sortedRanges []agerotate.Range) []*bucket {
	buckets := make([]*bucket, len(sortedRanges))
	for i := range sortedRanges {
		buckets[i] = newBucket(sortedRanges[i])
	}
	return buckets
}

func readObjects(objects agerotate.Objects, buckets []*bucket, now time.Time) ([]agerotate.Object, error) {
	overflow := []agerotate.Object{}
	oList, err := objects.List()
	if err != nil {
		return nil, err
	}

	for _, o := range oList {
		found := false
		for _, b := range buckets {
			if o.Age(now) < b.Age() {
				b.Add(o)
				found = true
			}
		}
		if !found {
			overflow = append(overflow, o)
		}
	}
	return overflow, nil
}

func cleanupBuckets(buckets []*bucket) error {
	for _, b := range buckets {
		if err = b.Cleanup(now); err != nil {
			return err
		}
	}
	return nil
}

func cleanupOverflow(overflow []Object) error {
	for _, o := range overflow {
		if err = o.Delete(); err != nil {
			return err
		}
	}
	return nil
}