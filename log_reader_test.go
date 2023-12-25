package tbft

import (
	"fmt"
	"github.com/hpcloud/tail"
	"strings"
	"testing"
	"time"
)

func TestReadSynSystemLog(t *testing.T) {
	cfg := tail.Config{
		ReOpen: true, // 当文件被移动或删除后，tail 将尝试重新打开文件
		Follow: true, // 是否跟随,类似tail -f 命令
		Location: &tail.SeekInfo{ // 从文件的哪个位置开始读，默认从末尾开始读
			Offset: 0, //偏移量
			Whence: 2, //偏移起始位置，0-文件开始，1-当前位置，2-文件末尾
		},
		MustExist: false, //文件是否必须存在
		Poll:      true,  //文件是否轮询，true为轮询，false为inotify
		//inotify 是 Linux 系统中的一个功能，它允许应用程序监视文件系统的变化。
		//当指定的文件或目录发生变化时（例如，文件被修改、删除、移动，或者有新的文件被创建），inotify 可以通知应用程序
	}
	tails, _ := tail.TailFile("./log.txt", cfg)
	fmt.Println("tailed")

	layout := "2006-01-02 15:04:05.000"
	nowTime := time.Now()
	// 定期清理过期数据
	go func() {
		for {
			fmt.Println("clean")
			time.Sleep(time.Second) // 每分钟检查一次过期数据

			Requestcache.mu.Lock()
			for key, item := range Requestcache.store {
				if time.Now().After(item.Expiration) {
					delete(Requestcache.store, key)
				}
			}
			Requestcache.mu.Unlock()

			ParamsMap.mu.Lock()
			for key, item := range ParamsMap.store {
				if time.Now().After(item.Expiration) {
					delete(ParamsMap.store, key)
				}
			}
			ParamsMap.mu.Unlock()
		}
	}()
	for line := range tails.Lines {
		fmt.Printf("line: %s\n", line.Text)
		linetr := line.Text
		//跳过空行
		if linetr == "" {
			continue
		}
		// 解析日志内容
		parts := strings.Split(linetr, "\t")
		//如果不为标准日志格式，那么跳过
		if len(parts) < 3 {
			continue
		}

		//解析时间，如果不为第一个参数不为日期格式并且日期小于当前时间，那么不解析该条日志，直接跳过
		timestr := parts[0]
		t, err := time.Parse(layout, timestr)
		if err != nil || t.Before(nowTime) {
			continue
		}
		level := parts[1]
		message := strings.Join(parts[3:], " ")
		// 解析日志级别
		if strings.Contains(level, "ERROR") && !isFilterError(message) {
			//处理日志,将错误日志打印出来,并找到匹配的对象加入corpors
			//client.Logger.Infof("%+v", message)
			fmt.Println("error!!!!!!!!!!!!")
		} else if strings.Contains(level, "DEBUG") { //记录上次请求的request
			fmt.Println("debug!!!!!!!!!!!!")
		}
	}
}

func TestFun(t *testing.T) {
	//str := "2023-11-28 01:01:07.065\t[INFO]\t[Consensus] \u001B[31;1m@chain1\u001B[0m\tv2@v2.3.3/consensus_tbft_impl.go:1251\t[QmUCUpe9mBA1GHM2Xs9DBHK4chu5b1rQh3GJsCqXg5Qfqy](1/170/PROPOSE) receive invalid round qc from [QmcXU5wLwpUX2qi4V6fQE6FKL8YBqEb4boQzPmTSXtSqtK](1/169/4e696c48617368)"
	//flag, _ := regexp.MatchString("\\([0-9]*/[0-9]*/\\w*\\)", "(1/1/PREVOTE)--1")
	//parts := strings.Split(str, "\t")
	//message := strings.Join(parts[3:], " ")
	//l := strings.Index(message, "(")
	//r := strings.Index(message, ")")
	//fmt.Printf("l:%d, r:%d, str:[%s]\n", l, r, message[l+1:r])
	//fmt.Println(flag)
	fmt.Println(state[0].Round)
}
