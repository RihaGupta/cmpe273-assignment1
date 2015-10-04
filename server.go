
package main

import "net/http"
import "fmt"
import "bufio"
import "strconv"
import "strings"
import "net"
import "io/ioutil"

type StockRes struct {
	tradeId int
	stocks string
	crntValue float64
	unvestedAmount float64
}

var trans [20]StockRes
var count = 0

func buyStock(req string) StockRes {
	var reqDtls, stockLst, stock []string
	var cmpny, cmpnyPrice string
	var str []string
	var quantity, percent int
	var price, partial, amount float64
	trans[count].tradeId = count;
	reqDtls = strings.Split(req, "^")
	balance, _ := strconv.ParseFloat(reqDtls[2], 64);
	amount = balance
	stockLst = strings.Split(reqDtls[1], ",")
	for index :=0; index< len(stockLst); index++{
		stock = strings.Split(stockLst[index], ":")
		cmpny = stock[0]
		percent, _ = strconv.Atoi(stock[1])
		partial = amount*float64(percent)/100
		response, _ := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+stock[0]+"/quote?format=json&view=detail")
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		str = strings.SplitAfter(string(contents), "price")
		str = strings.SplitAfter(str[1]," ")		
		str = strings.SplitAfter(str[2],",")
		cmpnyPrice = str[0][1:len(str[0])-2]
		price,_ = strconv.ParseFloat(cmpnyPrice, 64)
		quantity = int(partial/price)
		balance = balance - float64(quantity) * price
		pr := strconv.FormatFloat(price, 'f', 2, 64)
		if(trans[count].stocks != ""){
			trans[count].stocks = trans[count].stocks + "," + cmpny + ":" + strconv.Itoa(quantity) + ":" + pr
		} else {
			trans[count].stocks = cmpny + ":" + strconv.Itoa(quantity) + ":" + pr
		}
	}
	trans[count].unvestedAmount = balance;
	count++
	return trans[count-1]
}

func tranDetails(req string) StockRes {
	var stockLst, stock []string
	var cmpnyPrice string
	var str,up []string
	var price float64
	var tDetails []string
	tDetails = strings.Split(req,"^")
	i, _ := strconv.Atoi(tDetails[1]);
	trans[i].crntValue = 0
	stockLst = strings.Split(trans[i].stocks, ",")
	up = strings.Split(trans[i].stocks, ",");
	for index :=0; index< len(stockLst); index++{
		stock = strings.Split(stockLst[index], ":")
		quantity, _ := strconv.Atoi(stock[1])
		response, _ := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+stock[0]+"/quote?format=json&view=detail")
		defer response.Body.Close()
		contents, _ := ioutil.ReadAll(response.Body)
		response, _ = http.Get("http://query.yahooapis.com/v1/public/yql?q=select%20*%20from%20yahoo.finance.quotes%20where%20symbol%20IN%20(%22"+stock[0]+"%22)&format=json&env=http://datatables.org/alltables.env")
		defer response.Body.Close()
		content, _ := ioutil.ReadAll(response.Body)
		str = strings.SplitAfter(string(content), "Change")
		cmpnyPrice = str[2][3:]
		cmpnyPrice = cmpnyPrice[:1]
		str = strings.Split(up[index], ":");
		str[2] = cmpnyPrice + str[2];
		up[index] = strings.Join(str, ":")
		str = strings.SplitAfter(string(contents), "price")
		str = strings.SplitAfter(str[1]," ")		
		str = strings.SplitAfter(str[2],",")
		cmpnyPrice = str[0][1:len(str[0])-2]
		price,_ = strconv.ParseFloat(cmpnyPrice, 64)
		trans[i].crntValue = trans[i].crntValue + float64(quantity)*price
	}	
	trans[i].stocks = strings.Join(up, ",")
	return trans[i]

}

func main() {
	fmt.Println("Launching server...")
	var r StockRes
	ln, _ := net.Listen("tcp", ":8081")
	conn, _ := ln.Accept()
	for
	{
		message, _ := bufio.NewReader(conn).ReadString(byte('#'))
		message = message[:len(message)-1]		
		fmt.Print("Message Received:", string(message))
		if(message[0]=='1'){
			r = buyStock(message);
			conn.Write([]byte("T Id: " + strconv.Itoa(r.tradeId) + " Stocks: " + r.stocks + "  Balance: " + strconv.FormatFloat(r.unvestedAmount, 'f', 2, 64) + "\n"))
		} else {
			r= tranDetails(message);
			conn.Write([]byte("T Id: " + strconv.Itoa(r.tradeId) + " Stocks: " + r.stocks + " Current Value: " + strconv.FormatFloat(r.crntValue, 'f', 2, 64) + " Balance: " + strconv.FormatFloat(r.unvestedAmount, 'f', 2, 64) + "\n"))
		}
	}
}
