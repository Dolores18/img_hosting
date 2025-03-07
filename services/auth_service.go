package services

import (
	"errors"
	"fmt"
	"img_hosting/dao"
	"img_hosting/models"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct{}

// Login 处理用户登录
func (s *AuthService) Login(identifier, password string) (*models.UserInfo, error) {
	fmt.Println("开始处理登录服务")
	db := models.GetDB()
	var user *models.UserInfo
	var err error

	fmt.Printf("登录信息: identifier=%s, password=%s (长度: %d)\n",
		identifier, password, len(password))

	// 确定查询方式（邮箱或用户名）
	if isEmail(identifier) {
		fmt.Println("使用邮箱登录")
		user, err = dao.GetUserByEmail(db, identifier)
	} else {
		fmt.Println("使用用户名登录")
		// 这里有问题，应该使用 name 字段查询
		user, err = dao.GetUserByName(db, identifier)
	}

	fmt.Printf("查询用户结果: user=%+v, err=%v\n", user, err)

	if err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		return nil, err
	}
	if user == nil {
		fmt.Println("用户不存在")
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	fmt.Println("开始验证密码")
	fmt.Printf("数据库中的密码哈希: %s (长度: %d)\n", user.Password, len(user.Password))
	fmt.Printf("用户提供的密码: %s (长度: %d)\n", password, len(password))

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Printf("密码验证失败: %v\n", err)
		return nil, errors.New("密码错误")
	}
	fmt.Println("密码验证成功")

	// 更新登录信息
	if err := dao.UpdateLoginInfo(db, user.UserID, ""); err != nil {
		fmt.Printf("更新登录信息失败: %v\n", err)
		return nil, err
	}
	fmt.Println("更新登录信息成功")

	return user, nil
}

// Register 处理用户注册
func (s *AuthService) Register(input *models.UserInput) (*models.UserInfo, error) {
	fmt.Println("开始处理注册服务")
	db := models.GetDB()

	// 检查邮箱是否已存在
	existingUser, err := dao.GetUserByEmail(db, input.Email)
	fmt.Printf("检查邮箱存在: err=%v, existingUser=%+v\n", err, existingUser)

	// 如果用户已存在，返回错误
	if existingUser != nil {
		fmt.Println("邮箱已被使用")
		return nil, errors.New("邮箱已被使用")
	}

	// 如果是"用户不存在"错误，则继续注册流程
	if err != nil && err.Error() == "用户不存在" {
		fmt.Println("用户不存在，继续注册流程")
	} else if err != nil {
		// 其他错误
		fmt.Printf("数据库错误: %v\n", err)
		return nil, err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Psd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("密码加密失败: %v\n", err)
		return nil, err
	}

	// 创建用户
	user := &models.UserInfo{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Age:      input.Age,
		Status:   models.UserStatusActive,
	}

	if err := dao.CreateUser(db, user); err != nil {
		fmt.Printf("创建用户失败: %v\n", err)
		return nil, err
	}

	fmt.Printf("用户创建成功: %+v\n", user)
	return user, nil
}

// isEmail 检查字符串是否为邮箱格式
func isEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
