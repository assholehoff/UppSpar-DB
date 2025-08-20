package gui

import (
	"UppSpar/backend"
	"embed"
	"errors"
	"fmt"
	"log"
	"mime"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	midget "github.com/assholehoff/fyne-midget"
	theme "github.com/assholehoff/fyne-theme"
	tooltip "github.com/dweymouth/fyne-tooltip"
)

type App struct {
	app     fyne.App
	backend *backend.Backend
	gui     *gui
	window  fyne.Window
}

func NewApp(translation embed.FS) *App {
	mime.AddExtensionType(".db", "application/sqlite")
	mime.AddExtensionType(".xlsx", "application/excel")

	a := &App{
		app: app.NewWithID("se.antondahlen.uppspar-db"),
	}

	a.addTranslationFS(translation)
	a.app.Settings().SetTheme(&theme.Tight{})

	a.window = a.app.NewWindow("UppSpar DB")
	a.window.Resize(fyne.NewSize(800, 600))

	a.selectDatabase()
	return a
}

func (a *App) Run() {
	a.window.Show()
	a.app.Run()
}

func (a *App) addTranslationFS(t embed.FS) {
	err := lang.AddTranslationsFS(t, "translation")
	if err != nil {
		log.Println(fmt.Errorf("add translations error: %w", err))
	}
}
func (a *App) selectDatabase() {
	var file string
	var d *dialog.CustomDialog
	var n, o *dialog.FileDialog

	newFileSaveDialog := func() *dialog.FileDialog {
		n = dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if writer != nil {
				file = writer.URI().Path()
				a.openDatabase(file)
			} else {
				d.Show()
			}
		}, a.window)
		n.Resize(fyne.NewSize(640, 480))
		n.SetTitleText(lang.L("Create new database"))
		n.SetConfirmText(lang.L("Create"))
		n.SetDismissText(lang.L("Close"))
		n.SetFileName("uppspar.db")
		n.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/sqlite"}))
		return n
	}
	newFileOpenDialog := func() *dialog.FileDialog {
		o = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				file = reader.URI().Path()
				a.openDatabase(file)
			} else {
				d.Show()
			}
		}, a.window)
		o.Resize(fyne.NewSize(640, 480))
		o.SetTitleText(lang.L("Open database"))
		o.SetConfirmText(lang.L("Open"))
		o.SetDismissText(lang.L("Close"))
		o.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/sqlite"}))
		return o
	}

	d = dialog.NewCustomWithoutButtons(
		"UppSpar DB",
		midget.NewLabel(
			lang.L("Open existing or create new database?"),
			"",
			"",
		),
		a.window,
	)
	d.SetButtons([]fyne.CanvasObject{
		widget.NewButton(lang.L("New"), func() {
			d.Hide()
			n = newFileSaveDialog()
			n.Show()
		}),
		widget.NewButton(lang.L("Open"), func() {
			d.Hide()
			o = newFileOpenDialog()
			o.Show()
		}),
	})

	file = a.app.Preferences().String("file")
	_, err := os.Stat(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			d.Show()
		} else {
			log.Printf("err != nil; but not os.ErrNotExist")
			/*
			 * Weird error!
			 * // TODO fix dialog and do stuff
			 */
			panic(err)
		}
	} else if !a.app.Preferences().BoolWithFallback("resume", false) {
		d.Show()
	} else {
		// TODO what is happening here when there are Items in the DB??
		a.openDatabase(file)
		// d.Show()
	}
}

func (a *App) openDatabase(file string) {
	log.Printf("openDatabase(%s)", file)
	var err error
	a.backend, err = backend.NewBackend(file)
	if err != nil {
		panic(err)
	}
	a.app.Preferences().SetString("file", file)

	a.newGui()

	canvas := a.window.Canvas()
	content := a.gui.tabs // TODO <-- is this the culprit?
	tipped := tooltip.AddWindowToolTipLayer(content, canvas)

	a.window.Resize(fyne.NewSize(1200, 900))
	a.window.SetMaster()
	a.window.Show()
	a.window.SetContent(tipped) // TODO <-- this blows up when there are Items in the DB, *why* ???
}
