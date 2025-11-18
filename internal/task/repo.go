package task

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var ErrNotFound = errors.New("task not found")

type Repo struct {
	mu    sync.RWMutex
	seq   int64
	items map[int64]*Task
}

func NewRepo() *Repo {
	return &Repo{items: make(map[int64]*Task)}
}

func (r *Repo) List(page, limit int, doneFilter *bool) []*Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*Task, 0, len(r.items))
	for _, t := range r.items {
		if doneFilter != nil && t.Done != *doneFilter {
			continue //пропускаем если не done, продолжаем цикл не прерывая
		}
		filtered = append(filtered, t)
	}

	sort.Slice(filtered, func(i, j int) bool { //сортировка листа
		return filtered[i].ID < filtered[j].ID
	})

	total := len(filtered)

	// Вычисляем начальный и конечный индексы среза
	start := (page - 1) * limit
	end := start + limit

	// Проверяем границы
	if start >= total {
		return []*Task{} // Пустой срез, если страница за пределами
	}

	if end > total {
		end = total // Обрезаем 'end', если он выходит за пределы
	}

	// Возвращаем подмножество задач
	return filtered[start:end]
}

func (r *Repo) Get(id int64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (r *Repo) Create(title string) *Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	now := time.Now()
	t := &Task{ID: r.seq, Title: title, CreatedAt: now, UpdatedAt: now, Done: false}
	r.items[t.ID] = t
	return t
}

func (r *Repo) Update(id int64, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	return t, nil
}

func (r *Repo) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.items[id]; !ok {
		return ErrNotFound
	}
	delete(r.items, id)
	return nil
}
