package webview

import (
	"github.com/Oppodelldog/simpleci/build"
	"github.com/Oppodelldog/simpleci/webview/assets"
	"html/template"
	"io"
	"path"
	"path/filepath"
)

func RenderIndexPage(w io.Writer) error {

	t, err := newPageTemplate("about.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
		}{})
}

func RenderQueuePage(w io.Writer) error {

	t, err := newPageTemplate("queue.html")
	if err != nil {
		return err
	}
	return t.Execute(w,
		struct {
			Builds []build.Build
		}{
			Builds: build.GetBuildQueueList(),
		})
}

func RenderLogPage(w io.Writer, buildId string, logId int) error {
	t, err := newPageTemplate("log.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			BuildID string
			LogID   int
			Log     string
		}{
			BuildID: buildId,
			LogID:   logId,
			Log:     build.GetBuildLog(buildId, logId),
		})
}

func RenderAbortPage(w io.Writer, id string) error {
	t, err := newPageTemplate("abort.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			Error error
		}{
			Error: build.AbortBuild(id),
		})
}

func RenderBuildsPage(w io.Writer) error {
	t, err := newPageTemplate("builds.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			Repositories []build.Repository
		}{
			Repositories: build.GetRepositories(),
		})
}

func RenderBuildPage(w io.Writer, buildId string) error {
	t, err := newPageTemplate("build.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			Builds  []int
			BuildID string
		}{
			BuildID: buildId,
			Builds:  build.GetBuild(buildId),
		})
}

func getPartials() []string {
	partialsPath := path.Join("/templates", "partials")
	return []string{
		path.Join(partialsPath, "logo.html"),
		path.Join(partialsPath, "head.html")}
}

func newPageTemplate(pageFileName string) (*template.Template, error) {
	pageFilePath := path.Join("/templates", pageFileName)
	name := filepath.Base(pageFilePath)

	t := template.New(name)

	b, err := assets.Templates.ReadFile(pageFilePath)
	if err != nil {
		return nil, err
	}
	_, err = t.Parse(string(b))
	if err != nil {
		return nil, err
	}

	for _, partial := range getPartials() {
		b, err := assets.Templates.ReadFile(partial)
		if err != nil {
			return nil, err
		}
		tmpl := t.New(partial)
		_, err = tmpl.Parse(string(b))
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
