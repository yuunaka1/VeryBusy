package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	
	"verybusy/internal/sim"
	"verybusy/internal/ui"
)

var (
	themeFlag string
	seedFlag  int64
)

var rootCmd = &cobra.Command{
	Use:   "verybusy",
	Short: "VeryBusy - SOC Dashboard Simulator",
	Long:  "A terminal dashboard that simulates a security operations monitoring screen.",
}

func runMode(mode string) {
	engine := sim.NewEngine(seedFlag, themeFlag)
	
	// Fast forward simulation a bit to populate data
	for i := 0; i < 50; i++ {
		engine.Tick()
		engine.GenerateLogs(time.Now())
	}

	m := ui.NewMainModel(engine, mode)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

var socCmd = &cobra.Command{
	Use:   "soc",
	Short: "Run the split-screen SOC overview",
	Run:   func(cmd *cobra.Command, args []string) { runMode("soc") },
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Run live security event log mode",
	Run:   func(cmd *cobra.Command, args []string) { runMode("logs") },
}

var alertsCmd = &cobra.Command{
	Use:   "alerts",
	Short: "Run the live detection/alert stream mode",
	Run:   func(cmd *cobra.Command, args []string) { runMode("alerts") },
}

var graphsCmd = &cobra.Command{
	Use:   "graphs",
	Short: "Run telemetry graph mode",
	Run:   func(cmd *cobra.Command, args []string) { runMode("graphs") },
}

var assetsCmd = &cobra.Command{
	Use:   "assets",
	Short: "Run asset/endpoint status mode",
	Run:   func(cmd *cobra.Command, args []string) { runMode("assets") },
}

var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Run live network traffic mode",
	Run:   func(cmd *cobra.Command, args []string) { runMode("network") },
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&themeFlag, "theme", "t", "soc", "Simulation theme/scenario (soc, cloud, endpoint)")
	rootCmd.PersistentFlags().Int64VarP(&seedFlag, "seed", "s", 0, "Simulation seed (0 for random)")

	rootCmd.AddCommand(socCmd, logsCmd, alertsCmd, graphsCmd, assetsCmd, networkCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
