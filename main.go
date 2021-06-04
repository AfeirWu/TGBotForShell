package main

import (
    "time"
    "log"
    "os"
    "strconv"
    "strings"
    "os/exec"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

var clean_domain string
var tgbot_token string
var tgbot_chatid string
var chat_id string

func checkErr(err error) {
    if err != nil {
        log.Println("错误：",err.Error() )
    }
}

func removeCharacters(input string, characters string) string {
    filter := func(r rune) rune {
        if strings.IndexRune(characters, r) < 0 {
            return r
        }
        return -1
    }
    return strings.Map(filter, input)
}


func exec_shell(command string) string {
    log.Println("/bin/sh -c",command)
    out, err := exec.Command("/bin/sh","-c",command).Output()
    checkErr(err)
    return string(out)
}

func main() {
    f, err := os.OpenFile("bot.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    log.SetOutput(f)
    if os.Getenv("TGBOT_TOKEN") != "" && os.Getenv("TGBOT_CHATID") != ""  {
        tgbot_token = os.Getenv("TGBOT_TOKEN")
        tgbot_chatid =  os.Getenv("TGBOT_CHATID")
    } else {
        log.Println("没有设置TGBOT_TOKEN或TGBOT_CHATID")
    }
    bot, err := tgbotapi.NewBotAPI(tgbot_token)
    if err != nil {
        log.Panic(err)
    }
    //bot.Debug = true
    log.Printf("Bot名称：%s", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60
    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        log.SetOutput(f)
        log.Printf("UserName=%s,Message=%s", update.Message.From.UserName, update.Message.Text)

        var output string

        chat_id:= strconv.Itoa(int(update.Message.Chat.ID))

        output = "您无权访问！"

        if strings.Contains(chat_id , tgbot_chatid) {
            output = "无效命令！"
            if strings.Contains(update.Message.Text,"/sh") {
                cmd_raw:=update.Message.Text
                s := strings.Replace(cmd_raw, "/sh ", "", -1)
                if len(s) == 1 {
                    output = "请在/sh 后加命令内容！"
                } else {
                    output = exec_shell(s)
                }
            }
        }
        // telegramBot发送文本消息最大支持4096字节，所以发送文件
        if len([]byte(output)) > 4096 {
            logname := time.Now().Format("2006-01-02-15-04-05") + ".log"
            f1, err1 := os.OpenFile(logname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
            if err1 != nil {
                log.Fatal(err1)
            }
            defer f1.Close()
            log.SetOutput(f1)
            log.Println(output)
            msg := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, logname)
            bot.Send(msg)

            // 删除日志文件
            os.Remove(logname)
        } else {
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, output)
            bot.Send(msg)
        }
   }
}
