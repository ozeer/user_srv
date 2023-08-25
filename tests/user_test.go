package tests

import (
	"context"
	"fmt"
	"testing"
	"user_srv/proto"

	"google.golang.org/grpc"
)

var (
	userClient proto.UserClient
	conn       *grpc.ClientConn
)

func GetUserList() {
	resp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Page: 1,
		Size: 10,
	})

	if err != nil {
		panic(err)
	}

	for _, user := range resp.Data {
		fmt.Println(user.Id, "#", user.Mobile, user.Nickname, user.Password)
		checkResp, err := userClient.CheckPassword(context.Background(), &proto.CheckPasswordRequest{
			Password:          "admin",
			EncryptedPassword: user.Password,
		})

		if err != nil {
			panic(err)
		}

		fmt.Println("测试结果: ", checkResp.Success)
	}
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

func Test(t *testing.T) {
	Init()
	GetUserList()
	conn.Close()
}
