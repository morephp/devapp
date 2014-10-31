package main

import (
	"github.com/go-fsnotify/fsnotify"
	"log"
	"os"
)

func main() {
	paths := []string{
		"d:\\temp",
	}
	wait := make(chan bool)
	go NewWatchDirectory(paths)
	<-wait
}

func NewWatchDirectory(paths []string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("[Error] Can't create fsnotify watcher.")
		os.Exit(-1)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				// if ev.Op&fsnotify.Create == fsnotify.Create {
				// 	log.Println("[Event] Create ", ev.Name)
				// }
				// if ev.Op&fsnotify.Remove == fsnotify.Remove {
				// 	log.Println("[Event] Remove ", ev.Name)
				// }
				// if ev.Op&fsnotify.Rename == fsnotify.Rename {
				// 	log.Println("[Event] Rename ", ev.Name)
				// }
				// if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
				// 	log.Println("[Event] Chmod ", ev.Name)
				// }
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if IsExistDirectory(ev.Name) {
						log.Println("[Warn] ", ev.Name, " is Directory.")
						err = watcher.Add(ev.Name)
						if err != nil {
							log.Println("[Error] Watch Directory (", err, ")")
							os.Exit(-1)
						}
						log.Println("[Suecces] Add Direcotry ", ev.Name)
					}

					log.Println("[Event] Write ", ev.Name)
				}
			case err := <-watcher.Errors:
				log.Println("[Error] ", err)
			}

		}
	}()

	for _, v := range paths {
		log.Println("[Trac] Directory (", v, ")")
		err = watcher.Add(v)
		if err != nil {
			log.Println("[Error] Watch Directory (", err, ")")
			os.Exit(-1)
		}
	}
	<-done
}

func IsExistDirectory(path string) bool {
	fd, err := os.Stat(path)
	if err != nil {
		log.Println("[Error] ", err)
		return false
	}
	return fd.IsDir()
}
