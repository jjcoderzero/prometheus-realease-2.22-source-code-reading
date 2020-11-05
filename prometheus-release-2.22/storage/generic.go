package storage

import (
	"github.com/prometheus/prometheus/pkg/labels"
)

type genericQuerier interface {
	LabelQuerier
	Select(bool, *SelectHints, ...*labels.Matcher) genericSeriesSet
}

type genericSeriesSet interface {
	Next() bool
	At() Labels
	Err() error
	Warnings() Warnings
}

type genericSeriesMergeFunc func(...Labels) Labels

type genericSeriesSetAdapter struct {
	SeriesSet
}

func (a *genericSeriesSetAdapter) At() Labels {
	return a.SeriesSet.At()
}

type genericChunkSeriesSetAdapter struct {
	ChunkSeriesSet
}

func (a *genericChunkSeriesSetAdapter) At() Labels {
	return a.ChunkSeriesSet.At()
}

type genericQuerierAdapter struct {
	LabelQuerier

	// One-of. If both are set, Querier will be used.
	q  Querier
	cq ChunkQuerier
}

func (q *genericQuerierAdapter) Select(sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) genericSeriesSet {
	if q.q != nil {
		return &genericSeriesSetAdapter{q.q.Select(sortSeries, hints, matchers...)}
	}
	return &genericChunkSeriesSetAdapter{q.cq.Select(sortSeries, hints, matchers...)}
}

func newGenericQuerierFrom(q Querier) genericQuerier {
	return &genericQuerierAdapter{LabelQuerier: q, q: q}
}

func newGenericQuerierFromChunk(cq ChunkQuerier) genericQuerier {
	return &genericQuerierAdapter{LabelQuerier: cq, cq: cq}
}

type querierAdapter struct {
	genericQuerier
}

type seriesSetAdapter struct {
	genericSeriesSet
}

func (a *seriesSetAdapter) At() Series {
	return a.genericSeriesSet.At().(Series)
}

func (q *querierAdapter) Select(sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) SeriesSet {
	return &seriesSetAdapter{q.genericQuerier.Select(sortSeries, hints, matchers...)}
}

type chunkQuerierAdapter struct {
	genericQuerier
}

type chunkSeriesSetAdapter struct {
	genericSeriesSet
}

func (a *chunkSeriesSetAdapter) At() ChunkSeries {
	return a.genericSeriesSet.At().(ChunkSeries)
}

func (q *chunkQuerierAdapter) Select(sortSeries bool, hints *SelectHints, matchers ...*labels.Matcher) ChunkSeriesSet {
	return &chunkSeriesSetAdapter{q.genericQuerier.Select(sortSeries, hints, matchers...)}
}

type seriesMergerAdapter struct {
	VerticalSeriesMergeFunc
}

func (a *seriesMergerAdapter) Merge(s ...Labels) Labels {
	buf := make([]Series, 0, len(s))
	for _, ser := range s {
		buf = append(buf, ser.(Series))
	}
	return a.VerticalSeriesMergeFunc(buf...)
}

type chunkSeriesMergerAdapter struct {
	VerticalChunkSeriesMergeFunc
}

func (a *chunkSeriesMergerAdapter) Merge(s ...Labels) Labels {
	buf := make([]ChunkSeries, 0, len(s))
	for _, ser := range s {
		buf = append(buf, ser.(ChunkSeries))
	}
	return a.VerticalChunkSeriesMergeFunc(buf...)
}

type noopGenericSeriesSet struct{}

func (noopGenericSeriesSet) Next() bool { return false }

func (noopGenericSeriesSet) At() Labels { return nil }

func (noopGenericSeriesSet) Err() error { return nil }

func (noopGenericSeriesSet) Warnings() Warnings { return nil }
