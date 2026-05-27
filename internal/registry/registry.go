package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ── Types ─────────────────────────────────────────────────────────────────────

type Todo struct {
	Task    string `json:"task"`
	Done    bool   `json:"done"`
	Created string `json:"created"`
}

type Project struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	LastOpened  string `json:"last_opened"`
	RepoURL     string `json:"repo_url"`
	NetlifyURL  string `json:"netlify_url"`
	Description string `json:"description"`
	Todos       []Todo `json:"todos,omitempty"`
}

type store struct {
	Projects []Project `json:"projects"`
}

// Registry manages the projects.json file.
type Registry struct {
	path string
}

// New returns a Registry backed by the given JSON file path.
func New(path string) *Registry {
	return &Registry{path: path}
}

// ── I/O ───────────────────────────────────────────────────────────────────────

func (r *Registry) load() (*store, error) {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(r.path); os.IsNotExist(err) {
		s := &store{Projects: []Project{}}
		return s, r.save(s)
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("reading registry: %w", err)
	}
	var s store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parsing registry: %w", err)
	}
	// Migrate: fix null types
	for i := range s.Projects {
		if s.Projects[i].Type == "" {
			s.Projects[i].Type = "vanilla"
		}
	}
	return &s, nil
}

// save performs an atomic write (temp file + rename).
func (r *Registry) save(s *store) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	tmp := r.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, r.path)
}

// ── CRUD ──────────────────────────────────────────────────────────────────────

func (r *Registry) Add(p Project) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	today := time.Now().Format("2006-01-02")
	p.Status = "active"
	p.Created = today
	p.LastOpened = today
	s.Projects = append(s.Projects, p)
	return r.save(s)
}

func (r *Registry) Exists(name string) bool {
	s, err := r.load()
	if err != nil {
		return false
	}
	for _, p := range s.Projects {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (r *Registry) Get(name string) (Project, bool) {
	s, err := r.load()
	if err != nil {
		return Project{}, false
	}
	for _, p := range s.Projects {
		if p.Name == name {
			return p, true
		}
	}
	return Project{}, false
}

func (r *Registry) All() ([]Project, error) {
	s, err := r.load()
	if err != nil {
		return nil, err
	}
	return s.Projects, nil
}

func (r *Registry) Active() ([]Project, error) {
	all, err := r.All()
	if err != nil {
		return nil, err
	}
	var active []Project
	for _, p := range all {
		if p.Status == "active" {
			active = append(active, p)
		}
	}
	return active, nil
}

// LastActive returns the most recently opened active project name.
func (r *Registry) LastActive() string {
	active, err := r.Active()
	if err != nil || len(active) == 0 {
		return ""
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].LastOpened > active[j].LastOpened
	})
	return active[0].Name
}

func (r *Registry) UpdateField(name, field, value string) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	for i := range s.Projects {
		if s.Projects[i].Name == name {
			switch field {
			case "name":
				s.Projects[i].Name = value
			case "path":
				s.Projects[i].Path = value
			case "type":
				s.Projects[i].Type = value
			case "status":
				s.Projects[i].Status = value
			case "repo_url":
				s.Projects[i].RepoURL = value
			case "netlify_url":
				s.Projects[i].NetlifyURL = value
			case "description":
				s.Projects[i].Description = value
			case "last_opened":
				s.Projects[i].LastOpened = value
			case "visibility":
				// stored in description as a note; no separate field
			}
			return r.save(s)
		}
	}
	return fmt.Errorf("project %q not found", name)
}

func (r *Registry) UpdateLastOpened(name string) error {
	return r.UpdateField(name, "last_opened", time.Now().Format("2006-01-02"))
}

func (r *Registry) Delete(name string) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	filtered := s.Projects[:0]
	for _, p := range s.Projects {
		if p.Name != name {
			filtered = append(filtered, p)
		}
	}
	s.Projects = filtered
	return r.save(s)
}

func (r *Registry) SetStatus(name, status string) error {
	return r.UpdateField(name, "status", status)
}

// CountByStatus returns counts: active, completed, archived.
func (r *Registry) CountByStatus() (active, completed, archived int) {
	s, _ := r.load()
	for _, p := range s.Projects {
		switch p.Status {
		case "active":
			active++
		case "completed":
			completed++
		case "archived":
			archived++
		}
	}
	return
}

// AddTodo appends a todo to a project.
func (r *Registry) AddTodo(name, task string) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	for i := range s.Projects {
		if s.Projects[i].Name == name {
			s.Projects[i].Todos = append(s.Projects[i].Todos, Todo{
				Task:    task,
				Done:    false,
				Created: time.Now().Format("2006-01-02"),
			})
			return r.save(s)
		}
	}
	return fmt.Errorf("project %q not found", name)
}

// ToggleTodo marks a todo done/undone by index.
func (r *Registry) ToggleTodo(name string, idx int) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	for i := range s.Projects {
		if s.Projects[i].Name == name {
			todos := s.Projects[i].Todos
			if idx < 0 || idx >= len(todos) {
				return fmt.Errorf("todo index %d out of range", idx)
			}
			s.Projects[i].Todos[idx].Done = !s.Projects[i].Todos[idx].Done
			return r.save(s)
		}
	}
	return fmt.Errorf("project %q not found", name)
}

// Rename renames a project (name + path update atomically).
func (r *Registry) Rename(oldName, newName, newPath string) error {
	s, err := r.load()
	if err != nil {
		return err
	}
	for i := range s.Projects {
		if s.Projects[i].Name == oldName {
			s.Projects[i].Name = newName
			s.Projects[i].Path = newPath
			return r.save(s)
		}
	}
	return fmt.Errorf("project %q not found", oldName)
}
