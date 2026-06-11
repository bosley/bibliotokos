# bible

Embeddable Bible library with multiple translations. The database is embedded at compile time and extracted to the user's XDG data directory on first run.

## Setup

Call `Init(appName)` before using any other methods. The app name determines where the database is stored locally.

## Types

**Book** — abbreviation, full name, canonical order

**Version** — id (short code like "kjv", "asv", "web"), display name, language code. The id is what you pass as `versionID` to other methods.

**Verse** — book abbr, chapter, verse number, version id, text

**VersePage** — paginated result with total count, offset, and verse slice

**MultiVersionVerse** — single verse location with text from multiple versions (keyed by version id)

**PassageRange** — parsed reference with encoded start/end positions for comparison

**Collection** — `OldTestament`, `NewTestament`, or `Apocrypha`

## Methods

**GetVersions** — lists all available translations from the database; returns id, display name, and language for each. Use the id as versionID in other calls, and the name for UI display. Adding translations to the database makes them available here automatically.

**GetBooks** — lists all books in canonical order

**GetBooksByCollection(collection)** — books filtered to OT, NT, or Apocrypha only

**HasCollection(versionID, collection)** — whether a version contains any verses from OT, NT, or Apocrypha

**GetChapterCount(versionID, book)** — number of chapters in a book for a version

**GetVerseCount(versionID, book, chapter)** — number of verses in a chapter

**QueryPage(reference, versionID, offset, limit)** — fetch verses by reference with pagination (limit capped at 200)

**QueryMultiVersion(reference, versionIDs)** — same reference across multiple versions, returns combined results

**Search(versionID, book, phrase, offset, limit)** — substring search within a version; pass empty book to search all books

**RandomVerse(versionID, book)** — returns a random verse; pass empty book for any verse in the version

**ResolveRange(reference)** — parses a reference string into a PassageRange with encoded positions

## Reference Format

References are flexible and support common aliases:

- Book only: `Genesis`, `Gen`, `Gn`
- Chapter: `Genesis 1`, `Gen 1`
- Verse: `Genesis 1:1`, `Gen 1:1`
- Verse range: `Genesis 1:1-5`
- Cross-book range: `Genesis 1:1 - Exodus 2:10`

Numbered books accept various forms: `1 Samuel`, `1Sam`, `1sam`, `Sa1`
