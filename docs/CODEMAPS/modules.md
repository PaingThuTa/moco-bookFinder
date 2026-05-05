# Module Codemap

<!-- AUTO-GENERATED -->

## `cmd/bot/main.go`

Entry point. Wires config, HTTP client, scrapers, source manager, and bot handler together. Runs the long-polling update loop.

**Dependencies**: `internal/config`, `internal/bot`, `internal/source`

```
main()
  config.Load()
  telego.NewBot(token)
  botClient.GetMe()
  scrapers := [ZLibrary, OceanPDF, LibGen]
  source.NewSourceManager(scrapers)
  bot.NewHandler(cfg, mgr)
  botClient.UpdatesViaLongPolling()
  for update := range updates { route to handlers }
```

---

## `internal/config/config.go`

Environment variable loading and user authorization.

**Dependencies**: none (stdlib only)

```
Config struct { BotToken, AllowedUserIDs }
  Load() -> *Config, error    // reads TELEGRAM_BOT_TOKEN, ALLOWED_USER_IDS
  IsAllowed(userID int64) -> bool
```

**Tests**: `config_test.go` (4 tests: valid load, missing token, missing IDs, IsAllowed)

---

## `internal/bot/handler.go`

Telegram message and callback handler. Uses the `telego` SDK.

**Dependencies**: `internal/config`, `internal/source`, `github.com/mymmrac/telego`

```
BotHandler struct { cfg, mgr, results map, mutex }
  NewHandler(cfg, mgr) -> *BotHandler
  HandleStart(bot, update)              // /start command
  HandleSearch(bot, update)             // /search command
  HandleCallback(bot, update)           // inline button tap
  checkAuthorized(bot, update) -> bool  // user whitelist check
  parseSearchArgs(args) -> title, author // --author flag parsing
  extractCommandArgs(text) -> string     // strip command prefix
  storeResults(chatID, results)          // cache with 10m TTL
  getResults(chatID) -> []BookResult     // retrieve cached results
```

**Tests**: `handler_test.go` (5 tests: parse args title-only, with author, no author flag, mock search, empty query)

---

## `internal/source/source.go`

`Scraper` interface and `SourceManager` that chains scrapers with fallback.

**Dependencies**: `net/http`, `context`

```
Scraper interface { Name(), Search(ctx, title, author) }
BookResult struct { Title, Author, DownloadURL, Source }
SourceManager struct { scrapers, client }
  NewSourceManager(scrapers) -> *SourceManager
  HTTPClient() -> *http.Client
  Search(ctx, title, author) -> []BookResult, error
FormatBooks(results []BookResult) -> string
```

**Tests**: `source_test.go` (5 tests: first-success, all-fail, empty-skipped, FormatBooks, FormatBooks empty)

---

## `internal/source/zlibrary.go`

Z-Library scraper. Searches `z-library.bz/s/<query>` using goquery CSS selectors.

**Dependencies**: `github.com/PuerkitoBio/goquery`

```
ZLibraryScraper struct { client, base }
  NewZLibraryScraper(client) -> *ZLibraryScraper
  Name() -> "Z-Library"
  Search(ctx, title, author) -> []BookResult, error
resolveURL(base, relative) -> string
```

---

## `internal/source/oceanpdf.go`

Ocean of PDF scraper. Searches `oceanofpdf.com/?s=<query>`.

**Dependencies**: `github.com/PuerkitoBio/goquery`

```
OceanPDFScraper struct { client, base }
  NewOceanPDFScraper(client) -> *OceanPDFScraper
  Name() -> "Ocean of PDF"
  Search(ctx, title, author) -> []BookResult, error
```

---

## `internal/source/libgen.go`

LibGen scraper. Uses `libgen.is/search.php?req=<title>&column=<title|author>`.

**Dependencies**: `github.com/PuerkitoBio/goquery`

```
LibGenScraper struct { client, base }
  NewLibGenScraper(client) -> *LibGenScraper
  Name() -> "LibGen"
  Search(ctx, title, author) -> []BookResult, error
```
