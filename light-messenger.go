package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/usb-radiology/light-messenger/src/configuration"
	"github.com/usb-radiology/light-messenger/src/lmdatabase"
	"github.com/usb-radiology/light-messenger/src/server"
	"github.com/usb-radiology/light-messenger/src/version"
)

func main() {

	log.Printf("%s %s", version.Version, version.BuildTime)

	initConfig, err := configuration.LoadAndSetConfiguration("./config.json")
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}

	app := cli.NewApp()
	app.Name = "light-messenger"
	app.Usage = ""
	app.Version = version.Version + " " + version.BuildTime

	app.Commands = []cli.Command{
		{
			Name:  "web",
			Usage: "run web server (default)",
			Action: func(c *cli.Context) error {
				return actionWeb(initConfig)
			},
		},
		{
			Name:  "db-exec",
			Usage: "execute db script",
			Action: func(c *cli.Context) error {
				return actionDbExec(initConfig, c)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "script-path"},
			},
		},
	}

	app.Action = app.Commands[0].Action

	errRun := app.Run(os.Args)
	if errRun != nil {
		log.Fatalf("%+v", errors.WithStack(errRun))
	}

}

func actionWeb(initConfig *configuration.Configuration) error {
	httpServer := server.InitServer(initConfig)
	server.Start(httpServer, initConfig.Server.HTTPPort)
	return nil
}

func actionDbExec(initConfig *configuration.Configuration, c *cli.Context) error {
	db, errDb := lmdatabase.GetDB(initConfig)
	if errDb != nil {
		return errDb
	}

	scriptPath := c.String("script-path")

	statements, errReadStatements := lmdatabase.ReadStatementsFromSQL(scriptPath)
	if errReadStatements != nil {
		return errReadStatements
	}

	for _, statement := range *statements {
		log.Printf("%s", statement)

		execStatementResult, errExecStatement := lmdatabase.ExecStatement(db, statement)
		if errExecStatement != nil {
			return errExecStatement
		}

		if execStatementResult == nil {
			continue
		}

		rowsAffected, errRowsAffected := execStatementResult.RowsAffected()
		if errRowsAffected != nil {
			return errRowsAffected
		}

		log.Printf("rowsAffected %d", rowsAffected)
	}

	return nil
}
