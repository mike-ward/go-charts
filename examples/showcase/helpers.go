package main

import "github.com/mike-ward/go-gui/gui"

// demoWithCode wraps a chart view with its source code shown
// below as a markdown code block.
func demoWithCode(
	w *gui.Window, id string, chartView gui.View, code string,
) gui.View {
	source := "```go\n" + code + "\n```"
	return gui.Column(gui.ContainerCfg{
		Sizing:  gui.FillFit,
		Padding: gui.NoPadding,
		Spacing: gui.SomeF(12),
		Content: []gui.View{
			chartView,
			line(),
			gui.Text(gui.TextCfg{
				Text:      "Code",
				TextStyle: gui.CurrentTheme().B3,
			}),
			w.Markdown(gui.MarkdownCfg{
				ID:      "code-" + id,
				Source:  source,
				Padding: gui.NoPadding,
				Style:   gui.DefaultMarkdownStyle(),
			}),
		},
	})
}

func line() gui.View {
	t := gui.CurrentTheme()
	return gui.Column(gui.ContainerCfg{
		Sizing:     gui.FillFit,
		Padding:    gui.SomeP(3, 0, 0, 0),
		SizeBorder: gui.NoBorder,
		Radius:     gui.NoRadius,
		Content: []gui.View{
			gui.Row(gui.ContainerCfg{
				Sizing:     gui.FillFit,
				Padding:    gui.NoPadding,
				SizeBorder: gui.NoBorder,
				Radius:     gui.NoRadius,
				Color:      t.ColorActive,
				Height:     1,
			}),
		},
	})
}
