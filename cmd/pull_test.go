package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"astuart.co/goq"
	"github.com/araddon/dateparse"
)

func TestParse(t *testing.T) {
	html := `
        <li id="note_ffxxl1" class="note-item dropdown">
            <a href="/note/read/ffxxl1" title="阅读笔记" target="read" class="note-preview"><img src="https://cdn.anotepad.com/Images/notepad_lockopen.svg?v2" height="28" width="28" alt="Read Note" /></a>
            <a href="/notes/ffxxl1" title="最近更新时间: 2023-01-29 12:22:13 PM &#10;2023-01-29" class="note-title">
                2023-01-29
            </a>
            <span class="note-action topLinks">
                <a id="noteActionLabelffxxl1" data-target="#" href="#" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false" class="note-action-button"></a>
                <ul class="dropdown-menu" aria-labelledby="noteActionLabelffxxl1">
                    <li><a href="/notes/ffxxl1">编辑</a></li>
                    <li><a href="/note/read/ffxxl1" target="read">读</a></li>
                    <li><a href="Javascript:fnDeleteNote('ffxxl1');">删除</a></li>
                    <li><a href="Javascript:fnCopyNote('ffxxl1');">复制</a></li>
                </ul>
            </span>

        </li>
        <li id="note_ffxxl" class="note-item dropdown">
            <a href="/note/read/ffxxl" title="阅读笔记" target="read" class="note-preview"><img src="https://cdn.anotepad.com/Images/notepad_lockopen.svg?v2" height="28" width="28" alt="Read Note" /></a>
            <a href="/notes/ffxxl" title="最近更新时间: 2023-01-29 10:36:06 AM &#10;2023-01-29" class="note-title">
                2023-01-29
            </a>
            <span class="note-action topLinks">
                <a id="noteActionLabelffxxl" data-target="#" href="#" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false" class="note-action-button"></a>
                <ul class="dropdown-menu" aria-labelledby="noteActionLabelffxxl">
                    <li><a href="/notes/ffxxl">编辑</a></li>
                    <li><a href="/note/read/ffxxl" target="read">读</a></li>
                    <li><a href="Javascript:fnDeleteNote('ffxxl');">删除</a></li>
                    <li><a href="Javascript:fnCopyNote('ffxxl');">复制</a></li>
                </ul>
            </span>

        </li>`
	doc := goq.NewDecoder(strings.NewReader(html))

	type noteEntry struct {
		Title     string `goquery:".note-title"`
		Link      string `goquery:".note-title,[href]"`
		UpdatedAt string `goquery:".note-title,[title]"`
	}
	type notes struct {
		List []noteEntry `goquery:".note-item"`
	}
	var err error
	n := &notes{}
	err = doc.Decode(n)
	if err != nil {
		panic(err)
	}
	for _, entry := range n.List {
		pos := strings.Index(entry.UpdatedAt, ": ")
		if pos > 0 {
			entry.UpdatedAt = entry.UpdatedAt[pos+2:]
		}
		pos = strings.Index(entry.UpdatedAt, "M ")
		if pos > 0 {
			entry.UpdatedAt = entry.UpdatedAt[0 : pos-2]
		}
		updatedAt, _ := time.ParseInLocation("2006-01-02 15:04:05", entry.UpdatedAt, time.Local)
		fmt.Println(entry, updatedAt)
	}
}

func TestNote(t *testing.T) {
	html, err := os.Open("note.html")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = html.Close()
	}()

	doc := goq.NewDecoder(html)

	x := Note{}
	err = doc.Decode(&x)

	if err != nil {
		panic(err)
	}

	t.Log(x.Content)
}

func TestDay(t *testing.T) {
	n := today()
	t.Log(n)

	entries := []*NoteEntry{
		{Title: "1", UpdatedAt: "2023-01-29 12:00:00"},
		{Title: "2", UpdatedAt: "2023-01-07 12:00:00"},
		{Title: "3", UpdatedAt: "2023-02-09 12:00:00"},
	}
	now := today().Add(time.Duration(-2) * 24 * time.Hour)
	filterd := entries[:0]
	for _, note := range entries {
		tt, _ := dateparse.ParseLocal(note.UpdatedAt)
		if tt.Before(now) {
			continue
		}
		fmt.Println(note.UpdatedAt, now)
		filterd = append(filterd, note)
	}
	t.Log(filterd)
}
