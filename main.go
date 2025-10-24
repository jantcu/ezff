package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func checkFFmpeg() bool {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Warning: ffmpeg not installed or not in PATH. Commands will not work.")
		return false
	}
	return true
}

func getDuration(input string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-show_entries", "format=duration", "-of", "csv=p=0", input)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
}

func generateOutput(input, suffix string) string {
	dir := filepath.Dir(input)
	base := filepath.Base(input)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	out := filepath.Join(dir, fmt.Sprintf("%s%s%s", name, suffix, ext))
	i := 1
	for {
		if _, err := os.Stat(out); os.IsNotExist(err) {
			return out
		}
		out = filepath.Join(dir, fmt.Sprintf("%s%s%d%s", name, suffix, i, ext))
		i++
	}
}

func runCmd(name string, args []string) {
	c := exec.Command(name, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error running", name+":", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: ezff <command> [args]

Commands:
  trim-start <input> <seconds> [--output <output>]
  trim-end <input> <seconds> [--output <output>]
  trim-mid <input> <start_cut> <end_cut> [--output <output>]
  trim <input> <trim_start> <trim_end> [--output <output>]`)
	os.Exit(1)
}

func main() {
	checkFFmpeg()
	args := os.Args
	if len(args) < 2 {
		usage()
	}
	switch args[1] {
	case "trim-start":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: ezff trim-start <input> <seconds> [--output <output>]")
			os.Exit(1)
		}
		input := args[2]
		seconds, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid seconds")
			os.Exit(1)
		}
		output := ""
		if len(args) > 5 && args[4] == "--output" {
			output = args[5]
		}
		if output == "" {
			output = generateOutput(input, "_trim")
		}
		cmdArgs := []string{"-ss", fmt.Sprintf("%f", seconds), "-i", input, "-c", "copy", output}
		runCmd("ffmpeg", cmdArgs)
	case "trim-end":
		if len(args) < 4 {
			fmt.Fprintln(os.Stderr, "Usage: ezff trim-end <input> <seconds> [--output <output>]")
			os.Exit(1)
		}
		input := args[2]
		seconds, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid seconds")
			os.Exit(1)
		}
		dur, err := getDuration(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting duration:", err)
			os.Exit(1)
		}
		if dur <= seconds {
			fmt.Fprintln(os.Stderr, "Duration too short")
			os.Exit(1)
		}
		t := dur - seconds
		output := ""
		if len(args) > 5 && args[4] == "--output" {
			output = args[5]
		}
		if output == "" {
			output = generateOutput(input, "_trim")
		}
		cmdArgs := []string{"-i", input, "-c", "copy", "-t", fmt.Sprintf("%f", t), output}
		runCmd("ffmpeg", cmdArgs)
	case "trim-mid":
		if len(args) < 5 {
			fmt.Fprintln(os.Stderr, "Usage: ezff trim-mid <input> <start_cut> <end_cut> [--output <output>]")
			os.Exit(1)
		}
		input := args[2]
		start, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid start_cut")
			os.Exit(1)
		}
		end, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid end_cut")
			os.Exit(1)
		}
		if start >= end {
			fmt.Fprintln(os.Stderr, "start_cut must be less than end_cut")
			os.Exit(1)
		}
		output := ""
		if len(args) > 6 && args[5] == "--output" {
			output = args[6]
		}
		if output == "" {
			output = generateOutput(input, "_trim")
		}
		vf := fmt.Sprintf("select='not(between(t,%f,%f))',setpts=N/FRAME_RATE/TB", start, end)
		af := fmt.Sprintf("aselect='not(between(t,%f,%f))',asetpts=N/SR/TB", start, end)
		cmdArgs := []string{"-i", input, "-vf", vf, "-af", af, output}
		runCmd("ffmpeg", cmdArgs)
	case "trim":
		if len(args) < 5 {
			fmt.Fprintln(os.Stderr, "Usage: ezff trim <input> <trim_start> <trim_end> [--output <output>]")
			os.Exit(1)
		}
		input := args[2]
		trimStart, err := strconv.ParseFloat(args[3], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid trim_start")
			os.Exit(1)
		}
		trimEnd, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid trim_end")
			os.Exit(1)
		}
		dur, err := getDuration(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting duration:", err)
			os.Exit(1)
		}
		totalTrim := trimStart + trimEnd
		if dur <= totalTrim {
			fmt.Fprintln(os.Stderr, "Duration too short for trimming")
			os.Exit(1)
		}
		t := dur - totalTrim
		output := ""
		if len(args) > 6 && args[5] == "--output" {
			output = args[6]
		}
		if output == "" {
			output = generateOutput(input, "_trim")
		}
		cmdArgs := []string{"-ss", fmt.Sprintf("%f", trimStart), "-i", input, "-t", fmt.Sprintf("%f", t), "-c", "copy", output}
		runCmd("ffmpeg", cmdArgs)
	default:
		usage()
	}
}
