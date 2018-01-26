package disruptor

type CompositeBarrier []*Cursor

func NewCompositeBarrier(upsteam ...*Cursor) CompositeBarrier {
	if len(upsteam) == 0 {
		panic("At least one upstream cursor is required.")
	}

	cursors := make([]*Cursor, len(upsteam))
	copy(cursors, upsteam)
	return CompositeBarrier(cursors)
}

func (c CompositeBarrier) Read(noop int64) int64 {
	minimum := MaxSequenceValue
	for _, item := range c {
		sequence := item.Load()
		if sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
