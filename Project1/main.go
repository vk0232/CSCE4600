package main
//package scheduling

import (
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
    "sort"
    "github.com/olekukonko/tablewriter"
)

type (
    Process struct {
        ProcessID     int64
        Name string
        ArrivalTime   int64
        BurstDuration int64
        Priority      int64
        Completed bool
        Waiting bool
        TurnaroundTime int64
        WaitingTime int64
    }
    TimeSlice struct {
        PID   int64
        Start int64
        Stop  int64
    }
)


/*
type Process struct {
    Name   string
    Burst  int
    Arrival int
    Priority int
    Completed bool
    Turnaround int
    Waiting int
}
*/

func main() {
    
    // CLI args
    f, closeFile, err := openProcessingFile(os.Args...)
    if err != nil {
        log.Fatal(err)
    }
    defer closeFile()

    // Load and parse processes
    processes, err := loadProcesses(f)
    if err != nil {
        log.Fatal(err)
    }

    // First-come, first-serve scheduling
    FCFSSchedule(os.Stdout, "First-come, first-serve", processes)

    SJFSchedule(os.Stdout, "Shortest-job-first", processes)

    SJFPrioritySchedule(os.Stdout, "Priority", processes)

    RRSchedule(os.Stdout, "Round-robin", processes, 10)
}

func openProcessingFile(args ...string) (*os.File, func(), error) {
    if len(args) != 2 {
        return nil, nil, fmt.Errorf("%w: must give a scheduling file to process", ErrInvalidArgs)
    }
    // Read in CSV process CSV file
    f, err := os.Open(args[1])
    if err != nil {
        return nil, nil, fmt.Errorf("%v: error opening scheduling file", err)
    }
    closeFn := func() {
        if err := f.Close(); err != nil {
            log.Fatalf("%v: error closing scheduling file", err)
        }
    }

    return f, closeFn, nil
}




//region Schedulers

// FCFSSchedule outputs a schedule of processes in a GANTT chart and a table of timing given:
// • an output writer
// • a title for the chart
// • a slice of processes
func FCFSSchedule(w io.Writer, title string, processes []Process) {
    var (
        serviceTime     int64
        totalWait       float64
        totalTurnaround float64
        lastCompletion  float64
        waitingTime     int64
        schedule        = make([][]string, len(processes))
        gantt           = make([]TimeSlice, 0)
    )
    for i := range processes {
        if processes[i].ArrivalTime > 0 {
            waitingTime = serviceTime - processes[i].ArrivalTime
        }
        totalWait += float64(waitingTime)

        start := waitingTime + processes[i].ArrivalTime

        turnaround := processes[i].BurstDuration + waitingTime
        totalTurnaround += float64(turnaround)

        completion := processes[i].BurstDuration + processes[i].ArrivalTime + waitingTime
        lastCompletion = float64(completion)

        schedule[i] = []string{
            fmt.Sprint(processes[i].ProcessID),
            fmt.Sprint(processes[i].Priority),
            fmt.Sprint(processes[i].BurstDuration),
            fmt.Sprint(processes[i].ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(completion),
        }
        serviceTime += processes[i].BurstDuration

        gantt = append(gantt, TimeSlice{
            PID:   processes[i].ProcessID,
            Start: start,
            Stop:  serviceTime,
        })
    }

    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := count / lastCompletion

    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}

//func SJFPrioritySchedule(w io.Writer, title string, processes []Process) { }
//Function for SJFPriority Scheudle to remove the process
func removeProcess(processes []Process, process Process) []Process {
    var remaining []Process
    for i := range processes {
        if processes[i].ProcessID != process.ProcessID {
            remaining = append(remaining, processes[i])
        }
    }
    return remaining
}
//Function used to find the shortest path
func fndShrtPath(remaining []Process, serviceTime int64) *Process {
    var shortest *Process
    for i := range remaining {
        if remaining[i].ArrivalTime > serviceTime {

            break
        }
        if shortest == nil || remaining[i].BurstDuration < shortest.BurstDuration {
            shortest = &remaining[i]
        }
    }
    return shortest
}
//main function
func SJFSchedule(w io.Writer, title string, processes []Process) {
    var (
        serviceTime int64
        totalWait   float64
        totalTurnaround float64
        lastCompletion float64
        schedule    = make([][]string, len(processes))
        gantt       = make([]TimeSlice, 0)
    )
    remaining := make([]Process, len(processes))
    copy(remaining, processes)

    sort.SliceStable(remaining, func(p, q int) bool {  
      return remaining[p].ArrivalTime < remaining[q].ArrivalTime }) 

    for len(remaining) > 0 {
        next := fndShrtPath(remaining, serviceTime)
        if next == nil {
            // Used when there is no job
            serviceTime++
            continue
        }

        process := *next
        remaining = removeProcess(remaining, process)
		//Waiting time is the service time - the arrival time
        waitingTime := serviceTime - process.ArrivalTime
        if waitingTime < 0 {
            waitingTime = 0
        }
        totalWait += float64(waitingTime)

        start := serviceTime

        turnaround := process.BurstDuration + waitingTime
        totalTurnaround += float64(turnaround)

        completion := process.BurstDuration + serviceTime
        lastCompletion = float64(completion)

        schedule[process.ProcessID-1] = []string{
            fmt.Sprint(process.ProcessID),
            fmt.Sprint(process.Priority),
            fmt.Sprint(process.BurstDuration),
            fmt.Sprint(process.ArrivalTime),
            fmt.Sprint(waitingTime),
            fmt.Sprint(turnaround),
            fmt.Sprint(completion),
        }

        gantt = append(gantt, TimeSlice{
            PID:   process.ProcessID,
            Start: start,
            Stop:  completion,
        })

        serviceTime += process.BurstDuration
    }

    count := float64(len(processes))
    aveWait := totalWait / count
    aveTurnaround := totalTurnaround / count
    aveThroughput := count / lastCompletion
	//Outputs the function
    outputTitle(w, title)
    outputGantt(w, gantt)
    outputSchedule(w, schedule, aveWait, aveTurnaround, aveThroughput)
}


//func SJFSchedule(w io.Writer, title string, processes []Process) { }

func SJFPrioritySchedule(w io.Writer, title string, processes []Process) {
    fmt.Fprintf(w, "------ %s ------\n", title)
    var currTime int64

    finished := 0
    currTime = 0
    var waiting []Process
    var active *Process
	//If those finished is less than the length of the process
    for finished < len(processes) {
        for i := range processes {
            if !processes[i].Completed && processes[i].ArrivalTime <= currTime {
                waiting = append(waiting, processes[i])
            }
        }
        sort.Slice(waiting, func(i, j int) bool {
            return waiting[i].BurstDuration < waiting[j].BurstDuration
        })
        if active == nil && len(waiting) > 0 {
            active = &waiting[0]
            waiting = waiting[1:]
        }
		//if those active is not null
        if active != nil {
            active.BurstDuration--
            if active.BurstDuration == 0 {
                active.Completed = true
                finished++
                active.TurnaroundTime = currTime + 1 - active.ArrivalTime
                active.WaitingTime = active.TurnaroundTime - active.Priority
                active = nil
            }
        }
        currTime++
    }

    var totalTurnaround, totalWaiting int64
    for i := range processes {
        totalTurnaround += processes[i].TurnaroundTime
        totalWaiting += processes[i].WaitingTime
    }

    fmt.Fprintf(w, "Average turnaround time: %.2f\n", float64(totalTurnaround)/float64(len(processes)))
    fmt.Fprintf(w, "Average waiting time: %.2f\n", float64(totalWaiting)/float64(len(processes)))
    fmt.Fprintf(w, "Throughput: %.2f\n", float64(len(processes))/float64(currTime))
}

//func RRSchedule(w io.Writer, title string, processes []Process) { }


type Queue struct {
    processes []Process
    quantum int
}
func min(a, b int64) int64 {
    if a < b {
        return a
    }
    return b
}
func RRSchedule(w io.Writer, title string, processes []Process, quantum int64) {
    var (
        currTime, totalTurnaroundTime, totalWaitingTime int64
        n, finishedProcesses, qIndex, qSize = len(processes), 0, 0, 0
        readyQueue Queue
    )

    fmt.Fprintf(w, "==== %s ====\n\n", title)

    for len(processes) > 0 || qSize > 0 {
		//This for loop moves the new processes to the queue
        for len(processes) > 0 && processes[0].ArrivalTime <= currTime {
            readyQueue.processes = append(readyQueue.processes, processes[0])
            processes = processes[1:]
            qSize++
        }

        if qSize == 0 {
            currTime++
            continue
        }

		// This moves the next process in the queue
        process := readyQueue.processes[qIndex]

		//Runs the process for quantum it is in or until it is done
        executedTime := min(process.BurstDuration, quantum)
        currTime += executedTime
        process.BurstDuration -= executedTime

        // Update waiting time for all other processes in the queue
        for i := 0; i < qSize; i++ {
            if i == qIndex {
                continue
            }
            readyQueue.processes[i].WaitingTime += executedTime
        }

        // Remove finished processes
        if process.BurstDuration == 0 {
            finishedProcesses++
            qSize--
            totalTurnaroundTime += currTime - process.ArrivalTime
            totalWaitingTime += process.WaitingTime
            fmt.Fprintf(w, "Process %s finished at time %d (turnaround time %d, waiting time %d)\n", process.ProcessID, currTime, currTime-process.ArrivalTime, process.WaitingTime)
            for i := qIndex; i < qSize; i++ {
                readyQueue.processes[i] = readyQueue.processes[i+1]
            }
        } else {
            qIndex = (qIndex + 1) % qSize
        }
    }

    fmt.Fprintf(w, "\nAverage turnaround time: %f\n", float64(totalTurnaroundTime)/float64(n))
    fmt.Fprintf(w, "Average waiting time: %f\n", float64(totalWaitingTime)/float64(n))
    fmt.Fprintf(w, "Average throughput: %f\n", float64(n)/float64(currTime))
}

//endregion

//region Output helpers

func outputTitle(w io.Writer, title string) {
    _, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
    _, _ = fmt.Fprintln(w, strings.Repeat(" ", len(title)/2), title)
    _, _ = fmt.Fprintln(w, strings.Repeat("-", len(title)*2))
}

func outputGantt(w io.Writer, gantt []TimeSlice) {
    _, _ = fmt.Fprintln(w, "Gantt schedule")
    _, _ = fmt.Fprint(w, "|")
    for i := range gantt {
        pid := fmt.Sprint(gantt[i].PID)
        padding := strings.Repeat(" ", (8-len(pid))/2)
        _, _ = fmt.Fprint(w, padding, pid, padding, "|")
    }
    _, _ = fmt.Fprintln(w)
    for i := range gantt {
        _, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Start), "\t")
        if len(gantt)-1 == i {
            _, _ = fmt.Fprint(w, fmt.Sprint(gantt[i].Stop))
        }
    }
    _, _ = fmt.Fprintf(w, "\n\n")
}

func outputSchedule(w io.Writer, rows [][]string, wait, turnaround, throughput float64) {
    _, _ = fmt.Fprintln(w, "Schedule table")
    table := tablewriter.NewWriter(w)
    table.SetHeader([]string{"ID", "Priority", "Burst", "Arrival", "Wait", "Turnaround", "Exit"})
    table.AppendBulk(rows)
    table.SetFooter([]string{"", "", "", "",
        fmt.Sprintf("Average\n%.2f", wait),
        fmt.Sprintf("Average\n%.2f", turnaround),
        fmt.Sprintf("Throughput\n%.2f/t", throughput)})
    table.Render()
}

//endregion

//region Loading processes.

var ErrInvalidArgs = errors.New("invalid args")

func loadProcesses(r io.Reader) ([]Process, error) {
    rows, err := csv.NewReader(r).ReadAll()
    if err != nil {
        return nil, fmt.Errorf("%w: reading CSV", err)
    }

    processes := make([]Process, len(rows))
    for i := range rows {
        processes[i].ProcessID = mustStrToInt(rows[i][0])
        processes[i].BurstDuration = mustStrToInt(rows[i][1])
        processes[i].ArrivalTime = mustStrToInt(rows[i][2])
        if len(rows[i]) == 4 {
            processes[i].Priority = mustStrToInt(rows[i][3])
        }
    }

    return processes, nil
}

func mustStrToInt(s string) int64 {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        _, _ = fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    return i
}

//endregion
