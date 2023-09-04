/*
Copyright © 2023 Ken'ichiro Oyama <k1lowxb@gmail.com>

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
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	gotestbench "github.com/k1LoW/octocov-go-test-bench"
	"github.com/k1LoW/octocov-go-test-bench/version"
	"github.com/k1LoW/octocov/report"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/tools/benchmark/parse"
)

var (
	tee     bool
	targets []string
)

var rootCmd = &cobra.Command{
	Use:   "octocov-go-test-bench",
	Short: "Generate octocov custom metrics JSON from the output of 'go test -bench'",
	Long:  `Generate octocov custom metrics JSON from the output of 'go test -bench'.`,
	Args: func(cmd *cobra.Command, args []string) error {
		versionVal, err := cmd.Flags().GetBool("version")
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if versionVal {
			fmt.Println(version.Version)
			os.Exit(0)
		}

		if isatty.IsTerminal(os.Stdin.Fd()) {
			return errors.New("octocov-go-test-bench need STDIN. Please use pipe")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var in io.Reader = os.Stdin
		if tee {
			in = io.TeeReader(os.Stdin, os.Stderr)
		}
		set, err := parse.ParseSet(in)
		if err != nil {
			return err
		}
		cset := gotestbench.Convert(set)
		if len(targets) > 0 {
			filtered := []*report.CustomMetricSet{}
			for _, t := range targets {
				re, err := regexp.Compile(t)
				if err != nil {
					return err
				}
				for _, cs := range cset {
					if re.MatchString(cs.Key) {
						filtered = append(filtered, cs)
					}
				}
			}
			cset = filtered
		}
		if len(cset) == 0 {
			return errors.New("no benchmarks found")
		}
		b, err := json.MarshalIndent(cset, "", "  ")
		if err != nil {
			return err
		}
		cmd.Println(string(b))
		return nil
	},
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "print the version")
	rootCmd.Flags().BoolVarP(&tee, "tee", "", false, "print stdin to stderr")
	rootCmd.Flags().StringSliceVarP(&targets, "target", "", []string{}, "target benchmark name")
}
