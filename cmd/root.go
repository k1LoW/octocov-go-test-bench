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
	"os"

	gotestbench "github.com/k1LoW/octocov-go-test-bench"
	"github.com/k1LoW/octocov-go-test-bench/version"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"golang.org/x/tools/benchmark/parse"
)

var rootCmd = &cobra.Command{
	Use:   "octocov-go-test-bench",
	Short: "Generate custom metrics JSON from the output of 'go test -bench'",
	Long:  `Generate custom metrics JSON from the output of 'go test -bench'.`,
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
		set, err := parse.ParseSet(os.Stdin)
		if err != nil {
			return err
		}
		cset := gotestbench.Convert(set)
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
}
