// based on https://github.com/tomlazar/table

package dyn_table

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/mgutz/ansi"
)

const (
	AlignLeft = iota
	AlignRight
)

// Config is the
type Config struct {
	ShowIndex       bool     // shows the index/row number as the first column
	Color           bool     // use the color codes in the output
	AlternateColors bool     // alternate the colors when writing
	TitleColorCode  string   // the ansi code for the title row
	AltColorCodes   []string // the ansi codes to alternate between
}

// DefaultConfig returns the default config for table, if its ever left null in a method this will be the one
// used to display the table
func DefaultConfig() *Config {
	return &Config{
		ShowIndex:       false,
		Color:           true,
		AlternateColors: true,
		TitleColorCode:  ansi.ColorCode("yellow+buf"),
		AltColorCodes: []string{
			"",
			"\u001b[40m",
		},
	}
}

// DynTable is the struct used to define the structure, this can be used from a zero state, or inferred using the
// reflection based methods
type DynTable struct {
	Width          []int
	Align          []int
	Headers        []string
	config         *Config
	rowColor       bool
	idxColumnWidth int
	rowIdx         int
}

func (t *DynTable) WriteHeader(conf *Config, idxColumnWidth int) {

	if conf == nil {
		t.config = DefaultConfig()
	}

	t.idxColumnWidth = idxColumnWidth

	if t.config.Color {
		fmt.Print(t.config.TitleColorCode)
	}
	if t.config.ShowIndex {
		fmt.Printf(" [%*v]  ", t.idxColumnWidth, "#")
	}

	for i, header := range t.Headers {
		if len(t.Align) == 0 || t.Align[i] == AlignRight {
			fmt.Printf("  %s", runewidth.FillLeft(header, t.Width[i]))
		} else {
			fmt.Printf("  %s", runewidth.FillRight(header, t.Width[i]))
		}
	}

	if t.config.Color {
		fmt.Print(ansi.Reset)
	}
	fmt.Println()

	t.rowColor = t.config.Color && t.config.AlternateColors && len(t.config.AltColorCodes) > 1
}

func (t *DynTable) AppendRow(row []string) {
	t.AppendRowWithColor(row,"")
}

func (t *DynTable) AppendRowWithColor(row []string, altColor string) {

	if t.rowColor {
		fmt.Print(t.config.AltColorCodes[t.rowIdx%len(t.config.AltColorCodes)])
	}
	
	if len(altColor) != 0 {
		fmt.Print(ansi.ColorCode(altColor))
	}
	
	if t.config.ShowIndex {
		fmt.Printf(" [%*v]  ", t.idxColumnWidth, t.rowIdx+1)
	}
	
	for i, v := range row {

		if runewidth.StringWidth(v) > t.Width[i] {
			v = runewidth.Truncate(v, t.Width[i], "")
		}

		if len(t.Align) == 0 || t.Align[i] == AlignRight {
			fmt.Printf("  %s", runewidth.FillLeft(v, t.Width[i]))
		} else {
			fmt.Printf("  %s", runewidth.FillRight(v, t.Width[i]))
		}
	}

	if t.rowColor {
		fmt.Print(ansi.Reset)
	}

	fmt.Println()

	t.rowIdx += 1
}
