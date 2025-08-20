package bridge

import (
	"UppSpar/backend/journal"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	midget "github.com/assholehoff/fyne-midget"
	ttw "github.com/dweymouth/fyne-tooltip/widget"
)

var _ fyne.Widget = (*JournalEntry)(nil)

type JournalEntry struct {
	ttw.ToolTipWidget
	entry     *journal.Entry
	timestamp *midget.Label
	message   *midget.Label
}

func NewJournalEntry(e *journal.Entry) *JournalEntry {
	level := strings.ToUpper(lang.L(strings.ToLower(e.Level.String())))
	// event := e.Event.String()
	j := &JournalEntry{
		entry:     e,
		timestamp: midget.NewLabel("15:04:05", "02-01-2006", ""),
		message:   midget.NewLabel(e.Message, level, ""),
	}
	j.timestamp.SetTop()
	j.message.SetTop()
	j.ExtendBaseWidget(j)
	j.Update()
	return j
}

func (j *JournalEntry) Format(n int) {
	itemIdRegex := regexp.MustCompile(`<ItemId>(\d+(\.\d+)?)<\/ItemId>`)
	formatNumber := func(match string) string {
		subMatch := itemIdRegex.FindStringSubmatch(match)
		if len(subMatch) < 2 {
			return match
		}
		num, err := strconv.ParseInt(subMatch[1], 10, 64)
		if err != nil {
			return match
		}
		return fmt.Sprintf("%0*d", n, num)
	}

	result := itemIdRegex.ReplaceAllStringFunc(j.entry.Message, formatNumber)
	j.message.Text = regexp.MustCompile(`<ItemId>|<\/ItemId>`).ReplaceAllString(result, "")
	j.message.Text = regexp.MustCompile(`<CatId>|<\/CatId>`).ReplaceAllString(result, "")
	j.message.Refresh()
}

func (j *JournalEntry) Bind(e *journal.Entry) {
	j.entry = e
	j.Update()
}

/* Update labels */
func (j *JournalEntry) Update() {
	day := j.entry.Time.Format("02-01-2006")
	tid := j.entry.Time.Format("15:04:05")
	j.timestamp.SetSubtext(day)
	j.timestamp.SetText(tid)
	j.timestamp.Refresh()
	j.message.SetSubtext(strings.ToUpper(lang.L(strings.ToLower(j.entry.Level.String()))))
	j.message.SetText(j.entry.Message)
	j.message.Refresh()
}

/* CreateRenderer implements fyne.Widget. */
func (j *JournalEntry) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(j.timestamp, j.message)
	return widget.NewSimpleRenderer(c)
}

var _ fyne.WidgetRenderer = (*journalEntryRenderer)(nil)

type journalEntryRenderer struct {
	j *JournalEntry
}

/* Destroy implements fyne.WidgetRenderer. */
func (r *journalEntryRenderer) Destroy() {
}

/* Layout implements fyne.WidgetRenderer. */
func (r *journalEntryRenderer) Layout(s fyne.Size) {
	tsize := r.j.timestamp.MinSize()
	tsize.Width += theme.InnerPadding() / 2
	msize := s
	msize.Width -= theme.InnerPadding() / 2
	msize.Width -= tsize.Width

	pos := fyne.NewSquareOffsetPos(theme.InnerPadding() / 2)
	r.j.timestamp.Resize(tsize)
	r.j.timestamp.Move(pos)

	pos.X += theme.InnerPadding() / 2
	pos.X += r.j.timestamp.Size().Width
	r.j.message.Resize(msize)
	r.j.message.Move(pos)
}

/* MinSize implements fyne.WidgetRenderer. */
func (r *journalEntryRenderer) MinSize() fyne.Size {
	size := r.j.timestamp.MinSize()
	size.Width += theme.InnerPadding()
	size.Width += r.j.message.MinSize().Width
	return size
}

/* Objects implements fyne.WidgetRenderer. */
func (r *journalEntryRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.j.timestamp, r.j.message}
}

/* Refresh implements fyne.WidgetRenderer. */
func (r *journalEntryRenderer) Refresh() {
	r.j.Update()
}
