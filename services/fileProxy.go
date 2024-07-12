package services

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func ProxyFile(reqUrl string) ([]byte, error) {
    resp, err := http.Get(reqUrl)
    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("page %s returned status %d: %s", reqUrl, resp.StatusCode, resp.Status)
    }
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return ioutil.ReadAll(resp.Body)
}
