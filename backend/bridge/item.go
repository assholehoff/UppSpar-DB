package bridge

import (
	"UppSpar/backend"

	"fyne.io/fyne/v2"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.Widget = (*ItemLabel)(nil)

/* ItemLabel is a widget for the item list view displaying a few text objects */
type ItemLabel struct {
	ttw.ToolTipWidget
	item *backend.Item
}

func NewItemLabel() *ItemLabel {
	t := &ItemLabel{}
	t.ExtendBaseWidget(t)
	return t
}

func (t *ItemLabel) Bind(id backend.ItemID) {}
func (t *ItemLabel) Unbind()                {}

/* CreateRenderer implements fyne.Widget. */
func (t *ItemLabel) CreateRenderer() fyne.WidgetRenderer {
	panic("unimplemented")
}

var _ fyne.WidgetRenderer = (*itemLabelRenderer)(nil)

type itemLabelRenderer struct {
	label *ItemLabel
}

/* Destroy implements fyne.WidgetRenderer. */
func (i *itemLabelRenderer) Destroy() {
}

/* Layout implements fyne.WidgetRenderer. */
func (i *itemLabelRenderer) Layout(fyne.Size) {
	panic("unimplemented")
}

/* MinSize implements fyne.WidgetRenderer. */
func (i *itemLabelRenderer) MinSize() fyne.Size {
	panic("unimplemented")
}

/* Objects implements fyne.WidgetRenderer. */
func (i *itemLabelRenderer) Objects() []fyne.CanvasObject {
	panic("unimplemented")
}

/* Refresh implements fyne.WidgetRenderer. */
func (i *itemLabelRenderer) Refresh() {
	panic("unimplemented")
}

var _ fyne.Widget = (*ItemCondition)(nil)

/* ItemCondition is a widget box with checkboxes, choices and entries specific to an item category or unique to an item */
type ItemCondition struct {
	ttw.ToolTipWidget
	item *backend.Item
}

func NewItemCondition() *ItemCondition {
	return &ItemCondition{}
}

func (t *ItemCondition) Bind(id backend.ItemID) {}
func (t *ItemCondition) Unbind()                {}

/* CreateRenderer implements fyne.Widget. */
func (t *ItemCondition) CreateRenderer() fyne.WidgetRenderer {
	panic("unimplemented")
}

var _ fyne.Widget = (*ItemProperties)(nil)

/* ItemProperties is a widget box with checkboxes, choices and entries specific to an item category or unique to an item */
type ItemProperties struct {
	ttw.ToolTipWidget
	item *backend.Item
}

func NewItemProperties() *ItemProperties {
	return &ItemProperties{}
}

func (t *ItemProperties) Bind(id backend.ItemID) {}
func (t *ItemProperties) Unbind()                {}

/* CreateRenderer implements fyne.Widget. */
func (t *ItemProperties) CreateRenderer() fyne.WidgetRenderer {
	panic("unimplemented")
}
