package config

import (
    "os"

    "github.com/joho/godotenv"
)

var (
    BlogHost        string
    BlogKey         string
    GeminiHost      string
    GeminiPort      string
    TextDescription string
    TextFooter      string
    TextRead        string
    TextStrikeTag   string
)

func LoadConfig() {
    godotenv.Load()

    BlogHost = os.Getenv("BLOG_HOST")
    BlogKey = os.Getenv("BLOG_KEY")
    GeminiHost = os.Getenv("GEMINI_HOST")
    GeminiPort = os.Getenv("GEMINI_PORT")
    TextDescription = os.Getenv("TEXT_DESCRIPTION")
    TextFooter = os.Getenv("TEXT_FOOTER")
    TextRead = os.Getenv("TEXT_READ")
    TextStrikeTag = os.Getenv("TEXT_STRIKE_TAG")
}
