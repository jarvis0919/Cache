package information

import "github.com/jinzhu/gorm"

type Information struct {
	gorm.Model
	K string
	V string
}
