package bridge

import (
	"slices"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

type Checks map[string]*ttw.Check

func (c Checks) Bind(m map[string]binding.Bool) {
	for key, val := range m {
		if c[key] == nil {
			c[key] = ttw.NewCheck(key, func(b bool) {})
		}
		c[key].Bind(val)
	}
}
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
func (c Checks) Uncheck() {
	for _, val := range c {
		val.SetChecked(false)
	}
}

type Entries map[string]*midget.Entry

func (e Entries) Bind(m map[string]binding.String) {
	for key, val := range m {
		if e[key] == nil {
			e[key] = midget.NewEntry()
		}
		e[key].Bind(val)
	}
}
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

type RadioConfig struct {
	Options  binding.StringList
	Function binding.Untyped
}

type Radios map[string]*widget.RadioGroup

func (r Radios) Bind(m map[string]RadioConfig) {
	for key, val := range m {
		if r[key] == nil {
			r[key] = widget.NewRadioGroup([]string{}, func(string) {})
		}
		val.Options.AddListener(binding.NewDataListener(func() {
			opt, _ := val.Options.Get()
			r[key].Options = opt
			r[key].Refresh()
		}))
		val.Function.AddListener(binding.NewDataListener(func() {
			fn, _ := val.Function.Get()
			r[key].OnChanged = fn.(func(string))
			r[key].Refresh()
		}))
	}
}
func (r Radios) Setup(o map[string][]string, f map[string]func(string)) {
	for key := range o {
		if r[key] == nil {
			r[key] = widget.NewRadioGroup([]string{}, func(string) {})
		}
		r[key].Options = o[key]
		r[key].OnChanged = f[key]
		r[key].Refresh()
	}
}
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
func (r Radios) Uncheck() {
	for _, val := range r {
		val.SetSelected("")
	}
}

type SelectConfig struct {
	Options  binding.StringList
	Function binding.Untyped
}

type Selects map[string]*ttw.Select

func (s Selects) Bind(m map[string]SelectConfig) {
	for key, val := range m {
		if s[key] == nil {
			s[key] = ttw.NewSelect([]string{}, func(string) {})
		}
		val.Options.AddListener(binding.NewDataListener(func() {
			opt, _ := val.Options.Get()
			s[key].SetOptions(opt)
		}))
		val.Function.AddListener(binding.NewDataListener(func() {
			fun, _ := val.Function.Get()
			s[key].OnChanged = fun.(func(string))
			s[key].Refresh()
		}))
	}
}
func (s Selects) Setup(o map[string][]string, f map[string]func(string)) {
	for key, val := range o {
		if s[key] == nil {
			s[key] = ttw.NewSelect([]string{}, func(string) {})
		}
		s[key].SetOptions(val)
		s[key].OnChanged = f[key]
		s[key].Refresh()
	}
}
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

func Sieve(m map[string]binding.String, list []string) map[string]binding.String {
	for key := range m {
		if !slices.Contains(list, key) {
			delete(m, key)
		}
	}
	return m
}
