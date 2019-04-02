package main

import "fmt"

func main() {
	data := "0"

	firstRes := SingleHash(data)
	secondRes := MultiHash(firstRes)

	fmt.Println(firstRes)
	fmt.Println(secondRes)
}

func SingleHash(data string) string {
	return DataSignerCrc32(data) + "~" + DataSignerCrc32(DataSignerMd5(data))
}

func MultiHash(data string) string {
	steps := []string{"0", "1", "2", "3", "4", "5"}
	res := ""

	for _, val := range steps {
		res += DataSignerCrc32(val + data)
	}

	return res
}

func CombineResults() {
	// @TODO
}
