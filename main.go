package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"os/exec"
	"os"
)

func main() {
	// 设置Webhook的路由和处理函数
	http.HandleFunc("/webhook", handleWebhook)

	// 启动服务器并监听指定的端口
	port := ":80"
	fmt.Printf("Webhook服务器正在监听端口 %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}


type Message struct {
	Operation string `json:"Operation"`
	Pair string `json:"Pair"`
	Side string `json:"Side"`
}


type Trade struct {
	Trades []struct {
		TradeID                     int         `json:"trade_id"`
		Pair                        string      `json:"pair"`
		BaseCurrency                string      `json:"base_currency"`
		QuoteCurrency               string      `json:"quote_currency"`
		IsOpen                      bool        `json:"is_open"`
		Exchange                    string      `json:"exchange"`
		Amount                      float64     `json:"amount"`
		AmountRequested             float64     `json:"amount_requested"`
		StakeAmount                 float64     `json:"stake_amount"`
		MaxStakeAmount              float64     `json:"max_stake_amount"`
		Strategy                    string      `json:"strategy"`
		EnterTag                    string      `json:"enter_tag"`
		Timeframe                   int         `json:"timeframe"`
		FeeOpen                     float64     `json:"fee_open"`
		FeeOpenCost                 float64     `json:"fee_open_cost"`
		FeeOpenCurrency             string      `json:"fee_open_currency"`
		FeeClose                    float64     `json:"fee_close"`
		FeeCloseCost                float64     `json:"fee_close_cost"`
		FeeCloseCurrency            string      `json:"fee_close_currency"`
		OpenDate                    string      `json:"open_date"`
		OpenTimestamp               int64       `json:"open_timestamp"`
		OpenRate                    float64     `json:"open_rate"`
		OpenRateRequested           float64     `json:"open_rate_requested"`
		OpenTradeValue              float64     `json:"open_trade_value"`
		CloseDate                   string      `json:"close_date"`
		CloseTimestamp              int64       `json:"close_timestamp"`
		RealizedProfit              float64     `json:"realized_profit"`
		RealizedProfitRatio         float64     `json:"realized_profit_ratio"`
		CloseRate                   float64     `json:"close_rate"`
		CloseRateRequested          float64     `json:"close_rate_requested"`
		CloseProfit                 float64     `json:"close_profit"`
		CloseProfitPct              float64     `json:"close_profit_pct"`
		CloseProfitAbs              float64     `json:"close_profit_abs"`
		TradeDurationS              int         `json:"trade_duration_s"`
		TradeDuration               int         `json:"trade_duration"`
		ProfitRatio                 float64     `json:"profit_ratio"`
		ProfitPct                   float64     `json:"profit_pct"`
		ProfitAbs                   float64     `json:"profit_abs"`
		ExitReason                  string      `json:"exit_reason"`
		ExitOrderStatus             string      `json:"exit_order_status"`
		StopLossAbs                 float64     `json:"stop_loss_abs"`
		StopLossRatio               float64     `json:"stop_loss_ratio"`
		StopLossPct                 int         `json:"stop_loss_pct"`
		StoplossOrderID             interface{} `json:"stoploss_order_id"`
		StoplossLastUpdate          interface{} `json:"stoploss_last_update"`
		StoplossLastUpdateTimestamp interface{} `json:"stoploss_last_update_timestamp"`
		InitialStopLossAbs          float64     `json:"initial_stop_loss_abs"`
		InitialStopLossRatio        float64     `json:"initial_stop_loss_ratio"`
		InitialStopLossPct          int         `json:"initial_stop_loss_pct"`
		MinRate                     float64     `json:"min_rate"`
		MaxRate                     float64     `json:"max_rate"`
		Leverage                    float64     `json:"leverage"`
		InterestRate                float64     `json:"interest_rate"`
		LiquidationPrice            interface{} `json:"liquidation_price"`
		IsShort                     bool        `json:"is_short"`
		TradingMode                 string      `json:"trading_mode"`
		FundingFees                 float64     `json:"funding_fees"`
		AmountPrecision             float64     `json:"amount_precision"`
		PricePrecision              float64     `json:"price_precision"`
		PrecisionMode               int         `json:"precision_mode"`
		ContractSize                float64     `json:"contract_size"`
		HasOpenOrders               bool        `json:"has_open_orders"`
		Orders                      []struct {
			Amount               float64     `json:"amount"`
			SafePrice            float64     `json:"safe_price"`
			FtOrderSide          string      `json:"ft_order_side"`
			OrderFilledTimestamp int64       `json:"order_filled_timestamp"`
			FtIsEntry            bool        `json:"ft_is_entry"`
			Pair                 string      `json:"pair"`
			OrderID              string      `json:"order_id"`
			Status               string      `json:"status"`
			Average              float64     `json:"average"`
			Cost                 float64     `json:"cost"`
			Filled               float64     `json:"filled"`
			IsOpen               bool        `json:"is_open"`
			OrderDate            string      `json:"order_date"`
			OrderTimestamp       int64       `json:"order_timestamp"`
			OrderFilledDate      string      `json:"order_filled_date"`
			OrderType            string      `json:"order_type"`
			Price                float64     `json:"price"`
			Remaining            float64     `json:"remaining"`
			FtFeeBase            interface{} `json:"ft_fee_base"`
			FundingFee           float64     `json:"funding_fee"`
		} `json:"orders"`
	} `json:"trades"`
	TradesCount int `json:"trades_count"`
	Offset      int `json:"offset"`
	TotalTrades int `json:"total_trades"`
}


func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// 处理Webhook请求
	if r.Method == "POST" {
		// 解析JSON数据
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "无法解析JSON数据: %v", err)
			return
		}

		// 处理解析后的数据
		fmt.Println("接收到Webhook请求:")
		fmt.Println("operation:", message.Operation)
		fmt.Println("pair:", message.Pair)
		fmt.Println("side:", message.Side)
		fmt.Println()


		switch message.Operation {
			case "enter":
				err := exit(message)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "退出失败: %v", err)
					return
				}
				err = entry(message)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "进入失败“: %v", err)
					return
				}
			case "exit":
				err := exit(message)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, "退出失败: %v", err)
					return
				}
			default:
				fmt.Println("未知的操作")
		}


		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Webhook请求已接收")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "不支持的请求方法")
	}
}

func entry(message Message) error {
	// 构建要执行的命令
	command := fmt.Sprintf("source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceenter %s %s", message.Pair, message.Side)
	fmt.Println("执行的命令:", command)
	// 执行命令
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// cmd := exec.Command("/bin/bash", "-c", "source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceenter "+ message.Pair + " " +message.Side)
	// cmd := exec.Command("python", "../freqtrade/scripts/rest_client.py", "--config", "../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json", "forceenter", message.Pair, message.Side)
	return cmd.Run()
	
}

func exit(message Message) error{
	//query trade
	command := fmt.Sprintf("source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json trade")
	fmt.Println("执行的命令:", command)
	cmd := exec.Command("/bin/bash", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var trades Trade
	err=json.Unmarshal(output,&trades)
	if err != nil {
		return err
	}

	var tradeId int = -1
	for _,trade := range trades.Trades{
		if trade.Pair == message.Pair{
			tradeId = trade.TradeID
		}
	}
	if tradeId == -1{
		fmt.Println("无单，无需退出")
	}

	// 退出当前单
	command = fmt.Sprintf("source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceexit %d", tradeId)
	fmt.Println("执行的命令:", command)
	// 执行命令
	cmd = exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// cmd := exec.Command("/bin/bash", "-c", "source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceenter "+ message.Pair + " " +message.Side)
	// cmd := exec.Command("python", "../freqtrade/scripts/rest_client.py", "--config", "../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json", "forceenter", message.Pair, message.Side)
	return cmd.Run()
}
