package main

import (
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"image/color"
	"io/ioutil"
	"os"
)

func GetQuotes(path string) map[string]map[string]string {
	var jsonMap map[string]map[string]string
	plan, _ := ioutil.ReadFile(path)
	err := json.Unmarshal(plan, &jsonMap)
	if err != nil {
		jsonMap = map[string]map[string]string{}
	}
	return jsonMap
}

func SaveQuotes(path string, mapQ map[string]map[string]string) {
	jsonString, _ := json.Marshal(mapQ)
	f, _ := os.Create(path)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic("Close error")
		}
	}(f)
	_, err := f.Write(jsonString)
	if err != nil {
		panic("Write error")
	}
}

var mapQuotes = GetQuotes("quotes.json")
var mapStar = GetQuotes("star_quotes.json")

func CreateNewQuote() fyne.CanvasObject {
	QuoteField := widget.NewMultiLineEntry()
	QuoteField.SetPlaceHolder("Quote...")
	QuoteField.Wrapping = fyne.TextWrapBreak
	QuoteField.SetMinRowsVisible(10)

	AuthorField := widget.NewEntry()
	AuthorField.SetPlaceHolder("Author...")

	btnSave := container.New(
		layout.NewMaxLayout(),
		canvas.NewRectangle(color.White),
		widget.NewButton("Save", func() {
			if QuoteField.Text != "" && AuthorField.Text != "" {
				id := uuid.New().String()
				mapQuotes[id] = map[string]string{"author": AuthorField.Text, "text": QuoteField.Text}
				QuoteField.Text = ""
				AuthorField.Text = ""
				SaveQuotes("quotes.json", mapQuotes)
			}
		}),
	)

	box := container.NewVBox(QuoteField, AuthorField, btnSave)
	return box
}

func CreateMyQuote() fyne.CanvasObject {
	box := container.NewVBox()
	icStar, _ := fyne.LoadResourceFromPath("icon/star.png")
	icShare, _ := fyne.LoadResourceFromPath("icon/share.png")
	icDelete, _ := fyne.LoadResourceFromPath("icon/delete.png")
	//icStarBlack, _ := fyne.LoadResourceFromPath("icon/star_yel.png")

	for i, v := range mapQuotes {
		key := i
		value := v
		quote := widget.NewLabel(v["text"])
		author := widget.NewLabel(v["author"])

		quote.Wrapping = fyne.TextWrapBreak
		author.Wrapping = fyne.TextWrapBreak

		buttonStar := widget.NewButton("", func() {
			mapStar[key] = value
			SaveQuotes("star_quotes.json", mapStar)
			delete(mapQuotes, key)
			SaveQuotes("quotes.json", mapQuotes)
		})

		buttonShare := widget.NewButton("", func() {})

		buttonDel := widget.NewButton("", func() {
			delete(mapQuotes, key)
			SaveQuotes("quotes.json", mapQuotes)
			if _, ok := mapStar[key]; !ok {
				delete(mapStar, key)
				SaveQuotes("star_quotes.json", mapStar)
			}
		})

		buttonStar.SetIcon(icStar)
		buttonShare.SetIcon(icShare)
		buttonDel.SetIcon(icDelete)

		boxBtn := container.NewHBox(buttonStar, buttonShare, buttonDel)
		text := container.NewVBox(quote, author, boxBtn)
		card := widget.NewCard(
			"",
			"",
			text,
		)

		box.Add(card)
	}
	return container.NewVScroll(box)
}

func CreateFavourites() fyne.CanvasObject {
	box := container.NewVBox()
	icStar, _ := fyne.LoadResourceFromPath("icon/star_yel.png")
	icShare, _ := fyne.LoadResourceFromPath("icon/share.png")
	//icStarBlack, _ := fyne.LoadResourceFromPath("icon/star_yel.png")

	for i, v := range mapStar {
		key := i
		value := v
		quote := widget.NewLabel(v["text"])
		author := widget.NewLabel(v["author"])

		quote.Wrapping = fyne.TextWrapBreak
		author.Wrapping = fyne.TextWrapBreak

		buttonStar := widget.NewButton("", func() {
			mapQuotes[key] = value
			delete(mapStar, key)
			SaveQuotes("quotes.json", mapQuotes)
			SaveQuotes("star_quotes.json", mapStar)
		})
		buttonShare := widget.NewButton("", func() {})

		buttonStar.SetIcon(icStar)
		buttonShare.SetIcon(icShare)

		boxBtn := container.NewHBox(buttonStar, buttonShare)
		text := container.NewVBox(quote, author, boxBtn)
		card := widget.NewCard(
			"",
			"",
			text,
		)

		box.Add(card)
	}
	return container.NewVScroll(box)
}

func main() {
	a := app.New()
	//a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("MyApp")
	w.Resize(fyne.NewSize(380, 670))

	formNewQuote := CreateNewQuote()
	formMyQuotes := CreateMyQuote()
	formFavourites := CreateFavourites()

	tabs := container.NewAppTabs(
		container.NewTabItem("My quotes", formMyQuotes),
		container.NewTabItem("New quote", formNewQuote),
		container.NewTabItem("Favourites", formFavourites),
		//container.NewTabItem("Book", widget.NewLabel("Books")),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	w.SetContent(tabs)
	w.ShowAndRun()
	a.Run()
}
