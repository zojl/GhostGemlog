# GhostGemlog — a proxy for your Ghost blog to the Geminispace
The project is based on [ninedraft/gemax](https://github.com/ninedraft/gemax) and [LukeEmmet/html2gmi](https://github.com/LukeEmmet/html2gmi) libraries

### Setting up
Create an `.env` file or use docker with the following environment variables:
- BLOG_HOST — a full url of your blog with https
- BLOG_KEY — content API key from the integrations section of your blog admin panel 
- Gemini_HOST — host to listen on for the Gemini server, can be 0.0.0.0
- Gemini_PORT — port of the Gemini server
- TEXT_DESCRIPTION — custom text for the main page of your capsule
- TEXT_FOOTER — footer text for each page of your capsule
- TEXT_READ — text for post links after their publication date
- TEXT_STRIKE_TAG — name of the ironic tag that will be displayed around the strikethrough text

## Example docker-compose.yml
The project is cloned to the src directory.
``` yaml
services:
  app:
    build: src
    env_file: .env
    volumes:
     - ./volume/app/data/blog:/app/data/blog
     - ./volume/app/data/cert:/app/data/cert
     - ./volume/app/logs:/app/logs
    ports:
     - 1965:1965
```

## Dont forget about the certificate
Simply run `openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=<YOUR_DOMAIN>'` to generate the certificate with key and put them to your data/cert directory.