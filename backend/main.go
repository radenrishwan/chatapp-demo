package main

import (
	"chatapp/ws"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

var (
	PORT = flag.String("port", "8080", "app port")
)

func main() {
	flag.Parse()

	router := chi.NewRouter()

	pool := ws.NewPool()
	go pool.Run()

	router.Use(middleware.Logger)

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		// passing data to html
		files, err := template.ParseFiles("./../frontend/web/index.html")
		if err != nil {
			log.Println("Error on parse file : ", err.Error())
		}

		err = files.Execute(writer, map[string]any{
			"Rooms": pool.GetRooms(),
			"Len":   len(pool.GetRooms()),
		})
		if err != nil {
			log.Println("Error on execute file : ", err.Error())
		}
	})

	router.Get("/{room}", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "./../frontend/web/room.html")
	})

	router.Get("/ws/{room}", func(writer http.ResponseWriter, request *http.Request) {
		runApp(writer, request, pool)
	})

	log.Println("App running....")
	log.Fatalln(http.ListenAndServe(":"+*PORT, router))
}

func runApp(w http.ResponseWriter, r *http.Request, pool *ws.Pool) {
	// upgrade connection to ws
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})
	if err != nil {
		log.Println("Err : ", err.Error())
	}

	// close connection
	defer conn.Close(websocket.StatusInternalError, "internal error")

	client := ws.Client{
		Id:     uuid.NewString(),
		Conn:   conn,
		Pool:   pool,
		RoomId: r.URL.Path[len("/ws/"):],
	}

	// read
	pool.Join <- &client
	client.Read(r.Context())
}
