package main

import (
	"UppSpar/gui"
	"embed"
	"log"

	"fyne.io/fyne/v2/lang"
)

//go:embed translation
var translations embed.FS

func main() {
	err := lang.AddTranslationsFS(translations, "translation")
	if err != nil {
		log.Println(err)
	}

	a := gui.NewApp(translations)
	a.Run()
}
