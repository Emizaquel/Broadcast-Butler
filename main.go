package main

import (
	"fmt"
	"log"
	"os"

	hook "github.com/robotn/gohook"
	"go.yaml.in/yaml/v4"
)

type bbconfig struct {
	Port   string `yaml:"port"`
	Source string `yaml:"source"`
}

type stream_event struct {
	Platform string
	Event    int
	Data     *any
}

type message_part struct {
	Format  string
	Content string
}

type stream_user struct {
	Name string
	Kind int
}

type platform struct{
	Name string
	Messages []func()
}

type stream_message struct {
}

func add() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})

	fmt.Println("--- Please press w---")
	hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func low() {

	evChan := hook.Start()
	defer hook.End()
	for ev := range evChan {
		if ev.Kind == hook.KeyDown {
			fmt.Println("hook: ", hook.RawcodetoKeychar(ev.Rawcode))
		}
	}
}

func event() {
	ok := hook.AddEvents("q", "ctrl", "shift")
	if ok {
		fmt.Println("add events...")
	}

	keve := hook.AddEvent("k")
	if keve {
		fmt.Println("you press... ", "k")
	}

	mleft := hook.AddEvent("mleft")
	if mleft {
		fmt.Println("you press... ", "mouse left button")
	}
}

func main() {
	test, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalln("Error loading file: ", err)
	}
	var c bbconfig
	err = yaml.Unmarshal(test, &c)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}

	fmt.Println(c)

	add()
	low()
	event()
}
