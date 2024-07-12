package models

import (
    "sort"
    "time"
)

type FeedItem struct {
    Title       string
    Link        string
    Description string
    PubDate     string
    Guid        string
}

type FeedChannel struct {
    Title       string
    Link        string
    Description string
    Items       []FeedItem
}

type Feed struct {
    Channel FeedChannel
}

type ByPubDate []FeedItem

func (a ByPubDate) Len() int {
    return len(a)
}

func (a ByPubDate) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a ByPubDate) Less(i, j int) bool {
    layout := "Mon, 02 Jan 2006 15:04:05 GMT"
    dateI, err := time.Parse(layout, a[i].PubDate)
    if err != nil {
        return false
    }
    dateJ, err := time.Parse(layout, a[j].PubDate)
    if err != nil {
        return true
    }
    return dateI.After(dateJ)
}

func (channel *FeedChannel) SortItemsByPubDate() {
    sort.Sort(ByPubDate(channel.Items))
}
