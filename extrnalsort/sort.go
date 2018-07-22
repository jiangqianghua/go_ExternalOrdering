package main
import(
	"fmt"
	"os"
	"github.com/myproject/bingfa/pipeline"
	"bufio"
	"strconv"
	)
/**
func main(){
	infile := "small.in" 
	outfile := "small.out"
	p := createPipeline(infile,512,4)
	writeToFile(p,outfile)
	//printFile(outfile)
}
**/

func main(){
	infile := "small.in" 
	outfile := "small.out"
	p := createNetworkPipeline(infile,512,4)
	writeToFile(p,outfile)
	printFile(outfile)
}

func printFile(filename string){
	file , err := os.Open(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file,-1)

	for v := range p{
		fmt.Println(v)
	}
}

func writeToFile(p <-chan int , filename string){
	file ,err := os.Create(filename)
	if err != nil{
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer,p)
}

func createPipeline(filename string ,fileSize,chunkCount int) <-chan int{
	chunkSize := fileSize / chunkCount
	pipeline.Init()
	sortResults := []<-chan int{}
	for i := 0 ; i < chunkCount ; i++{
		file ,err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize),0)
		source := pipeline.ReaderSource(
			bufio.NewReader(file),chunkSize)
		sortResults = append(sortResults,
			pipeline.InMemSort(source))
	}
	return pipeline.MergeN(sortResults...)
}

// 以下是网络版本
func createNetworkPipeline(filename string ,fileSize,chunkCount int) <-chan int{
	chunkSize := fileSize / chunkCount
	pipeline.Init()
	sortAddr := []string{}
	for i := 0 ; i < chunkCount ; i++{
		file ,err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		file.Seek(int64(i * chunkSize),0)
		source := pipeline.ReaderSource(
			bufio.NewReader(file),chunkSize)
		addr := ":"+strconv.Itoa(7000+i)
		// 塞给网络服务器
		pipeline.NetworkSink(addr,pipeline.InMemSort(source))
		sortAddr = append(sortAddr,addr)
	}	
	// 从网络服务器取
	sortResults :=[] <-chan int{}
	for _,addr := range sortAddr{
		sortResults = append(sortResults,
			pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortResults...)
}