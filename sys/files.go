package sys

import (
    "os"
)

/*
 * Do a stat on file_name. Return true if it exists, false otherwise
 */
func FileExists (file_name string) (exists bool) {
    _, err := os.Stat(file_name)
    exists = err == nil
    return
}
