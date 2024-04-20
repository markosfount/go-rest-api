package data

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrRecordExists = errors.New("record already exists")