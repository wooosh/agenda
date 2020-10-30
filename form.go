package main

import (
    "strconv"
    "time"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var form *tview.Form

func createForm() {
    form = tview.NewForm().
        SetFieldBackgroundColor(tcell.ColorBlack).
        SetItemPadding(0).
        AddCheckbox("Done", false, nil).
        AddInputField("Due", "", 0, nil, nil).
        AddInputField("Est. Time", "", 0, nil, nil).
        AddInputField("Difficulty", "", 0, nil, nil).
        AddInputField("Description", "", 0, nil, nil)

    form.SetCancelFunc(func() {
        row, _ := table.GetSelection()
        var task task
        task.done = form.GetFormItemByLabel("Done").(*tview.Checkbox).
            IsChecked()

        var err error
        task.date, err = time.Parse(
            "2006.01.02",
            form.GetFormItemByLabel("Due").(*tview.InputField).GetText())
        check(err)

        task.estimatedTime, err = strconv.Atoi(form.GetFormItemByLabel("Est. Time").(*tview.InputField).GetText())
        check(err)

        task.difficulty, err = strconv.ParseFloat(form.GetFormItemByLabel("Difficulty").(*tview.InputField).GetText(), 64)
        check(err)

        task.desc = form.GetFormItemByLabel("Description").(*tview.InputField).GetText()

        tasks[row-1] = task
        displayTable()
        writeTasks()
        app.SetFocus(table)
    })
}
