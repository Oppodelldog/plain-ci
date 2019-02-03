package webview

import (
	"github.com/Oppodelldog/simpleci/build"
	"html/template"
	"io"
	"path"
)

const templatesDir = "webview/assets/templates"
const ImagesDir = "webview/assets/images"

func RenderIndexPage(w io.Writer) error {
	templateFile := path.Join(templatesDir, "about.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	return tpl.Execute(w,
		struct {
		}{})
}

func RenderQueuePage(w io.Writer) error {

	templateFile := path.Join(templatesDir, "queue.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	return tpl.Execute(w,
		struct {
			Builds []build.Build
		}{
			Builds: build.GetBuildQueueList(),
		})
}

func RenderLogPage(w io.Writer, buildId string, logId int) error {
	templateFile := path.Join(templatesDir, "log.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	return tpl.Execute(w,
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
	templateFile := path.Join(templatesDir, "abort.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	return tpl.Execute(w,
		struct {
			Error error
		}{
			Error: build.AbortBuild(id),
		})
}

func RenderBuildsPage(w io.Writer) error {
	templateFile := path.Join(templatesDir, "builds.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	return tpl.Execute(w,
		struct {
			Repositories []build.Repository
		}{
			Repositories: build.GetRepositories(),
		})
}

func RenderBuildPage(w io.Writer, buildId string) error {
	templateFile := path.Join(templatesDir, "build.html")
	templates := append([]string{templateFile}, getPartials()...)
	tpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

	return tpl.Execute(w,
		struct {
			Builds  []int
			BuildID string
		}{
			BuildID: buildId,
			Builds:  build.GetBuild(buildId),
		})
}
func getPartials() []string {
	partialsPath := path.Join(templatesDir, "partials")
	return []string{
		path.Join(partialsPath, "logo.html"),
		path.Join(partialsPath, "head.html")}
}
