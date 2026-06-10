package notes

import (
	crand "crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"

	"bibliotokos/platform"
)

type NoteHeader struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	UpdatedAt string   `json:"updatedAt"`
	Tags      []string `json:"tags"`
}

type Note struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	UpdatedAt string   `json:"updatedAt"`
	Tags      []string `json:"tags"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NotesService struct {
	db *sql.DB
}

func (s *NotesService) Init() error {
	dir := filepath.Join(xdg.DataHome, platform.GetInstallAppName())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	dbPath := filepath.Join(dir, "notes.db")

	existed := true
	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		existed = false
	}
	_ = existed

	db, err := sql.Open("sqlite", "file:"+dbPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	s.db = db

	if _, err := s.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf("enable foreign keys: %w", err)
	}

	schema := `
CREATE TABLE IF NOT EXISTS notes (
    id         TEXT PRIMARY KEY,
    title      TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS note_content (
    note_id TEXT PRIMARY KEY REFERENCES notes(id) ON DELETE CASCADE,
    content TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS tags (
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS note_tags (
    note_id TEXT NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
    tag_id  TEXT NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (note_id, tag_id)
);`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}

func (s *NotesService) noteTagNames(noteID string) ([]string, error) {
	rows, err := s.db.Query(`
		SELECT t.name FROM tags t
		JOIN note_tags nt ON nt.tag_id = t.id
		WHERE nt.note_id = ?
		ORDER BY t.name`, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tags = append(tags, name)
	}
	if tags == nil {
		tags = []string{}
	}
	return tags, rows.Err()
}

func (s *NotesService) syncTags(tx *sql.Tx, noteID string, tagNames []string) error {
	if _, err := tx.Exec("DELETE FROM note_tags WHERE note_id = ?", noteID); err != nil {
		return err
	}
	for _, name := range tagNames {
		var tagID string
		err := tx.QueryRow("SELECT id FROM tags WHERE name = ?", name).Scan(&tagID)
		if errors.Is(err, sql.ErrNoRows) {
			tagID = newID()
			if _, err := tx.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", tagID, name); err != nil {
				return err
			}
		} else if err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT OR IGNORE INTO note_tags (note_id, tag_id) VALUES (?, ?)", noteID, tagID); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotesService) ListNotes() ([]NoteHeader, error) {
	rows, err := s.db.Query(`SELECT id, title, updated_at FROM notes ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var headers []NoteHeader
	for rows.Next() {
		var h NoteHeader
		if err := rows.Scan(&h.ID, &h.Title, &h.UpdatedAt); err != nil {
			return nil, err
		}
		tags, err := s.noteTagNames(h.ID)
		if err != nil {
			return nil, err
		}
		h.Tags = tags
		headers = append(headers, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if headers == nil {
		headers = []NoteHeader{}
	}
	return headers, nil
}

func (s *NotesService) GetNote(id string) (Note, error) {
	var n Note
	err := s.db.QueryRow(`SELECT n.id, n.title, n.updated_at, COALESCE(c.content, '')
		FROM notes n
		LEFT JOIN note_content c ON c.note_id = n.id
		WHERE n.id = ?`, id).Scan(&n.ID, &n.Title, &n.UpdatedAt, &n.Content)
	if err != nil {
		return Note{}, fmt.Errorf("get note: %w", err)
	}
	tags, err := s.noteTagNames(id)
	if err != nil {
		return Note{}, err
	}
	n.Tags = tags
	return n, nil
}

func (s *NotesService) SaveNote(note Note) error {
	now := time.Now().Format(time.RFC3339)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO notes (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET title = excluded.title, updated_at = excluded.updated_at`,
		note.ID, note.Title, now, now)
	if err != nil {
		return fmt.Errorf("upsert note: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO note_content (note_id, content) VALUES (?, ?)
		ON CONFLICT(note_id) DO UPDATE SET content = excluded.content`,
		note.ID, note.Content)
	if err != nil {
		return fmt.Errorf("upsert content: %w", err)
	}

	tags := note.Tags
	if tags == nil {
		tags = []string{}
	}
	if err := s.syncTags(tx, note.ID, tags); err != nil {
		return fmt.Errorf("sync tags: %w", err)
	}

	return tx.Commit()
}

func (s *NotesService) DeleteNote(id string) error {
	if _, err := s.db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return err
	}
	_, err := s.db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

func (s *NotesService) GetTags() ([]Tag, error) {
	rows, err := s.db.Query("SELECT id, name FROM tags ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	if tags == nil {
		tags = []Tag{}
	}
	return tags, rows.Err()
}

func (s *NotesService) CreateTag(name string) (Tag, error) {
	id := newID()
	_, err := s.db.Exec("INSERT INTO tags (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return Tag{}, fmt.Errorf("create tag: %w", err)
	}
	return Tag{ID: id, Name: name}, nil
}

func (s *NotesService) DeleteTag(id string) error {
	_, err := s.db.Exec("DELETE FROM tags WHERE id = ?", id)
	return err
}

func newID() string {
	b := make([]byte, 16)
	if _, err := crand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
