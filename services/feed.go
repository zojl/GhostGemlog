package services

import (
    "encoding/json"
    "fmt"
    "ghost2gemini/models"
    "ghost2gemini/models/config"
    "io/ioutil"
    "net/url"
    "os"
)

const jsonFilePath = "data/blog/posts.json"

func SaveFeed() error {
    rssFeed, err := GetGhostPostList()
    if err != nil {
        return err
    }

    existingItemsMap := make(map[string]models.FeedItem)
    existingFeed, _ := readRssFromJson()
    for _, item := range existingFeed.Channel.Items {
        existingItemsMap[item.Guid] = item
    }

    for _, item := range rssFeed.Channel.Items {
        existingItemsMap[item.Guid] = item
    }

    newItems := make([]models.FeedItem, 0, len(existingItemsMap))
    for _, item := range existingItemsMap {
        newItems = append(newItems, item)
    }
    rssFeed.Channel.Items = newItems

    rssFeed.Channel.SortItemsByPubDate()

    file, err := os.Create(jsonFilePath)
    if err != nil {
        return err
    }
    defer file.Close()

    jsonData, err := json.MarshalIndent(rssFeed, "", "  ")
    if err != nil {
        return err
    }

    _, err = file.Write(jsonData)
    if err != nil {
        return err
    }

    return nil
}

func GetFeedAsGemini() (string, error) {
    rssFeed, err := readRssFromJson()
    if err != nil {
        return "", err
    }

    geminiContent := fmt.Sprintf(
        "# %s\n\n## %s\n\n%s\n\n",
        rssFeed.Channel.Title,
        rssFeed.Channel.Description,
        config.TextDescription,
    )
    for _, item := range rssFeed.Channel.Items {
        postLink, _ := url.Parse(item.Link)
        geminiContent += fmt.Sprintf(
            "### %s\n\n%s\n\n=> %s %s | %s\n\n",
            item.Title,
            item.Description,
            postLink.Path,
            item.PubDate,
            config.TextRead,
        )
    }

    return geminiContent, nil
}

func readRssFromJson() (*models.Feed, error) {
    file, err := os.Open(jsonFilePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    byteValue, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, err
    }
    var rssFeed models.Feed
    err = json.Unmarshal(byteValue, &rssFeed)
    return &rssFeed, err
}
