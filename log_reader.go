package tbft

import (
	"chainmaker.org/chainmaker/protocol/v2"
	"github.com/hpcloud/tail"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var filterErrorList = []string{}

type Args struct {
	intArgs       [6]uint64
	stringArgs    [10]string
	methodIndex   int
	contractIndex int
}

type HashSet struct {
	m map[string]bool
}

var Address = [5]string{"../../chainmaker-v2.3.2-wx-org1.chainmaker.org/log/system.log",
	"../../chainmaker-v2.3.2-wx-org1.chainmaker.org/log/system.log",
	"../../chainmaker-v2.3.2-wx-org2.chainmaker.org/log/system.log",
	"../../chainmaker-v2.3.2-wx-org3.chainmaker.org/log/system.log",
	"../../chainmaker-v2.3.2-wx-org4.chainmaker.org/log/system.log",
}

type NodeState struct {
	time time.Time
	// current height
	Height uint64
	// current round
	Round int32
	// current step
	Step string
}

var state = &[5]NodeState{
	{
		time:   time.Time{},
		Height: 0,
		Round:  0,
		Step:   "",
	},
	{
		time:   time.Time{},
		Height: 0,
		Round:  0,
		Step:   "",
	},
	{
		time:   time.Time{},
		Height: 0,
		Round:  0,
		Step:   "",
	},
	{
		time:   time.Time{},
		Height: 0,
		Round:  0,
		Step:   "",
	},
	{
		time:   time.Time{},
		Height: 0,
		Round:  0,
		Step:   "",
	},
}

type KeyValue struct {
	Value      string
	Expiration time.Time
}

type Cache struct {
	mu    sync.Mutex
	store map[string]KeyValue
}

type ArgsValue struct {
	Value      Args
	Expiration time.Time
}

type ParamCache struct {
	mu    sync.Mutex
	store map[string]ArgsValue
}

func (c *Cache) Set(key string, value string, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = KeyValue{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, _ := c.store[key]

	return item.Value, true
}

func (c *ParamCache) SetParam(key string, value Args, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = ArgsValue{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
}

func (c *ParamCache) GetParam(key string) (Args, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, _ := c.store[key]

	return item.Value, true
}

// 保存txId和请求参数的对应关系
var ParamsMap = &ParamCache{
	store: make(map[string]ArgsValue),
}

// 保存txId和请求参数的对应关系
var errorSet = &HashSet{m: make(map[string]bool)}

// 保存txId和入参日志的对应关系
var Requestcache = &Cache{
	store: make(map[string]KeyValue),
}

func ReadSystemLog(id int, logger protocol.Logger) error {
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
	tails, _ := tail.TailFile(Address[id], cfg)

	layout := "2006-01-02 15:04:05.000"
	nowTime := time.Now()
	//// 定期清理过期数据
	//go func() {
	//	for {
	//		time.Sleep(time.Second) // 每分钟检查一次过期数据
	//
	//		Requestcache.mu.Lock()
	//		for key, item := range Requestcache.store {
	//			if time.Now().After(item.Expiration) {
	//				delete(Requestcache.store, key)
	//			}
	//		}
	//		Requestcache.mu.Unlock()
	//
	//		ParamsMap.mu.Lock()
	//		for key, item := range ParamsMap.store {
	//			if time.Now().After(item.Expiration) {
	//				delete(ParamsMap.store, key)
	//			}
	//		}
	//		ParamsMap.mu.Unlock()
	//	}
	//}()

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
			logger.Errorf("Error occurred:%s", message)
		} else {
			flag, _ := regexp.MatchString("\\([0-9]*/[0-9]*/\\w*\\)", message)
			//记录上次请求的request
			if flag {
				stateStr := message[strings.Index(message, "(")+1 : strings.Index(message, ")")]
				parts := strings.Split(stateStr, "/")
				height, _ := strconv.Atoi(parts[0])
				round, _ := strconv.Atoi(parts[1])
				step := parts[2]
				n := state[id]
				if n.time.Equal(time.Time{}) { // 第一次状态存储
					n.time = t
					n.Height = uint64(height)
					n.Round = int32(round)
					n.Step = step
				} else if strings.EqualFold(n.Step, step) && uint64(height) == n.Height && int32(round) == n.Round { // 状态未改变
					// 与进入该状态的时间差
					delta := t.Sub(n.time).Minutes()
					if delta > 5 { // 超过5min未更新状态
						logger.Errorf("Long time not changed: node:%d state:(%d/%d/%s)", id, height, round, step)
					}
				} else {
					n.time = t
					n.Height = uint64(height)
					n.Round = int32(round)
					n.Step = step
				}
			}
		}

	}
	return nil
}

/*
*
增加corpus
*/
//func addCorpus(f *testing.F, args Args) {
//	intparams := args.intArgs
//	stringparams := args.stringArgs
//	f.Add(args.contractIndex, args.methodIndex, stringparams[0], stringparams[1], stringparams[2], stringparams[3], stringparams[4], stringparams[5], stringparams[6], stringparams[7], stringparams[8], stringparams[9], intparams[0], intparams[1], intparams[2], intparams[3], intparams[4], intparams[5])
//}

/*
*
同步读取日志
*/
/*
func ReadSynSystemLog(client *sdk2.ChainClient, nowtime time.Time) (bool, error) {
	cfg := tail.Config{
		ReOpen: false, // 重新打开, 在单个日志文件写满做切隔时, 重新打开新一个文件
		Follow: false, // 开启不跟随，
		Location: &tail.SeekInfo{ // 从文件的哪个位置开始读，默认从末尾开始读，读最新的20行
			Offset: -30,
			Whence: 2,
		},
		MustExist: false,
		Poll:      true,
	}
	tails, _ := tail.TailFile(LoggerAddress, cfg)

	layout := "2006-01-02 15:04:05.000"
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

		//解析时间，如果不为第一个参数不为日期格式并且日期小于当前时间，那么不解析该条日志，直接跳过
		timestr := parts[0]
		t, err := time.Parse(layout, timestr)
		if err != nil || t.Before(nowtime) {
			continue
		}
		level := parts[1]
		message := strings.Join(parts[3:], " ")
		// 解析日志级别,判断是否
		if strings.Contains(level, "ERROR") && !isFilterError(message) {
			//处理日志,将错误日志打印出来

			client.Logger.Errorf("%+v", message)
			return true, nil
		}

	}

	return false, nil
}
*/

/*
*
先判断是否包含重复的错误，如果包含重复的错误，那么不再添加
判断是否包含过滤中的错误，如果是返回true,如果不是返回false
*/
func isFilterError(message string) bool {
	//过滤重复的错误
	if errorSet.Contains(message) {
		return true
	}
	for _, v := range filterErrorList {
		if strings.Contains(message, v) {
			return true
		}
	}
	errorSet.Add(message)
	return false
}

// Add 方法Add会返回一个bool类型的结果值，以表示添加元素值的操作是否成功。
// 方法Add的声明中的接收者类型是*HashSet。
func (set *HashSet) Add(e string) bool {
	if !set.m[e] { // 当前的m的值中还未包含以e的值为键的键值对
		set.m[e] = true // 将键为e(代表的值)、元素为true的键值对添加到m的值当中
		return true     // 添加成功
	}
	return false // 添加失败
}

// Contains 用于判断其值是否包含某个元素值。
// 这里判断结果得益于元素类型为bool的字段m
func (set *HashSet) Contains(e string) bool {
	return set.m[e]
}
