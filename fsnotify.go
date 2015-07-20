package main

import (
	"github.com/howeyc/fsnotify"
	"log"
)

func fswatch(dir string, cb func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-watcher.Event:
				cb()
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(dir)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}
