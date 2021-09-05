package model

import (
	"chatroom/common/message"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

//实现对数据库的操作

//我们在服务器启动后就初始化一个userDao实例
//把他做成全局变量，在需要和redis操作时直接使用
//因为整个程序只需要建立一次，所以定义成全局变量
var (
	MyUserDao *UserDao
)

//定义一个UserDao的结构体，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

//使用工厂模式创建UserDao的实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

//1.根据用户id返回User实例或err
func (this *UserDao) GetUserById(conn redis.Conn, id int) (user *message.User, err error) {
	//通过给定的id，在redis中查询用户
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil { //这种错误表示没有在数据库中找到对应id的数据
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &message.User{}
	//将res反序列化成user实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("反序列化失败：", err)
	}
	return
}

//完成登录的校验
//1.完成对用户的验证
//如果用户的id和密码都正确，返回一个user实例
//如果有错误，返回一个错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *message.User, err error) {
	//先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.GetUserById(conn, userId)
	if err != nil {
		return
	}

	//这时证明这个用户是获取到
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

//处理注册信息的函数
func (this *UserDao) Register(user *message.User) (err error) {
	//先从UserDao的连接池中取出一个连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.GetUserById(conn, user.UserId)
	if err == nil { //没有错误返回说明取到了用户，说明已经存在
		err = ERROR_USER_EXISTS
		return
	}

	//这时，说明id在redis中还没有，可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		fmt.Println("序列化失败")
		return
	}

	//入库
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("注册信息入库失败，可能是数据库开启保护模式")
		return
	}
	return

}
