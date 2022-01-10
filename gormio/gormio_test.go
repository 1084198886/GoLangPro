package gormio

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func TestDB(t *testing.T) {
	SetLogOutput()
	storage := newStorage()
	if err := storage.connect("test"); err != nil {
		t.Errorf("connnect database err %v", err)
		return
	}
	fmt.Println("connect database success")
	doInsert(storage.cfg())
	doUpdate(storage.cfg())
	doDelete(storage.cfg())
	doQuery(storage.cfg())
	doTransaction(storage.cfg())
}

/*
  开启事务操作表
*/
func doTransaction(db *gorm.DB) {
	fmt.Println("doTransaction start")
	tx := db.Begin()
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("defer use Transaciton Update err:%v", err)
			tx.Rollback()
		} else {
			fmt.Println("defer use Transaciton Update success")
		}
	}()
	if err := tx.Error; err != nil {
		panic(fmt.Sprintf("connnect database err %v", err))
	}
	updateRet := tx.Table("users").Where("Stuempno=?", "002").Updates(map[string]interface{}{"Name": "fds"})
	if updateRet.Error != nil {
		panic(fmt.Sprintf("Update  err %v", updateRet.Error))
	}
	if err := tx.Commit().Error; err != nil {
		panic(fmt.Sprintf("Commit  err %v", updateRet.Error))
	}
	fmt.Println("use Transaciton Update success")

}

/*删除*/
func doDelete(db *gorm.DB) {
	ret := db.Debug().Table("users").
		Where("Stuempno=?", "002").
		Delete(&User{})

	if err := ret.Error; err != nil {
		fmt.Printf("delete user err :%v \n", err)
		return
	}
	fmt.Printf("delete user success \n")
}

/*更新*/
func doUpdate(db *gorm.DB) {
	// 更新一个字段
	//ret := db.Debug().Table("users").
	//	Where("Stuempno=?", "002").
	//	Update("name", "zhangsan2")
	// 更新多个字段
	ret := db.Debug().Table("users").
		Where("Stuempno=?", "002").
		Updates(map[string]interface{}{"name": "zhangsan2"})

	if err := ret.Error; err != nil {
		fmt.Printf("Updates record err :%v \n", err)
		return
	}
	fmt.Println("Updates user success")
	fmt.Println("raw Update start")
	execRet := db.Exec("Update users set Name=? where Stuempno=?", "NewName", "002")
	if execRet.Error != nil {
		fmt.Printf("raw Update err:%v", execRet.Error)
		return
	}
	if execRet.RowsAffected == 0 {
		fmt.Printf("raw Update err:%v", execRet.Error)
		return
	}
	fmt.Printf("raw Update success  RowsAffected:%v", execRet.RowsAffected)
}

/*
 查询
*/
func doQuery(db *gorm.DB) {
	//user:=&User{}
	var count int64
	users := make([]User, 10)
	ret := db.
		Debug().
		Table("users").
		Select([]string{"Id", "Stuempno", "Name", "Fee"}).
		//Where("Stuempno=?", "002").
		//Where("Stuempno<> ?", "001").
		//Where("Stuempno in (?)", []string{"001","002"}).
		//Where("Stuempno between ? and ?", "000","001").
		//Where("Stuempno like ?", "%02%").
		//Where("Stuempno =? and Name=?", "002", "zhangsan2").
		//First(user).
		//Offset(1).
		//Limit(1).
		Order("Stuempno desc").
		Find(&users).
		Count(&count)
	if err := ret.Error; err != nil {
		fmt.Printf("select record err :%v \n", err)
		return
	}
	if ret.RowsAffected == 0 {
		fmt.Printf("select record no data \n")
		return
	}
	fmt.Printf("count=%d \n", count)

	// 遍历查询结果
	for index, user := range users {
		jsonStr, _ := json.Marshal(user)
		fmt.Printf("发现 Id=%d bean=%s \n", index, string(jsonStr))
	}

	//obj:= struct {
	// Stuempno string
	//}{}
	/*
	  使用sql原始语句查询
	*/
	fmt.Println("raw query start...")
	rawRet := db.Debug().Raw(`select * from users where Stuempno=?`, "001").Scan(&users)
	//rawRet := db.Debug().Raw(`select * from users where Stuempno=?`, "001").Scan(&obj)
	if rawRet.Error != nil {
		fmt.Printf("raw query err:%v \n", rawRet.Error)
		return
	}
	//jsonStr, _ := json.Marshal(obj)
	//fmt.Printf("raw query 发现bean=%s \n", string(jsonStr))

	for index, user := range users {
		jsonStr, _ := json.Marshal(user)
		fmt.Printf("raw query 发现 Id=%d bean=%s \n", index, string(jsonStr))
	}

}

func doInsert(db *gorm.DB) {
	user := &User{Stuempno: "001", Name: "zhangsan", Fee: 0.20}
	/*
	 插入
	*/
	if err := db.Debug().Create(user).Error; err != nil {
		fmt.Printf("insert record err :%v \n", err)
		return
	}
	user = &User{Stuempno: "002", Name: "zhangsan", Fee: 0.20}
	if err := db.Debug().Create(user).Error; err != nil {
		fmt.Printf("insert record err :%v \n", err)
		return
	}
	fmt.Printf("insert user success \n")
}

func newStorage() *gormioStorage {
	return &gormioStorage{
		mu: sync.Mutex{},
	}
}
