package main

import (
	"embed"
	"log"
	neturl "net/url"

	"github.com/wailsapp/wails/v3/pkg/application"

	"bibliotokos/platform"
	"bibliotokos/services/bible"
	"bibliotokos/services/notes"
	"bibliotokos/services/system"
)

func eventString(data any) string {
	switch v := data.(type) {
	case string:
		return v
	case []any:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	bibleService := &bible.BibleService{}
	if err := bibleService.Init(platform.GetInstallAppName()); err != nil {
		log.Fatal(err)
	}

	notesService := notes.New(bibleService)
	if err := notesService.Init(); err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name:        "BiblioTokos",
		Description: "BiblioTokos is a desktop application for managing your study",
		Services: []application.Service{
			application.NewService(&system.SystemService{}),
			application.NewService(bibleService),
			application.NewService(notesService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:  "BiblioTokos",
		Width:  1200,
		Height: 900,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 38,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	app.Event.On("open-notes", func(e *application.CustomEvent) {
		url := "/?view=notes"
		if noteID := eventString(e.Data); noteID != "" {
			url += "&note=" + neturl.QueryEscape(noteID)
		}
		app.Window.NewWithOptions(application.WebviewWindowOptions{
			Title:  "Notes — BiblioTokos",
			Width:  900,
			Height: 680,
			Mac: application.MacWindow{
				InvisibleTitleBarHeight: 38,
				Backdrop:                application.MacBackdropTranslucent,
				TitleBar:                application.MacTitleBarHiddenInset,
			},
			BackgroundColour: application.NewRGB(27, 38, 54),
			URL:              url,
		})
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
