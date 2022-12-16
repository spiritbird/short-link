package short_link

import "errors"

var ErrRedisConnect = errors.New("redis connect has gone")
var ErrNotFound = errors.New("not found match long url")
