# Bibliotokos

Bibliotokos is a lightweight desktop application for reading and comparing Bible texts. It sits adjacent to tools like Logos — focused on scripture lookup, side-by-side translation comparison, and personal study notes rather than a full theological library.


## Features

**Scripture** — Search by book name, alias, chapter, verse, range, or cross-book range (e.g. `John 3:16`, `Gen 1`, `Psalm 23:1-6`, `Gen 1 - Rev 22`). Open multiple translation panes to compare texts side by side. Scroll lock keeps panes aligned as you read.

**Notes** — A dedicated notes workspace for study notes with tags. Use **Link passage** on a note to associate Bible references with your writing. When you view a linked passage in Scripture, matching notes appear in a bar at the bottom of the window — click to jump straight to the note.

## Supported Texts

| Code | Name | Language | Coverage | Notes |
|------|------|----------|----------|-------|
| ASV | American Standard Version | English | Old & New Testament | Public domain |
| WEB | World English Bible | English | Old & New Testament | Public domain |
| KJV | King James Version | English | Old & New Testament | Public domain |
| GNT | Greek New Testament | Greek | New Testament only | Original-language NT |
| TAN | Tanach | Hebrew | Old Testament only | Hebrew Bible |
| SEP | Septuagint | Greek | Old Testament | Greek OT translation |
| VUL | Vulgate | Latin | Old & New Testament | Latin Bible |
| ENV | The Emergent Noetic Version | English | Old & New Testament | Exclusive to Bibliotokos; translated from Hebrew Tanach + Greek NT |
| VENV | The Vulgate Emergent Noetic Version | English | Old & New Testament | Exclusive to Bibliotokos; translated from the Latin Vulgate |

### ENV & VENV

ENV and VENV are exclusive to Bibliotokos. They were created by iteratively translating original-language texts through multiple LLMs to preserve meaning as closely as possible.

- **ENV** — sourced from Hebrew Tanach and Greek New Testament
- **VENV** — sourced from the Latin Vulgate

Source texts are available in the translator.

### Other Sources

Standard translations (ASV, WEB, KJV, GNT, TAN, SEP, VUL) are drawn from public domain and Creative Commons sources, including [sacred-texts.com](https://sacred-texts.com/bib/osrc/index.htm) and [github.com/seven1m/open-bibles](https://github.com/seven1m/open-bibles).

## Getting Started

The application lives in `bibliotokos-app/`. From that directory:

```bash
wails3 dev
```

To build a production binary:

```bash
wails3 build
```
