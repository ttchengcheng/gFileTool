/*
TODO:
. skip file customization: .syncignore
. test case
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"flag"

	"github.com/ttchengcheng/file/fsync"
	"github.com/ttchengcheng/file/fsync/local"
)

type Args struct {
	source, destination, ignore string
	keep, watch                 bool
}

func parseCommandLine() (args Args) {
	flagSet := flag.NewFlagSet("foldersync", flag.ContinueOnError)

	flagSet.StringVar(&args.source, "src", ".", "Path of the source directory")
	flagSet.StringVar(&args.destination, "des", ".", "Path of the destination directory")
	flagSet.StringVar(&args.ignore, "ignore", "", "Path of the ignore file")
	flagSet.BoolVar(&args.keep, "keep", false, "Keep the files that in destination but not in source")
	flagSet.BoolVar(&args.keep, "k", false, "Keep the files that in destination but not in source(shorthand)")
	flagSet.BoolVar(&args.watch, "watch", false, "Watch both dir and sync when any changes")
	flagSet.BoolVar(&args.watch, "w", false, "Watch both dir and sync when any changes(shorthand)")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		fmt.Println()

		flagSet.PrintDefaults()
		os.Exit(-1)
	}
	return
}

func main() {
	args := parseCommandLine()

	parseCommandLine()

	srcDir, err := filepath.Abs(args.source)
	if err != nil {
		fmt.Println(err)
		return
	}

	desDir, err := filepath.Abs(args.destination)
	if err != nil {
		fmt.Println(err)
		return
	}

	filter := &fsync.IgnoreSetting{}
	filter.Parse(".git/\n.sync-cache")

	src := local.File{Path: srcDir, IgnoreSetting: filter}
	des := local.File{Path: desDir, IgnoreSetting: filter}

	sc := &fsync.Sync{Source: &src, Destination: &des, KeepOther: args.keep}
	sc.Run()
}
