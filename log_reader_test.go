package tbft

import (
	"fmt"
	"github.com/hpcloud/tail"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestReadSynSystemLog(t *testing.T) {
	id := 0
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
	nowTime := time.Now()

	for line := range tails.Lines {
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
		level := parts[1]
		message := strings.Join(parts[3:], " ")
		fmt.Printf("enter fuzzing state\n")
		// 解析日志级别
		if strings.Contains(level, "ERROR") && !isFilterError(message) {
			//处理日志,将错误日志打印出来,并找到匹配的对象加入corpors
			fmt.Printf("fuzzing: Error occurred:%s\n", message)
		} else {
			n := LastState(id)
			var height, round int
			var step string
			parts = strings.Split(message, " ")
			if strings.Contains(message, "attempt enter new height to") {
				str := parts[len(parts)-1]
				str = str[1 : len(str)-2]
				fmt.Println(str)
				height, _ = strconv.Atoi(str)
				round = -1
				step = "NEW_HEIGHT"
				if uint64(height) != n.Height+1 || n.Step != "COMMIT" {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else if strings.Contains(message, "attempt enterNewRound to") {
				str := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(str, "/")
				height, _ = strconv.Atoi(parts[0])
				round, _ = strconv.Atoi(parts[1])
				step = "NEW_ROUND"
				if uint64(height) != n.Height || int32(round) != n.Round+1 || (n.Step != "PROPOSE" && n.Step != "PREVOTE" && n.Step != "PRECOMMIT" && n.Step != "NEW_HEIGHT") {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else if strings.Contains(message, "attempt enterPropose to") {
				str := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(str, "/")
				height, _ = strconv.Atoi(parts[0])
				round, _ = strconv.Atoi(parts[1])
				step = "PROPOSE"
				if uint64(height) != n.Height || int32(round) != n.Round || n.Step != "NEW_ROUND" {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else if strings.Contains(message, "enter prevote") {
				str := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(str, "/")
				height, _ = strconv.Atoi(parts[0])
				round, _ = strconv.Atoi(parts[1])
				step = "PREVOTE"
				if uint64(height) != n.Height || int32(round) != n.Round || n.Step != "PROPOSE" {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else if strings.Contains(message, "enter precommit") {
				str := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(str, "/")
				height, _ = strconv.Atoi(parts[0])
				round, _ = strconv.Atoi(parts[1])
				step = "PRECOMMIT"
				if uint64(height) != n.Height || int32(round) != n.Round || n.Step != "PREVOTE" {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else if strings.Contains(message, "enter commit") {
				str := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(str, "/")
				height, _ = strconv.Atoi(parts[0])
				round, _ = strconv.Atoi(parts[1])
				step = "COMMIT"
				if uint64(height) != n.Height || int32(round) != n.Round || n.Step != "PRECOMMIT" {
					fmt.Printf("fuzzing: Wrong state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
				}
			} else {
				continue
			}
			n = NodeState{
				time:   nowTime,
				Height: uint64(height),
				Round:  int32(round),
				Step:   step,
			}
			if n.time.After(LastState(id).time.Add(time.Minute * 5)) {
				fmt.Printf("fuzzing: Delayed state: [%d/%d/%s] after %s\n", height, round, step, n.ToString())
			}
			fmt.Printf("fuzzing: %d:%s\n", id, n.ToString())
			StateLists[id] = append(StateLists[id], n)
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
	TestReadSynSystemLog(t)
}
