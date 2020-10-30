package main

import (
    "fmt"
    "strconv"
    "io/ioutil"
    "strings"
    "time"
    "sort"
    "os"

    "github.com/gdamore/tcell/v2"
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
var table *tview.Table

func parseTasks(filename string) {
    raw, err := ioutil.ReadFile(filename)
    check(err)

    lines := strings.Split(string(raw), "\n")

    for _, line := range lines {
        if line == "" {
            continue
        }
        var task task
        fields := strings.Split(line,"\t")

        if (fields[0] == "X") {
            task.done = true
        }

        task.date, err = time.Parse("2006.01.02", fields[1])
        check(err)

        task.estimatedTime, err = strconv.Atoi(fields[2])
        check(err)

        task.difficulty, err = strconv.ParseFloat(fields[3], 64)
        check(err)

        task.desc = fields[4]

        tasks = append(tasks, task)
    }
}

func writeTasks() {
    f, err := os.Create("tasklist")
    check(err)
    defer f.Close()

    for _, task := range tasks {
        if task.done {
            f.WriteString("X")
        }
        f.WriteString(fmt.Sprintf("\t%s\t%d\t%.1f\t%s\n",
            task.date.Format("2006.01.02"),
            task.estimatedTime,
            task.difficulty,
            task.desc))
    }
}

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

func main() {
    parseTasks("tasklist")

    app := tview.NewApplication()


    form := tview.NewForm().
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

    flex := tview.NewFlex().
        SetDirection(tview.FlexRow).
        AddItem(table, 0, 1, false).
        AddItem(form, 7, 0, false)

    displayTable()
    if err := app.SetRoot(flex, true).SetFocus(table).Run(); err != nil {
        panic(err)
    }
}
