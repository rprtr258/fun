package iter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/fun"
)

func assertStream[T any](t *testing.T, s Seq[T], expected []T) {
	t.Helper()
	assert.Equal(t, expected, ToSlice(s))
}

var (
	nats   = FromGenerator(0, func(s int) int { return s + 1 })
	nats10 = FromMany(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	mul2   = func(i int) int { return i * 2 }
	isEven = func(i int) bool { return i%2 == 0 }
)

func TestStream(t *testing.T) {
	assertStream(t, Map(FromMany(10, 11, 12), mul2), []int{20, 22, 24})
}

func TestGenerate(t *testing.T) {
	powers2 := FromGenerator(1, mul2)

	a, ok := Head(powers2)
	assert.True(t, ok)
	assert.Equal(t, 1, a)

	b, ok := Head(Skip(powers2, 9))
	assert.True(t, ok)
	assert.Equal(t, 512, b)
}

func TestRepeat(t *testing.T) {
	base := FromMany(0, 1, 2)
	assertStream(t, Take(Repeat(base), 7), []int{
		0, 1, 2,
		0, 1, 2,
		0,
	})
}

func TestNats10(t *testing.T) {
	assertStream(t, nats10, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestSum(t *testing.T) {
	sum := Sum(nats10)
	assert.Equal(t, 45, sum)
}

func TestFlatMap(t *testing.T) {
	nats3 := func(yield func(int) bool) bool {
		return yield(0) &&
			yield(1) &&
			yield(2)
	}
	assertStream(t, FlatMap(nats3, func(i int) Seq[int] {
		return func(yield func(int) bool) bool {
			return yield(i*3) &&
				yield(i*4) &&
				yield(i*5)
		}
	}), []int{
		0, 0, 0,
		3, 4, 5,
		6, 8, 10,
	})
}

func TestFlatMap2(t *testing.T) {
	floatsNested := FlatMap(nats10, func(i int) Seq[float32] {
		return FromMany(float32(i), float32(2*i))
	})
	floats := Sum(floatsNested)
	assert.Equal(t, float32(45+45*2), floats)
}

func TestChunks(t *testing.T) {
	// chunks cant be retained which doesnt allow to use assertStream
	i := 0
	_ = Chunked(Take(nats, 19), 10)(func(chunk []int) bool {
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
	powers2 := FromGenerator(1, mul2)
	is := []int{}
	ForEach(Take(powers2, 5), func(i int) {
		is = append(is, i)
	})
	assert.Equal(t, []int{1, 2, 4, 8, 16}, is)
}

func TestFilter(t *testing.T) {
	sumEven := Sum(Filter(nats10, isEven))
	assert.Equal(t, 20, sumEven)
}

func TestUnique(t *testing.T) {
	intsSet := Unique(FromMany(0, 0, 1, 1, 1, 1, 2, 0, 1, 2, 2, 1, 0))
	assert.Equal(t, 3, Count(intsSet))
}

func TestGroupBy(t *testing.T) {
	intsDuplicated := FlatMap(nats10, func(i int) Seq[int] {
		return Map(nats10, func(j int) int { return i + j })
	})
	intsGroups := Group(intsDuplicated, fun.Identity[int])
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
	assertStream(t, Chunked(nats10, 3), [][]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{9},
	})
}

func TestGroupByMapCount(t *testing.T) {
	counter := ToCounterBy(nats10, isEven)
	assert.Equal(t, 5, counter[false])
	assert.Equal(t, 5, counter[true])
}

func TestChain(t *testing.T) {
	got := ToSlice(Concat(
		FromMany(1, 2),
		FromMany(3, 4, 5),
		FromMany(6),
	))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestFlatten(t *testing.T) {
	got := ToSlice(Flatten(FromMany([]Seq[int]{
		FromMany(1, 2),
		FromMany(3, 4, 5),
		FromMany(6),
	}...)))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestIntersperse(t *testing.T) {
	got := ToSlice(Intersperse(FromMany(1, 2, 3, 4, 5), 0))
	assert.Equal(t, []int{1, 0, 2, 0, 3, 0, 4, 0, 5}, got)
}

func TestIntersperseEmpty(t *testing.T) {
	got := ToSlice(Intersperse(FromNothing[int](), 0))
	assert.Equal(t, []int{}, got)
}

func TestIntersperseTwoElems(t *testing.T) {
	got := ToSlice(Intersperse(FromMany(1, 2), 0))
	assert.Equal(t, []int{1, 0, 2}, got)
}

func TestSkip(t *testing.T) {
	assertStream(t, Skip(FromMany(1, 2, 3), 2), []int{3})
}

func TestSkipToEmpty(t *testing.T) {
	got := ToSlice(Skip(FromMany(1, 2, 3), 100))
	assert.Equal(t, []int{}, got)
}

func TestFind(t *testing.T) {
	got := Find(
		FromMany(1, 2, 3, 4, 5),
		func(x int) bool { return x%4 == 0 },
	)
	assert.Equal(t, fun.Some(4), got)
}

func TestFindNotFound(t *testing.T) {
	got := Find(
		FromMany(1, 2, 3),
		func(x int) bool { return x%4 == 0 },
	)
	assert.Equal(t, fun.None[int](), got)
}

func TestTakeWhile(t *testing.T) {
	stream := TakeWhile(
		FromMany(2, 4, 6, 7, 8),
		func(x int) bool { return x%2 == 0 },
	)
	assertStream(t, stream, []int{2, 4, 6})
}

func TestFilterMap(t *testing.T) {
	got := ToSlice(MapFilter(
		FromMany(2, 4, 6, 7, 8),
		func(x int) (int, bool) {
			if x%2 == 1 {
				return 0, false
			}
			return x / 2, true
		},
	))
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestPaged(t *testing.T) {
	assertStream(t, Paged(FromMany(
		[]int{1, 2, 3},
		[]int{4, 5, 6},
		[]int{7},
	)), []int{1, 2, 3, 4, 5, 6, 7})
}

// func collectStreamsConcurrently(streams []Stream[int]) [][]int {
// 	var (
// 		mu sync.Mutex
// 		wg sync.WaitGroup
// 	)
// 	n := len(streams)
// 	slices := make([][]int, n)
// 	wg.Add(n)
// 	for i, stream := range streams {
// 		stream := stream
// 		i := i
// 		go func() {
// 			defer wg.Done()

// 			slice := CollectToSlice(stream)

// 			mu.Lock()
// 			defer mu.Unlock()
// 			slices[i] = slice
// 		}()
// 	}
// 	wg.Wait()
// 	return slices
// }

// func TestScatterEvenly(t *testing.T) {
// 	n := 4
// 	streams := ScatterEvenly(nats10, n)
// 	got := collectStreamsConcurrently(streams)
// 	assert.Equal(t, [][]int{
// 		{0, 4, 8},
// 		{1, 5, 9},
// 		{2, 6},
// 		{3, 7},
// 	}, got)
// }

// func TestScatter(t *testing.T) {
// 	n := 4
// 	streams := Scatter(nats10, n)
// 	slices := collectStreamsConcurrently(streams)
// 	got := make([]int, 0)
// 	for _, slice := range slices {
// 		got = append(got, slice...)
// 	}
// 	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
// }

// func TestScatterCopy(t *testing.T) {
// 	n := 4
// 	streams := ScatterCopy(nats10, n)
// 	got := collectStreamsConcurrently(streams)
// 	assert.Equal(t, [][]int{
// 		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
// 		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
// 		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
// 		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
// 	}, got)
// }

// func TestScatterRoute(t *testing.T) {
// 	streams := ScatterRoute(nats10, []func(int, int) bool{
// 		func(i int, _ int) bool { return i < 3 },    // first three to first stream
// 		func(_ int, x int) bool { return x%2 == 0 }, // evens to second stream
// 		func(_ int, x int) bool { return x%3 == 0 }, // multiples of three to third stream
// 		// rest to fourth stream
// 	})
// 	got := collectStreamsConcurrently(streams)
// 	assert.Equal(t, [][]int{
// 		{0, 1, 2},
// 		{4, 6, 8},
// 		{3, 9},
// 		{5, 7},
// 	}, got)
// }

// func TestGather(t *testing.T) {
// 	got := CollectToSlice(Gather([]Stream[int]{nats10, nats10}))
// 	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
// }

func TestRange(t *testing.T) {
	assertStream(t, FromRange(0, 10, 3), []int{0, 3, 6, 9})
}

func TestNewStream(t *testing.T) {
	assertStream(t, FromInfiniteGenerator(func(yield func(int)) {
		yield(1)
		yield(2)
		yield(3)
		for i := 4; i <= 9; i++ {
			yield(i)
		}
	}), []int{1, 2, 3, 4, 5, 6, 7, 8, 9})
}
