package iter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun/iter"
)

func assertStream[T any](t *testing.T, s iter.Seq[T], expected []T) {
	t.Helper()
	assert.Equal(t, expected, iter.ToSlice(s))
}

var (
	nats   = iter.FromGenerator(0, func(s int) int { return s + 1 })
	nats10 = iter.FromMany(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	mul2   = func(i int) int { return i * 2 }
	isEven = func(i int) bool { return i%2 == 0 }
)

func TestStream(t *testing.T) {
	assertStream(t, iter.Map(iter.FromMany(10, 11, 12), mul2), []int{20, 22, 24})
}

func TestGenerate(t *testing.T) {
	powers2 := iter.FromGenerator(1, mul2)

	a, ok := iter.Head(powers2)
	assert.True(t, ok)
	assert.Equal(t, 1, a)

	b, ok := iter.Head(iter.Skip(powers2, 9))
	assert.True(t, ok)
	assert.Equal(t, 512, b)
}

func TestRepeat(t *testing.T) {
	base := iter.FromMany(0, 1, 2)
	assertStream(t, iter.Take(iter.Repeat(base), 7), []int{
		0, 1, 2,
		0, 1, 2,
		0,
	})
}

func TestNats10(t *testing.T) {
	assertStream(t, nats10, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestSum(t *testing.T) {
	sum := iter.Sum(nats10)
	assert.Equal(t, 45, sum)
}

func TestFlatMap(t *testing.T) {
	nats3 := func(yield func(int) bool) {
		_ = !yield(0) || !yield(1) || !yield(2)
	}
	assertStream(t, iter.FlatMap(nats3, func(i int) iter.Seq[int] {
		return func(yield func(int) bool) {
			_ = !yield(i*3) || !yield(i*4) || !yield(i*5)
		}
	}), []int{
		0, 0, 0,
		3, 4, 5,
		6, 8, 10,
	})
}

func TestFlatMap2(t *testing.T) {
	floatsNested := iter.FlatMap(nats10, func(i int) iter.Seq[float32] {
		return iter.FromMany(float32(i), float32(2*i))
	})
	floats := iter.Sum(floatsNested)
	assert.Equal(t, float32(45+45*2), floats)
}

func TestChunks(t *testing.T) {
	// chunks cant be retained which doesnt allow to use assertStream
	i := 0
	iter.Chunked(iter.Take(nats, 19), 10)(func(chunk []int) bool {
		switch i {
		case 0:
			assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, chunk)
		case 1:
			assert.Equal(t, []int{10, 11, 12, 13, 14, 15, 16, 17, 18}, chunk)
		default:
			assert.Fail(t, "unexpected chunk")
		}
		i++
		return true
	})
}

func TestForEach(t *testing.T) {
	powers2 := iter.FromGenerator(1, mul2)
	is := []int{}
	iter.ForEach(iter.Take(powers2, 5), func(i int) {
		is = append(is, i)
	})
	assert.Equal(t, []int{1, 2, 4, 8, 16}, is)
}

func TestFilter(t *testing.T) {
	sumEven := iter.Sum(iter.Filter(nats10, isEven))
	assert.Equal(t, 20, sumEven)
}

func TestUnique(t *testing.T) {
	intsSet := iter.Unique(iter.FromMany(0, 0, 1, 1, 1, 1, 2, 0, 1, 2, 2, 1, 0))
	assert.Equal(t, 3, iter.Count(intsSet))
}

func TestGroupBy(t *testing.T) {
	intsDuplicated := iter.FlatMap(nats10, func(i int) iter.Seq[int] {
		return iter.Map(nats10, func(j int) int { return i + j })
	})
	intsGroups := iter.Group(intsDuplicated, func(i int) int { return i })
	assert.Equal(t, 19, len(intsGroups))
	for k, as := range intsGroups {
		assert.Equal(t, k, as[0])
	}
}

func TestGrouped(t *testing.T) {
	/* chunks by 3, taking first 3 chunks:
	[0, 1, 2] <-
	[3, 4, 5] <-
	[6, 7, 8] <-
	[9]
	*/
	s := iter.
		Chunked(nats10, 3).
		Map(func(chunk []int) []int {
			return append([]int(nil), chunk...)
		})
	assertStream(t, s, [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{9},
	})
}

func TestGroupByMapCount(t *testing.T) {
	counter := iter.ToCounterBy(nats10, isEven)
	assert.Equal(t, 5, counter[false])
	assert.Equal(t, 5, counter[true])
}

func TestChain(t *testing.T) {
	got := iter.ToSlice(iter.Concat(
		iter.FromMany(1, 2),
		iter.FromMany(3, 4, 5),
		iter.FromMany(6),
	))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestFlatten(t *testing.T) {
	got := iter.ToSlice(iter.Flatten(iter.FromMany([]iter.Seq[int]{
		iter.FromMany(1, 2),
		iter.FromMany(3, 4, 5),
		iter.FromMany(6),
	}...)))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestIntersperse(t *testing.T) {
	got := iter.ToSlice(iter.Intersperse(iter.FromMany(1, 2, 3, 4, 5), 0))
	assert.Equal(t, []int{1, 0, 2, 0, 3, 0, 4, 0, 5}, got)
}

func TestIntersperseEmpty(t *testing.T) {
	got := iter.ToSlice(iter.Intersperse(iter.FromNothing[int](), 0))
	assert.Equal(t, []int{}, got)
}

func TestIntersperseTwoElems(t *testing.T) {
	got := iter.ToSlice(iter.Intersperse(iter.FromMany(1, 2), 0))
	assert.Equal(t, []int{1, 0, 2}, got)
}

func TestSkip(t *testing.T) {
	assertStream(t, iter.Skip(iter.FromMany(1, 2, 3), 2), []int{3})
}

func TestSkipToEmpty(t *testing.T) {
	got := iter.ToSlice(iter.Skip(iter.FromMany(1, 2, 3), 100))
	assert.Equal(t, []int{}, got)
}

func TestFind(t *testing.T) {
	got, ok := iter.Find(
		iter.FromMany(1, 2, 3, 4, 5),
		func(x int) bool { return x%4 == 0 },
	)
	assert.True(t, ok)
	assert.Equal(t, 4, got)
}

func TestFindNotFound(t *testing.T) {
	_, ok := iter.Find(
		iter.FromMany(1, 2, 3),
		func(x int) bool { return x%4 == 0 },
	)
	assert.False(t, ok)
}

func TestTakeWhile(t *testing.T) {
	stream := iter.TakeWhile(
		iter.FromMany(2, 4, 6, 7, 8),
		func(x int) bool { return x%2 == 0 },
	)
	assertStream(t, stream, []int{2, 4, 6})
}

func TestFilterMap(t *testing.T) {
	got := iter.
		FromMany(2, 4, 6, 7, 8).
		MapFilter(func(x int) (int, bool) {
			return x / 2, x%2 == 0
		}).
		ToSlice()
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestPaged(t *testing.T) {
	assertStream(t, iter.Paged(iter.FromMany(
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		[]int{7},
	)), []int{1, 2, 3, 4, 5, 6, 7})
}

func TestRange(t *testing.T) {
	assertStream(t, iter.FromRange(0, 10, 3), []int{0, 3, 6, 9})
}

func TestNewStream(t *testing.T) {
	assertStream(t, iter.FromInfiniteGenerator(func(yield func(int)) {
		yield(1)
		yield(2)
		yield(3)
		for i := 4; i <= 9; i++ {
			yield(i)
		}
	}), []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestCount_lenToSlice(t *testing.T) {
	for i := 0; i < 100; i++ {
		arr := [100]struct{}{}
		seq1 := iter.FromMany(arr[:i]...)
		seq2 := iter.FromMany(arr[:i]...)
		assert.Equal(t, iter.Count(seq1), len(seq2.ToSlice()))
	}
}
