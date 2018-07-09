package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/googollee/go-socket.io"
	"time"
	"html/template"
)
type Todo struct {
	Text string
	Date  string
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}
func main() {
	server, err := socketio.NewServer(nil)
	if(err != nil) {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		log.Println("Someone Connected!!")
		so.Join("chat_room")
		so.On("chat message", func(msg string) {
			log.Println("emit:", so.Emit("chat message", msg))

			add(msg)
			so.BroadcastTo("chat_room", "chat message", msg)
		})
	})

	http.Handle("/socket.io/", server)
	http.HandleFunc("/", htmlHandle)

	log.Println("Serving at :5000 rapbeh...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}



func htmlHandle(w http.ResponseWriter, r *http.Request)  {
	db, err := sql.Open("mysql", "root:123456@/test_db?charset=utf8")
	checkErr(err)

	rows, err := db.Query("SELECT text, date FROM chat")
	checkErr(err)

	data := TodoPageData{
		PageTitle: "My Chat",
		Todos: []Todo{
			{Text: "Task 1", Date: "2006-01-02 15:04:05"},
			{Text: "Task 2", Date: "2006-01-02 15:04:05"},
			{Text: "Task 3", Date: "2006-01-02 15:04:05"},
		},
	}

	for rows.Next() {
		var text string
		var date string

		err = rows.Scan(&text, &date)
		checkErr(err)

		log.Println(text, date)

		mss := Todo{
			Text: text,
			Date: date,
		}

		data.Todos = append(data.Todos, mss)
	}


	log.Println(data)

	tmpl := template.Must(template.ParseFiles("public/index.html"))
	tmpl.Execute(w, data)
	http.ListenAndServe(":5000", nil)
}

func add(text string) {
	db, err := sql.Open("mysql", "root:123456@/test_db?charset=utf8")
	checkErr(err)

	t := time.Now()
	var tres string = t.Format("2006-01-02 15:04:05")

	result, err := db.Query("INSERT INTO chat VALUES(NULL,?,?)", text, tres)

	checkErr(err)
	log.Println(result)
	log.Println("step 3")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}