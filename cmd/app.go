package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

// App FlyDB command app
var App = grumble.New(&grumble.Config{
	Name:                  "FlyDB Cli",
	Description:           "A command of FlyDB",
	HistoryFile:           path.Join(os.TempDir(), ".aDrive.history"),
	HistoryLimit:          10000,
	ErrorColor:            color.New(color.FgRed, color.Bold, color.Faint),
	HelpHeadlineColor:     color.New(color.FgGreen),
	HelpHeadlineUnderline: false,
	HelpSubCommands:       true,
	Prompt:                "aDrive $> ",
	PromptColor:           color.New(color.FgBlue, color.Bold),
	Flags:                 func(f *grumble.Flags) {},
})

func init() {
	App.OnInit(func(a *grumble.App, fm grumble.FlagMap) error {
		return nil
	})
	App.SetPrintASCIILogo(func(a *grumble.App) {
		fmt.Println(strings.Join([]string{`
Welcome to aDrive
`}, "\r\n"))
	})
	register(App)
}

func register(app *grumble.App) {
	app.AddCommand(&grumble.Command{
		Name: "get",
		Help: "get file",
		Run:  get,
		Args: func(a *grumble.Args) {
			a.String("file", "file to get")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "put",
		Help: "put file",
		Run:  put,
		Args: func(a *grumble.Args) {
			a.String("file", "file to put")
			a.String("directory", "directory to put")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "ls",
		Help: "list files",
		Run:  ls,
		Args: func(a *grumble.Args) {
			a.String("directory", "directory to list")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "rm",
		Help: "remove file",
		Run:  rm,
		Args: func(a *grumble.Args) {
			a.String("file", "file to remove")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "mkdir",
		Help: "make directory",
		Run:  mkdir,
		Args: func(a *grumble.Args) {
			a.String("directory", "directory to make")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "rename",
		Help: "rename file",
		Run:  rename,
		Args: func(a *grumble.Args) {
			a.String("old", "old file name")
			a.String("new", "new file name")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "addNameNode",
		Help: "add name node",
		Run:  addNameNode,
		Args: func(a *grumble.Args) {
			a.String("addr", "address of name node")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "addDataNode",
		Help: "add data node",
		Run:  addDataNode,
		Args: func(a *grumble.Args) {
			a.String("addr", "address of data node")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "rmNameNode",
		Help: "remove name node",
		Run:  rmNameNode,
		Args: func(a *grumble.Args) {
			a.String("addr", "address of name node")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "rmDataNode",
		Help: "remove data node",
		Run:  rmDataNode,
		Args: func(a *grumble.Args) {
			a.String("addr", "address of data node")
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "lsNameNode",
		Help: "list name nodes",
		Run:  lsNameNode,
	})
	app.AddCommand(&grumble.Command{
		Name: "lsDataNode",
		Help: "list data nodes",
		Run:  lsDataNode,
	})
}
