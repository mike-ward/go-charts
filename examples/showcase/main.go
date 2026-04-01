package main

import (
	"github.com/mike-ward/go-gui/gui"
	"github.com/mike-ward/go-gui/gui/backend"
)

const (
	scrollCatalog uint32 = iota + 1
	scrollDetail
	focusSearch
)

const catalogWidth float32 = 250

func main() {
	gui.SetTheme(gui.ThemeDarkBordered)

	w := gui.NewWindow(gui.WindowCfg{
		State:  newShowcaseApp(),
		Title:  "Charts Showcase",
		Width:  800,
		Height: 600,
		OnInit: func(w *gui.Window) {
			w.UpdateView(mainView)
		},
	})

	backend.Run(w)
}

func mainView(w *gui.Window) gui.View {
	ww, wh := w.WindowSize()
	return gui.Row(gui.ContainerCfg{
		Width:   float32(ww),
		Height:  float32(wh),
		Sizing:  gui.FixedFixed,
		Padding: gui.NoPadding,
		Spacing: gui.NoSpacing,
		Content: []gui.View{catalogPanel(w), detailPanel(w)},
	})
}
