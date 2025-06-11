package utils

import (
    "io/fs"
    "os"
    "path/filepath"
    "sort"
)

func CleanupDownloads(maxSizeMB int64) {
    dir := "downloads"
    var totalSize int64
    files := []fs.FileInfo{}

    _ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
        if path == dir || err != nil || !info.IsDir() {
            return nil
        }
        files = append(files, info)
        return nil
    })

    sort.Slice(files, func(i, j int) bool {
        return files[i].ModTime().Before(files[j].ModTime())
    })

    for _, f := range files {
        p := filepath.Join(dir, f.Name())
        size := getDirSize(p)
        totalSize += size / (1024 * 1024)

        if totalSize > maxSizeMB {
            os.RemoveAll(p)
        }
    }
}

func getDirSize(path string) int64 {
    var size int64
    _ = filepath.Walk(path, func(_ string, info fs.FileInfo, _ error) error {
        if !info.IsDir() {
            size += info.Size()
        }
        return nil
    })
    return size
}
