package storage

// lazyGenericSeriesSet is a wrapped series set that is initialised on first call to Next().
type lazyGenericSeriesSet struct {
	init func() (genericSeriesSet, bool)

	set genericSeriesSet
}

func (c *lazyGenericSeriesSet) Next() bool {
	if c.set != nil {
		return c.set.Next()
	}
	var ok bool
	c.set, ok = c.init()
	return ok
}

func (c *lazyGenericSeriesSet) Err() error {
	if c.set != nil {
		return c.set.Err()
	}
	return nil
}

func (c *lazyGenericSeriesSet) At() Labels {
	if c.set != nil {
		return c.set.At()
	}
	return nil
}

func (c *lazyGenericSeriesSet) Warnings() Warnings {
	if c.set != nil {
		return c.set.Warnings()
	}
	return nil
}

type warningsOnlySeriesSet Warnings

func (warningsOnlySeriesSet) Next() bool           { return false }
func (warningsOnlySeriesSet) Err() error           { return nil }
func (warningsOnlySeriesSet) At() Labels           { return nil }
func (c warningsOnlySeriesSet) Warnings() Warnings { return Warnings(c) }

type errorOnlySeriesSet struct {
	err error
}

func (errorOnlySeriesSet) Next() bool         { return false }
func (errorOnlySeriesSet) At() Labels         { return nil }
func (s errorOnlySeriesSet) Err() error       { return s.err }
func (errorOnlySeriesSet) Warnings() Warnings { return nil }
