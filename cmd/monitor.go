package cmd

import (
	"github.com/rivo/tview"
	"fmt"
	"sort"
	"time"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/urfave/cli/v2"
)

func cmdMonitor() *cli.Command {
	return &cli.Command{
		Name:   "profile",
		Action: profile,
		Usage:  "Show reousrce usage of system",
		Description: `
Resource monitoring tools

Examples:
# Monitor real time operations
$ goup profile 

`,
	}
}
func profile(ctx *cli.Context) error {

	app := tview.NewApplication()
	textView := tview.NewTextView().SetDynamicColors(true)

	go func() {
		for {
			// Fetch metrics
			cpuPercent, _ := cpu.Percent(0, false)
			memStats, _ := mem.VirtualMemory()

			// Fetch processes
			procs, _ := process.Processes()
			var procStats []string
			for _, p := range procs {
				name, _ := p.Name()
				cpu, _ := p.CPUPercent()
				procStats = append(procStats, fmt.Sprintf("%-20s %.2f%%", name, cpu))
			}

			// Sort processes by CPU usage
			sort.Slice(procStats, func(i, j int) bool {
				return procStats[i] > procStats[j]
			})

			// Display top 5 processes
			topProcs := ""
			for i := 0; i < 5 && i < len(procStats); i++ {
				topProcs += procStats[i] + "\n"
			}

			// Update display
			stats := fmt.Sprintf(
				"[yellow]CPU Usage:[white] %.2f%%\n[yellow]Memory Usage:[white] %.2f%% of %.2fGB\n\n[yellow]Top Processes:\n[white]%s",
				cpuPercent[0],
				memStats.UsedPercent,
				float64(memStats.Total)/1e9,
				topProcs,
			)

			app.QueueUpdateDraw(func() {
				textView.SetText(stats)
			})
			time.Sleep(1 * time.Second)
		}
	}()

	if err := app.SetRoot(textView, true).Run(); err != nil {
		return err
	}
	return nil
}
