package main

import "net"
import "fmt"
import "bufio"
import "strconv"

func main() {
	var op int;
	var bal float64;
	var id, text string;
	conn, _ := net.Dial("tcp", "127.0.0.1:8081")
	for
	{
		fmt.Print("Enter 1 for buying stock and 2 for checking your portfolio: ")
		fmt.Scan(&op)
		if(op==1) {
			fmt.Print("Enter stock symbol and percentage (eg=GOOG:50,YHOO:50) : ");
			fmt.Scan(&text)
			fmt.Print("Enter balance: ");
			fmt.Scan(&bal)
			// send to socket
			balance := strconv.FormatFloat(bal, 'f', 2, 64)
			fmt.Fprintf(conn, "1^" + text + "^" + balance + "#")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println("Message from server: "+message)
		} else if(op==2) {
			fmt.Print("Enter transaction ID: ")
			fmt.Scan(&id)
			// send to socket
			fmt.Fprintf(conn, "2^" + id + "#")
			// listen for reply
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Println("Message from server: "+message)
		} else {
			fmt.Println("Please Enter appropriate value");
		}
	}
}
