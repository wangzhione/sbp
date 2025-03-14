package mysql

import (
	"context"
	"database/sql"
	"runtime"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wangzhione/sbp/util/chain"
	"github.com/wangzhione/sbp/util/jsou"
)

var connects = "mysql -u root -p123456 -h 127.0.0.1 -P 3306 demo"

func TestNewDB(t *testing.T) {
	connects = "mysql -u root -p123456 demo"

	s, err := NewDB(chain.Background, connects)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	if s != nil {
		t.Log("Success")
	}

	_ = s

	runtime.AddCleanup(s, func(fd uintptr) { println("runtime.AddCleanup close", fd) }, 0)

	s.Close(chain.Background)
	// 必须主动 close 后才能被回收, 所以 runtime.SetFinalizer(s, (*DB).Close) 永远不会执行
	s = nil

	for i := 0; i < 6; i++ {
		runtime.GC()
		t.Log(i+1, "runtime.GC()")
		time.Sleep(time.Second)
	}
}

func TestQueryRow(t *testing.T) {
	s, err := NewDB(chain.Background, connects)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	var count int
	err = s.QueryRow(chain.Background, "select count(*) from t_user", nil, &count)
	if err != nil {
		t.Fatal("s.QueryRow fatal", err)
	}
	t.Log("count = ", count)

	err = s.QueryRow(chain.Background, "select count(*) from t_user where id > ?", []any{6}, &count)
	if err != nil {
		t.Fatal("s.QueryRow fatal", err)
	}
	t.Log("count = ", count)
}

// User 结构体与 t_user 表字段对应
type User struct {
	ID               uint64
	UserName         string
	Password         string
	PasswordSalt     string
	EmailNotVerified string
	UserEmail        string
	UpdateTime       string
	CreateTime       string
	DeleteTime       uint64
}

func TestDB_QueryCallBack(t *testing.T) {
	s, err := NewDB(chain.Background, connects)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	// 现代 Go 开发, 一定要集合 AI, 例如 ChatGPT 辅助, 很多工作都好节省
	var users []User
	query := "SELECT id, user_name, password, password_salt, email_not_verified, user_email, update_time, create_time, delete_time FROM t_user WHERE delete_time = 0"
	err = s.QueryCallBack(chain.Background, func(ctx context.Context, rows *sql.Rows) error {
		// 遍历查询结果
		for rows.Next() {
			var user User
			// 扫描当前行数据到 user 结构体
			if err := rows.Scan(&user.ID, &user.UserName, &user.Password, &user.PasswordSalt,
				&user.EmailNotVerified, &user.UserEmail, &user.UpdateTime, &user.CreateTime, &user.DeleteTime); err != nil {
				t.Errorf("Failed to scan row: %v", err)
				return err
			}

			// 打印用户信息
			t.Logf("User ID: %d, UserName: %s, UserEmail: %s", user.ID, user.UserName, user.UserEmail)

			users = append(users, user)
		}

		return nil
	}, query)
	if err != nil {
		t.Fatal("QueryCallBack fatal", err)
	}
}

func TestDB_QueryAll(t *testing.T) {
	s, err := NewDB(chain.Background, connects)
	if err != nil {
		t.Fatal("NewDB fatal", err)
	}

	query := "SELECT id, user_name, password, password_salt, email_not_verified, user_email, update_time, create_time, delete_time FROM t_user WHERE delete_time = 0"
	results, err := s.QueryAll(chain.Background, query)
	if err != nil {
		t.Fatal("s.QueryAll fatal", err)
	}

	t.Log("Success", "\n", jsou.String(results))
}
