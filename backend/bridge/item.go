package bridge

import (
	"UppSpar/backend"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"

	// mtheme "github.com/assholehoff/fyne-theme"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.SecondaryTappable = (*ItemLabel)(nil)
var _ fyne.Widget = (*ItemLabel)(nil)

/* ItemLabel is a widget for the item list view displaying a few text objects */
type ItemLabel struct {
	ttw.ToolTipWidget
	ItemID backend.ItemID

	cat, id, mfr, model, name string
}

func NewItemLabel() *ItemLabel {
	t := &ItemLabel{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *ItemLabel) Bind(id backend.ItemID) {
	id.Item().ItemIDString.AddListener(binding.NewDataListener(func() {
		s, err := id.Item().ItemIDString.Get()
		if err != nil {
			log.Println(err)
		}
		t.id = s
		t.Refresh()
	}))
	id.Item().Category.AddListener(binding.NewDataListener(func() {
		s, err := id.Item().Category.Get()
		if err != nil {
			log.Println(err)
		}
		t.cat = s
		t.Refresh()
	}))
	id.Item().Manufacturer.AddListener(binding.NewDataListener(func() {
		s, err := id.Item().Manufacturer.Get()
		if err != nil {
			log.Println(err)
		}
		t.mfr = s
		t.Refresh()
	}))
	id.Item().ModelName.AddListener(binding.NewDataListener(func() {
		s, err := id.Item().ModelName.Get()
		if err != nil {
			log.Println(err)
		}
		t.model = s
		t.Refresh()
	}))
	id.Item().Name.AddListener(binding.NewDataListener(func() {
		s, err := id.Item().Name.Get()
		if err != nil {
			log.Println(err)
		}
		t.name = s
		t.Refresh()
	}))
}

// func (t *ItemLabel) Unbind() {}

/* TappedSecondary implements fyne.SecondaryTappable. */
func (t *ItemLabel) TappedSecondary(*fyne.PointEvent) {
	panic("unimplemented")
}

/* CreateRenderer implements fyne.Widget. */
func (t *ItemLabel) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	th := t.Theme()
	vt := fyne.CurrentApp().Settings().ThemeVariant()

	id := canvas.NewText(t.id, th.Color(theme.ColorNameForeground, vt))
	id.TextSize = theme.Size(theme.SizeNameCaptionText)
	id.TextStyle = fyne.TextStyle{Bold: true}

	cat := canvas.NewText(t.cat, th.Color(theme.ColorNameForeground, vt))
	cat.TextSize = theme.Size(theme.SizeNameCaptionText)

	mfr := canvas.NewText(strings.ToUpper(t.mfr), th.Color(theme.ColorNameForeground, vt))
	mfr.TextSize = theme.Size(theme.SizeNameCaptionText)

	model := canvas.NewText(t.model, th.Color(theme.ColorNameForeground, vt))
	model.TextSize = theme.Size(theme.SizeNameCaptionText)

	name := canvas.NewText(t.name, th.Color(theme.ColorNameForeground, vt))

	r := &itemLabelRenderer{
		label: t,
		cat:   cat,
		id:    id,
		mfr:   mfr,
		model: model,
		name:  name,
	}

	r.applyTheme()
	return r
}

var _ fyne.WidgetRenderer = (*itemLabelRenderer)(nil)

type itemLabelRenderer struct {
	label *ItemLabel

	cat, id, name, mfr, model *canvas.Text
}

/* Destroy implements fyne.WidgetRenderer. */
func (r *itemLabelRenderer) Destroy() {
}

/* Layout implements fyne.WidgetRenderer. */
func (r *itemLabelRenderer) Layout(fyne.Size) {
	// TODO layout
	panic("unimplemented")
}

/* MinSize implements fyne.WidgetRenderer. */
func (r *itemLabelRenderer) MinSize() fyne.Size {
	hasMfr := r.label.mfr != ""
	hasModel := r.label.model != ""

	padding := r.label.Theme().Size(theme.SizeNamePadding)
	if !hasMfr && !hasModel {
		return fyne.NewSize(
			fyne.Max(r.name.MinSize().Width, r.id.MinSize().Width+r.cat.MinSize().Width)+padding,
			r.id.MinSize().Height+r.name.MinSize().Height+padding,
		)
	}

	min := fyne.NewSquareSize(padding)

	if hasMfr || hasModel {
		min = min.Add(fyne.NewSize(0, padding))
		if hasMfr {
		}
		if hasModel {
		}
	}

	return min
}

/* Objects implements fyne.WidgetRenderer. */
func (r *itemLabelRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.cat, r.id, r.mfr, r.model, r.name}
}

/* Refresh implements fyne.WidgetRenderer. */
func (r *itemLabelRenderer) Refresh() {
	r.cat.Text = r.label.cat
	r.cat.Refresh()

	r.id.Text = r.label.id
	r.id.Refresh()

	r.mfr.Text = r.label.mfr
	r.mfr.Refresh()

	r.model.Text = r.label.model
	r.model.Refresh()

	r.name.Text = r.label.name
	r.name.Refresh()

	r.applyTheme()
	r.Layout(r.label.Size())
	canvas.Refresh(r.label)
}

func (r *itemLabelRenderer) applyTheme() {
	th := r.label.Theme()
	vt := fyne.CurrentApp().Settings().ThemeVariant()

	if r.cat != nil {
		r.cat.TextSize = th.Size(theme.SizeNameCaptionText)
		r.cat.Color = th.Color(theme.ColorNameForeground, vt)
	}
	if r.id != nil {
		r.id.TextSize = th.Size(theme.SizeNameCaptionText)
		r.id.Color = th.Color(theme.ColorNameForeground, vt)
	}
	if r.mfr != nil {
		r.mfr.TextSize = th.Size(theme.SizeNameCaptionText)
		r.mfr.Color = th.Color(theme.ColorNameForeground, vt)
	}
	if r.model != nil {
		r.model.TextSize = th.Size(theme.SizeNameCaptionText)
		r.model.Color = th.Color(theme.ColorNameForeground, vt)
	}
	if r.name != nil {
		// r.name.TextSize = th.Size(theme.SizeNameText)
		r.name.Color = th.Color(theme.ColorNameForeground, vt)
	}
}
