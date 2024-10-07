// Filename: internal/data/comments.go
package data

import (
	"time"
)

// each name begins with uppercase so that they are exportable/public
type Comment struct {
	ID        int64     // unique value for each comment
	Content   string    // the comment data
	Author    string    // the person who wrote the comment
	CreatedAt time.Time // database timestamp
	Verson    int32     // incremented on each update
}
