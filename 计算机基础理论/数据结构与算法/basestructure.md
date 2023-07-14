# 1. 排序



## 1.1  归并排序

**自顶向下递归**，每次一分为二处理再合并

```go
func MergeSort(arr []int) {
	merge := func(l, r int) {
		temp := make([]int, r-l+1)
		mid := l + (r-l)/2
		i, j := l, mid+1
		p := 0
		for i <= mid && j <= r {
			if arr[i] <= arr[j] {
				temp[p] = arr[i]
				i++
			} else {
				temp[p] = arr[j]
				j++
			}
			p++
		}
		for i <= mid {
			temp[p] = arr[i]
			p++
			i++
		}
		for j <= r {
			temp[p] = arr[j]
			p++
			j++
		}

		copy(arr[l:r+1], temp)
	}
	var mergeSort func(l, r int)
	mergeSort = func(l, r int) {
		if l >= r {
			return
		}
		mid := l + (r-l)/2
		mergeSort(l, mid)
		mergeSort(mid+1, r)
		//优化
		if arr[mid] > arr[mid+1] {
			merge(l, r)
		}
	}
	mergeSort(0, len(arr)-1)
}
```

时间复杂度：O(NlogN)

空间复杂度：O(N)

### 题目

**Leetcode 51 逆序对**

在合并的时候，一旦右半部分较小，计数器累加当前左半部分剩余的数量

```go
var cnt int
	merge := func(l, r int) {
		temp := make([]int, r-l+1)
		mid := l + (r-l)/2
		i, j := l, mid+1
		p := 0
		for i <= mid && j <= r {
			if nums[i] <= nums[j] {
				temp[p] = nums[i]
				i++
			} else {
				temp[p] = nums[j]
				j++
                //只对归并排序做了一点点补充
				cnt += mid - i + 1
			}
			p++
		}
		for i <= mid {
			temp[p] = nums[i]
			p++
			i++
		}
		for j <= r {
			temp[p] = nums[j]
			p++
			j++
		}

		copy(nums[l:r+1], temp)
	}
	var mergeSort func(l, r int)
	mergeSort = func(l, r int) {
		if l >= r {
			return
		}
		mid := l + (r-l)/2
		mergeSort(l, mid)
		mergeSort(mid+1, r)
		//优化
		if nums[mid] > nums[mid+1] {
			merge(l, r)
		}

	}
	mergeSort(0, len(nums)-1)
	return cnt
```



## 1.2 快速排序

快排跟归并看起来都是将数组一分为二，只是切分不同

快排的一个子过程能确定靶元素的**最终位置**，这个特性能很好完成**寻找第K大元素**问题

### 单路

维护一个连个指针 `i`和`j`

`i`：用于遍历数组，记录当前正在检查的位置

`j`：用于分割两类数组段，满足`arr[l+1...j] < v ; arr[j+1...i) >= v`

```go
func QuickSort1(nums []int) {

	var quick func(l, r int) int
	var quickSort func(l, r int)

	//返回p, 使得arr[l...p-1] < arr[p] ; arr[p+1...r] > arr[p]
	quick = func(l, r int) int {

		//随机化
		ri := l + rand.Intn(r-l+1)
		nums[l], nums[ri] = nums[ri], nums[l]

		v := nums[l]
		j := l // arr[l+1...j] < v ; arr[j+1...i) >= v
		for i := l + 1; i <= r; i++ {
			if nums[i] < v {
				j++
				nums[i], nums[j] = nums[j], nums[i]
			}
		}
		nums[l], nums[j] = nums[j], nums[l]
		return j
	}
	quickSort = func(l, r int) {
		if l >= r {
			return
		}
		//先做一个子过程找到分割点，然后左右递归
		p := quick(l, r)
		quickSort(l, p-1)
		quickSort(p+1, r)
	}

	quickSort(0, len(nums)-1)
}
```



### 双路

维护两个指针`i`和`j`，分别从左和右往中间靠拢，满足`arr[l+1...i-1]<v arr[j+1...r]>v`直到一方<font color=red>越过</font>另一方，子过程停止

`i`：小于组的右开边界

`j`：大于组的左开边界

```go
func QuickSort2(nums []int) {
	var quick func(l, r int) int
	var quickSort func(l, r int)

	//返回p, 使得arr[l...p-1] < arr[p] ; arr[p+1...r] > arr[p]
	quick = func(l, r int) int {

		//随机化
		ri := l + rand.Intn(r-l+1)
		nums[l], nums[ri] = nums[ri], nums[l]

		v := nums[l]
		i, j := l+1, r
		for true {
			for i <= j && nums[j] >= v {
				j--
			}
			for i <= j && nums[i] <= v {
				i++
			}
			if i > j {
				break
			}
			nums[i], nums[j] = nums[j], nums[i]
		}
		nums[l], nums[j] = nums[j], nums[l]
		return j
	}
	quickSort = func(l, r int) {
		if l >= r {
			return
		}
		//先做一个子过程找到分割点，然后左右递归
		p := quick(l, r)
		quickSort(l, p-1)
		quickSort(p+1, r)
	}

	quickSort(0, len(nums)-1)
}
```





### 三路

最高效，每次都能将相等的元素一次性处理

维护三个指针`i`、`lt`和`gt`，满足`arr[l+1...lt]<v  arr[lt+1...i)=v  arr(gt...r]>v`

`i`：用于遍历，表示当前正在检查的元素

`lt`：小于组的右闭边界

`gt`：大于组的左开边界



```go
func QuickSort3(nums []int) {
	var quick func(l, r int) int
	var quickSort func(l, r int)

	//返回p, 使得arr[l...p-1] < arr[p] ; arr[p+1...r] > arr[p]
	quick = func(l, r int) int {

		//随机化
		ri := l + rand.Intn(r-l+1)
		nums[l], nums[ri] = nums[ri], nums[l]

		v := nums[l]
		i := l
		lt, gt := l, r
		for i <= gt {
			if nums[i] < v {
				lt++
				nums[lt], nums[i] = nums[i], nums[lt]
			} else if nums[i] > v {
				nums[gt], nums[i] = nums[i], nums[gt]
				gt--
				continue
			}
			i++
		}
		nums[l], nums[lt] = nums[lt], nums[l]
		return lt
	}
	quickSort = func(l, r int) {
		if l >= r {
			return
		}
		//先做一个子过程找到分割点，然后左右递归
		p := quick(l, r)
		quickSort(l, p-1)
		quickSort(p+1, r)
	}

	quickSort(0, len(nums)-1)
}
```

### 题目

**Leetcode 215  数组中第k大元素**

第一次找到第p大元素位置。若k>p，则从arr[p+1...n-1]中找；若k<p，则从arr[0...p-1]中找；直到k==p。

```go
func findKthLargest(nums []int, k int) int {
    var quick func(l, r int) (int, int)
	//返回p, 使得arr[l...p-1] < arr[p] ; arr[p+1...r] > arr[p]
	quick = func(l, r int) (int, int) {
        
		//随机化
		ri := l + rand.Intn(r-l+1)
		nums[l], nums[ri] = nums[ri], nums[l]

		v := nums[l]
		i := l
		lt, gt := l, r
		for i <= gt {
			if nums[i] < v {
				lt++
				nums[lt], nums[i] = nums[i], nums[lt]
			} else if nums[i] > v {
				nums[gt], nums[i] = nums[i], nums[gt]
				gt--
				continue
			}
			i++
		}
		nums[l], nums[lt] = nums[lt], nums[l]
		return lt, gt
	}
	n := len(nums)
    l,r:=0,n-1
	p, q := quick(l,r)
	k = n - k
	for true {
		if k < p {
            r=p
			p, q = quick(l, r)
		} else if k > q {
            l=q+1
			p, q = quick(l, r)
		} else {
			break
		}
	}
	return nums[q]

}
```

# 2. 堆（优先队列）

## 2.1 普通堆

最大堆：所有父节点的值都大于子节点

最小堆：所有父节点的值都小于子节点



<font color=red>最大堆</font>

```go
type Heap struct {
	tables   []int
	len, cap int
}

func NewHeap(cap int) *Heap {
	return &Heap{
		cap:    cap,
		tables: make([]int, cap+1),
	}
}
//shift up
func (h *Heap) Push(val int) {
	h.len++
	p := h.len
	h.tables[p] = val
	for p > 1 && h.tables[p] > h.tables[p/2] {
		h.tables[p], h.tables[p/2] = h.tables[p/2], h.tables[p]
		p = p / 2
	}
}

//shift down
func (h *Heap) Pop() int {
	val := h.tables[1]
	h.tables[1], h.tables[h.len] = h.tables[h.len], h.tables[1]
	h.len--
	p := 1
	for 2*p <= h.len {
		j := 2 * p
		if j+1 <= h.len && h.tables[j] < h.tables[j+1] {
			j = j + 1
		}
		if h.tables[p] < h.tables[j] {
			h.tables[p], h.tables[j] = h.tables[j], h.tables[p]
			p = j
		} else {
			break
		}
	}
	return val
}
```



## 2.2 索引堆

每个值都有一个唯一整型索引值

push或pop时，比较的是值，交换的是值对应的索引

索引堆比普通堆多一个功能：修改`change`

允许修改任意点的优先级并快速调整

```go
package heap

type IndexHeap struct {
	len, cap int
	tables   []int
	indexes  []int
	vers     []int
}

func NewIndexHeap(cap int) *IndexHeap {
	return &IndexHeap{
		cap:     cap,
		tables:  make([]int, cap+1),
		indexes: make([]int, cap+1),
		vers:    make([]int, cap+1),
	}
}
func (h *IndexHeap) Push(idx int, val int) {
	idx = idx + 1

	h.tables[idx] = val

	h.len++
	p := h.len
	h.indexes[p] = idx
	h.vers[h.indexes[p]] = p
	h.shiftup(p)
}
func (h *IndexHeap) Pop() (int, int) {
	idx, val := h.indexes[1], h.tables[h.indexes[1]]

	h.indexes[1], h.indexes[h.len] = h.indexes[h.len], h.indexes[1]
	h.vers[h.indexes[1]], h.vers[h.indexes[h.len]] = 1, h.len
	h.len--
	p := 1
	h.shiftdown(p)

	return idx - 1, val
}
func (h *IndexHeap) Contain(idx int) bool {
	return h.vers[idx+1] != 0
}
func (h *IndexHeap) GetValue(idx int) int {
	return h.tables[idx+1]
}
func (h *IndexHeap) Change(idx int, val int) {
	idx = idx + 1
	h.tables[idx] = val
	p := h.vers[idx]
	h.shiftup(p)
	h.shiftdown(p)
}

func (h *IndexHeap) shiftup(p int) {
	for p > 1 && h.tables[h.indexes[p]] > h.tables[h.indexes[p/2]] {
		h.indexes[p], h.indexes[p/2] = h.indexes[p/2], h.indexes[p]
		h.vers[h.indexes[p]], h.vers[h.indexes[p/2]] = p, p/2
		p = p / 2
	}
}

func (h *IndexHeap) shiftdown(p int) {
	for 2*p <= h.len {
		j := 2 * p
		if j+1 <= h.len && h.tables[h.indexes[j]] < h.tables[h.indexes[j+1]] {
			j = j + 1
		}
		if h.tables[h.indexes[p]] < h.tables[h.indexes[j]] {
			h.indexes[p], h.indexes[j] = h.indexes[j], h.indexes[p]
			h.vers[h.indexes[p]], h.vers[h.indexes[j]] = p, j
			p = j
		} else {
			break
		}
	}
}

```

# 5. 最小生成树

<font color=red>prime</font>

- 基于点
- 新加点，更新与该点相连的切边
- 数据结构：最小索引堆







```go
type Edge struct {
	s, d int
	w    int
}
type Prime struct {
	v         int
	graph     [][]Edge
	indexHeap *heap.IndexHeap
	visited   []bool
}

func NewPrime(g [][]Edge) *Prime {
	return &Prime{
		v:         len(g),
		graph:     g,
		indexHeap: heap.NewIndexHeap(len(g)),
		visited:   make([]bool, len(g)),
	}
}

func (p Prime) MinTree() int {

	minCost := 0

	for _, e := range p.graph[0] {
		p.indexHeap.Push(e.d, e.w)
	}
	p.visited[0] = true
	for p.indexHeap.Len() > 0 {
		idx, w := p.indexHeap.Pop()
		p.visited[idx] = true
		minCost += w
		for _, e := range p.graph[idx] {
			if !p.visited[e.d] {
				if !p.indexHeap.Contain(e.d) {
					p.indexHeap.Push(e.d, e.w)
				} else if e.w < p.indexHeap.GetValue(e.d) {
					p.indexHeap.Change(e.d, e.w)
				}
			}
		}
	}
	return minCost
}
```



<font color=red>krush</font>

- 基于边
- 新加不成环的最小边
- 数据结构：最小堆、并查集



```go
type Dijkstra struct {
	v     int
	graph [][]Edge
	indexHeap *heap.IndexHeap
	visited   []bool
	path      []int
}

func NewDijkstra(g [][]Edge) *Dijkstra {
	v := len(g)
	path := make([]int, v)
	for i := 0; i < v; i++ {
		path[i] = -1
	}
	return &Dijkstra{
		v:     v,
		graph: g,
		indexHeap: heap.NewIndexHeap(v),
		visited:   make([]bool, v),
		path:      path,
	}
}

func (d *Dijkstra) MinPath(start int) []int {
	d.path[start] = start
	d.visited[start] = true

	for _, e := range d.graph[start] {
		d.path[e.d] = start
		d.indexHeap.Push(e.d, e.w)
	}
	for d.indexHeap.Len() > 0 {
		id, cost := d.indexHeap.Pop()
		d.visited[id] = true
		//松弛
		for _, ee := range d.graph[id] {
			if !d.visited[ee.d] {
				if !d.indexHeap.Contain(ee.d) {
					d.indexHeap.Push(ee.d, ee.w+cost)
					d.path[ee.d] = id
				} else if ee.w+cost < d.indexHeap.GetValue(ee.d) {
					d.indexHeap.Change(ee.d, ee.w+cost)
					d.path[ee.d] = id
				}
			}
		}
	}
	return d.path
}
```



# 6. 单源最短路径

<font color=red>dijkstra</font>

- 基于点
- 新增点，更新以该点为中转点的未访问点的消耗cost以及前继点
- 数据结构：最小索引堆



```go
type Dijkstra struct {
	v         int
	graph     [][]Edge
	indexHeap *heap.IndexHeap
	visited   []bool
	path      []int
}

func NewDijkstra(g [][]Edge) *Dijkstra {
	v := len(g)
	path := make([]int, v)
	for i := 0; i < v; i++ {
		path[i] = -1
	}
	return &Dijkstra{
		v:         v,
		graph:     g,
		indexHeap: heap.NewIndexHeap(v),
		visited:   make([]bool, v),
		path:      path,
	}
}

func (d *Dijkstra) MinPath(start int) []int {
	d.path[start] = start
	d.visited[start] = true

	for _, e := range d.graph[start] {
		d.path[e.d] = start
		d.indexHeap.Push(e.d, e.w)
	}
	for d.indexHeap.Len() > 0 {
		id, cost := d.indexHeap.Pop()
		d.visited[id] = true
		//松弛
		for _, ee := range d.graph[id] {
			if !d.visited[ee.d] {
				if !d.indexHeap.Contain(ee.d) {
					d.indexHeap.Push(ee.d, ee.w+cost)
					d.path[ee.d] = id
				} else if ee.w+cost < d.indexHeap.GetValue(ee.d) {
					d.indexHeap.Change(ee.d, ee.w+cost)
					d.path[ee.d] = id
				}
			}
		}
	}
	return d.path
}
```

