# 64行代码实现零拷贝go的TCP拆包粘包

### 前言

这段时间想用go写一个简单IM系统，就思考了一下go语言TCP的拆包粘包。TCP的拆包粘包有一般有三种解决方案。

#### 使用定长字节

实际使用中，少于固定字长的，要用字符去填充，空间使用率不够高。

#### 使用分隔符

一般用文本传输的，使用分隔符，IM系统一般对性能要求高，不推荐使用文本传输。

#### 用消息的头字节标识消息内容的长度

可以使用二进制传输，效率高，推荐。下面看看怎么实现。

### 尝试使用系统库自带的bytes.Buffer实现

代码实现：

```go
package tcp

import (
	"fmt"
	"net"
	"log"
	"bytes"
	"encoding/binary"
)

const (
	BYTES_SIZE uint16 = 1024
	HEAD_SIZE  int    = 2
)

func StartServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error listening", err.Error())
		return
	}
	for {
		conn, err := listener.Accept()
		fmt.Println(conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doConn(conn)
	}
}

func doConn(conn net.Conn) {
	var (
		buffer           = bytes.NewBuffer(make([]byte, 0, BYTES_SIZE))
		bytes            = make([]byte, BYTES_SIZE);
		isHead      bool = true
		contentSize int
		head        = make([]byte, HEAD_SIZE)
		content     = make([]byte, BYTES_SIZE)
	)
	for {
		readLen, err := conn.Read(bytes);
		if err != nil {
			log.Println("Error reading", err.Error())
			return
		}
		_, err = buffer.Write(bytes[0:readLen])
		if err != nil {
			log.Println("Error writing to buffer", err.Error())
			return
		}

		for {
			if isHead {
				if buffer.Len() >= HEAD_SIZE {
					_, err := buffer.Read(head)
					if err != nil {
						fmt.Println("Error reading", err.Error())
						return
					}
					contentSize = int(binary.BigEndian.Uint16(head))
					isHead = false
				} else {
					break
				}
			}
			if !isHead {
				if buffer.Len() >= contentSize {
					_, err := buffer.Read(content[:contentSize])
					if err != nil {
						fmt.Println("Error reading", err.Error())
						return
					}
					fmt.Println(string(content[:contentSize]))
					isHead = true
				} else {
					break
				}
			}
		}
	}
}

```

测试用例：

```go
package tcp

import (
	"testing"
	"net"
	"fmt"
	"encoding/binary"
)

func TestStartServer(t *testing.T) {
	StartServer("localhost:50002")
}

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", "localhost:50002")
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}
	var headSize int
	var headBytes = make([]byte, 2)
	s := "hello world"
	content := []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello go"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)

	s = "hello tcp"
	content = []byte(s)
	headSize = len(content)
	binary.BigEndian.PutUint16(headBytes, uint16(headSize))
	conn.Write(headBytes)
	conn.Write(content)
}

```

执行结果

```
127.0.0.1:51062
hello world
hello go
hello tcp
```

用go系统库的buffer，是不是感觉代码特别别扭，两大缺点

1.要写大量的逻辑代码，来弥补buffer对这个场景的不适用。

2.性能不高，有三次次内存拷贝，coon->[]byte->Buffer->[]byte。

### 自己实现

既然轮子不合适，就自己造轮子，首先实现一个自己的Buffer,很简单，只有六十几行代码，所有过程只有一次byte数组的拷贝，conn->buffer,剩下的全部操作都在原buffer的字节数组里面操作

```go
package tcp

import (
	"errors"
	"io"
)

type buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func newBuffer(reader io.Reader, len int) buffer {
	buf := make([]byte, len)
	return buffer{reader, buf, 0, 0}
}

func (b *buffer) Len() int {
	return b.end - b.start
}

//将有用的字节前移
func (b *buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0;
}

//从reader里面读取数据，如果reader阻塞，会发生阻塞
func (b *buffer) readFromReader() (int, error) {
	b.grow()
	n, err := b.reader.Read(b.buf[b.end:])
	if (err != nil) {
		return n, err
	}
	b.end += n
	return n, nil
}

//返回n个字节，而不产生移位
func (b *buffer) seek(n int) ([]byte, error) {
	if b.end-b.start >= n {
		buf := b.buf[b.start:b.start+n]
		return buf, nil
	}
	return nil, errors.New("not enough")
}

//舍弃offset个字段，读取n个字段
func (b *buffer) read(offset, n int) ([]byte) {
	b.start += offset
	buf := b.buf[b.start:b.start+n]
	b.start += n
	return buf
}
```

再看看怎样使用它,将上面的doConn函数改成这样就行了。

```go
func doConn(conn net.Conn) {
	var (
		buffer      = newBuffer(conn, 16)
		headBuf     []byte
		contentSize int
		contentBuf  []byte
	)
	for {
		_, err := buffer.readFromReader()
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			headBuf, err = buffer.seek(HEAD_SIZE);
			if err != nil {
				break
			}
			contentSize = int(binary.BigEndian.Uint16(headBuf))
			if (buffer.Len() >= contentSize-HEAD_SIZE) {
				contentBuf = buffer.read(HEAD_SIZE, contentSize)
				fmt.Println(string(contentBuf))
				continue
			}
			break
		}
	}
}
```

跑下测试用例，看下结果

```
127.0.0.1:51062
hello world
hello go
hello tcp
```

源码地址：https://github.com/alberliu/goim