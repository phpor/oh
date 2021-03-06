// Released under an MIT license. See LICENSE.

package ui

import (
	"github.com/michaelmacinnis/oh/pkg/cell"
	"github.com/michaelmacinnis/oh/pkg/common"
	"github.com/michaelmacinnis/oh/pkg/system"
	"github.com/michaelmacinnis/oh/pkg/task"
	"github.com/peterh/liner"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type cli struct {
	*liner.State
}

var (
	cooked   liner.ModeApplier
	uncooked liner.ModeApplier
)

func CtrlCPressed(err error) bool {
	return err == liner.ErrPromptAborted
}

func New(args []string) *cli {
	if len(args) > 1 {
		return nil
	}

	// We assume the terminal starts in cooked mode.
	cooked, _ = liner.TerminalMode()
	if cooked == nil {
		return nil
	}

	i := &cli{liner.NewLiner()}

	if history_path, err := system.GetHistoryFilePath(); err == nil {
		if f, err := os.Open(history_path); err == nil {
			i.ReadHistory(f)
			f.Close()
		}
	}

	uncooked, _ = liner.TerminalMode()

	i.SetCtrlCAborts(true)
	i.SetTabCompletionStyle(liner.TabPrints)
	i.SetShouldRestart(restart)
	i.SetWordCompleter(complete)

	return i
}

func (i *cli) Close() error {
	if i.Exists() {
		if history_path, err := system.GetHistoryFilePath(); err == nil {
			if f, err := os.Create(history_path); err == nil {
				i.WriteHistory(f)
				f.Close()
			} else {
				println("Error writing history: " + err.Error())
			}
		}
	}
	return i.State.Close()
}

func (i *cli) Exists() bool {
	return i != nil
}

func (i *cli) ReadString(delim byte) (line string, err error) {
	system.SetForegroundGroup(system.Pgid())

	uncooked.ApplyMode()
	defer cooked.ApplyMode()

	command := cell.List(
		cell.Cons(
			cell.NewSymbol("_sys_"),
			cell.NewSymbol("get-prompt"),
		),
		cell.NewSymbol("> "),
	)
	prompt := task.Call(nil, command, "")

	if line, err = i.State.Prompt(prompt); err == nil {
		i.AppendHistory(line)
		if task.ForegroundTask().Job.Command == "" {
			task.ForegroundTask().Job.Command = line
		}
		line += "\n"
	}

	if err == liner.ErrPromptAborted {
		return line, common.CtrlCPressed
	}

	return
}

func complete(line string, pos int) (string, []string, string) {
	head := line[:pos]
	tail := line[pos:]

	fields := strings.Fields(head)

	if len(fields) == 0 {
		return head, []string{"    "}, tail
	}

	word := fields[len(fields)-1]
	if !strings.HasSuffix(head, word) {
		return head, []string{}, tail
	}

	head = head[0 : len(head)-len(word)]

	completions := task.ForegroundTask().Complete(word)
	completions = append(completions, files(word)...)
	if len(fields) == 1 {
		completions = append(completions, executables(word)...)
	}

	if len(completions) == 0 {
		return head, []string{word}, tail
	}

	unique := make(map[string]bool)
	for _, completion := range completions {
		unique[completion] = true
	}

	completions = make([]string, 0, len(unique))
	for completion := range unique {
		completions = append(completions, completion)
	}

	return head, completions, tail
}

func executables(word string) []string {
	completions := []string{}

	if strings.Contains(word, string(os.PathSeparator)) {
		return completions
	}

	pathenv := os.Getenv("PATH")
	for _, dir := range strings.Split(pathenv, string(os.PathListSeparator)) {
		if dir == "" {
			dir = "."
		} else {
			dir = path.Clean(dir)
		}

		stat, err := os.Stat(dir)
		if err != nil || !stat.IsDir() {
			continue
		}

		max := strings.Count(dir, "/") + 1
		filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
			depth := strings.Count(p, "/")
			if depth > max {
				if i.IsDir() {
					return filepath.SkipDir
				}
				return nil
			} else if depth < max {
				return nil
			}

			_, basename := filepath.Split(p)

			if strings.HasPrefix(basename, word) {
				completions = append(completions, basename)
			}

			return nil
		})
	}

	return completions
}

func files(word string) []string {
	completions := []string{}

	candidate := word
	if candidate[:1] == "~" {
		candidate = filepath.Join(os.Getenv("HOME"), candidate[1:])
	}

	candidate = path.Clean(candidate)
	if !path.IsAbs(candidate) {
		ft := task.ForegroundTask()
		n := cell.NewSymbol("$PWD")
		ref, _ := task.Resolve(ft.Lexical, ft.Frame, n)
		cwd := ref.Get().String()

		candidate = path.Join(cwd, candidate)
	}

	dirname, basename := filepath.Split(candidate)
	if strings.HasSuffix(word, "/") {
		dirname, basename = path.Join(dirname, basename)+"/", ""
	}

	stat, err := os.Stat(dirname)
	if err != nil {
		return completions
	} else if len(basename) == 0 && !stat.IsDir() {
		return completions
	}

	max := strings.Count(dirname, "/")

	filepath.Walk(dirname, func(p string, i os.FileInfo, err error) error {
		depth := strings.Count(p, "/")
		if depth > max {
			if i.IsDir() {
				return filepath.SkipDir
			}
			return nil
		} else if depth < max {
			return nil
		}

		full := path.Join(dirname, basename)
		if len(basename) == 0 {
			if p == dirname {
				return nil
			}
			full += "/"
		} else if !strings.HasPrefix(p, full) {
			return nil
		}

		if i.IsDir() {
			p += "/"
		}

		if len(full) >= len(p) {
			return nil
		}

		completion := word + p[len(full):]
		completions = append(completions, completion)

		return nil
	})

	return completions
}

func restart(err error) bool {
	if err == io.EOF {
		return false
	}

	return system.ResetForegroundGroup(os.Stdin)
}
