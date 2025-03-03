package main

import (
	"net"
	"sync"
	"time"
)

var rooms *Rooms
var laddrUDP *net.UDPAddr
var laddrTCP *net.TCPAddr

var tracker map[string]chan int
var playerQueue chan *Player

var locker *sync.Mutex

var serverStart time.Time

const maxRooms = 32

func init() {
	laddrTCP = &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8000,
	}
	laddrUDP = &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}

	ls := make([]chan int, maxRooms)
	for i := 0; i < maxRooms; i++ {
		ls[i] = make(chan int, 1)
	}

	rs := make([]*Room, maxRooms, maxRooms)

	lo := make(chan int, 1)

	rooms = &Rooms{
		LockSlice:  ls,
		RoomSlice:  rs,
		Offset:     0,
		LockOffset: lo,
	}

	tracker = make(map[string]chan int, maxRooms*2)
	playerQueue = make(chan *Player, 64)
	locker = &sync.Mutex{}

	serverStart = time.Now()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go UDPServer(&wg)
	go TCPServer(&wg)
	go sweeper(&wg)

	wg.Wait()

}
