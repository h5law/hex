package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	path := os.Args[1]
	fileName, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("[hex] error getting file path")
		os.Exit(1)
	}

	addresses, hex, text := dumpFile(fileName)

	app := tview.NewApplication()
	app.SetInputCapture(appKeyEvents(app))

	addyBox := tview.NewList().SetWrapAround(false)
	addyBox.SetBorder(true).
		SetTitle("Address").
		SetTitleAlign(tview.AlignLeft).
		SetBorderPadding(1, 1, 5, 5)
	hexBox := tview.NewList().SetWrapAround(false)
	hexBox.SetBorder(true).
		SetTitle("Hex").
		SetBorderPadding(1, 1, 5, 5)
	textBox := tview.NewList().SetWrapAround(false)
	textBox.SetBorder(true).
		SetTitle("Text").
		SetTitleAlign(tview.AlignRight).
		SetBorderPadding(1, 1, 5, 5)

	hexBox.SetChangedFunc(func(index int, _, _ string, _ rune) {
		addyBox.SetCurrentItem(index)
		textBox.SetCurrentItem(index)
	})

	hexBox.SetInputCapture(listKeyEvents(hexBox))

	flex := tview.NewFlex().
		AddItem(addyBox, 0, 1, false).
		AddItem(hexBox, 0, 2, true).
		AddItem(textBox, 0, 1, false)
	flex.SetBorderPadding(1, 1, 5, 5)

	for i := 0; i < len(addresses); i++ {
		addyBox.AddItem(addresses[i], "", 0, nil)
		hexBox.AddItem(hex[i], "", 0, nil)
		textBox.AddItem(text[i], "", 0, nil)
	}

	if err := app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

func appKeyEvents(app *tview.Application) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
		}
		return event
	}
}

func listKeyEvents(list *tview.List) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'j':
			list.SetCurrentItem(list.GetCurrentItem() + 1)
		case 'k':
			if list.GetCurrentItem() != 0 {
				list.SetCurrentItem(list.GetCurrentItem() - 1)
			}
		case '[':
			if list.GetCurrentItem() > 20 {
				list.SetCurrentItem(list.GetCurrentItem() - 20)
			} else {
				list.SetCurrentItem(0)
			}
		case ']':
			if list.GetCurrentItem() < list.GetItemCount()-1 {
				list.SetCurrentItem(list.GetCurrentItem() + 20)
			} else {
				list.SetCurrentItem(list.GetItemCount() - 1)
			}
		case 'G':
			list.SetCurrentItem(list.GetItemCount() - 1)
		case 'g':
			list.SetCurrentItem(0)
		}
		return event
	}
}
