/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"notepad/biz"
	"os"
	"path/filepath"
	"strings"
	"time"

	"astuart.co/goq"
	"github.com/araddon/dateparse"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var todayOnly int16

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "pull notes from remote.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pull notes...")
		entries, err := getNoteEntries()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		now := today().Add(time.Duration(-todayOnly) * 24 * time.Hour)
		filterd := entries[:0]
		for _, note := range entries {
			t, _ := dateparse.ParseLocal(note.UpdatedAt)
			if t.Before(now) {
				continue
			}
			filterd = append(filterd, note)
		}
		entries = filterd

		notes, err := getNotes(entries)
		if err != nil {
			cmd.PrintErr(err)
			return
		}

		var key = viper.GetString("key")
		home, err := os.UserHomeDir()
		if err != nil {
			cmd.PrintErr(err)
			return
		}
		var vault = strings.ReplaceAll(viper.GetString("vault"), "~", home)

		for _, note := range notes {
			cmd.Println("note: ", note.Entry)
			noteFile := filepath.Join(vault, note.Entry.Title+".md")
			f, err := os.Create(noteFile)
			if err != nil {
				panic(err)
			}
			content, err := biz.Decrypt([]byte(key), []byte(note.Content))
			if err != nil {
				panic(err)
			}
			_, err = f.Write(content)
			if err != nil {
				fmt.Println(err)
			}
		}
	},
}

func today() time.Time {
	y, m, d := time.Now().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func init() {
	rootCmd.AddCommand(pullCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pullCmd.PersistentFlags().String("foo", "", "A help for foo")

	pullCmd.Flags().Int16VarP(&todayOnly, "day", "d", 0, "pull x day note from the repository")
}

const host = "https://cn.anotepad.com"

func getNoteEntries() ([]NoteEntry, error) {

	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/note/list", host), nil)

	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("cookie", viper.GetString("cookie"))
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	list := Notes{}

	doc := goq.NewDecoder(res.Body)
	err = doc.Decode(&list)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(list.List); i++ {
		entry := list.List[i]
		pos := strings.Index(entry.UpdatedAt, ": ")
		if pos > 0 {
			entry.UpdatedAt = entry.UpdatedAt[pos+2:]
		}
		pos = strings.Index(entry.UpdatedAt, "M ")
		if pos > 0 {
			entry.UpdatedAt = entry.UpdatedAt[0 : pos-2]
		}
		list.List[i] = entry
	}

	return list.List, nil
}

type NoteEntry struct {
	Title     string `goquery:".note-title"`
	Link      string `goquery:".note-title,[href]"`
	UpdatedAt string `goquery:".note-title,[title]"`
}
type Notes struct {
	List []NoteEntry `goquery:".note-item"`
}

type Note struct {
	Entry   *NoteEntry
	Content string `goquery:"div.note_content .plaintext"`
}

func getNoteDetail(entry *NoteEntry) (*Note, error) {
	r, err := http.Get(host + entry.Link)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	doc := goq.NewDecoder(r.Body)
	var note = &Note{}
	err = doc.Decode(note)
	if err != nil {
		return nil, err
	}
	d, err := url.QueryUnescape(string(note.Content))
	if err != nil {
		return nil, err

	}
	note.Content = string(d)
	return note, nil
}

func getNotes(entries []NoteEntry) ([]*Note, error) {
	var notes []*Note
	for _, entry := range entries {
		item := entry
		note, err := getNoteDetail(&item)
		if err != nil {
			return nil, err
		}

		note.Entry = &item
		notes = append(notes, note)
	}
	return notes, nil
}
