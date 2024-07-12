package services

import (
    "io"
    "os"
)

func GetFile(urlPath string) ([]byte, bool, error) {
    postPath := "data/blog/contents/" + urlPath + "/post.gmi"
    postFile, err := os.Open(postPath)
    defer postFile.Close()
    if err == nil {
        content, err := io.ReadAll(postFile)
        if err != nil {
            return nil, false, err
        }
        return content, false, nil
    }

    mediaPath := "data/blog/contents/" + urlPath
    mediaFile, err := os.Open(mediaPath)
    defer mediaFile.Close()
    if err == nil {
        content, err := io.ReadAll(mediaFile)
        if err != nil {
            return content, true, nil
        }
        return content, true, nil
    }
    
    return nil, false, nil
}
