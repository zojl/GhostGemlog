package models

type GhostItem struct {
    Title               string `json:"title"`
    Html                string `json:"html"`
    FeatureImage        string `json:"feature_image"`
    Visibility          string `json:"visibility"`
    Exceprt             string `json:"excerpt"`
    CustomExceprt       string `json:"custom_excerpt"`
    FeatureImageCaption string `json:"feature_image_caption"`
    Id                  string `json:"id"`
    PubDate             string `json:"published_at"`
    Slug                string `json:"slug"`
}

type GhostPosts struct {
    Posts []GhostItem `json:"posts"`
}

type GhostPages struct {
    Pages []GhostItem `json:"pages"`
}

type GhostSettings struct {
    Content GhostSettingsContent `json:"settings"`
}

type GhostSettingsContent struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Url         string `json:"url"`
}
