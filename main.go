package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	wa := application.New(application.Options{
		Name:        "悦速进销存",
		Description: "悦速进销存-by:悦速互联-定制开发:17337771173(微信同号)",
		Services: []application.Service{
			application.NewService(app),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	win := wa.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:                  "悦速进销存",
		Width:                  1580,
		Height:                 968,
		MinWidth:               800,
		MinHeight:              600,
		URL:                    "/",
		BackgroundColour:       application.NewRGB(27, 38, 54),
		OpenInspectorOnStartup: false,
	})

	// 窗口关闭时退出应用
	win.OnWindowEvent(events.Common.WindowClosing, func(event *application.WindowEvent) {
		wa.Quit()
	})

	err := wa.Run()
	if err != nil {
		log.Fatal(err)
	}
}
