package bridge

import (
	"fyne.io/fyne/v2/widget"
	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type Checks map[string]*ttw.Check

func (c Checks) Disable() {
	for _, val := range c {
		val.Disable()
	}
}
func (c Checks) Enable() {
	for _, val := range c {
		val.Enable()
	}
}
func (c Checks) Unbind() {
	for _, val := range c {
		val.Unbind()
	}
}
func (c Checks) Hide() {
	for _, val := range c {
		val.Hide()
	}
}
func (c Checks) Show() {
	for _, val := range c {
		val.Show()
	}
}

type Entries map[string]*midget.Entry

func (e Entries) Clear() {
	for _, val := range e {
		val.SetText("")
	}
}
func (e Entries) Disable() {
	for _, val := range e {
		val.Disable()
	}
}
func (e Entries) Enable() {
	for _, val := range e {
		val.Enable()
	}
}
func (e Entries) Unbind() {
	for _, val := range e {
		val.Unbind()
	}
}
func (e Entries) Hide() {
	for _, val := range e {
		val.Hide()
	}
}
func (e Entries) Show() {
	for _, val := range e {
		val.Show()
	}
}

type Labels map[string]*ttw.Label

func (l Labels) Clear() {
	for _, val := range l {
		val.SetText("")
	}
}
func (l Labels) Unbind() {
	for _, val := range l {
		val.Unbind()
	}
}
func (l Labels) Hide() {
	for _, val := range l {
		val.Hide()
	}
}
func (l Labels) Show() {
	for _, val := range l {
		val.Show()
	}
}

type Radios map[string]*widget.RadioGroup

func (r Radios) Disable() {
	for _, val := range r {
		val.Disable()
	}
}
func (r Radios) Enable() {
	for _, val := range r {
		val.Enable()
	}
}
func (r Radios) Hide() {
	for _, val := range r {
		val.Hide()
	}
}
func (r Radios) Show() {
	for _, val := range r {
		val.Show()
	}
}

type Selects map[string]*ttw.Select

func (s Selects) Clear() {
	for _, val := range s {
		val.ClearSelected()
	}
}
func (s Selects) Disable() {
	for _, val := range s {
		val.Disable()
	}
}
func (s Selects) Enable() {
	for _, val := range s {
		val.Enable()
	}
}
func (s Selects) Unbind() {
	for _, val := range s {
		val.Unbind()
	}
}
func (s Selects) Hide() {
	for _, val := range s {
		val.Hide()
	}
}
func (s Selects) Show() {
	for _, val := range s {
		val.Show()
	}
}
