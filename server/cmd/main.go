package main

import "gitlab.com/kirasmir2/vogo/server/run"

func main() {
	app := run.NewApp()

	app.Init().Start()
}
