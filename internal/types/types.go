package types

import (
	"github.com/google/uuid"
)

type DropdownItem struct {
	Name  string
	Count string
	Href  string
}

type DataCount struct {
	ClassicLevelCount int
	PlatformerLevelCount int
	ClassicRecordCount int
	PlatformerRecordCount int
}

type Pagination struct {
	Start      int
	End        int
	Page       int
	Total      int
	TotalPages int
}

func (p Pagination) Pages() []int {
	if p.TotalPages <= 7 {
		pages := make([]int, p.TotalPages)
		for i := range pages {
			pages[i] = i + 1
		}
		return pages
	}
	pages := []int{1}
	if p.Page > 3 {
		pages = append(pages, 0)
	}
	start := max(p.Page-1, 2)
	end := min(p.Page+1, p.TotalPages-1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	if p.Page < p.TotalPages-2 {
		pages = append(pages, 0)
	}
	pages = append(pages, p.TotalPages)
	return pages
}

type SessionData struct {
	PlayerID       uuid.UUID
	PlayerUsername string
	PlayerAvatar   string
	UserRoles      []string
}

func (s SessionData) HasRole(roleName string) bool {
	for _, role := range s.UserRoles {
		if role == roleName {
			return true
		}
	}
	return false
}
