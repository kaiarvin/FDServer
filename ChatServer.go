// ChatServer
package ChatServer

type Client struct {
	ID       int
	Name     string
	Icon     int
	Level    int
	Area     int
	Viplevel int
	Msg      chan string
	UserRelation
}

type UserRelation struct {
	FriendInfo map[uint64]string
	IgnorID    map[uint64]string
}

type ChatServer struct {
	ID          int
	Client_Pool []Client
}
