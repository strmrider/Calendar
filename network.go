package calendar

import (
	"encoding/binary"
	"fmt"
	"net"
)

const( 
	host = "localhost"
	protocol = "tcp"
)

type Network struct {
	listener net.Listener
	// a callback to recieve updated serialized calendar
	getCalendar getCalendarInBytes

}

type getCalendarInBytes func() ([]byte)

func createNetwork(calendarHandler getCalendarInBytes) *Network{
	network:= new(Network)
	network.getCalendar = calendarHandler
	return network
}

func (network *Network) Listen(port int){
	var err error
	network.listener, err = net.Listen(protocol, host+":"+fmt.Sprint(port))
	if err == nil{
		fmt.Println("Created listener...")
		defer network.listener.Close()
		for{
			fmt.Println("Listening...")
			var conn net.Conn
			conn, err= network.listener.Accept()
			fmt.Println("Accepted connection from " + conn.RemoteAddr().String())
			syncwg.operate(wgAdd)
			go network.handleConnection(conn)
		}
	}
	syncwg.operate(wgDone)
}

func (network *Network) handleConnection(conn net.Conn){
	calendar := network.getCalendar()
	if calendar != nil{
		var size []byte = make([]byte, 4)
		binary.LittleEndian.PutUint32(size, uint32(len(calendar)))
		conn.Write(size)
		conn.Write(calendar)
	}
	syncwg.operate(wgDone)
}

func readData(conn net.Conn) []byte {
	var buffer []byte = make([]byte, 4)
	_, err := conn.Read(buffer)
	if err == nil{
		size := int32(binary.LittleEndian.Uint32(buffer))
		buffer = make([]byte, size)
		_, err := conn.Read(buffer)
		if err == nil{
			return buffer
		}
	}
	return nil
}

func ReadCalendarFromNetwork(host string, port int, alert reminderAlert) *Calendar{
	var server net.Conn
	var err error
	server, err = net.Dial(protocol, host+":"+fmt.Sprint(port))
	if err == nil{
		defer server.Close()
		data := readData(server)
		return Deserialize(data, alert)
	}
	return nil
}
