package core

import (
	"github.com/oklog/ulid/v2"
)

type Site struct {
	UID     ulid.ULID `json:"uid"`
	Title   string    `json:"title"`
	Version ulid.ULID `json:"version"`
}

type Version struct {
	UID   ulid.ULID   `json:"uid"`
	Pages []ulid.ULID `json:"pages"`
	Theme ulid.ULID   `json:"theme"`
}

type Page struct {
	UID          ulid.ULID `json:"uid"`
	Title        string    `json:"title"`
	Path         string    `json:"path"`
	Content      ulid.ULID `json:"content"`
	TemplateName string    `json:"templateName"`
	Tags         []string  `json:"tags"`
	Excerpt      string    `json:"excerpt"`
}

type Theme struct {
	UID       ulid.ULID            `json:"uid"`
	Templates map[string]ulid.ULID `json:"templates"`
	CSSAssets []ulid.ULID          `json:"cssAssets"`
	JSAssets  []ulid.ULID          `json:"jsAssets"`
}
