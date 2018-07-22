package main 
// how run
// go run ***.go  编译切执行，内存不会保存编译后的文件
// go build ***.go 或则 go build 会编译到当前文件下
// go install 把当前用到的所有go编译安装到go的bin目录下
import(
	"fmt"
	"time"
	)
func main(){	
	ch := make(chan string)
	// 调用函数前面加上go表示启动并行机制
	for i := 1 ; i < 10 ; i++{
		go HelloWordFunc(i,ch)
	}
	for{
		fmt.Println("开始等待ch的值")
		// 只要ch有值，会自动取出复制给msg，然后等待下一次的值
		msg := <- ch
		fmt.Println(msg)
	}
	// 延迟1毫秒
	time.Sleep(time.Millisecond)

}

// ch 类型是chan string
func HelloWordFunc(i int, ch chan string){
	//time.Sleep(5000*time.Millisecond)
	// 单独的for相当于while，不限循环
	for {
		// 把格式化的string复制给ch，ch可以一次性复制多个，外部取出会按照复制的顺序取
		fmt.Printf("开始复制给chan %d \n",i)
		ch <- fmt.Sprintf("hello world goroutine %d!\n",i);
	}
}