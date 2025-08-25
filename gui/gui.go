package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

type gui struct {
	items    *itemView
	journal  *journalView
	metadata *metadataView
	settings *settingsView
	wishlist *wishlistView
	tabs     *container.AppTabs
}

func (a *App) newGui() {
	a.gui = &gui{}
	a.gui.items = newItemView(a)
	a.gui.journal = newJournalView(a.backend)
	a.gui.metadata = newMetadataView(a.backend)
	a.gui.settings = newSettingsView(a.backend)
	a.gui.wishlist = newWishlistView(a.backend)
	a.newAppTabs()
}
func (a *App) newAppTabs() {
	a.gui.tabs = container.NewAppTabs(
		container.NewTabItemWithIcon(lang.L("Items"), theme.ListIcon(), a.gui.items.container),
		container.NewTabItemWithIcon(lang.L("Metadata"), theme.StorageIcon(), a.gui.metadata.tabs),
		container.NewTabItemWithIcon(lang.L("Journal"), theme.InfoIcon(), a.gui.journal.container),
		container.NewTabItemWithIcon(lang.L("Wishlist"), theme.MenuIcon(), a.gui.wishlist.container),
		container.NewTabItemWithIcon(lang.L("Settings"), theme.SettingsIcon(), a.gui.settings.container),
	)
	a.gui.tabs.SetTabLocation(container.TabLocationLeading)
}
