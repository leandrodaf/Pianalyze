package store

import (
	"sync"
)

// State mantém o estado compartilhado do pipeline, como notas pressionadas.
type State struct {
	mu           sync.RWMutex
	PressedNotes []int // Slice para manter a ordem das notas pressionadas
	LastNoteTime uint64
}

// NewPipelineState inicializa o estado do pipeline.
func NewPipelineState() *State {
	return &State{
		PressedNotes: []int{},
	}
}

// AddNote adiciona uma nota pressionada.
func (ps *State) AddNote(note int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	// Verifica se a nota já está pressionada para evitar duplicações
	for _, n := range ps.PressedNotes {
		if n == note {
			return
		}
	}
	ps.PressedNotes = append(ps.PressedNotes, note)
}

// RemoveNote remove uma nota que foi solta.
func (ps *State) RemoveNote(note int) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	for i, n := range ps.PressedNotes {
		if n == note {
			// Remove a nota mantendo a ordem
			ps.PressedNotes = append(ps.PressedNotes[:i], ps.PressedNotes[i+1:]...)
			break
		}
	}
}

// GetPressedNotes retorna uma cópia das notas atualmente pressionadas.
func (ps *State) GetPressedNotes() []int {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	notesCopy := make([]int, len(ps.PressedNotes))
	copy(notesCopy, ps.PressedNotes)
	return notesCopy
}

// UpdateLastNoteTime atualiza o timestamp da última nota se houver uma alteração.
func (ps *State) UpdateLastNoteTime(timestamp uint64) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if ps.LastNoteTime != timestamp {
		ps.LastNoteTime = timestamp
	}
}

// GetLastNoteTime retorna o timestamp da última nota.
func (ps *State) GetLastNoteTime() uint64 {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	return ps.LastNoteTime
}
