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


		switch message.operation {
			case "enter":
				entry()
			case "exit":
				exit()
			default:
				fmt.Println("未知的操作")
		}


		// 构建要执行的命令
		command := fmt.Sprintf("source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceenter %s %s", message.Pair, message.Side)
		fmt.Println("执行的命令:", command)
		// 执行命令
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		// cmd := exec.Command("/bin/bash", "-c", "source activate freqtrade && python ../freqtrade/scripts/rest_client.py --config ../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json forceenter "+ message.Pair + " " +message.Side)
		// cmd := exec.Command("python", "../freqtrade/scripts/rest_client.py", "--config", "../freqtrade/user_data/config_future_dryrun_3MA_1H_DCA.json", "forceenter", message.Pair, message.Side)
		err = cmd.Run()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "无法执行本地脚本: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Webhook请求已接收")
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "不支持的请求方法")
	}
}

func entry(msg Message) error {


}

func exit(msg Message) error{


}
