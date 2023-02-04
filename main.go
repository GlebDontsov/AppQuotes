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

func Remove(list []fyne.CanvasObject, item fyne.CanvasObject) []fyne.CanvasObject {
	for i, value := range list {
		if value == item {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func CreateCardF(key string, value map[string]string, mapWidget map[string]fyne.CanvasObject,
	box *[]fyne.CanvasObject, box1 *[]fyne.CanvasObject,
	icStarYel *fyne.Resource, icShare *fyne.Resource,
	icDel *fyne.Resource, icStar *fyne.Resource) fyne.CanvasObject {

	quote := widget.NewLabel(value["text"])
	author := widget.NewLabel(value["author"])

	quote.Wrapping = fyne.TextWrapBreak
	author.Wrapping = fyne.TextWrapBreak

	buttonStar := widget.NewButton("", func() {
		mapQuotes[key] = value
		delete(mapStar, key)
		SaveQuotes("quotes.json", mapQuotes)
		SaveQuotes("star_quotes.json", mapStar)
		*box = Remove(*box, mapWidget[key])
		*box1 = append(*box1, CreateCardM(key, value, mapWidget, box1, box, icStarYel, icShare, icDel, icStar))
	})

	buttonShare := widget.NewButton("", func() {})

	buttonStar.SetIcon(*icStarYel)
	buttonShare.SetIcon(*icShare)

	boxBtn := container.NewHBox(buttonStar, buttonShare)
	text := container.NewVBox(quote, author, boxBtn)
	card := widget.NewCard(
		"",
		"",
		text,
	)
	mapWidget[key] = card
	return card
}

func CreateCardM(key string, value map[string]string, mapWidget map[string]fyne.CanvasObject,
	box *[]fyne.CanvasObject, box1 *[]fyne.CanvasObject,
	icStarYel *fyne.Resource, icShare *fyne.Resource,
	icDel *fyne.Resource, icStar *fyne.Resource) fyne.CanvasObject {

	quote := widget.NewLabel(value["text"])
	author := widget.NewLabel(value["author"])

	quote.Wrapping = fyne.TextWrapBreak
	author.Wrapping = fyne.TextWrapBreak

	buttonStar := widget.NewButton("", func() {
		mapStar[key] = value
		delete(mapQuotes, key)
		SaveQuotes("quotes.json", mapQuotes)
		SaveQuotes("star_quotes.json", mapStar)
		*box = Remove(*box, mapWidget[key])
		*box1 = append(*box1, CreateCardF(key, value, mapWidget, box1, box, icStarYel, icShare, icDel, icStar))
	})

	buttonDel := widget.NewButton("", func() {
		delete(mapQuotes, key)
		SaveQuotes("quotes.json", mapQuotes)
		if _, ok := mapStar[key]; !ok {
			delete(mapStar, key)
			SaveQuotes("star_quotes.json", mapStar)
		}
		*box1 = Remove(*box1, mapWidget[key])
		delete(mapWidget, key)
	})

	buttonShare := widget.NewButton("", func() {})

	buttonStar.SetIcon(*icStar)
	buttonShare.SetIcon(*icShare)
	buttonDel.SetIcon(*icDel)

	boxBtn := container.NewHBox(buttonStar, buttonShare, buttonDel)
	text := container.NewVBox(quote, author, boxBtn)
	card := widget.NewCard(
		"",
		"",
		text,
	)
	mapWidget[key] = card
	return card
}

func main() {
	a := app.New()
	//a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("MyApp")
	w.Resize(fyne.NewSize(380, 670))

	icStar, _ := fyne.LoadResourceFromPath("icon/star.png")
	icStarYel, _ := fyne.LoadResourceFromPath("icon/star_yel.png")
	icShare, _ := fyne.LoadResourceFromPath("icon/share.png")
	icDelete, _ := fyne.LoadResourceFromPath("icon/delete.png")

	box := container.NewVBox()
	box1 := container.NewVBox()

	mapWidget := map[string]fyne.CanvasObject{}

	for i, v := range mapStar {
		card := CreateCardF(i, v, mapWidget, &box.Objects, &box1.Objects, &icStarYel, &icShare, &icDelete, &icStar)
		box.Add(card)
	}

	for i, v := range mapQuotes {
		card := CreateCardM(i, v, mapWidget, &box1.Objects, &box.Objects, &icStarYel, &icShare, &icDelete, &icStar)
		box1.Add(card)
	}

	QuoteField := widget.NewMultiLineEntry()
	AuthorField := widget.NewEntry()
	QuoteField.SetPlaceHolder("Quote...")
	QuoteField.SetMinRowsVisible(10)
	AuthorField.SetPlaceHolder("Author...")
	QuoteField.Wrapping = fyne.TextWrapBreak
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
				box1.Add(CreateCardM(id, mapQuotes[id], mapWidget, &box1.Objects,
					&box.Objects, &icStarYel, &icShare, &icDelete, &icStar))
			}
		}),
	)
	box2 := container.NewVBox(QuoteField, AuthorField, btnSave)

	tabs := container.NewAppTabs(
		container.NewTabItem("My quotes", container.NewVScroll(box1)),
		container.NewTabItem("New quote", box2),
		container.NewTabItem("Favourites", container.NewVScroll(box)),
		//container.NewTabItem("Book", widget.NewLabel("Books")),
	)
	tabs.SetTabLocation(container.TabLocationBottom)

	w.SetContent(tabs)
	w.ShowAndRun()
	a.Run()
}
