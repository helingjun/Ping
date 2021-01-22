package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

var wg sync.WaitGroup
var prefix string

func init() {

	flag.StringVar(&prefix, "prefix", "192.168", "ip prefix ")
}

// Ping 测试目标是否能达到
func Ping(dest string) bool {
	pinger, err := ping.NewPinger(dest)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Second * 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return false
	}
	return true
}

// IPchan 列表通道
var IPchan = make(chan string, 100)

var ipstats = make(chan string)

// GetIPs 获取IP列表
func GetIPs() {
	for x := 0; x < 255; x++ {
		for y := 0; y < 255; y++ {
			IPchan <- fmt.Sprintf("%s.%d.%d", prefix, x, y)
		}
	}
	close(IPchan)
}
func scan() {
	flag.Parse()
	// 并发数设置
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			for ip := range IPchan {
				stat := Ping(ip)
				if stat {
					ipstats <- fmt.Sprintf("IP:%s is online", ip)
				} else {
					ipstats <- fmt.Sprintf("IP:%s is offline", ip)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ipstats)
}
func main() {
	go GetIPs()
	go scan()
	for v := range ipstats {
		fmt.Println(v)
	}

}
