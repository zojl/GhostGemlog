package handlers

import (
    "context"
    "ghost2gemini/models/config"
    "ghost2gemini/services"
    "github.com/ninedraft/gemax/gemax"
    "github.com/ninedraft/gemax/gemax/status"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

func BlogHandler(_ context.Context, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
    err := logRequest(req)
    if err != nil {
        log.Println(err)
    }

    urlPath := req.URL().Path
    urlPath = strings.TrimPrefix(urlPath, "/")
    urlPath = strings.TrimSuffix(urlPath, "/")
    urlPath = strings.ReplaceAll(urlPath, "/../", "/")
    if urlPath == "" {
        handleIndex(rw)
        return
    }

    if urlPath == "updateBlogPosts" {
        handleUpdateFeed(rw)
        return
    }

    fileContent, isMedia, err := services.GetFile(urlPath)
    if err != nil {
        rw.WriteStatus(status.ServerUnavailable, "Can't display contents of: "+req.URL().Path)
        log.Println(err)
        return
    }

    if fileContent != nil {
        handleLocalFile(isMedia, fileContent, rw)
        return
    }

    reqUrl := config.BlogHost + req.URL().Path
    if isImageURL(urlPath) {
        handleFileProxy(reqUrl, rw, req)
        return
    }

    handlePostProxy(urlPath, rw, req)
}

func handleLocalFile(isMedia bool, fileContent []byte, rw gemax.ResponseWriter) {
    if isMedia {
        contentType := http.DetectContentType(fileContent)
        rw.WriteStatus(status.Success, contentType)
    }
    rw.Write(fileContent)
    if !isMedia {
        rw.Write([]byte(config.TextFooter))
    }
}

func logRequest(request gemax.IncomingRequest) error {
    logFilename := "logs/access.log"
    file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    logger := log.New(file, "", log.LstdFlags)

    currentTime := time.Now().Format("2006-01-02 15:04:05")
    logger.Printf(
        "%s | IP %s | %s \n",
        currentTime,
        request.RemoteAddr(),
        request.URL(),
    )

    return nil
}

func handlePostProxy(reqUrl string, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
    postText, err := services.GetGhostPost(reqUrl)
    if err == nil {
        io.WriteString(rw, postText)
        return
    }

    pageText, err := services.GetGhostPage(reqUrl)
    if err == nil {
        io.WriteString(rw, pageText+config.TextFooter)
        return
    }

    rw.WriteStatus(status.NotFound, "Not found: "+req.URL().Path)
}

func handleFileProxy(reqUrl string, rw gemax.ResponseWriter, req gemax.IncomingRequest) {
    file, err := services.ProxyFile(reqUrl)
    if err != nil {
        log.Println(err)
        rw.WriteStatus(status.NotFound, "Not found: "+req.URL().Path)
    }
    contentType := http.DetectContentType(file)
    rw.WriteStatus(status.Success, contentType)
    _, err = rw.Write(file)
    if err != nil {
        log.Println(err)
    }
}

func handleUpdateFeed(rw gemax.ResponseWriter) {
    err := services.SaveFeed()
    if err != nil {
        log.Println(err)
    }
    io.WriteString(rw, "Updated")
}

func handleIndex(rw gemax.ResponseWriter) {
    geminiContent, _ := services.GetFeedAsGemini()
    io.WriteString(rw, geminiContent+config.TextFooter)
}

func isImageURL(urlPath string) bool {
    suffixes := []string{".jpg", ".jpeg", ".webp", ".png", ".gif"}
    lowerUrlPath := strings.ToLower(urlPath)
    for _, suffix := range suffixes {
        if strings.HasSuffix(lowerUrlPath, suffix) {
            return true
        }
    }

    return false
}
