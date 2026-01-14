package repositories

import "errors"

var ErrReadOnly = errors.New("write not supported in JSON mode")