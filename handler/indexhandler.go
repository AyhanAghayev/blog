package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/renderer/html"
)

type Post struct {
	Name    string
	Content template.HTML
}

var mdRenderer goldmark.Markdown = goldmark.New(
	goldmark.WithExtensions(
		highlighting.NewHighlighting(
			highlighting.WithStyle("dracula"),
		),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

func ServeAbout(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func ServeProjects(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/projects.html")
}

func ServePosts(w http.ResponseWriter, r *http.Request) {
	var posts []Post
	dirPath := "./contents"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		f, err := os.ReadFile(filepath.Join(dirPath, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		var buf bytes.Buffer
		r := []byte(f)
		if len(r) > 80 {
			f = f[:80]
		}

		if err := mdRenderer.Convert(f, &buf); err != nil {
			panic(err)
		}

		posts = append(posts, Post{
			Name:    fileName,
			Content: template.HTML(buf.String()),
		})
	}

	tmpl := template.Must(template.ParseFiles("./static/posts.html"))

	// Execute template
	if err := tmpl.Execute(w, posts); err != nil {
		log.Println("Template execution error:", err)
	}
}

func ServePost(w http.ResponseWriter, r *http.Request) {
	slugName := chi.URLParam(r, "slug")
	path := fmt.Sprintf("./contents/%s.md", slugName)
	f, err := os.ReadFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer


	if err := mdRenderer.Convert(f, &buf); err != nil {
		panic(err)
	}

	// Parse template
	tmpl := template.Must(template.ParseFiles("./static/post.html"))

	// Render page
	post := Post{
		Name:    slugName,
		Content: template.HTML(buf.String()),
	}

	tmpl.Execute(w, post)
}
