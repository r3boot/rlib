package sys

import (
    "errors"
    "os"
    "strings"
)

/*
 * Do a stat on file_name. Return true if it exists, false otherwise
 */
func FileExists (file_name string) (exists bool) {
    _, err := os.Stat(file_name)
    exists = err == nil
    return
}

/*
 * Find the absolute path of cmd by traversing PATH. Return an error if the
 * absolute path is not found, else, path contains the absolute path to cmd
 */
func BinaryPrefix (cmd string) (path string, err error) {
    raw_PATH := os.Getenv("PATH")
    paths := strings.Split(raw_PATH, ":")

    for _, p := range paths {
        path = p + "/" + cmd
        if FileExists(path) {
            return
        }
    }

    path = ""
    err = errors.New("Cannot find " + cmd)
    return
}

/*
 * Look in various directories to find the prefix of a configuration file.
 * Returns the absolute path to the config file if it is found and an error
 * if not
 */
func ConfigPrefix (cfg_file string) (path string, err error) {
    for _, p := range ETC_PATHS {
        path = p + "/" + cfg_file
        if FileExists(path) {
            return
        }
    }

    path = ""
    err = errors.New("Cannot find " + cfg_file)
    return
}
