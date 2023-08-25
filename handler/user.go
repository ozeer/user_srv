package handler

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
	"user_srv/tool"

	"gorm.io/gorm"
)

type UserServer struct {
	proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// GetUserList 用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected)

	global.DB.Scopes(Paginate(int(req.Page), int(req.Size))).Find(&users)

	for _, user := range users {
		userInfoResp := proto.UserInfoResponse{
			Id:       user.ID,
			Password: user.Password,
			Nickname: user.NickName,
			Mobile:   user.Mobile,
			Gender:   user.Gender,
			Role:     int32(user.Role),
		}
		if user.Birthday != nil {
			userInfoResp.Birthday = uint64(user.Birthday.Unix())
		}
		rsp.Data = append(rsp.Data, &userInfoResp)
	}

	return rsp, nil
}

// GetUserByMobile 通过手机号码查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Nickname: user.NickName,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix())
	}

	return &userInfoResp, nil
}

// GetUserById 通过Id查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("id = ?", req.Id).First(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Nickname: user.NickName,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix())
	}

	return &userInfoResp, nil
}

// CreateUser 新建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where("mobile = ?", req.Mobile).First(&user)

	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	user.Mobile = req.Mobile
	user.NickName = req.Nickname
	user.Password = tool.Encrypt(tool.KEY, req.Password)
	result = global.DB.Create(&user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	userInfoResp := proto.UserInfoResponse{
		Id:       user.ID,
		Password: user.Password,
		Nickname: user.NickName,
		Mobile:   user.Mobile,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoResp.Birthday = uint64(user.Birthday.Unix())
	}

	return &userInfoResp, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfoRequest) (*empty.Empty, error) {
	var user model.User
	result := global.DB.Where("id = ?", req.Id).First(&user)

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户不存在")
	}

	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.Nick
	user.Birthday = &birthday
	user.Gender = req.Gender
	result = global.DB.Save(user)

	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}

	return &empty.Empty{}, nil
}

// CheckPassword 检查密码
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.CheckPasswordRequest) (*proto.CheckResponse, error) {
	// 校验密码
	encryptedPassword := tool.Encrypt(tool.KEY, req.Password)
	check := encryptedPassword == req.EncryptedPassword

	return &proto.CheckResponse{Success: check}, nil
}
