package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	modFile       string
	graphFile     string
	ignoreVersion bool
)

func initTree() (string, error) {
	buf := bytes.NewBuffer(nil)
	var source io.Reader
	if graphFile == "" {
		source = os.Stdin
	} else {
		if f, err := os.Open(graphFile); err != nil {
			return "", err
		} else {
			source = f
		}
	}

	if _, err := io.Copy(buf, source); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func registerGraph(root *cobra.Command) {
	graphCmd := &cobra.Command{
		Use: "graph",
		RunE: func(cmd *cobra.Command, args []string) error {
			content, err := initTree()
			if err != nil {
				return err
			}
			t := NewTree(content)
			t.Graph()
			return nil
		},
	}
	root.AddCommand(graphCmd)
}

func registerDep(root *cobra.Command) {
	tail := false
	depCmd := &cobra.Command{
		Use: "dep",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no target")
			}
			content, err := initTree()
			if err != nil {
				return err
			}
			t := NewTree(content)
			t.Dep(args[0], tail)
			return nil
		},
	}
	depCmd.PersistentFlags().BoolVar(&tail, "t", false, "print modules dependency to tail")
	root.AddCommand(depCmd)
}

func registerRef(root *cobra.Command) {
	refCmd := &cobra.Command{
		Use: "ref",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no target")
			}
			content, err := initTree()
			if err != nil {
				return err
			}
			t := NewTree(content)
			t.Ref(args[0])
			return nil
		},
	}
	root.AddCommand(refCmd)
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "mgraph",
		Short: "mgraph transform the result of 'go mod graph' to a better view",
	}
	rootCmd.PersistentFlags().StringVar(&graphFile, "f", "", "read module graph from file")
	rootCmd.PersistentFlags().StringVar(&modFile, "m", "./go.mod", "read go.mod file")
	rootCmd.PersistentFlags().BoolVar(&ignoreVersion, "i", false, "parse module dependency tree without version") // TODO

	registerGraph(rootCmd)
	registerDep(rootCmd)
	registerRef(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
