package main

import (
	//"encoding/json"
	//"fmt"
	"log"

	"git.sr.ht/~rockorager/vaxis"
	"git.sr.ht/~rockorager/vaxis/vxfw"

	"ash/bgm/widgets/root"
	"ash/bgm/remote"
)

var app *vxfw.App

func main() {

	app, err := vxfw.NewApp(vaxis.Options{})
	if err != nil {
		log.Fatalf("Couldn't create a new app: %v", err)
	}

	// The remote needs the app to post vaxis.Events in the mainloop
	mpd_remote := remote.New(app)

	// Connect to the MPD server
	mpd_remote.Dial()
	defer mpd_remote.HangUp()

	// The widgets need a way to send commands to the remote
	newroot := root.New(mpd_remote, app)

	// Kick off the main loop
	app.Run(newroot)
}
