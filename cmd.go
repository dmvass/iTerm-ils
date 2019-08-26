package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"
)

const (
	defaultLocation = "."
	imgProto        = " \033]1337;File=inline=1;height=1:%s\a"
	tabSize         = 8
)

// Command is a command to list files in Unix and Unix-like operating systems
type Command struct {
	// Target directory location
	location string

	// Command flags
	LongFormat, NotSort, Revealing, ListAll, Recursion, TimeSort, HumanSize bool

	theme *Theme
}

// NewCommand is a command constructor
func NewCommand(t *Theme, args []string) (*Command, error) {
	cmd := Command{theme: t}
	// Parse flags and target directory
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			cmd.parseFlag(arg)
			continue
		}
		cmd.location = arg
	}
	if cmd.location == "" {
		cmd.location = defaultLocation
	}
	return &cmd, nil
}

// parseFlag is a command flags parser
//
// This method parse flags which starts with '-' and can parse as single
// '-l -a' and as grouped flags '-la'.
func (c *Command) parseFlag(args string) int {
	matched := 0
	for _, f := range args[1:] {
		switch string(f) {
		// LongFormat (-l): displaying Unix file types, permissions, number of
		// hard links, owner, group, size, last-modified date and filename.
		case "l":
			c.LongFormat = true

		// NotSort (-f): useful for directories containing large numbers
		// of files.
		case "f":
			c.NotSort = true

		// Revealing (-F): appends a character revealing the nature of a file,
		// for example, * for an executable, or / for a directory. Regular
		// files have no suffix.
		case "F":
			c.Revealing = true

		// ListAll (-a): lists all files in the given directory, including
		// those whose names start with "." (which are hidden files in Unix).
		// By default, these files are excluded from the list.
		case "a":
			c.ListAll = true

		// Recursion (-R): recursively lists subdirectories. The command
		// ls -R / would therefore list all files.
		case "R":
			c.Recursion = true

		// TimeSort (-t): sort the list of files by modification time.
		case "t":
			c.TimeSort = true

		// HumanSize (-h): print sizes in human readable format.
		case "h":
			c.HumanSize = true

		default:
			continue
		}
		matched++
	}
	return matched
}

func (c Command) execute(dir string) error {
	files, err := c.readdir(dir)
	if err != nil {
		return err
	}
	// display files in this directory
	c.display(files)

	if !c.Recursion {
		return nil
	}

	// walkthrough all directories and display files If Recursion is true
	for _, file := range files {
		if file == nil || !file.IsDir() {
			continue
		}
		nextDir := filepath.Join(dir, file.Name())

		// print displayed directory path
		fmt.Printf("\n%s\n", nextDir)

		err = c.execute(nextDir)
		if err != nil {
			return err
		}
	}

	return nil
}

// Execute command and display results
func (c Command) Execute() error {
	return c.execute(c.location)
}

// sort files in directory by last modification date or filename
func (c Command) sort(files []os.FileInfo) {
	switch {
	// sort files by last-modified date
	case !c.TimeSort:
		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[i].ModTime())
		})
	// sort files by filename
	default:
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() < files[j].Name()
		})
	}
}

func (c Command) readdir(dirname string) ([]os.FileInfo, error) {
	dir, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	files, err := dir.Readdir(-1)
	dir.Close()
	if err != nil {
		return nil, err
	}
	if !c.NotSort {
		c.sort(files)
	}
	if c.ListAll {
		return files, nil
	}
	for i, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			files[i] = nil
		}
	}
	return files, nil
}

// displayLong is display Unix file types, permissions, number of hard
// links, owner, group, size, last-modified date and filename.
func (c Command) displayLong(files []os.FileInfo) {
	fmt.Println("total ", len(files))
	for _, f := range files {
		if f == nil {
			continue
		}
		stat, _ := f.Sys().(*syscall.Stat_t)

		// Print file permissions
		fmt.Print(Permissions(f.Mode()))
		fmt.Print("\t")

		// Print file number of hard links
		fmt.Printf("%4d", getNlink(stat))
		fmt.Print("\t")

		// Print file owner
		fmt.Printf("%8s", getUser(stat))
		fmt.Print("\t")

		// Print file group
		fmt.Printf("%8s", getGroup(stat))
		fmt.Print("\t")

		// Print file size
		fmt.Printf("%10s", c.getSize(f))
		fmt.Print("\t")

		// Print file icon and filename
		fmt.Print(iTermIcon(c.theme.GetIcon(f)))
		fmt.Print(c.getFilename(f))
		fmt.Print("\n")
	}
}

// displayBare is display only Unix file types and filenames
func (c Command) displayBare(files []os.FileInfo) {
	var printedSize int
	size := lineSize()
	for _, f := range files {
		if f == nil {
			continue
		}
		filename := c.getFilename(f)
		icon := iTermIcon(c.theme.GetIcon(f))
		// check that line is not full
		if size > 0 {
			// add filename characters count
			printedSize += utf8.RuneCountInString(filename)
			// add printed icon characters count
			printedSize += 4

			if printedSize > size {
				fmt.Print("\n")
				printedSize = 0
			} else {
				// Add tabulation characters count
				printedSize += tabSize
			}
		}
		// Print file icon and filename
		fmt.Print(icon, filename)
		fmt.Print("\t")
	}
	fmt.Print("\n")
}

// display prints to stdout file result in a short or long style
func (c Command) display(files []os.FileInfo) {
	if c.LongFormat {
		c.displayLong(files)
	} else {
		c.displayBare(files)
	}
}

// getSize returns file size representation in bytes on human redable format
func (c Command) getSize(f os.FileInfo) string {
	if c.HumanSize {
		return Bytes(f.Size())
	}
	return strconv.FormatInt(f.Size(), 10)
}

// getFilename returns file name with a character revealing the nature
func (c Command) getFilename(f os.FileInfo) string {
	if c.Revealing && f.IsDir() {
		return fmt.Sprintf("%s/", f.Name())
	}
	return f.Name()
}

func lineSize() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return -1
	}

	r, err := regexp.Compile("[0-9]+")
	if err != nil {
		return -1
	}

	result := r.FindAllString(string(out), -1)
	size, err := strconv.Atoi(result[1])
	if err != nil {
		return -1
	}

	return size
}

// getUser returns username by user id
func getUser(s *syscall.Stat_t) string {
	if s == nil {
		return ""
	}
	u, err := user.LookupId(strconv.FormatUint(uint64(s.Uid), 10))
	if err != nil {
		return ""
	}
	return u.Username
}

// getGroup returns group name by group id
func getGroup(s *syscall.Stat_t) string {
	if s == nil {
		return ""
	}
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(s.Gid), 10))
	if err != nil {
		return ""
	}
	return g.Name
}

// getNlink returns file number of hard links
func getNlink(s *syscall.Stat_t) uint16 {
	if s == nil {
		return 0
	}
	return uint16(s.Nlink)
}

// iTermIcon format and returns icon regarding iTerm display image protocol
func iTermIcon(icon string) string {
	return fmt.Sprintf(imgProto, icon)
}
