package main

import (
	"flag"
	// "fmt"

	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var (
	AppName string
	BuiltAt string
	debug = flag.Bool("d", true, "enables the debug mode")
	w *astilectron.Window
	StaticPath string
)

func main() {
	flag.Parse()
	astilog.FlagInit()

	astilog.Debugf("Running app built at %s", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:          Asset,
		RestoreAssets:  RestoreAssets,
		AssetDir:       AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			// AppIconDarwinPath:  "resources/icon.icns",
			// AppIconDefaultPath: "resources/icon.png",
		},
		Debug:    *debug,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astilectron.PtrStr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				// {Label: astilectron.PtrStr("About")},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(app *astilectron.Astilectron, iw []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = iw[0]
			StaticPath = app.Paths().DataDirectory()
			return nil
		},
		Windows: []*bootstrap.Window{{
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(700),
				Width:           astilectron.PtrInt(1300),
			},
			MessageHandler: HandleMessages,
			Homepage: "index.html",
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}