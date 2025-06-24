package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "time"
)

type IDWithAccess struct {
    ID           string
    LastAccessed time.Time
    DirPath      string
}

func RunSimpleCleanup() {
    fmt.Println("🧹 Iniciando cleanup automático (50% dos mais antigos)...")

    dirs, err := os.ReadDir(cfg.DownloadDir)
    if err != nil {
        fmt.Printf("❌ Erro ao ler diretório: %v\n", err)
        return
    }
    
    var idsWithAccess []IDWithAccess
    
    // Coleta todos os IDs com seus últimos acessos
    for _, dir := range dirs {
        if !dir.IsDir() {
            continue
        }
        
        id := dir.Name()
        lastAccess := getLastAccess(id)
        
        idsWithAccess = append(idsWithAccess, IDWithAccess{
            ID:           id,
            LastAccessed: lastAccess,
            DirPath:      filepath.Join(cfg.DownloadDir, id),
        })
    }
    
    totalCount := len(idsWithAccess)
    if totalCount <= 1 {
        fmt.Println("⚠️  Menos de 2 downloads, pulando cleanup")
        return
    }
    
    // Ordena por último acesso (mais antigo primeiro)
    sort.Slice(idsWithAccess, func(i, j int) bool {
        return idsWithAccess[i].LastAccessed.Before(idsWithAccess[j].LastAccessed)
    })
    
    // Calcula quantos remover (50%)
    toRemoveCount := totalCount / 2
    if toRemoveCount == 0 {
        toRemoveCount = 1 // Remove pelo menos 1 se tiver mais de 1
    }
    
    // Remove os mais antigos
    removedCount := 0
    for i := 0; i < toRemoveCount && i < len(idsWithAccess); i++ {
        idToRemove := idsWithAccess[i]
        
        if err := os.RemoveAll(idToRemove.DirPath); err == nil {
            fmt.Printf("🗑️  Removido: %s (último acesso: %s)\n", 
                idToRemove.ID, 
                formatTimeAgo(idToRemove.LastAccessed))
            removedCount++
        } else {
            fmt.Printf("❌ Erro ao remover %s: %v\n", idToRemove.ID, err)
        }
    }
    
    fmt.Printf("✅ Cleanup concluído: %d/%d downloads removidos\n", removedCount, totalCount)
}

// Formata tempo para exibição amigável
func formatTimeAgo(t time.Time) string {
    if t.IsZero() {
        return "nunca"
    }
    
    duration := time.Since(t)
    hours := int(duration.Hours())
    
    if hours < 1 {
        return "menos de 1h"
    } else if hours < 24 {
        return fmt.Sprintf("%dh atrás", hours)
    } else {
        days := hours / 24
        return fmt.Sprintf("%dd atrás", days)
    }
}