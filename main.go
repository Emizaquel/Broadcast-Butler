package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"database/sql"

	goobs "github.com/andreykaipov/goobs"
	goobs_sceneitems "github.com/andreykaipov/goobs/api/requests/sceneitems"
	_ "github.com/mattn/go-sqlite3"
	hook "github.com/robotn/gohook"
	"go.yaml.in/yaml/v4"
)

// implement DB Stuff with sqlite3, it's stupid to do anything else

var platforms []platform

type bbconfig struct {
	Port      string  `yaml:"Broadcast Butler Port"`
	bbdb      string  `yaml:"Broadcast Butler Database"`
	Source    string  `yaml:"Site Folder"`
	DbVersion float32 `yaml:"DB Version"`
	ObsPort   string  `yaml:"OBS Port"`
	ObsPass   string  `yaml:"OBS Password"`
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
	id   string
	Kind int
}

type platform struct {
	Name       string
	Url        string
	Operations []func()
}

type stream_message struct {
	user stream_user
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

var (
	bbconf bbconfig
)

func api_serve(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.RequestURI, "/api") {
		http.ServeFile(w, r, bbconf.Source+"/home.html")
		return
	}
}

func base_serve(w http.ResponseWriter, r *http.Request) {
	log.Printf(r.RequestURI)

	if strings.HasPrefix(r.RequestURI, "/api") {
		api_serve(w, r)
		return
	}

	if r.RequestURI == "/" {
		http.ServeFile(w, r, bbconf.Source+"/home.html")
		return
	}

	http.ServeFile(w, r, bbconf.Source+r.RequestURI)
}

func backend_loop() {
	for {
		time.Sleep(2 * time.Minute)
		log.Printf("test")
	}
}

func serve_http() {
	http.HandleFunc("/", base_serve)
	log.Println("Running http server on http://localhost" + bbconf.Port)
	http.ListenAndServe(bbconf.Port, nil)
}

func migrate_db(db *sql.DB, config bbconfig) error {
	var err error
	if config.DbVersion < 0.1 {
		_, err = db.Exec(`ALTER TABLE "tasks" RENAME TO "process"`)
		if err != nil {
			return err
		}
		_, err = db.Exec(`create table tasks (id integer not null primary key, name text);delete from tasks;`)
		if err != nil {
			return err
		}
		var names []string
		query, err := db.Query(`SELECT name FROM process`)
		if err != nil {
			return err
		}
		defer query.Close()
		for query.Next() {
			var name string
			query.Scan(&name)
			log.Println("name: ", name)
			names = append(names, name)
		}
		stmt, err := db.Prepare(`INSERT INTO tasks(name) values(?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		for idex := range names {
			_, err = stmt.Exec(names[idex])
			if err != nil {
				return err
			}
		}
		_, err = db.Exec(`DROP TABLE process`)
		if err != nil {
			return err
		}
	}
	return nil
}

func setup_db(db *sql.DB, config bbconfig) error {

	rows, err := db.Query(`SELECT name FROM sqlite_schema WHERE type='table' AND name='tasks';`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	rcount := 0
	for rows.Next() {
		rcount += 1
		var name string
		rows.Scan(&name)
		log.Println("name: ", name)
	}
	if rcount < 1 {
		log.Println("no tables D:")
		_, err = db.Exec(`create table tasks (id integer not null primary key, name text);`)
		if err != nil {
			return err
		}
		config.DbVersion = 0.1
	}

	err = migrate_db(db, config)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(`INSERT INTO tasks(name) values(?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec("test 1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec("test 2")
	if err != nil {
		return err
	}

	log.Println("tables", rcount)
	return nil
}

func obs_thread(port string, password string) {
	port = "localhost:" + port
	log.Println(port)
	client, err := goobs.New(port, goobs.WithPassword(obs_password))
	if err != nil {
		log.Fatalln("OBS Error: ", err)
	}
	defer client.Disconnect()

	version, err := client.General.GetVersion()
	if err != nil {
		panic(err)
	}

	fmt.Printf("OBS Studio version: %s\n", version.ObsVersion)
	fmt.Printf("Server protocol version: %s\n", version.ObsWebSocketVersion)
	fmt.Printf("Client protocol version: %s\n", goobs.ProtocolVersion)
	fmt.Printf("Client library version: %s\n", goobs.LibraryVersion)

	params := goobs_sceneitems.NewGetSceneItemListParams().WithSceneName("Primary Layout")
	sceneList, err := client.SceneItems.GetSceneItemList(params)
	if err != nil {
		panic(err)
	}

	for _, item := range sceneList.SceneItems {
		log.Println(item.SourceName)
	}
}

func main() {

	test, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalln("Error loading file: ", err)
	}
	err = yaml.Unmarshal(test, &bbconf)
	if err != nil {
		log.Fatalf("cannot unmarshal config file: %v", err)
	}
	// log.Println(bbconf)
	db, err := sql.Open("sqlite3", bbconf.bbdb)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = setup_db(db, bbconf)
	if err != nil {
		log.Fatal("Error setting up DB: ", err)
	}

	bbconf.Port = ":" + bbconf.Port

	youtube_init()

	// fmt.Println(bbconf)
	var work_group sync.WaitGroup

	work_group.Add(1)
	go backend_loop()
	work_group.Add(1)
	go serve_http()
	work_group.Add(1)
	log.Println("obs port:", bbconf.ObsPort)
	go obs_thread(bbconf.ObsPort, bbconf.ObsPass)

	time.Sleep(2000)

	log.Println("test")

	// // add()
	// // low()
	// // event()

	work_group.Wait()
}
