package main

import (
    "fmt"
    "strconv"
    "sort"

    "github.com/gdamore/tcell/v2"
    "github.com/rivo/tview"
)

var table *tview.Table

func displayTable() {
    table.Clear()

    headers := []string{"DUE", "EST TIME", "DIFFICULTY", "DESCRIPTION"}
    for i, header := range headers {
        table.SetCell(0, i+1, tview.NewTableCell(header).
            SetTextColor(tcell.ColorYellow).
            SetSelectable(false))
    }

    table.GetCell(0, 4).SetExpansion(1)

    sort.Slice(tasks, func(i, j int) bool {
        return tasks[j].done
    })

    for i, task := range tasks {
        if (task.done) {
            table.SetCellSimple(i+1, 0, "X")
        }
        table.SetCellSimple(i+1, 1, task.date.Format("2006.01.02"))
        table.SetCellSimple(i+1, 2, strconv.Itoa(task.estimatedTime))
        table.SetCellSimple(i+1, 3, fmt.Sprintf("%.1f", task.difficulty))
        table.SetCellSimple(i+1, 4, task.desc)
    }
}

func createTable() {
    table = tview.NewTable().
        SetSeparator(tview.Borders.Vertical).
        SetSelectable(true, false).
        SetFixed(1,5).
        SetDoneFunc(func(key tcell.Key) {
            if key == tcell.KeyEscape {
                app.Stop()
            }
        }).SetSelectionChangedFunc(func(row, column int) {
            // do nothing if header is selected
            if row == 0 {
                return
            }
            task := tasks[row-1] // -1 to account for header
            form.GetFormItemByLabel("Done").(*tview.Checkbox).
                SetChecked(task.done)
            form.GetFormItemByLabel("Due").(*tview.InputField).
                SetText(task.date.Format("2006.01.02"))
            form.GetFormItemByLabel("Est. Time").(*tview.InputField).
                SetText(strconv.Itoa(task.estimatedTime))
            form.GetFormItemByLabel("Difficulty").(*tview.InputField).
                SetText(fmt.Sprintf("%.1f", task.difficulty))
            form.GetFormItemByLabel("Description").(*tview.InputField).
                SetText(task.desc)
        }).SetSelectedFunc(func(row, column int) {
            app.SetFocus(form)
        })
}
