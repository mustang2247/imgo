package im

import (
    "code.google.com/p/go-uuid/uuid"
    "fmt"
    "log"
    "net/http"
    "strings"
    "im-go/im/common"
    "im-go/im/model"
)

// 启动HTTP服务
func StartHttpServer(config common.IMConfig) error {
    log.Printf("HttpServer starting...")

    // 设置请求映射地址及对应处理方法
    http.HandleFunc("/login", handleLogin)
    //打印监听端口
    log.Printf("HTTP 开始监听端口: %d", config.HttpPort)
    // 设置监听地址及端口
    addr := fmt.Sprintf("0.0.0.0:%d", config.HttpPort)
    if err := http.ListenAndServe(addr, nil); err != nil {
        return fmt.Errorf("监听Http失败: %s", err)
    }
    return nil
}

//登录请求处理方法
func handleLogin(resp http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
        handlePost(resp, req)
    } else {
        resp.Write(NewIMResponseSimple(404, "Not Found: " + req.Method, "").Encode())
    }
}

// POST登录请求
func handlePost(resp http.ResponseWriter, req *http.Request) {
    ip := common.GetIp(req)
    device := req.FormValue("device")
    account := req.FormValue("account")
    password := req.FormValue("password")

    log.Printf("ip %s", ip)
    log.Printf("device %s", device)
    log.Printf("account %s", account)
    log.Printf("password %s", password)

    login(resp, account, password, device, ip)
}

// 登录主方法
func login(resp http.ResponseWriter, account string, password string, device string, ip string) {
    if account == "" {
        resp.Write(NewIMResponseSimple(101, "账号不能为空", "").Encode())
    } else if password == "" {
        resp.Write(NewIMResponseSimple(102, "密码不能为空", "").Encode())
    } else if device == "" {
        resp.Write(NewIMResponseSimple(103, "设备名不能空", "").Encode())
    } else {
        var user model.IMUser
        num := CheckAccount(account)
        if num > 0 {
            user = LoginUser(account, password)
            if !strings.EqualFold(user.Id, "") {
                token := uuid.New()
                if SaveLogin(user.Id, token, ip) > 0 {
                    returnData := make(map[string]string)
                    returnData["id"] = user.Id
                    returnData["nick"] = user.Nick
                    returnData["avatar"] = user.Avatar
                    returnData["status"] = user.Status
                    returnData["token"] = token //token uuid 带 横杠
                    resp.Write(NewIMResponseData(common.GetJson("user", returnData), "LOGIN_RETURN").Encode())
                } else {
                    resp.Write(NewIMResponseSimple(105, "保存登录记录错误,请稍后再试", "").Encode())
                }

            } else {
                resp.Write(NewIMResponseSimple(104, "密码错误", "").Encode())
            }
        } else {
            resp.Write(NewIMResponseSimple(103, "账户不存在", "").Encode())
        }
    }
}