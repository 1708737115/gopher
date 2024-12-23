/*
@author: fengxu
@since: 2024-12-23
@desc: //TODO
*/
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopher/src/grpc/grpc1/services"
)

func main() {
	user := &services.User{
		Name:  "fengxu",
		Email: "fengxu@gmail.com",
		Phone: "123456789",
	}

	//序列化
	newUser := &services.User{}
	marshal, err := proto.Marshal(user)
	if err != nil {
		panic(err)
	}

	//反序列化
	err = proto.Unmarshal(marshal, newUser)
	if err != nil {
		panic(err)
	}

	fmt.Println(newUser)
}
