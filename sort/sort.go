package main

import(
	"sort"
	"fmt"
	)
// 一个简单的内排序
func main(){
	a := []int{3,4,2,1,5,6,4,3};
	sort.Ints(a)
	// for i, v := range a{
	// 	fmt.Println(i,v)
	// }

	for _, v := range a{
		fmt.Println(v)
	}
}