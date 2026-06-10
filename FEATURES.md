# BiblioTokos Features

BiblioTokos is a lightweight personal Bible study desktop application built with Wails v3 (Go) and Svelte. It serves as a focused, simpler alternative to complex Bible software like Logos.

## Scripture Reading

### Bible Navigation & Search
- Search and navigate to any Bible passage using natural query syntax
  - Book/chapter references: `Gen 1`, `John 3`, `Psalm 23`
  - Specific verses: `John 3:16`, `Romans 8:28`
  - Verse ranges: `Psalm 23:1-6`, `Romans 8:28-39`
  - Cross-book ranges: `Gen 1 - Rev 22`
- Clean, readable verse-by-verse display with generous spacing
- Smart verse reference parsing with support for book name aliases and abbreviations
- Real-time query updates via Enter key

### Multiple Translation Comparison
- Side-by-side display of multiple Bible translations/versions
- Dynamically add or remove translation panes with `+` and `×` buttons
- Resizable panes with drag-to-adjust width
- Each pane independently selects which translation to display
- All panes display the same passage for easy comparison

### Synchronized Scrolling
- **Scroll lock toggle** in the top search bar (lock/unlock icon)
- **Locked mode**: All translation panes scroll together, automatically aligned by verse position
- **Unlocked mode**: Each pane scrolls independently
- Re-locking automatically aligns all panes to the last-scrolled pane's position
- Verse-level synchronization (not pixel-based) ensures alignment across translations with different text lengths

### Performance & UX
- Virtualized scrolling with intelligent chunk loading
- Smooth scroll anchoring during dynamic content loads
- Intersection Observer-based visibility tracking
- Automatic scroll position preservation when adding/removing panes
- Copy buttons for verse references and text

## Notes System

### Note Management
- Dedicated notes window (separate from main scripture window)
- Create, edit, save, and delete personal study notes
- Each note has a title, content, and timestamps
- Notes list shows all notes sorted by most recently updated
- Full-text editing with automatic save

### Tagging & Organization
- Tag-based organization system
- Multiple tags per note
- Create and delete tags
- Tag display in note headers for quick scanning

### Linked Passages Feature
- **Link Bible passages to notes**: Associate specific verses or verse ranges with any note
- **Many-to-many relationships**: 
  - One note can link to multiple passages
  - One passage can appear in multiple notes
- **Automatic surfacing**: When viewing scripture, see all notes linked to the visible verses
- **Linked Notes Bar**: Persistent bar at bottom of scripture workspace showing relevant linked notes
- Real-time updates as you scroll through different passages
- Click linked note to open in notes window

## User Interface

### Navigation & Layout
- Clean, minimalist top navigation bar
- Scripture and Notes icon buttons with tooltips
- Centered search input with lock toggle
- Multi-window architecture:
  - Main scripture reading workspace
  - Separate notes windows (can open multiple)
- macOS-native window styling:
  - Translucent backdrop
  - Hidden inset titlebar
  - Double-click titlebar to maximize

### Theme Support
- Light and dark theme toggle (sun/moon icon)
- Theme preference persisted across sessions
- Carefully chosen color variables for optimal readability

### Empty States & Feedback
- Helpful empty state messages when no query entered
- Loading indicators during data fetches
- Query examples in empty states
- Smooth transitions and hover effects

## Technical Architecture

### Backend (Go)
- **Bible Service**: Embedded SQLite database with Bible texts; query parsing and verse resolution
- **Notes Service**: SQLite-based storage with full CRUD operations, tags, and linked passages
- **System Service**: Theme and system-level operations
- Foreign key constraints with cascading deletes for data integrity
- Position-based range queries for efficient linked passage lookups

### Frontend (Svelte)
- Reactive state management with Svelte stores
- Component-based architecture
- Wails v3 auto-generated bindings for Go<->JavaScript communication
- Custom events for inter-component communication
- Efficient DOM updates and virtualized rendering

### Data Storage
- Bible texts embedded in application binary (read-only)
- User notes stored in local SQLite database at `~/.local/share/BiblioTokos/notes.db` (XDG-compliant)
- Schema includes tables for notes, note content, tags, tag associations, and linked passages

## Design Philosophy

- **Lightweight and focused**: Only essential features, no bloat
- **Calm and readable**: Generous spacing, excellent typography, low cognitive load
- **Fast and responsive**: Virtualized scrolling, intelligent caching, minimal latency
- **Personal**: Your notes, your connections, your study system
- **Purposeful**: Every feature serves the core purpose of reading and studying Scripture
