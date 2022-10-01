package stream

import (
	"sync"
	"testing"

	"github.com/rprtr258/go-flow/fun"
	"github.com/stretchr/testify/assert"
)

func nats() Stream[int] {
	return Generate(0, func(s int) int { return s + 1 })
}

func nats10() Stream[int] {
	return Take(nats(), 10)
}

var mul2 = func(i int) int { return i * 2 }

var isEven = func(i int) bool {
	return i%2 == 0
}

func TestStream(t *testing.T) {
	empty := NewStreamEmpty[int]()
	DrainAll(empty)

	res := CollectToSlice(Map(FromMany(10, 11, 12), mul2))
	assert.Equal(t, []int{20, 22, 24}, res)
}

func TestGenerate(t *testing.T) {
	powers2 := Generate(1, mul2)

	res := Head(powers2).Unwrap()
	assert.Equal(t, 1, res)

	res = Head(Skip(powers2, 9)).Unwrap()
	assert.Equal(t, 1024, res)
}

func TestRepeat(t *testing.T) {
	results := CollectToSlice(Take(Repeat(nats10()), 21))
	assert.Equal(t, results, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0})
}

func TestSum(t *testing.T) {
	sum := Sum(nats10())
	assert.Equal(t, 45, sum)
}

func TestFlatMap(t *testing.T) {
	pipe := FlatMap(nats10(), func(i int) Stream[int] {
		return Map(nats10(), func(j int) int {
			return i + j
		})
	})
	sum := Sum(Filter(pipe, func(i int) bool {
		return i%2 == 0
	}))
	assert.Equal(t, 450, sum)
}

func TestFlatMap2(t *testing.T) {
	floatsNested := FlatMap(nats10(), func(i int) Stream[float32] {
		return FromMany(float32(i), float32(2*i))
	})
	floats := Sum(floatsNested)
	assert.Equal(t, float32(45+45*2), floats)
}

func TestChunks(t *testing.T) {
	natsBy10 := Chunked(Take(nats(), 19), 10)
	nats10to19 := Head(Skip(natsBy10, 1)).Unwrap()
	assert.Equal(t, []int{10, 11, 12, 13, 14, 15, 16, 17, 18}, nats10to19)
}

func TestForEach(t *testing.T) {
	powers2 := Generate(1, mul2)
	is := []int{}
	ForEach(Take(powers2, 5), func(i int) {
		is = append(is, i)
	})
	assert.Equal(t, []int{1, 2, 4, 8, 16}, is)
}

func TestFilter(t *testing.T) {
	sumEven := Sum(Filter(nats10(), isEven))
	assert.Equal(t, 20, sumEven)
}

func TestSet(t *testing.T) {
	intsDuplicated := FlatMap(nats10(), func(i int) Stream[int] {
		return Map(
			nats10(),
			func(j int) int { return i + j },
		)
	})
	intsSet := Unique(intsDuplicated)
	assert.Equal(t, 19, Count(intsSet))
}

func TestGroupBy(t *testing.T) {
	intsDuplicated := FlatMap(nats10(), func(i int) Stream[int] {
		return Map(nats10(), func(j int) int { return i + j })
	})
	intsGroups := Group(intsDuplicated, fun.Identity[int])
	assert.Equal(t, 19, len(intsGroups))
	for k, as := range intsGroups {
		assert.Equal(t, k, as[0])
	}
}

func TestGrouped(t *testing.T) {
	lastWindow := CollectToSlice(Skip(Chunked(nats10(), 3), 3))
	assert.Equal(t, lastWindow, [][]int{{9}})
}

func TestGroupByMapCount(t *testing.T) {
	counter := ToCounterBy(nats10(), isEven)
	assert.Equal(t, uint(5), counter[false])
	assert.Equal(t, uint(5), counter[true])
}

func TestChain(t *testing.T) {
	got := CollectToSlice(Chain(
		FromSlice([]int{1, 2}),
		FromSlice([]int{3, 4, 5}),
		FromSlice([]int{6}),
	))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestFlatten(t *testing.T) {
	got := CollectToSlice(Flatten(FromSlice([]Stream[int]{
		FromSlice([]int{1, 2}),
		FromSlice([]int{3, 4, 5}),
		FromSlice([]int{6}),
	})))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, got)
}

func TestIntersperse(t *testing.T) {
	got := CollectToSlice(Intersperse(FromSlice([]int{1, 2, 3, 4, 5}), 0))
	assert.Equal(t, []int{1, 0, 2, 0, 3, 0, 4, 0, 5}, got)
}

func TestIntersperseEmpty(t *testing.T) {
	got := CollectToSlice(Intersperse(NewStreamEmpty[int](), 0))
	assert.Equal(t, []int{}, got)
}

func TestIntersperseTwoElems(t *testing.T) {
	got := CollectToSlice(Intersperse(FromSlice([]int{1, 2}), 0))
	assert.Equal(t, []int{1, 0, 2}, got)
}

func TestSkip(t *testing.T) {
	got := CollectToSlice(Skip(FromSlice([]int{1, 2, 3}), 2))
	assert.Equal(t, []int{3}, got)
}

func TestSkipToEmpty(t *testing.T) {
	got := CollectToSlice(Skip(FromSlice([]int{1, 2, 3}), 100))
	assert.Equal(t, []int{}, got)
}

func TestFind(t *testing.T) {
	got := Find(
		FromSlice([]int{1, 2, 3, 4, 5}),
		func(x int) bool { return x%4 == 0 },
	)
	assert.Equal(t, fun.Some(4), got)
}

func TestFindNotFound(t *testing.T) {
	got := Find(
		FromSlice([]int{1, 2, 3}),
		func(x int) bool { return x%4 == 0 },
	)
	assert.Equal(t, fun.None[int](), got)
}

func TestTakeWhile(t *testing.T) {
	stream := TakeWhile(
		FromSlice([]int{2, 4, 6, 7, 8}),
		func(x int) bool { return x%2 == 0 },
	)
	got := CollectToSlice(stream)
	assert.Equal(t, []int{2, 4, 6}, got)
	assert.Equal(t, fun.None[int](), Head(stream))
}

func TestFilterMap(t *testing.T) {
	got := CollectToSlice(MapFilter(
		FromSlice([]int{2, 4, 6, 7, 8}),
		func(x int) fun.Option[int] {
			if x%2 == 1 {
				return fun.None[int]()
			}
			return fun.Some(x / 2)
		},
	))
	assert.Equal(t, []int{1, 2, 3, 4}, got)
}

func TestPaged(t *testing.T) {
	got := CollectToSlice(Paged(FromSlice([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7},
	})))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, got)
}

func collectStreamsConcurrently(streams []Stream[int]) [][]int {
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	n := len(streams)
	slices := make([][]int, n)
	wg.Add(n)
	for i, stream := range streams {
		stream := stream
		i := i
		go func() {
			defer wg.Done()

			slice := CollectToSlice(stream)

			mu.Lock()
			defer mu.Unlock()
			slices[i] = slice
		}()
	}
	wg.Wait()
	return slices
}

func TestScatterEvenly(t *testing.T) {
	n := uint(4)
	streams := ScatterEvenly(nats10(), n)
	got := collectStreamsConcurrently(streams)
	assert.Equal(t, [][]int{
		{0, 4, 8},
		{1, 5, 9},
		{2, 6},
		{3, 7},
	}, got)
}

func TestScatter(t *testing.T) {
	n := uint(4)
	streams := Scatter(nats10(), n)
	slices := collectStreamsConcurrently(streams)
	got := make([]int, 0)
	for _, slice := range slices {
		got = append(got, slice...)
	}
	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
}

func TestScatterCopy(t *testing.T) {
	n := 4
	streams := ScatterCopy(nats10(), n)
	got := collectStreamsConcurrently(streams)
	assert.Equal(t, [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
	}, got)
}

func TestScatterRoute(t *testing.T) {
	streams := ScatterRoute(nats10(), []func(uint, int) bool{
		func(i uint, _ int) bool { return i < 3 },    // first three to first stream
		func(_ uint, x int) bool { return x%2 == 0 }, // evens to second stream
		func(_ uint, x int) bool { return x%3 == 0 }, // multiples of three to third stream
		// rest to fourth stream
	})
	got := collectStreamsConcurrently(streams)
	assert.Equal(t, [][]int{
		{0, 1, 2},
		{4, 6, 8},
		{3, 9},
		{5, 7},
	}, got)
}

func TestGather(t *testing.T) {
	got := CollectToSlice(Gather([]Stream[int]{nats10(), nats10()}))
	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, got)
}

func TestRange(t *testing.T) {
	got := CollectToSlice(Range(0, 10, 3))
	assert.Equal(t, []int{0, 3, 6, 9}, got)
}
