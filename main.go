package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hokaccha/go-prettyjson"
	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func CountingServer(ws *websocket.Conn) {
	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		ws.Write([]byte(fmt.Sprint(
			i,
		)))
	}
}

func BindToJsonServer(ws *websocket.Conn) {
	asMap := map[string]any{}
	if err := json.NewDecoder(ws).Decode(&asMap); err != nil {
		log.Println(err)
	}

	s, _ := prettyjson.Marshal(asMap)
	log.Println(string(s))

	CountingServer(ws) //? This? I just want this server to sending something back to the client
}

// This example demonstrates a trivial echo server.
func main() {
	// ? Fixing 403: https://stackoverflow.com/a/23324279
	// wsServer := websocket.Server{Handler: websocket.Handler(EchoServer)}
	// wsServer := websocket.Server{Handler: websocket.Handler(CountingServer)}
	// wsServer := websocket.Server{Handler: websocket.Handler(BindToJsonServer)}

	//* -------------------------------- NET HTTP -------------------------------- */
	// ? Ref: https://pkg.go.dev/golang.org/x/net/websocket
	// http.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
	// 	wsServer.ServeHTTP(w, req)
	// })
	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	panic("ListenAndServe: " + err.Error())
	// }
	//* ----------------------------------- GIN ---------------------------------- */
	//? Ref: https://github.com/gin-gonic/gin/issues/51#issuecomment-48201747
	r := gin.New()
	r.GET("/ws", func(c *gin.Context) {
		wsServer := websocket.Server{Handler: websocket.Handler(func(ws *websocket.Conn) {
			for {
				select {
				case <-c.Request.Context().Done(): //? Disconnect if the client is done connecting
				default:
					param := map[string]any{}
					if err := json.NewDecoder(ws).Decode(&param); err != nil {
						log.Println(err)
						return
					}
					ws.Write([]byte(param)) //? Send Response
				}
			}
		})}
		wsServer.ServeHTTP(c.Writer, c.Request)
	})
	r.Run(":8080")
}
