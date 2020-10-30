package main

import (
    "fmt"
    "strconv"
    "io/ioutil"
    "strings"
    "time"
    "os"
)

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
