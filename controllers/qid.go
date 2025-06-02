package controllers

import (
	"encoding/json"
	"sync"
	"time"

	log "github.com/astaxie/beego/logs"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"

	mf "github.com/gtlyy/myfun"
	"github.com/gtlyy/mytime"

	okex "github.com/gtlyy/myokx"

	"github.com/gtlyy/myrabbitmq"
)

// Error 处理函数
func IfError(msg string, err error) {
	if err != nil {
		log.Error("%s: %s\n", msg, err)
	}
}

// 客户端发回给后端用的数据结构
type SetupData struct {
	Usdt0     string `json:"usdt0"`
	Goal      string `json:"goal"`
	Fail      string `json:"fail"`
	Strategy  string `json:"strategy"`
	Bar       string `json:"bar"`
	AutoOpen  bool   `json:"autoopen"`
	AutoClose bool   `json:"autoclose"`
	AA        bool   `json:"aa"`
	BB        bool   `json:"bb"`
}

// 界面控制器
type QidController struct {
	beego.Controller
}

// 生成qid
func (c *QidController) Get() {
	log.Info("In QidController: Get().")
	qid := uuid.New().String()
	c.Data["json"] = map[string]interface{}{
		"qid": qid,
	}
	log.Info("Create qid:", qid)
	// 核心一步：每个会话（进入或刷新），根据qid开启一个单独的协程！
	go tradegame(qid)
	c.ServeJSON()
}

// tradegame 主程序
func tradegame(qid string) {
	var wg sync.WaitGroup
	log.Info("In qid: tradegame()")

	// 一些常量、变量
	usdt0 := 1000.0
	usdt := usdt0
	coin := 0.0
	perF := 0.95       // 每次投入的比例
	price := 0.0       // 开仓价或平仓价
	price_start := 0.0 // 开始价格
	price_now := 0.0   // 实时价格
	fee := 0.0
	per := 0.01
	feeRate := 0.001
	canBuy := true
	canSell := false
	goal := 0.01 * 0.6
	fail := -0.01 * 0.6
	slow, fast, signal := 9, 23, 14
	exchange := "amq.topic"
	exchange_type := "topic"
	timeout_minutes := 3 * time.Minute
	// 策略相关：
	strategy := "macd1"
	bar := "15m"
	autoopen := false
	autoclose := false
	condition := false
	isA := true
	isB := true

	// macd
	var m okex.MyMacdClass

	// 建立 rabbitmq, mariadb 连接
	var rabbit *myrabbitmq.RabbitMqClass
	rabbit = &myrabbitmq.RabbitMqClass{}
	err := rabbit.Init("userdw01", "pq328hu7", "127.0.0.1", "5672")
	IfError("Rabbitmq Init().", err)
	// mariadb
	var maria *okex.MyMariaDBClass
	maria = &okex.MyMariaDBClass{}
	maria.Init("tgame", "tc_3469Uk", "127.0.0.1", "3306", "testdb")

	// 设置rabbit，创建一系列的queque
	// setupstart
	queue_name, err := rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 1:"+queue_name, err)
	routing_key_setupstart := qid + "-setupstart"
	err = rabbit.Bind(exchange, exchange_type, routing_key_setupstart, queue_name)
	IfError("In qid tradegame(): Bind() 1:"+queue_name, err)
	queue_name_setupstart := queue_name
	// account
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 2:"+queue_name, err)
	routing_key_account := qid + "-account"
	err = rabbit.Bind(exchange, exchange_type, routing_key_account, queue_name)
	IfError("In qid tradegame(): Bind() 2:"+queue_name, err)
	// kline
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 3:"+queue_name, err)
	routing_key_kline := qid + "-kline"
	err = rabbit.Bind(exchange, exchange_type, routing_key_kline, queue_name)
	IfError("In qid tradegame(): Bind() 3:"+queue_name, err)
	// trade
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 4:"+queue_name, err)
	routing_key_trade := qid + "-trade"
	err = rabbit.Bind(exchange, exchange_type, routing_key_trade, queue_name)
	mf.IfError("In qid tradegame(): Bind() 4:"+queue_name, err)
	// tips
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 5:"+queue_name, err)
	routing_key_tips := qid + "-tips"
	err = rabbit.Bind(exchange, exchange_type, routing_key_tips, queue_name)
	IfError("In qid tradegame(): Bind() 5:"+queue_name, err)
	// over
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 6:"+queue_name, err)
	routing_key_over := qid + "-over"
	err = rabbit.Bind(exchange, exchange_type, routing_key_over, queue_name)
	IfError("In qid tradegame(): Bind() 6:"+queue_name, err)
	// tradeReply
	queue_name, err = rabbit.CreateQueueReturnName()
	IfError("In qid tradegame(): CreateQueue3() 6:"+queue_name, err)
	routing_key_tradeReply := qid + "-tradeReply"
	err = rabbit.Bind(exchange, exchange_type, routing_key_tradeReply, queue_name)
	IfError("In qid tradegame(): Bind() 6:"+queue_name, err)
	queue_name_tradeReply := queue_name

	log.Info("Bind routing_keys Done. keys: account, kline, trade, over, tips, tradeReply")

	// 开始主程序
	breakFor := false
	wg.Add(1)
	go func() {
		log.Info("Waiting for setupstart......")
		msgs, err := rabbit.Receive(queue_name_setupstart, true)
		IfError("In qid tradegame(): Receive() : queue_name_setupstart: ", err)
		msgs1, err1 := rabbit.Receive(queue_name_tradeReply, true) // for open, pass, close
		IfError("In qid tradegame(): Receive(): queue_name_tradeReply", err1)
		shouldExit := false
		for !shouldExit {
			log.Info("In for{}.")
			price_start = 0.0
			select {
			case d, ok := <-msgs:
				// 初始化账户：
				log.Info("Receive msgs: setupstart.")
				if !ok {
					log.Error("msgs 通道已关闭，退出 goroutine")
					return
				}

				var dataSetup SetupData
				err := json.Unmarshal(d.Body, &dataSetup)
				if err != nil {
					log.Info("In tradegame(): Unmarshal():", err)
					continue
				}

				usdt0, goal, fail = mf.StringToFloat64(dataSetup.Usdt0), mf.StringToFloat64(dataSetup.Goal), mf.StringToFloat64(dataSetup.Fail)
				usdt = usdt0
				coin = 0.0
				// 策略相关：
				strategy = dataSetup.Strategy
				bar = dataSetup.Bar
				autoopen = dataSetup.AutoOpen
				autoclose = dataSetup.AutoClose
				isA = dataSetup.AA
				isB = dataSetup.BB
				log.Info("Server setup Done.", usdt0, usdt, goal, fail)
				log.Info("bar, autoopen, autoclose:", bar, autoopen, autoclose)

				// 同步，返回给前端（其实可以直接在前端修改的啦）
				msg0 := `{"total":` + mf.Float64ToString(usdt+coin*price, 2) +
					`,"usdt":` + mf.Float64ToString(usdt, 2) +
					`,"coin":` + mf.Float64ToString(coin, mf.CountFloat(price)) + `}`
				// rabbit.Send("amq.topic", routing_key_tips, "hello")
				rabbit.Send("amq.topic", routing_key_account, msg0)
				log.Info("Init account Done.")

				// 获取测试数据
				r, stock, name := maria.CreateTradeGameData3(isA, isB, bar)
				// Big A
				if stock[0] >= '0' && stock[0] <= '9' {
					bar = "1D"
				}
				log.Info("trade id:", stock+"-"+bar)
				if name == "" {
					rabbit.Send(exchange, qid+"-title", stock+"-"+bar)
				} else {
					rabbit.Send(exchange, qid+"-title", stock+"-"+name+"-"+bar)
				}

				close1 := make([]float64, len(r))
				for i, v := range r {
					close1[i] = mf.StringToFloat64(v.C)
					price_now = close1[i]
					if price_start == 0.0 {
						price_start = close1[0]
					}
					m.Init(close1, fast, slow, signal)
					m.CalMacd()

					// 发送到html的部分：
					str_dif := mf.Float64ToString(m.Dif[i], 5)
					str_dea := mf.Float64ToString(m.Dea[i], 5)
					str_macd := mf.Float64ToString(m.Hist[i], 5)
					log.Debug(str_dif, str_dea, str_macd)

					msg := `{"dd":` + `"` + mytime.TsToStrCST(v.Ts, "2006-01-02T15:04") + `"` +
						`,"kk":` + `[` + v.O + `,` + v.C + `,` + v.L + `,` + v.H + `]` +
						`,"vv":` + mf.Float64ToString(mf.StringToFloat64(v.VolCcy), 0) +
						`,"macd":` + str_macd +
						`,"dif":` + str_dif +
						`,"dea":` + str_dea +
						`}`
					time.Sleep(2000 * time.Microsecond)
					rabbit.Send(exchange, routing_key_kline, msg)

					// 有33个数据后，才开始启用macd大法。
					if i < 33 {
						continue
					}

					if strategy == "strategy1" {
						condition = canBuy && m.CrossOver(i) && usdt > 300
					} else if strategy == "strategy2" {
						condition = canBuy && m.CrossOverTwice(i) && usdt > 300
					} else if strategy == "strategy3" {
						condition = canBuy && m.Up3Hist(i) && usdt > 300
					}

					// 策略测试部分
					// if canBuy && m.CrossOverTwice(i) && usdt > 300 {
					if condition {
						log.Info(mytime.TsToStrCST(v.Ts, "2006-01-02T15:04"), "canBuy!!!")
						err = rabbit.Send(exchange, routing_key_trade, "enable_open_button")
						log.Info("Send: enable_open_button")
						IfError("routing_key_trade: enable_open_button", err)

						// 开仓提示
						err = rabbit.Send(exchange, routing_key_tips, "开仓提示！")

						// 接受是否开仓的消息
						response := ""
						log.Info("To receive msgs: canBuy")
						select {
						case d1, ok1 := <-msgs1:
							if !ok1 {
								log.Error("break select: canBuy.", ok1)
								breakFor = true
								break // select
							}
							response = string(d1.Body)
							log.Info("receive msgs : canBuy: ", response)
							break
						case <-time.After(timeout_minutes):
							// 超时时间到达，退出 for
							rabbit.Send(exchange, routing_key_tips, "Timeout. Game over.")
							log.Error("Timeout. Game over. canBuy.")
							breakFor = true
							break // select
						}

						if breakFor {
							log.Error("break here......开仓 timeout")
							break // break: for i, v := range r {
						}

						// 自动点击，模拟自动运行，效果一样的。
						if autoopen {
							response = "openDone"
						}

						// 根据回复消息判断是否继续执行后续的逻辑
						if response == "openDone" {
							log.Info("response: openDone")
							per = perF * usdt / m.Close[i]
							price = m.Close[i]
							fee = per * price * feeRate
							coin = coin + per
							usdt = usdt - per*price - fee
							canSell = true
							canBuy = false
						}
					} else if canSell && (((m.Close[i]-price)/price >= goal) || ((m.Close[i]-price)/price <= fail ||
						((m.Close[i]-price)/price <= fail*0.7 && m.CrossUnder(i)))) {

						log.Info(mytime.TsToStrCST(v.Ts, "2006-01-02T15:04"), "canSell!!!")
						err = rabbit.Send(exchange, routing_key_trade, "enable_close_button")

						// 平仓提示
						if (m.Close[i]-price)/price >= goal {
							err = rabbit.Send(exchange, routing_key_tips, "可止盈"+
								// mf.Float64ToString(usdt+coin*m.Close[i], 2)+
								" "+mf.Float64ToString((m.Close[i]-price)/(0.01*price), 2)+"%")
						} else {
							err = rabbit.Send(exchange, routing_key_tips, "应止损"+
								// mf.Float64ToString(usdt+coin*m.Close[i], 2)+
								" "+mf.Float64ToString((m.Close[i]-price)/(0.01*price), 2)+"%")

						}

						// 接受是否平仓的消息
						response := ""
						select {
						case d, ok := <-msgs1:
							if !ok {
								log.Error("break select: canSell.", ok)
								breakFor = true
								break
							}
							response = string(d.Body)
							log.Info("Receive msg: canSell: ", response)
							break
						case <-time.After(timeout_minutes):
							// 超时时间到达，退出 for
							log.Error("Timeout. Game over.")
							rabbit.Send(exchange, routing_key_tips, "Timeout. Game over.")
							breakFor = true
							break // select
						}

						if breakFor {
							log.Error("break here......平仓 timeout")
							break // break: for i, v := range r {
						}

						// 自动点击，模拟自动运行，效果一样的。
						if autoclose {
							response = "closeDone"
						}

						// pass，不平仓！
						if response == "passDone" {
							log.Info("response: passDone")
							continue
						}

						// 开始平仓：
						if (m.Close[i]-price)/price >= goal {
							price = m.Close[i]
							fee = per * price * feeRate
							usdt = usdt + coin*price - fee
							coin = coin - per
							canBuy = true
							canSell = false
							msg1 := `{"total":` + mf.Float64ToString(usdt+coin*price, 2) +
								`,"usdt":` + mf.Float64ToString(usdt, 2) +
								`,"coin":` + mf.Float64ToString(coin, mf.CountFloat(price)) + `}`
							rabbit.Send(exchange, routing_key_account, msg1)
						} else if (m.Close[i]-price)/price <= fail ||
							((m.Close[i]-price)/price <= fail*0.7 && m.CrossUnder(i)) { // 0.7 死叉，且跌到fail的70%
							price = m.Close[i]
							fee = per * price * feeRate
							usdt = usdt + coin*price - fee
							coin = coin - per
							canBuy = true
							canSell = false
							msg2 := `{"total":` + mf.Float64ToString(usdt+coin*price, 2) +
								`,"usdt":` + mf.Float64ToString(usdt, 2) +
								`,"coin":` + mf.Float64ToString(coin, mf.CountFloat(price)) + `}`
							rabbit.Send(exchange, routing_key_account, msg2)
						}
					}
				} // End for. 结束后，更新一下账户余额。

				if breakFor {
					shouldExit = true // Game over.
					break             // break: select
				}

				// 更新账户金额
				msg3 := `{"total":` + mf.Float64ToString(usdt+coin*price_now, 2) +
					`,"usdt":` + mf.Float64ToString(usdt, 2) +
					`,"coin":` + mf.Float64ToString(coin, mf.CountFloat(price_now)) + `}`
				rabbit.Send(exchange, routing_key_account, msg3)

				// 最终盈利率
				err = rabbit.Send(exchange, routing_key_tips, "最终："+"   "+mf.Float64ToString((usdt+coin*price_now-usdt0)/(0.01*usdt0), 2)+"%")
				// err = rabbit.Send(exchange, routing_key_tips, mf.Float64ToString((usdt+coin*price_now-usdt0)/(0.01*usdt0), 2)+"%")
				err = rabbit.Send(exchange, routing_key_over, "")
				IfError("over : ", err)
				canBuy = true
				canSell = false
				log.Info("最终盈利率：" + "   " + mf.Float64ToString((usdt+coin*price_now-usdt0)/(0.01*usdt0), 2) + "%")
				msgTitle := stock + "-" + bar + " " + mf.Float64ToString((price_now-price_start)/(0.01*price_start), 2) + "%"
				log.Info("price_start:", price_start)
				log.Info("price_now:", price_now)
				log.Info("Send title msg:", msgTitle)
				err = rabbit.Send(exchange, qid+"-title", msgTitle)

			case <-time.After(timeout_minutes):
				// 超时，退出 goroutine
				rabbit.Send(exchange, routing_key_tips, "Timeout. Game over.")
				err = rabbit.Send(exchange, routing_key_over, "")
				log.Error("Timeout. Game over.")
				shouldExit = true // Game over.
				break
			} // End Select.
		}
		wg.Done()
	}()

	// 等待所有协程完成
	log.Info("Waiting for all goroutines exit...")
	wg.Wait()
	log.Info("All goroutines are done. ")
	log.Info("Game over.")
}
