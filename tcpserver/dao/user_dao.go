package dao

import (
	"database/sql"
	"entrytask1/tcpserver/conf"
	"entrytask1/tcpserver/model"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


type UserDao struct {
	DB *sql.DB
}


// 新建连接
func NewUserDB() UserDao {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", conf.USERNAME, conf.PASSWORD,
		conf.NETWORK, conf.SERVER, conf.PORT, conf.DATABASE)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return UserDao{}
	}

	DB.SetConnMaxLifetime(100*time.Second)  //最大连接周期，超时的连接就close
	DB.SetMaxOpenConns(100) //最大连接数


	// 转成UserDao类型，可以使用其方法
	UserDB := UserDao{DB}

	return UserDB
}

// 关闭连接
func (u *UserDao) Close() error {
	err := u.DB.Close()
	if err != nil {
		return err
	}
	return nil
}

// 根据User的Id查
func (u *UserDao) UserQueryById(user *model.User) (int, error) {
	row := u.DB.QueryRow("select username,password,nickname from users where id=?", user.Id)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	//传入的是user结构体的成员的地址
	if err := row.Scan(&user.Id, &user.Password, &user.Nickname); err != nil {
		if err == sql.ErrNoRows {
			//查不到数据
			return 1, err
		} else {
			//其他问题
			return 2, err
		}
	}
	return 0, nil
}

// 根据User的username查
func (u *UserDao) UserQueryByName(user *model.User) (int, error) {
	row := u.DB.QueryRow("select id,password,nickname from users where username=?", user.Username)
	//row.scan中的字段必须是按照数据库存入字段的顺序，否则报错
	//传入的是user结构体的成员的地址
	if err := row.Scan(&user.Id, &user.Password, &user.Nickname); err != nil {
		if err == sql.ErrNoRows {
			//查不到数据
			return 1, err
		} else {
			//其他问题
			return 2, err
		}
	}
	return 0, nil
}

// 根据Id修改Nickname
func (u *UserDao) UserUpdateById(user *model.User) (int, error) {
	result,err := u.DB.Exec("UPDATE users set nickname=? where id=?", user.Nickname, user.Id)
	if err != nil{
		return 1, err
	}
	fmt.Println("update data successd:", result)

	rowsaffected, err := result.RowsAffected()
	if err != nil {
		return 1, err
	}
	fmt.Println("Affected rows:", rowsaffected)

	return 0, nil
}

