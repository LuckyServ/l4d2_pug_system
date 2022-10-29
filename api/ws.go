package api

import (
	//"fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


func HttpReqWebSocket(c *gin.Context) {
	wshandler(c.Writer, c.Request);
	//fmt.Printf("%d : HttpReqWebSocket(c *gin.Context) ended\n", time.Now().Unix());
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {return true;};
	ws, err := upgrader.Upgrade(w, r, nil);
	if (err != nil) {
		//fmt.Printf("%s\n", err.Error());
		return;
	}
	//fmt.Printf("%d : Client connected\n", time.Now().Unix());
	/*err = ws.WriteMessage(1, []byte("Hi Client!"));
	if (err != nil) {
		fmt.Printf("%s\n", err.Error());
		return;
	}*/
	//reader(ws);
	for {
		err = ws.WriteMessage(1, []byte("trigger_status"));
		if (err != nil) {
			//fmt.Printf("%s\n", err.Error());
			return;
		}
		time.Sleep(1 * time.Second);
	}
}

/*func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage();
		if err != nil {
			//fmt.Printf("%s\n", err.Error());
			return;
		}
		// print out that message for clarity
		//fmt.Printf("%s\n", string(p));

		if err := conn.WriteMessage(messageType, p); err != nil {
			//fmt.Printf("%s\n", err.Error());
			return;
		}
	}
}*/