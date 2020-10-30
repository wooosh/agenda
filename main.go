package main

import (
    "time"

    "github.com/rivo/tview"
)

type task struct {
    done bool
    date time.Time
    estimatedTime int
    difficulty float64
    desc string
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

var tasks []task
var app *tview.Application

func main() {
    parseTasks("tasklist")

    app = tview.NewApplication()

    createForm()
    createTable()

    flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(table, 0, 1, false).
        AddItem(form, 7, 0, false)

    displayTable()
    if err := app.SetRoot(flex, true).SetFocus(table).Run(); err != nil {
        panic(err)
    }
}
