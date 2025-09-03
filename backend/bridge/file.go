package bridge

import (
	"UppSpar/backend"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
)

func NewSaveFileDialog(b *backend.Backend, w fyne.Window) *dialog.FileDialog {
	d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if writer != nil {
			b.Items.ExportExcel(writer.URI().Path())
		} else {
			return
		}
	}, w)
	d.Resize(fyne.NewSize(900, 600))
	d.SetTitleText(lang.X("dialog.save.excel.title", "dialog.save.excel.title"))
	d.SetConfirmText(lang.L("Export"))
	d.SetDismissText(lang.L("Close"))
	d.SetFileName("UppSpar-" + time.Now().Format("20060102-150405") + ".xlsx")
	d.SetFilter(storage.NewMimeTypeFileFilter([]string{"application/excel"}))
	return d
}
