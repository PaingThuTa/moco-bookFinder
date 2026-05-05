# Architecture Codemap

<!-- AUTO-GENERATED -->

Book Finder is a Telegram bot that searches multiple book sources and returns download links as inline keyboard buttons. It uses long polling (no webhook) and runs as a single binary.

## High-Level Architecture

```
User (Telegram)
    |
    v
+------------------+
|   Telegram API   |
+------------------+
    |
    | long polling
    v
+------------------+     +------------------+     +------------------+
|   cmd/bot/main   |---->|  BotHandler      |---->|  SourceManager   |
|   (entry point)  |     |  (telego SDK)    |     |  (fallback chain)|
+------------------+     +------------------+     +------------------+
                                |                        |
                                v                        v
                         +-------------+     +------------------------+
                         |   Config    |     |  Scraper Implementations|
                         |  (env vars) |     |  - ZLibraryScraper     |
                         +-------------+     |  - OceanPDFScraper     |
                                             |  - LibGenScraper       |
                                             +------------------------+
```

## Data Flow

1. **User sends `/search <query>`** via Telegram
2. **Main loop** receives update, routes to `HandleSearch`
3. **BotHandler** parses args, sends "Searching..." message
4. **SourceManager.Search()** iterates scrapers in order (Z-Library -> Ocean of PDF -> LibGen)
5. **First scraper with results** returns immediately (fallback pattern)
6. **BotHandler** formats results as markdown with inline keyboard buttons
7. **Results stored in memory** with 10-minute TTL for callback lookup
8. **User taps "Download"** button -> callback resolved -> direct link sent

## Key Design Decisions

- **Long polling** over webhook: simpler deployment, no HTTPS/reverse proxy needed
- **Sequential fallback**: sources tried in priority order, first success wins
- **In-memory result cache**: chat-scoped, 10-minute TTL, avoids stale callbacks
- **User authorization**: whitelist-based via `ALLOWED_USER_IDS` env var
- **Shared HTTP client**: single client with 30s timeout injected into all scrapers
