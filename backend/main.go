package main

import (
	"chatapp/ws"
	"flag"
	"log"
	"net/http"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

var (
	PORT = flag.String("port", "8080", "app port")
)

func main() {
	flag.Parse()

	pool := ws.NewPool()
	go pool.Run()

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "./../frontend/web/index.html")
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		runApp(writer, request, pool)
	})

	log.Println("App running....")
	log.Fatalln(http.ListenAndServe(":"+*PORT, nil))
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
		Id:   uuid.NewString(),
		Conn: conn,
		Pool: pool,
	}

	// read
	pool.Join <- &client
	client.Read(r.Context())
}
