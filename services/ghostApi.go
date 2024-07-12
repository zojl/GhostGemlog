package services

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "ghost2gemini/models"
    "ghost2gemini/models/config"
    "github.com/LukeEmmet/html2gemini"
    "golang.org/x/net/html"
    "io"
    "net/http"
    "strings"
    "time"
)

func GetGhostPost(slug string) (string, error) {
    rawPost, err := getFromApi("posts/slug/" + slug)
    if err != nil {
        return "", err
    }

    var posts models.GhostPosts
    err = json.Unmarshal(rawPost, &posts)
    if err != nil {
        return "", err
    }

    return parseGhostItem(posts.Posts)
}

func GetGhostPage(slug string) (string, error) {
    rawPage, err := getFromApi("pages/slug/" + slug)
    if err != nil {
        return "", err
    }

    var pages models.GhostPages
    err = json.Unmarshal(rawPage, &pages)
    if err != nil {
        return "", err
    }

    return parseGhostItem(pages.Pages)
}

func GetGhostPostList() (*models.Feed, error) {
    rawPosts, err := getFromApi("posts")
    if err != nil {
        return nil, err
    }

    var posts models.GhostPosts
    err = json.Unmarshal(rawPosts, &posts)
    if err != nil {
        return nil, err
    }

    rawSettings, err := getFromApi("settings/")
    if err != nil {
        return nil, err
    }

    var settings models.GhostSettings
    err = json.Unmarshal(rawSettings, &settings)
    if err != nil {
        return nil, err
    }

    feed := models.Feed{
        Channel: models.FeedChannel{
            Title:       settings.Content.Title,
            Link:        settings.Content.Url,
            Description: settings.Content.Description,
            Items:       make([]models.FeedItem, 0, len(posts.Posts)),
        },
    }

    for _, post := range posts.Posts {
        parsedTime, _ := time.Parse(time.RFC3339, post.PubDate)
        pubDate := parsedTime.UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")

        item := models.FeedItem{
            Title:       post.Title,
            Link:        post.Slug,
            Description: post.Exceprt,
            PubDate:     pubDate,
            Guid:        post.Id,
        }

        feed.Channel.Items = append(feed.Channel.Items, item)
    }

    return &feed, nil
}

func parseGhostItem(posts []models.GhostItem) (string, error) {
    if len(posts) == 0 {
        return "", errors.New("not found")
    }
    post := posts[0]

    if post.Visibility != "public" {
        return "", errors.New("forbidden")
    }

    var excerpt string
    if len(post.CustomExceprt) > 0 {
        excerpt = "\n\n### " + post.CustomExceprt
    }

    var imageCaption string
    if len(post.FeatureImageCaption) > 0 {
        imageCaption, _ = htmlToGemini(post.FeatureImageCaption)
        imageCaption = "\n\n" + imageCaption
    }

    var image string
    if len(post.FeatureImage) > 0 {
        image = fmt.Sprintf(
            "\n\n=> %s 🖼️ %s",
            strings.ReplaceAll(post.FeatureImage, config.BlogHost, ""),
            post.Title,
        )
    }

    postText, _ := htmlToGemini(post.Html)

    return "# " + post.Title + excerpt + image + imageCaption + "\n\n" + postText, nil
}

func getFromApi(method string) ([]byte, error) {
    url := config.BlogHost + "/ghost/api/content/" + method + "?page=2&key=" + config.BlogKey
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return respBody, nil
}

func htmlToGemini(html string) (string, error) {
    html, _ = replaceTags(html, "s", config.TextStrikeTag)
    html, _ = replaceTags(html, "strike", config.TextStrikeTag)
    html = strings.ReplaceAll(html, "</p>", "</p><br><br>")
    html = strings.ReplaceAll(html, "</div>", "</div><br><br>")

    options := html2gemini.NewOptions()
    options.CitationStart = 4
    options.LinkEmitFrequency = 1
    options.CitationMarkers = true
    options.NumberedLinks = true
    options.EmitImagesAsLinks = true
    options.ImageMarkerPrefix = "🖼️"

    options.PrettyTables = false
    options.PrettyTablesOptions.HeaderLine = true
    options.PrettyTablesOptions.RowLine = true
    options.PrettyTablesOptions.CenterSeparator = " "
    options.PrettyTablesOptions.ColumnSeparator = " "
    options.PrettyTablesOptions.RowSeparator = " "

    ctx := html2gemini.NewTraverseContext(*options)

    text, err := html2gemini.FromString(html, *ctx)
    if err != nil {
        return "", err
    }

    text = strings.ReplaceAll(text, "=> "+config.BlogHost, "=> ")

    return text, nil
}

func replaceTags(input string, from string, to string) (string, error) {
    doc, err := html.Parse(strings.NewReader(input))
    if err != nil {
        return "", err
    }

    var traverse func(*html.Node)
    traverse = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == from {
            n.Data = to
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            traverse(c)
        }
    }
    traverse(doc)

    var buf bytes.Buffer
    if err := html.Render(&buf, doc); err != nil {
        return "", err
    }
    return buf.String(), nil
}
