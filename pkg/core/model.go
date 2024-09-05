package core

import (
	"crypto/sha256"
	"slices"

	"github.com/oklog/ulid/v2"
)

type Site struct {
	UID         ulid.ULID `json:"uid"`
	Title       string    `json:"title"`
	Version     ulid.ULID `json:"version"`
	NextVersion ulid.ULID `json:"nextVersion"`
}

type Version struct {
	UID   ulid.ULID            `json:"uid"`
	Pages map[string]ulid.ULID `json:"pages"`
	Theme ulid.ULID            `json:"theme"`
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

type WritablePage struct {
	Title        string   `json:"title"`
	Path         string   `json:"path"`
	TemplateName string   `json:"templateName"`
	Tags         []string `json:"tags"`
	Excerpt      string   `json:"excerpt"`
}

type Blob struct {
	UID  ulid.ULID `json:"uid"`
	Hash []byte    `json:"hash"`
}

// Equal compares the hash of the input bytes to the Hash on the blob
func (b Blob) Equal(input []byte) (bool, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(input); err != nil {
		return false, err
	}
	return slices.Equal(hasher.Sum(nil), b.Hash), nil
}

type Theme struct {
	UID       ulid.ULID            `json:"uid"`
	Templates map[string]ulid.ULID `json:"templates"`
	CSSAssets []ulid.ULID          `json:"cssAssets"`
	JSAssets  []ulid.ULID          `json:"jsAssets"`
}

type TargetVersion struct {
	uid     ulid.ULID
	current bool
	next    bool
}

func (v TargetVersion) UID() ulid.ULID {
	return v.uid
}

func (v TargetVersion) IsCurrent() bool {
	return v.current
}

func (v TargetVersion) IsNext() bool {
	return v.next
}

func NewTargetVersion(uid ulid.ULID) *TargetVersion {
	return &TargetVersion{uid: uid}
}

func NewCurrentTargetVersion(uid ulid.ULID) *TargetVersion {
	return &TargetVersion{
		uid:     uid,
		current: true,
	}
}

func NewNextTargetVersion(uid ulid.ULID) *TargetVersion {
	return &TargetVersion{
		uid:  uid,
		next: true,
	}
}
