package impl

type FibonacciRPC struct{}

func (t *FibonacciRPC) Fibo(n *int, reply *int) error {
	*reply = fibonacci(*n)
	return nil
}

func fibonacci(n int) int {
	ans := 1
	prev := 0
	for i := 1; i < n; i++ {
		temp := ans
		ans = ans + prev
		prev = temp
	}
	return ans
}
