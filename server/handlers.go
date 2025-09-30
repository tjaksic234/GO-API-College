package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"College/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Server struct {
	Hub *Hub
}

func NewServer(hub *Hub) *Server {
	return &Server{Hub: hub}
}

func (s *Server) HandleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}
	defer ws.Close()

	client := &Client{conn: ws}

	for {
		var in IncomingMessage
		if err := ws.ReadJSON(&in); err != nil {
			// client likely disconnected; remove if present
			if client.RoomName != "" {
				s.Hub.RemoveClient(client)
			}
			log.Printf("read error: %v", err)
			return
		}

		switch in.Type {
		case "join":
			// ensure user & room exist in DB
			user, err := services.EnsureUser(in.Username)
			if err != nil {
				log.Printf("ensure user: %v", err)
				continue
			}
			room, err := services.EnsureRoom(in.Room)
			if err != nil {
				log.Printf("ensure room: %v", err)
				continue
			}
			// bind client
			client.Username = user.Username
			client.UserID = user.ID
			client.RoomName = room.Name
			client.RoomID = room.ID

			s.Hub.AddClient(client)

			// notify room (optional)
			s.Hub.broadcast <- OutgoingMessage{
				Type:     "message",
				Room:     room.Name,
				Username: "system",
				Message:  user.Username + " joined",
			}

		case "message":
			if client.RoomID == 0 || client.UserID == 0 {
				// not joined yet
				continue
			}
			// persist message
			_ = services.SaveMessage(client.RoomID, client.UserID, in.Message)

			// broadcast to room
			s.Hub.broadcast <- OutgoingMessage{
				Type:     "message",
				Room:     client.RoomName,
				Username: client.Username,
				Message:  in.Message,
			}

		case "leave":
			if client.RoomName != "" {
				s.Hub.RemoveClient(client)
				s.Hub.broadcast <- OutgoingMessage{
					Type:     "message",
					Room:     client.RoomName,
					Username: "system",
					Message:  client.Username + " left",
				}
			}
			return
		default:
			// ignore unknown
		}
	}
}

// ---------- REST endpoints ----------

// GET /rooms          -> all rooms from DB (persistent)
func (s *Server) GetAllRooms(c *gin.Context) {
	rooms, err := services.ListRooms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list rooms"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

func (s *Server) GetActiveRooms(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"rooms": s.Hub.ActiveRooms()})
}

func (s *Server) GetRoomMessages(c *gin.Context) {
	roomName := c.Param("room")
	room, err := services.FindRoomByName(roomName)
	if err != nil || room.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "room not found"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	msgs, err := services.GetRoomMessages(room.ID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
		return
	}

	type msgDTO struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		Room      string `json:"room"`
		Content   string `json:"content"`
		CreatedAt string `json:"created_at"`
	}
	out := make([]msgDTO, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, msgDTO{
			ID:        m.ID,
			Username:  m.User.Username,
			Room:      m.Room.Name,
			Content:   m.Content,
			CreatedAt: m.CreatedAt.UTC().Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, gin.H{"messages": out})
}

// GET /users          -> list users from DB
func (s *Server) GetAllUsers(c *gin.Context) {
	users, err := services.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}
