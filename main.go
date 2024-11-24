package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	if ("gorutine" == os.Args[1]){
		main_gorutine()

	} else if ("channel" == os.Args[1]) {
		main_channel()

	} else if ("channel_buffer" == os.Args[1]) {
		main_channel_buffer()

	} else if ("channel_buffer2" == os.Args[1]) {
		main_channel_buffer2()

	} else if ("channel_direction" == os.Args[1]) {
		main_channel_direction()

	} else if ("channel_multiple" == os.Args[1]) {
		main_channel_multiple()

	} else {
		fmt.Println("ERROR! parameter 1:gorutine or channel or channel_buffer or main_channel_buffer2 or channel_direction or channel_multiple") 
	}
}

/// ゴールーチン（並行実行関数）
func main_gorutine() {
    start := time.Now()

    apis := []string{
        "https://management.azure.com",
        "https://dev.azure.com",
        "https://api.github.com",
        "https://outlook.office.com/",
        "https://api.somewhereintheinternet.com/",
        "https://graph.microsoft.com",
    }

	// update
	for _, api := range apis {
		go checkAPI_gorutine(api) // ここが同時実行される！
	}
    // for _, api := range apis {
    //     _, err := http.Get(api)
    //     if err != nil {
    //         fmt.Printf("ERROR: %s is down!\n", api)
    //         continue
    //     }

    //     fmt.Printf("SUCCESS: %s is up and running!\n", api)
    // }

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}

func checkAPI_gorutine(api string) {
    _, err := http.Get(api)
    if err != nil {
        fmt.Printf("ERROR: %s is down!\n", api)
        return
    }

    fmt.Printf("SUCCESS: %s is up and running!\n", api)
}

/// バッファーなしのチャネル
func main_channel() {
    start := time.Now()

    apis := []string{
        "https://management.azure.com",
        "https://dev.azure.com",
        "https://api.github.com",
        "https://outlook.office.com/",
        "https://api.somewhereintheinternet.com/",
        "https://graph.microsoft.com",
    }

	// チャンネル
	ch := make(chan string)

	for _, api := range apis {
		go checkAPI_channel(api,ch)
	}

	// update
	for i := 0; i < len(apis); i++ {
        fmt.Print(<-ch)
    }
	// fmt.Print(<-ch)
	// // add
	// fmt.Print(<-ch)
	// fmt.Print(<-ch)
	// fmt.Print(<-ch)
	// fmt.Print(<-ch)
	// fmt.Print(<-ch)

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}

func checkAPI_channel(api string, ch chan string) {
    _, err := http.Get(api)
    if err != nil {
        ch <- fmt.Sprintf("ERROR: %s is down!\n", api) // チャンネルを通してPrint
        return
    }

    ch <- fmt.Sprintf("SUCCESS: %s is up and running!\n", api) // チャンネルを通してPrint
}

/// バッファー有りのチャネル
func main_channel_buffer() {
	// size := 4
	size := 2
    ch := make(chan string, size)
    send_buffer(ch, "one")
    send_buffer(ch, "two")
    // send_buffer(ch, "three")
    // send_buffer(ch, "four")
    go send_buffer(ch, "three")
    go send_buffer(ch, "four")
    fmt.Println("All data sent to the channel ...")

    for i := 0; i < size; i++ {
        fmt.Println(<-ch)
    }

    fmt.Println("Done!")
}

func send_buffer(ch chan string, message string) {
    ch <- message
}

/// バッファーありのチャネル
func main_channel_buffer2() {
    start := time.Now()

    apis := []string{
        "https://management.azure.com",
        "https://dev.azure.com",
        "https://api.github.com",
        "https://outlook.office.com/",
        "https://api.somewhereintheinternet.com/",
        "https://graph.microsoft.com",
    }

	// チャンネル
	ch := make(chan string,10) // ここでサイズ指定

	for _, api := range apis {
		go checkAPI_channel(api,ch)
	}

	for i := 0; i < len(apis); i++ {
        fmt.Print(<-ch)
    }

    elapsed := time.Since(start)
    fmt.Printf("Done! It took %v seconds!\n", elapsed.Seconds())
}

/// チャネルの方向
func send_direction(ch chan<- string, message string) {
    fmt.Printf("Sending: %#v\n", message)
    ch <- message
}

func read_direction(ch <-chan string) {
    
	fmt.Printf("Receiving: %#v\n", <-ch)

	// チャンネルの方向が違うので、コンパイルエラーになる
	// fmt.Printf("Receiving: %#v\n", <-ch)
    // ch <- "Bye!"
}

func main_channel_direction() {
    ch := make(chan string, 1)
    send_direction(ch, "Hello World!")
    read_direction(ch)
}


/// 多重化
func main_channel_multiple() {
	ch1 := make(chan string)
    ch2 := make(chan string)
    go process_multiple(ch1)
    go replicate_multiple(ch2)

    for i := 0; i < 2; i++ {
        select {
        case process := <-ch1:
            fmt.Println(process)
        case replicate := <-ch2:
            fmt.Println(replicate)
        }
    }
}

func process_multiple(ch chan string) {
    time.Sleep(3 * time.Second)
    ch <- "Done processing!"
}

func replicate_multiple(ch chan string) {
    time.Sleep(1 * time.Second)
    ch <- "Done replicating!"
}