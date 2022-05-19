package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("服务已启动")
	for {
		conn, err := listener.Accept()
		fmt.Println("收到客户请求")
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

type readerError struct {
	io.Reader
	err error
}

func (r *readerError) Read(b []byte) (n int, err error) {
	if r.err != nil {
		return 0, err
	}

	n, r.err = r.Reader.Read(b)
	return n, r.err
}

func handleConn(conn net.Conn) {
	re := new(readerError)
	re.Reader = conn
	defer func() {
		if re.err != nil {
			conn.Close()
		}
	}()
	pkgLenBytes := make([]byte, 4)
	re.Read(pkgLenBytes)

	pkgLen := binary.BigEndian.Uint32(pkgLenBytes)

	//读取到整个包大小后，我们就可以整体把包读取出来
	pkgBytes := make([]byte, pkgLen-4) //因为已经读过4个字节了，所以前面四个字节被排除了
	re.Read(pkgBytes)

	headerLen := binary.BigEndian.Uint16(pkgBytes[:2])
	headerLen = headerLen - 4 // pkglength 就占头部四个字节，而它早就读出来了，所以头部长度先减去4

	protocolVersion := binary.BigEndian.Uint16(pkgBytes[2:4])

	operation := binary.BigEndian.Uint32(pkgBytes[4:8])

	seqId := binary.BigEndian.Uint32(pkgBytes[8:12])

	// 头部剩余字段
	extraHeader := pkgBytes[12:headerLen]

	body := pkgBytes[headerLen:]

	handleGoim(protocolVersion, operation, seqId, extraHeader, body)
}

func handleGoim(version uint16, operation, seqId uint32, extraHeader, body []byte) {
	// dosomethings
}
