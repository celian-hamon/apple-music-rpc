package ipc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

var socket net.Conn

// Choose the right directory to the ipc socket and return it
func GetIpcPath() string {
	return "/private/var/folders/dh/z_mh8_c17lgg536h52vvvqcr0000gn/T/"
}

func CloseSocket() error {
	if socket != nil {
		socket.Close()
		socket = nil
	}
	return nil
}

// Read the socket response
func Read() string {
	buf := make([]byte, 512)
	payloadlength, err := socket.Read(buf)
	if err != nil {
		//fmt.Println("Nothing to read")
	}

	buffer := new(bytes.Buffer)
	for i := 8; i < payloadlength; i++ {
		buffer.WriteByte(buf[i])
	}

	return buffer.String()
}

// Send opcode and payload to the unix socket
func Send(opcode int, payload string) string {
	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(opcode))
	if err != nil {
		fmt.Println(err)
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(payload)))
	if err != nil {
		fmt.Println(err)
	}

	buf.Write([]byte(payload))
	_, err = socket.Write(buf.Bytes())
	if err != nil {
		fmt.Println(err)
	}

	return Read()
}
