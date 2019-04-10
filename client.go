/*
 * @Author: atony2099
 * @Date: 2019-04-10 02:59:54
 * @Last Modified by: atony2099
 * @Last Modified time: 2019-04-10 12:32:41
 */

package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type IdArgs struct {
	BizTag string
}

func main() {

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	idArgs := &IdArgs{"test"}
	var idNum int64
	err = client.Call("Arith.GetId", idArgs, &idNum)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %v*%d\n", idArgs, idNum)

}
