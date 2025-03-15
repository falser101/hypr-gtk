package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Binding struct {
	FullLine    string
	Flags       string
	Modifiers   string
	Key         string
	Description string
	Command     string
	Args        string
	LineNumber  int
}

type Config struct {
	Lines    []string
	Bindings []Binding
}

func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		config.Lines = append(config.Lines, line)

		if isBindLine(line) {
			binding := parseBindLine(line, lineNum)
			if binding != nil {
				config.Bindings = append(config.Bindings, *binding)
			}
		}
		lineNum++
	}
	return config, scanner.Err()
}

func isBindLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "bind")
}

func parseBindLine(line string, lineNum int) *Binding {
	trimmed := strings.TrimSpace(line)
	parts := strings.SplitN(trimmed, "=", 2)
	if len(parts) != 2 {
		return nil
	}

	left := strings.TrimSpace(parts[0])
	if !strings.HasPrefix(left, "bind") {
		return nil
	}
	flags := strings.TrimPrefix(left, "bind")

	rightParts := strings.Split(strings.TrimSpace(parts[1]), ",")
	for i := range rightParts {
		rightParts[i] = strings.TrimSpace(rightParts[i])
	}

	if len(rightParts) < 4 {
		return nil
	}

	binding := &Binding{
		FullLine:   line,
		Flags:      flags,
		LineNumber: lineNum,
	}

	if len(rightParts) > 0 {
		binding.Modifiers = rightParts[0]
	}
	if len(rightParts) > 1 {
		binding.Key = rightParts[1]
	}
	if len(rightParts) > 2 {
		binding.Description = rightParts[2]
	}
	if len(rightParts) > 3 {
		binding.Command = rightParts[3]
	}
	if len(rightParts) > 4 {
		binding.Args = strings.Join(rightParts[4:], ", ")
	}

	return binding
}

func (c *Config) Save(filename string) error {
	content := []byte(strings.Join(c.Lines, "\n"))
	return os.WriteFile(filename, content, 0644)
}

func (c *Config) findInsertPosition(flags string) int {
	lastLine := -1
	for _, b := range c.Bindings {
		if b.Flags == flags && b.LineNumber > lastLine {
			lastLine = b.LineNumber
		}
	}
	if lastLine != -1 {
		return lastLine + 1
	}

	for i := len(c.Lines) - 1; i >= 0; i-- {
		if isBindLine(c.Lines[i]) {
			return i + 1
		}
	}
	return len(c.Lines)
}

func (c *Config) AddBinding(flags, modifiers, key, description, command, args string) {
	var rightParts []string
	rightParts = append(rightParts, modifiers, key, description, command)
	if args != "" {
		rightParts = append(rightParts, args)
	}

	newLine := fmt.Sprintf("bind%s = %s", flags, strings.Join(rightParts, ", "))
	insertPos := c.findInsertPosition(flags)

	c.Lines = append(c.Lines[:insertPos], append([]string{newLine}, c.Lines[insertPos:]...)...)
	c.Bindings = append(c.Bindings, Binding{
		FullLine:    newLine,
		Flags:       flags,
		Modifiers:   modifiers,
		Key:         key,
		Description: description,
		Command:     command,
		Args:        args,
		LineNumber:  insertPos,
	})

	for i := range c.Bindings {
		if c.Bindings[i].LineNumber >= insertPos {
			c.Bindings[i].LineNumber++
		}
	}
}

func (c *Config) UpdateBinding(index int, modifiers, key, description, command, args string) {
	if index < 0 || index >= len(c.Bindings) {
		return
	}

	b := &c.Bindings[index]
	var rightParts []string
	rightParts = append(rightParts, modifiers, key, description, command)
	if args != "" {
		rightParts = append(rightParts, args)
	}

	newLine := fmt.Sprintf("bind%s = %s", b.Flags, strings.Join(rightParts, ", "))
	c.Lines[b.LineNumber] = newLine
	b.FullLine = newLine
	b.Modifiers = modifiers
	b.Key = key
	b.Description = description
	b.Command = command
	b.Args = args
}

func (c *Config) DeleteBinding(index int) {
	if index < 0 || index >= len(c.Bindings) {
		return
	}

	b := c.Bindings[index]
	c.Lines = append(c.Lines[:b.LineNumber], c.Lines[b.LineNumber+1:]...)

	c.Bindings = append(c.Bindings[:index], c.Bindings[index+1:]...)
	for i := range c.Bindings {
		if c.Bindings[i].LineNumber > b.LineNumber {
			c.Bindings[i].LineNumber--
		}
	}
}
