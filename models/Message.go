package models

import (
	"encoding/json"
	"fmt"
	//"github.com/gin-gonic/gin/render"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
	"sync"
)

// Message 消息类型结构体
type Message struct {
	gorm.Model
	FromID   int64  //发送者
	TargetID int64  //接收者
	Type     int    //聊天类型：群聊、私聊、广播
	Media    int    //消息类型：文字、图片、音频
	Content  string //消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	//校验Tocken
	query := request.URL.Query()
	userId := query.Get("userId")
	userid, _ := strconv.ParseInt(userId, 10, 64)
	fmt.Println("userid=", userid)
	//tocken := query.Get("tocken")
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	//msgType := query.Get("type")
	isvalida := true //需从数据库获取进行tocken校验
	conn, err := (&websocket.Upgrader{
		//tocken校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//获取连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//获取用户关系
	//userid和node绑定并加锁
	rwLocker.Lock()
	clientMap[userid] = node
	rwLocker.Unlock()

	//完成发送的逻辑
	go sendProc(node)

	//完成接收的逻辑
	go recvProc(node)
}

func sendProc(node *Node) {
	for {
		fmt.Println("111发送消息")
		select {
		case data := <-node.DataQueue:
			fmt.Println("发送消息")
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("sendProc1", err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<  ", data)
	}
}

var udpsendChan = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成udp数据发送的协程
func udpSendProc() {
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(192, 168, 0, 6),
		Port: 3000,
	})
	if err != nil {
		fmt.Println("udpSendProc1", err)
		return
	}
	defer con.Close()
	for {
		select {
		case data := <-udpsendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println("udpSendProc2", err)
				return
			}

		}
	}
}

// 完成udp数据接收的协程
func udpRecvProc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println("udpRecvProc1", err)
		return
	}
	defer con.Close()
	for {
		var buffer [512]byte
		_, err := con.Read(buffer[0:])
		if err != nil {
			fmt.Println("udpRecvProc2", err)
			return
		}
		fmt.Println("udpRecvProc")
		dispatch(buffer[0:])
		fmt.Println("after dispatch......")
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("Unmarshal ", err)
		return
	}
	switch msg.Type {
	case 1: //私信
		fmt.Println("私信")
		sendMsg(msg.TargetID, data)
	case 2: //群发消息
		//sendGroupMsg()
	case 3: //广播消息
		//sendAllMsg()
	default:
		return
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLocker()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
