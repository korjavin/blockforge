package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/cli"
	"github.com/zserge/webview"
	"gitlab.com/jgillich/autominer/miner"
)

func init() {
	Commands["gui"] = func() (cli.Command, error) {
		return GuiCommand{}, nil
	}
}

type GuiCommand struct {
}

func (c GuiCommand) Run(args []string) int {
	flags := flag.NewFlagSet("miner", flag.PanicOnError)
	flags.Usage = func() { ui.Output(c.Help()) }

	var configPath = flags.String("config", "miner.hcl", "Config file path")

	buf, err := ioutil.ReadFile(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Config file not found. Set '-config' argument or run 'coin miner -init' to generate.")
			return 1
		}
		fmt.Println(err)
		return 1
	}

	var config miner.Config
	err = hcl.Decode(&config, string(buf))
	if err != nil {
		fmt.Println(err)
		return 1
	}

	http.Handle("/", http.FileServer(rice.MustFindBox("../gui").HTTPBox()))

	listener, err := net.Listen("tcp", "127.0.0.1:3333")
	//listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	go func() {
		log.Fatal(http.Serve(listener, nil))
	}()

	view := webview.New(webview.Settings{
		URL:       "http://" + listener.Addr().String(),
		Title:     "CoinStack",
		Width:     1024,
		Height:    768,
		Resizable: true,
		Debug:     true,
	})
	defer view.Exit()

	view.Dispatch(func() {
		view.Bind("miner", &GuiMiner{config})
		view.Eval("init()")
	})

	view.Run()

	return 0
}

func (c GuiCommand) Help() string {

	helpText := `
Usage: coin miner [options]

	Launch the graphical user interface.

	General Options:

	-config=<path>          Config file path.
`
	return strings.TrimSpace(helpText)
}

func (c GuiCommand) Synopsis() string {
	return "Launch the graphical user interface"
}

type GuiMiner struct {
	Config miner.Config `json:"config"`
}

func (b *GuiMiner) Log(s string) {
	fmt.Println(s)
}
