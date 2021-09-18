package utility

import "github.com/bwmarrin/snowflake"

// GenerateId 使用方法：
// n 可以填写任意的数字,保证每个业务填写的n不一样
// gen := utility.GenerateId(1)
// id := <- gen
func GenerateId(n int64) <-chan string {
	out := make(chan string, 10)
	node, _ := snowflake.NewNode(n)
	go func() {
		for {
			id := node.Generate()
			out <- id.String()
		}
	}()
	return out
}

func GenerateIdInt(n int64) <-chan int64 {
	out := make(chan int64, 10)
	node, _ := snowflake.NewNode(n)
	go func() {
		for {
			id := node.Generate()
			out <- id.Int64()
		}
	}()
	return out
}
