/*
Copyright © 2025 Brian Ketelsen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

// Package surgeon is a tool for modifying Go code
// using code mods. It allows you to define a set of code modifications
// and apply them to a Go project. The code mods are defined in a configuration
// file and can be applied to the source code using the surgeon tool.
package surgeon

type Config struct {
	Upstream   string
	ModsDir    string    `mapstructure:"modsdir"`
	CodeMods   []CodeMod `mapstructure:"codemods"`
	IgnoreList []Ignore  `mapstructure:"ignorelist"`
}

type Ignore struct {
	Prefix string
}

type CodeMod struct {
	Description string
	Mod         string
	Match       string // glob https://pkg.go.dev/path/filepath#Glob
	Args        []string
}
