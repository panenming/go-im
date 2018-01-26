package disruptor

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	cpuCacheLinePadding        = 7 // 用在cache line 填充，防止两个变量分配到一个cache line
)

type Cursor struct {
	sequence int64
	padding  [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{
		sequence: InitialSequenceValue,
	}
}
