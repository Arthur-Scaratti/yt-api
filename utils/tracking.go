package utils

import (
    "encoding/json"
    "os"
    "path/filepath"
    "time"
)

type AccessInfo struct {
    LastAccessed time.Time `json:"last_accessed"`
}

func UpdateLastAccess(id string) {
    accessPath := filepath.Join(cfg.DownloadDir, id, ".access")
    access := AccessInfo{
        LastAccessed: time.Now(),
    }
    
    data, _ := json.Marshal(access)
    os.WriteFile(accessPath, data, 0644)
}

func getLastAccess(id string) time.Time {
    accessPath := filepath.Join(cfg.DownloadDir, id, ".access")
    data, err := os.ReadFile(accessPath)
    if err != nil {
        return time.Time{}
    }
    
    var access AccessInfo
    if err := json.Unmarshal(data, &access); err != nil {
        return time.Time{}
    }
    return access.LastAccessed
}