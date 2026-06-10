package notes

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"

	"bibliotokos/platform"
)

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt"`
}

type NotesService struct {
	filePath string
}

func (s *NotesService) Init() error {
	dir := filepath.Join(xdg.DataHome, platform.GetInstallAppName())
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	s.filePath = filepath.Join(dir, "notes.json")
	return nil
}

func (s *NotesService) load() ([]Note, error) {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return []Note{}, nil
	}
	if err != nil {
		return nil, err
	}
	var notes []Note
	if err := json.Unmarshal(data, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *NotesService) persist(notes []Note) error {
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *NotesService) GetNotes() ([]Note, error) {
	return s.load()
}

func (s *NotesService) SaveNote(note Note) error {
	notes, err := s.load()
	if err != nil {
		return err
	}
	note.UpdatedAt = time.Now().Format(time.RFC3339)
	for i, n := range notes {
		if n.ID == note.ID {
			notes[i] = note
			return s.persist(notes)
		}
	}
	notes = append([]Note{note}, notes...)
	return s.persist(notes)
}

func (s *NotesService) DeleteNote(id string) error {
	notes, err := s.load()
	if err != nil {
		return err
	}
	filtered := make([]Note, 0, len(notes))
	for _, n := range notes {
		if n.ID != id {
			filtered = append(filtered, n)
		}
	}
	return s.persist(filtered)
}
