package webview

import (
	"html/template"
	"io"
	"path"
	"path/filepath"

	"github.com/Oppodelldog/plainci/build"
	"github.com/Oppodelldog/plainci/webview/assets"
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

func RenderAbortPage(abort func(id string) error, w io.Writer, id string) error {
	t, err := newPageTemplate("abort.html")
	if err != nil {
		return err
	}

	return t.Execute(w,
		struct {
			Error error
		}{
			Error: abort(id),
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

// (has a leading slash to work with static assets libray)
const templatesBasePath = "/templates"

func getPartials() []string {

	partialsPath := path.Join(templatesBasePath, "partials")
	return []string{
		path.Join(partialsPath, "logo.html"),
		path.Join(partialsPath, "head.html")}
}

func newPageTemplate(pageFileName string) (*template.Template, error) {
	pageFilePath := path.Join(templatesBasePath, pageFileName)
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
