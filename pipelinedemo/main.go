package main

import(
	"fmt"
	"os"
	"github.com/myproject/bingfa/pipeline"
	"bufio"
	)
/**
func main(){
	p := pipeline.ArraySource(3,2,6,7,4)
	for{
		// 错误写法 从chan p中取数据，
		// num := <- p 
		// fmt.Println(num)

		// if num,ok := <-p ; ok{
		// 	fmt.Println(num)
		// }else{
		// 	break;
		// }
		p <- 1;
		// 另外一种写法，上面if写法比较复杂
		for num := range p{
			fmt.Println(num)
		}
	}
}
**/
/**
func main(){
	arr1 := pipeline.ArraySource(3,2,6,7,4)
	arr2 := pipeline.ArraySource(0,2,7,2,5,1,2,9,4)
	arrSort1 :=pipeline.InMemSort(arr1)
	arrSort2 := pipeline.InMemSort(arr2)

	p := pipeline.Merge(arrSort1,arrSort2)

	// p := pipeline.InMemSort(
	// 	pipeline.ArraySource(3,2,6,7,4))

	for v := range p {
		fmt.Println(v)
	}
}**/

func main(){
	fmt.Println("start")
	const n = 64
	const filename = "small.in"
	// 创建一个输入的文件
	file , err := os.Create(filename)
	if err != nil{
		panic(err)
	}

	// 随机生成50个数，存放到p的管道中
	p := pipeline.RandomSource(n)
	// 写入文件
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer,p)
	writer.Flush()

	//  从文件中读取数据
	file , err = os.Open(filename)
	p = pipeline.ReaderSource(file,-1)
	count := 0 
	for v := range p{
		fmt.Println(v)
		count++;
		if count >= 100{
			break
		}
	}
	defer file.Close();
}