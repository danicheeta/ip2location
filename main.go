package main

import (
	"database/sql"

	"fmt"

	"bufio"

	"strconv"

	"runtime"
	"time"

	"os"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
)

var (
	head  *node
	dbmap *gorp.DbMap
	slice []int
)

type sourceMap []Ip2loc

type node struct {
	nodes []*node
	data  string
}

type Ip2loc struct {
	Ip_from      int
	Ip_to        int
	Country_Name string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	channel := make(chan sourceMap, 1)
	go queryDB(channel)
	counter := 0
	for source := range channel {
		for _, data := range source {
			for i := data.Ip_from; i < data.Ip_to; i++ {
				go initialize(head, digitable(i), data.Country_Name)
			}
		}
		fmt.Println(fmt.Sprintf("%d/1000000\n", counter))
		counter += 1000
	}
	user()
}

func queryDB(ch chan sourceMap) {
	for i := 0; i < 100000; i += 1000 {
		var source sourceMap
		query := `SELECT ip_from, ip_to, country_name FROM ip2location LIMIT 1000 OFFSET %d`
		_, err := dbmap.Select(&source, fmt.Sprintf(query, i)) //care
		if err != nil {
			fmt.Println(err)
		}
		ch <- source
	}
	close(ch)

}

func user() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		now := time.Now()
		fmt.Println(findPath(text))
		fmt.Println(time.Since(now))
	}
}

func initialize(head *node, slice []int, data string) {
	first := head
	for _, i := range slice {
		next := getNodeIndex(first, i)
		first = next
	}
	first.data = data
}

func findPath(address string) string {
	ip64, err := strconv.ParseInt(address[:len(address)-1], 10, 0)
	if err != nil {
		panic(err)
	}
	slice := digitable(int(ip64))
	next := head
	for i := range slice {
		next = next.nodes[slice[i]]
	}
	return next.data
}

func getNodeIndex(head *node, index int) (result *node) {
	if head.nodes[index] == nil {
		node := &node{
			nodes: make([]*node, 10),
		}
		head.nodes[index] = node
		return node
	} else {
		result = head.nodes[index]
		return result
	}

}

func digitable(num int) []int {
	result := []int{}
	for num > 0 {
		result = append(result, num%10)
		num /= 10
	}
	return result
}

func init() {
	db, _ := sql.Open("mysql", "smth")
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{}}
	head = &node{
		nodes: make([]*node, 10),
	}
}

func log(slice []int, target int) int {
	var mid int
	mid = len(slice) / 2
	if len(slice) == 1 {
		return slice[0]
	}

	if target < slice[mid-1] {
		fmt.Println("111")
		return log(slice[:mid], target)
	} else {
		fmt.Println("222")
		return log(slice[mid:], target)
	}
}
