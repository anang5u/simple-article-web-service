package helper

import "github.com/gosimple/slug"

// Slugify
func Slugify(str string) string {
	return slug.Make(str)
}
