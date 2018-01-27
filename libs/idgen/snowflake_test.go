package idgen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/panenming/go-im/libs/datastruct"
)

func TestNewIDGenerator(t *testing.T) {
	b := "\t\t\t"
	b2 := "\t\t\t\t\t"
	d := "====================================="

	//第一个生成器
	gentor1, err := NewIDGenerator().SetWorkerId(100).Init()
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	//第二个生成器
	gentor2, err := NewIDGenerator().
		SetTimeBitSize(48).
		SetSequenceBitSize(10).
		SetWorkerIdBitSize(5).
		SetWorkerId(30).Init()
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	fmt.Printf("%s%s%s\n", d, b, d)
	fmt.Printf("workerId=%d lastTimestamp=%d %s workerId=%d lastTimestamp=%d\n",
		gentor1.workerId, gentor1.lastMsTimestamp, b,
		gentor2.workerId, gentor2.lastMsTimestamp)
	fmt.Printf("sequenceBitSize=%d timeBitSize=%d %s sequenceBitSize=%d timeBitSize=%d\n",
		gentor1.sequenceBitSize, gentor1.timeBitSize, b,
		gentor2.sequenceBitSize, gentor2.timeBitSize)
	fmt.Printf("workerBitSize=%d sequenceBitSize=%d %s workerBitSize=%d sequenceBitSize=%d\n",
		gentor1.workerIdBitSize, gentor1.sequenceBitSize, b,
		gentor2.workerIdBitSize, gentor2.sequenceBitSize)
	fmt.Printf("%s%s%s\n", d, b, d)

	var ids []int64
	for i := 0; i < 100; i++ {
		id1, err := gentor1.NextId()
		if err != nil {
			t.Fatal(err)
			return
		}
		id2, err := gentor2.NextId()
		if err != nil {
			t.Fatal(err)
			return
		}
		ids = append(ids, id2)
		fmt.Printf("%d%s%d\n", id1, b2, id2)
	}

	//解析ID
	for _, id := range ids {
		ts, workerId, seq, err := gentor2.Parse(id)
		fmt.Printf("id=%d\ttimestamp=%d\tworkerId=%d\tsequence=%d\terr=%v\n",
			id, ts, workerId, seq, err)
	}
}

//多线程测试
func TestSnowFlakeIdGenerator_MultiThread(t *testing.T) {
	f := "snowflake.txt"
	//准备写入的文件
	fp, err := os.OpenFile(f, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	//初始化ID生成器，采用默认参数
	gentor, err := NewIDGenerator().SetWorkerId(100).Init()
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}

	//启动10个线程，出错就报出来
	for i := 0; i < 10; i++ {
		go func() {
			for {
				gid, err := gentor.NextId()
				if err != nil {
					t.Fatal(err)
				}
				n, err := fp.WriteString(fmt.Sprintf("%d\n", gid))
				if err != nil || n <= 0 {
					t.Fatal(err)
				}
			}
		}()
	}

	// 10s 生成9w+ id
	time.Sleep(60 * time.Second)
	//time.Sleep(600 * time.Second)
}

func TestSameInFile(t *testing.T) {
	// 记录一共生成多少id
	count := 0
	f := "snowflake.txt"
	//准备写入的文件
	fp, err := os.OpenFile(f, os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	// 将file中的数据装入set中
	set := datastruct.NewSet()
	rd := bufio.NewReader(fp)
	for {
		line, err := rd.ReadString('\n') //以'\n'为结束符读入一行

		if err != nil || io.EOF == err {
			break
		}
		set.Add(line)
		count++
	}

	if set.Len() != count {
		fmt.Println("set len = ", set.Len(), " count = ", count, " 两者不匹配，说明有相同的id")
		t.Fatal("有相同的id")
	}
}
