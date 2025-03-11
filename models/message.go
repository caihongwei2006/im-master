package models

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"encoding/json"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FromId   int //发送者
	TargetId int //接收者
	Type     int //消息类型  群聊 私聊 广播
	Media    int //消息类型 文字，图片，音频
	Content  string
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计

}

func (table *Message) TableName() string {
	return "message" //写死表名
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 将map的键类型从int64改为int
var clientMap map[int]*Node = make(map[int]*Node, 0)

var rwLocker sync.RWMutex

func Chat(w http.ResponseWriter, r *http.Request) {
	//111111111111获取参数，校验token合法性
	query := r.URL.Query()
	userIdStr := query.Get("userId")
	// 将userId从字符串转换为整型
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		fmt.Println("无效的userId:", err)
		return
	}
	//token:=query.Get("token")
	//messageType := query.Get("type")
	//targetIdStr := query.Get("targetId")
	//targetId, err := strconv.Atoi(targetIdStr)
	//context := query.Get("context")
	isvalid := true //这是checktoken的结果，先写死为true
	conn, err := (&websocket.Upgrader{
		//token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//2222222222222222222获取connection
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	//333333333用户关系
	//44444444userid和node绑定
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	//555555555555发送逻辑
	go sendProc(node)
	//6666666666666接收逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天室"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
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
		fmt.Println(string(data))
	}
}

var upsendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	upsendChan <- data
}

// 两个调度协程
func init() {
	go udpSendProc()
	go udprecevProc()
}

// 完成udp数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(172, 19, 223, 255), //这ip在wsl只有以太网，没有无线网
		Port: 3000,
	})
	if err != nil {
		fmt.Println("UDP连接失败:", err)
		return
	}
	defer conn.Close()
	for {
		select {
		case data := <-upsendChan:
			conn.Write(data)
		}
	}
	// 这里应该添加发送UDP数据的代码
}
func udprecevProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(172, 19, 223, 255),
		Port: 3000,
	})
	if err != nil {
		fmt.Println("UDP监听失败:", err)
		return
	}
	defer conn.Close()
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println("接收数据失败:", err)
			return
		}
		fmt.Println("收到数据:", string(buf[0:n]))
	}
}

// 后端调度逻辑
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("json解析失败:", err)
		return
	}
	switch msg.Type {
	case 1:
		sendMsg(msg.FromId, []byte{})
	case 2:
		sendGroupMsg(msg)
	case 3:
		sendAllMsg(msg)
	}
}
func sendMsg(userId int, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}

}
func sendGroupMsg(msg Message) {}
func sendAllMsg(msg Message)   {}
