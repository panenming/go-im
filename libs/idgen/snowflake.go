package idgen

/**
 * SnowFlake算法，生成的共64位，用uint64即可表示，各位情况说明：
 * 第1位（不提供调整）：
 *		二进制中最高位为1的都是负数，但是我们生成的id一般都使用整数，所以这个最高位固定是0
 * 2-42位（41位，本库可调整）：
 *		用来记录时间戳（毫秒）。
 *		41位可以表示2^41−1个数字，如果只用来表示正整数（计算机中正数包含0）
 *		可以表示的数值范围是：0 至 2^41−1，减1是因为可表示的数值范围是从0开始算的，而不是1。
 *		也就是说41位可以表示2^41−1个毫秒的值，转化成单位年则是(2^41−1)/(1000∗60∗60∗24∗365)=69年
 * 43-52位（10位，本库可调整）：
 *		用来记录工作机器id，可以部署在2^10=1024个节点，包括5位datacenterId和5位workerId
 *		5位（bit）可以表示的最大正整数是25−1=31，即可以用0、1、2、3、....31这32个数字，
 *		来表示不同的datecenterId或workerId
 * 53-64位（12位，本库可调整）：
 *		序列号，用来记录同毫秒内产生的不同id。12位（bit）可以表示的最大正整数是212−1=4096
 *		即可以用0、1、2、3、....4095这4096个数字，来表示同一机器同一时间截（毫秒)内产生的4096个ID序号
**/

import (
	"fmt"
	"sync"
	"time"
)

// snowflake 结构体
type snowFlakeIdGenerator struct {
	workerId           int64 // 当前的workerId
	workerIdAfterShift int64 // 移位后的workerId，可直接和时间戳、序号取位或操作
	lastMsTimestamp    int64 // 上次用的时间戳
	curSequence        int64 // 当前的序号
	timeBitSize        uint8 // 时间戳占的位数，默认41位，最大不超过60位
	workerIdBitSize    uint8 // workerId 占的位数，默认是10，最大不超过60位
	sequenceBitSize    uint8 //序号占的位数，默认12，最大不超过60位

	lock               *sync.Mutex // 同步锁
	isHaveInit         bool        // 是否已经初始化
	maxWorkerId        int64       // workerId的最大值，初始化时计算得出
	maxSequence        int64       // 序列号的最大值，初始化时计算得出
	workerIdLeftShift  uint8       // 生成的workerId只取最低的几位，这里要左移，给序列号腾位，初始化时计算得出
	timestampLeftShift uint8       // 生成的时间戳左移几位，给workId、序列号腾位，初始化时计算得出
}

// 实例化一个ID生成器
func NewIDGenerator() *snowFlakeIdGenerator {
	return &snowFlakeIdGenerator{
		workerId:           0,
		lastMsTimestamp:    0,
		curSequence:        0,
		timeBitSize:        41,
		workerIdBitSize:    10,
		sequenceBitSize:    12,
		maxWorkerId:        0,
		maxSequence:        0,
		workerIdLeftShift:  0,
		timestampLeftShift: 0,
		lock:               new(sync.Mutex),
		isHaveInit:         false,
	}
}

//设置worker id
func (sfg *snowFlakeIdGenerator) SetWorkerId(w int64) *snowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.workerId = w
	return sfg
}

//设置时间戳占的位数
func (sfg *snowFlakeIdGenerator) SetTimeBitSize(n uint8) *snowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.timeBitSize = n
	return sfg
}

//设置worker id占的位数
func (sfg *snowFlakeIdGenerator) SetWorkerIdBitSize(n uint8) *snowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.workerIdBitSize = n
	return sfg
}

//设置序号占的位数
func (sfg *snowFlakeIdGenerator) SetSequenceBitSize(n uint8) *snowFlakeIdGenerator {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()
	sfg.isHaveInit = false
	sfg.sequenceBitSize = n
	return sfg
}

// 初始化操作
func (s *snowFlakeIdGenerator) Init() (*snowFlakeIdGenerator, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	// 如果已经初始化
	if s.isHaveInit {
		return s, nil
	}

	if s.sequenceBitSize < 1 || s.sequenceBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid sequence bit size, should (1,60)")
	}
	if s.timeBitSize < 1 || s.timeBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid time bit size, should (1,60)")
	}
	if s.workerIdBitSize < 1 || s.workerIdBitSize > 60 {
		return nil, fmt.Errorf("Init failed:\tinvalid worker id bit size, should (1,60)")
	}
	if s.workerIdBitSize+s.sequenceBitSize+s.timeBitSize != 63 {
		return nil, fmt.Errorf("Init failed:\tinvalid sum of all bit size, should eq 63")
	}

	// 确定移位数
	s.workerIdLeftShift = s.sequenceBitSize
	s.timestampLeftShift = s.sequenceBitSize + s.workerIdBitSize

	//确定序列号及workerid最大值
	s.maxWorkerId = -1 ^ (-1 << s.workerIdBitSize)
	s.maxSequence = -1 ^ (-1 << s.sequenceBitSize)

	//移位之后的workerId，返回结果时可直接跟时间戳、序号取或操作即可
	s.workerIdAfterShift = s.workerId << s.workerIdLeftShift

	// 判断当前workerid是否合法
	if s.workerId > s.maxWorkerId {
		return nil, fmt.Errorf("Init failed:\tinvalid worker id, should not greater than %d", s.maxWorkerId)
	}

	// 初始化完毕
	s.isHaveInit = true
	s.lastMsTimestamp = 0
	s.curSequence = 0
	return s, nil

}

//生成时间戳，根据bit size设置取高几位
//即，生成的时间戳先右移几位，再左移几位，就保留了最高的指定位数
func (sfg *snowFlakeIdGenerator) genTs() int64 {
	rawTs := time.Now().UnixNano()
	diff := 64 - sfg.timeBitSize
	ret := (rawTs >> diff) << diff
	return ret
}

//生成下一个时间戳，如果时间戳的位数较小，且序号用完时此处等待的时间会较长
func (sfg *snowFlakeIdGenerator) genNextTs(last int64) int64 {
	for {
		cur := sfg.genTs()
		if cur > last {
			return cur
		}
	}
}

// 生成下一个id
func (sfg *snowFlakeIdGenerator) NextId() (int64, error) {
	sfg.lock.Lock()
	defer sfg.lock.Unlock()

	//如果还没有初始化
	if !sfg.isHaveInit {
		return 0, fmt.Errorf("Gen NextId failed:\tplease execute Init() first")
	}

	//先判断当前的时间戳，如果比上一次的还小，说明出问题了
	curTs := sfg.genTs()
	if curTs < sfg.lastMsTimestamp {
		return 0, fmt.Errorf("Gen NextId failed:\tunknown error, the system clock occur some wrong")
	}

	//如果跟上次的时间戳相同，则增加序号
	if curTs == sfg.lastMsTimestamp {
		sfg.curSequence = (sfg.curSequence + 1) & sfg.maxSequence
		//序号又归0即用完了，重新生成时间戳
		if sfg.curSequence == 0 {
			curTs = sfg.genNextTs(sfg.lastMsTimestamp)
		}
	} else {
		//如果两个的时间戳不一样，则归0序号
		sfg.curSequence = 0
	}

	sfg.lastMsTimestamp = curTs

	//将处理好的各个位组装成一个int64型
	curTs = curTs | sfg.workerIdAfterShift | sfg.curSequence
	return curTs, nil
}

//解析生成的ID
func (sfg *snowFlakeIdGenerator) Parse(id int64) (int64, int64, int64, error) {
	//如果还没有初始化
	if !sfg.isHaveInit {
		return 0, 0, 0, fmt.Errorf("Parse failed:\tplease execute Init() first")
	}

	//先提取时间戳部分
	shift := sfg.sequenceBitSize + sfg.sequenceBitSize
	timestamp := (id & (-1 << shift)) >> shift

	//再提取workerId部分
	shift = sfg.sequenceBitSize
	workerId := (id & (sfg.maxWorkerId << shift)) >> shift

	//序号部分
	sequence := id & sfg.maxSequence

	//解析错误
	if workerId != sfg.workerId || workerId > sfg.maxWorkerId {
		fmt.Printf("workerBitSize=%d\tMaxWorkerId=%d\n", sfg.workerIdBitSize, sfg.maxWorkerId)
		return 0, 0, 0, fmt.Errorf("parse failed：invalid id, originWorkerId=%d\tparseWorkerId=%d\n",
			sfg.workerId, workerId)
	}
	if sequence < 0 || sequence > sfg.maxSequence {
		fmt.Printf("sequesnceBitSize=%d\tMaxSequence=%d\n", sfg.sequenceBitSize, sfg.maxSequence)
		return 0, 0, 0, fmt.Errorf("parse failed：invalid id, parseSequence=%d\n", sequence)
	}

	return timestamp, workerId, sequence, nil
}
