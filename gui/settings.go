package gui

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type settingsView struct {
	container *fyne.Container
}

func newSettingsView(b *backend.Backend) *settingsView {
	SettingsTitle := `# ` + lang.X("settings.title", "settings.title")

	ItemIDText := lang.X("settings.itemid.text", "settings.itemid.text")
	ItemIDSubtext := lang.X("settings.itemid.subtext", "settings.itemid.subtext")
	ItemIDTooltip := lang.X("settings.itemid.tooltip", "settings.itemid.tooltip")

	ResumeText := lang.X("settings.resume.text", "settings.resume.text")
	// ResumeSubtext := lang.X("settings.resume.subtext", "settings.resume.subtext")
	// ResumeTooltip := lang.X("settings.resume.tooltip", "settings.resume.tooltip")

	f := container.New(layout.NewFormLayout(),
		layout.NewSpacer(),
		ttw.NewCheckWithData(ResumeText, b.Settings.ResumeLastSession),
		midget.NewLabel(ItemIDText, ItemIDSubtext, ItemIDTooltip),
		midget.NewIntEntryWithData(b.Settings.ItemIDWidth),
	)

	f.Objects[1].(*ttw.Check).Disable() // TODO fix the crash before enabling !

	g := container.NewBorder(widget.NewRichTextFromMarkdown(SettingsTitle), nil, nil, nil, f)

	return &settingsView{
		container: g,
	}
}
