package main

import (
	"fmt"
	"quickstart-chat/module_bindings"

	"github.com/alexanderbh/bubbleapp/app"
	"github.com/alexanderbh/bubbleapp/component/box"
	"github.com/alexanderbh/bubbleapp/component/button"
	"github.com/alexanderbh/bubbleapp/component/context"
	"github.com/alexanderbh/bubbleapp/component/divider"
	"github.com/alexanderbh/bubbleapp/component/stack"
	"github.com/alexanderbh/bubbleapp/component/text"
	"github.com/alexanderbh/bubbleapp/component/textfield"
	"github.com/alexanderbh/spacetimedb-go-sdk"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type MainApp struct {
	Db     *spacetimedb.DBConnection
	logs   []string
	logger func(format string, args ...any)
}

var mainApp *MainApp = &MainApp{
	Db:   nil,
	logs: make([]string, 0),
}

func Logger(format string, args ...any) {
	if mainApp.logger == nil {
		mainApp.logs = append(mainApp.logs, fmt.Sprintf(format, args...))
		return
	}
	mainApp.logger(format, args...)
}

var AppDataContext = context.Create(mainApp)

func NewRoot(db *spacetimedb.DBConnection) func(c *app.Ctx) *app.C {
	return func(c *app.Ctx) *app.C {

		app.UseEffect(c, func() {
			mainApp.Db = db
			mainApp.logger = func(format string, args ...any) {
				mainApp.logs = append(mainApp.logs, fmt.Sprintf(format, args...))
				c.Update()
			}
		}, app.RunOnceDeps)

		return context.NewProvider(c, AppDataContext, mainApp, func(c *app.Ctx) *app.C {
			return NewMain(c)
		})
	}
}

func NewMain(c *app.Ctx) *app.C {

	data := context.UseContext(c, AppDataContext)

	return stack.New(c, func(c *app.Ctx) []*app.C {
		return []*app.C{
			box.New(c, func(c *app.Ctx) *app.C {
				return NewChat(c)
			}),
			box.NewEmpty(c, box.WithWidth(1), box.WithBg(c.Theme.Colors.InfoDark)),
			box.New(c, func(c *app.Ctx) *app.C {
				return stack.New(c, func(c *app.Ctx) []*app.C {
					views := []*app.C{
						text.New(c, "Logs:", text.WithFg(c.Theme.Colors.InfoFg)),
						divider.New(c),
					}

					for _, log := range data.logs {
						views = append(views, text.New(c, log))
					}

					return views
				}, stack.WithGrow(true))
			}),
		}
	}, stack.WithDirection(app.Horizontal))
}

func NewChat(c *app.Ctx) *app.C {

	hasNameSet, setHasNameSet := app.UseState(c, false)
	name, setName := app.UseState(c, "")

	message, setMessage := app.UseState(c, "")

	submit := func(n string) bool {
		if n != "" {
			err := module_bindings.SetName(mainApp.Db, n)
			if err != nil {
				mainApp.logger("Error calling set_name reducer: %v", err)
				return false
			} else {
				setHasNameSet(true)
				return true
			}
		}
		return false
	}

	sendMessage := func(msg string) bool {
		if msg != "" {
			mainApp.logger("Sending message: %s", msg)
		}
		return false
	}

	app.UseGlobalKeyHandler(c, func(keyMsg tea.KeyMsg) bool {
		switch keyMsg.String() {
		case "enter":
			if !hasNameSet {
				return submit(name)
			}
			return sendMessage(message)
		}
		return false
	})

	app.UseEffect(c, func() {
		c.FocusNext()
	}, app.RunOnceDeps)

	if !hasNameSet {
		return stack.New(c, func(c *app.Ctx) []*app.C {
			return []*app.C{
				textfield.New(c, func(text string) {
					setName(text)
				}, name),
				button.New(c, "Set Name", func() {
					submit(name)
				}),
			}
		})
	}

	return stack.New(c, func(c *app.Ctx) []*app.C {
		return []*app.C{
			box.New(c, func(c *app.Ctx) *app.C {
				return text.New(c, "messages")
			}),
			divider.New(c),
			stack.New(c, func(c *app.Ctx) []*app.C {
				return []*app.C{
					textfield.New(c, func(text string) {
						setMessage(text)
					}, message),
				}
			}, stack.WithDirection(app.Horizontal), stack.WithGrowY(false)),
		}

	})
}
