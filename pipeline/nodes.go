package pipeline

import(
	"sort"
	"io"
	"encoding/binary"
	"math/rand"
	"fmt"
	"time"
	)

var startTime time.Time

func Init(){
	startTime = time.Now()
}
//<- chan int 表示返回chan int类型， <- 表示返回出去后，只能取值，不能复值
func ArraySource(a ...int) <-chan int{
	out := make(chan int)
	go func(){
		for _, v := range a{
			out <- v
		}
		// 一定要close，close后，外面会用if或则range取判断取失败
		close(out)
	}()
	return out
}

// 内排序
//  <- 表示 in 只读,只出不进
func InMemSort(in <-chan int) <-chan int{
	out := make(chan int);
	go func(){
		a := []int{}
		for v := range in{
			a = append(a,v)
		}
		fmt.Println("Read done",time.Now().Sub(startTime))
		sort.Ints(a)
		fmt.Println("InMemSort done",time.Now().Sub(startTime))
		for _, v := range a{
			out <- v
		}
		close(out)
	}()
	return out
}
// 合并
func Merge(in1, in2 <-chan int) <-chan int{
	out :=make(chan int)
	go func(){
		v1,ok1 := <- in1
		v2,ok2 := <- in2
		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2){
				out <- v1
				v1 , ok1 = <- in1
			} else {
				out <- v2
				v2 , ok2 = <- in2
			}
		}
		close(out)
		fmt.Println("Merge done",time.Now().Sub(startTime))
		}()
	return out 
}
// 从文件读取数据
func ReaderSource(reader io.Reader ,chunkSize int) <- chan int{
	out := make(chan int)
	go func(){
		buffer := make([]byte,8)
		bytesRead := 0 
		for{
			n , err := reader.Read(buffer)
			bytesRead += n
			if n > 0{
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}

			if err != nil ||
			 (chunkSize != -1 &&
			 	bytesRead >= chunkSize){
				fmt.Println(err)
				break
			}
		}
		close(out)
	}()
	return out
}

// 向文件写数据
func WriterSink(writer io.Writer , in <-chan int){
	for v := range in {
		buffer := make([]byte ,8)
		binary.BigEndian.PutUint64(buffer,uint64(v))
		writer.Write(buffer)
	}
}


//生成随机数据源
func RandomSource(count int) <-chan int{
	out := make(chan int)
	go func(){
		for i := 0 ; i < count ; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out 
}

func MergeN(inputs ...<-chan int) <-chan int{
	if len(inputs) == 1{
		return inputs[0]
	}

	m := len(inputs) / 2
	return Merge(
		MergeN(inputs[:m]...),
		MergeN(inputs[m:]...))
}
