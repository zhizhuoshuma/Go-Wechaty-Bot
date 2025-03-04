package main

import (
	"fmt"
	"github.com/spf13/viper"
	. "github.com/wechaty/go-wechaty/wechaty"
	wp "github.com/wechaty/go-wechaty/wechaty-puppet"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	. "github.com/wechaty/go-wechaty/wechaty/interface"
	"github.com/wechaty/go-wechaty/wechaty/user"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"
	. "wechatBot/bot-service"
	. "wechatBot/data"
)

var (
	err error
)

type (
	MessageInfo struct {
		Status   bool
		RoomName string
		RoomID   string
		UserName string
		UserID   string
		AutoInfo string
	}
)

func init() {
	// 设置日志格式
	log.SetPrefix("[xrsec] [\033[01;33m➜\033[0m] ") // 设置日志前缀
	log.SetFlags(log.Ltime | log.Lshortfile)

	// 初始化配置文件
	rootPath, _ := os.Getwd()
	exePath, _ := os.Executable()
	log.Printf("rootPath: %s, exePath: %s", rootPath, filepath.Dir(exePath))

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filepath.Dir(exePath))
	viper.AddConfigPath(rootPath)
	viper.Set("rootPath", rootPath)
	viper.Set("exePath", exePath)
}

func onScan(context *Context, qrCode string, status schemas.ScanStatus, data string) {
	log.Printf("%s[Scan] %s %s %s\n", viper.GetString("info"), qrCode, status, data)
}

/*
	@method onLogin 当机器人成功登陆后，会触发事件，并会在事件中传递当前登陆机器人的信息
	@param {*} user
*/
func onlogin(ctx *Context, user *user.ContactSelf) {
	log.Printf(`
                           //
               \\         //
                \\       //
        ##DDDDDDDDDDDDDDDDDDDDDD##
        ## DDDDDDDDDDDDDDDDDDDD ##      
        ## DDDDDDDDDDDDDDDDDDDD ##      
        ## hh                hh ##      ##         ## ## ## ##   ## ## ## ###   ##    ####     ##     
        ## hh    //    \\    hh ##      ##         ##       ##   ##             ##    ## ##    ##
        ## hh   //      \\   hh ##      ##         ##       ##   ##             ##    ##   ##  ##
        ## hh                hh ##      ##         ##       ##   ##     ##      ##    ##    ## ##
        ## hh      wwww      hh ##      ##         ##       ##   ##       ##    ##    ##     ####
        ## hh                hh ##      ## ## ##   ## ## ## ##   ## ## ## ###   ##    ##      ###
        ## MMMMMMMMMMMMMMMMMMMM ##    
        ##MMMMMMMMMMMMMMMMMMMMMM##      微信机器人: [%s] 已经登录成功了。
        `, user.Name())
	viper.Set("bot.name", user.Name())
}

/**
@method onLogout 当机器人检测到登出的时候，会触发事件，并会在事件中传递机器人的信息。
@param {*} user
*/
func onLogout(ontext *Context, user *user.ContactSelf, reason string) {
	log.Println("========================onLogout👇========================")
	DingMessage(user.Name() + "账号已退出登录, 请检查账号状态!" + reason)
}

/*
  @method onRoomInvite 当收到群邀请的时候，会触发这个事件。
  @param {*} user
*/
func onRoomInvite(ontext *Context, roomInvitation *user.RoomInvitation) {
	log.Println("========================onRoomInvite👇========================")
	if err = roomInvitation.Accept(); err != nil {
		ErrorFormat("Accept Room Invitation", err)
		// TODO	好像有点问题，群聊设置了邀请确认就用不了
	}
	log.Println(roomInvitation.String())
}

/*
	@method onRoomTopic 当有人修改群名称的时候会触发这个事件。
	@param {*} user
*/
func onRoomTopic(context *Context, room *user.Room, newTopic string, oldTopic string, changer IContact, date time.Time) {
	log.Println("========================onRoomTopic👇========================")
	//	暂时未测试
}

/*
	进入房间监听回调 room-群聊 inviteeList-受邀者名单 inviter-邀请者
	判断配置项群组id数组中是否存在该群聊id
*/
func onRoomJoin(context *Context, room *user.Room, inviteeList []IContact, inviter IContact, date time.Time) {
	// TODO 用不了
}

/*
	@method onRoomleave 当机器人把群里某个用户移出群聊的时候会触发这个时间。用户主动退群是无法检测到的。
	@param {*} user
*/
func onRoomleave(context *Context, room *user.Room, leaverList []IContact, remover IContact, date time.Time) {
	log.Println("========================onRoomleave👇========================")
	log.Printf("用户[%s]被踢出去聊", remover.Name())
	// 用不了
}

func onFriendship(context *Context, friendship *user.Friendship) {
	// 用不了
	switch friendship.Type() {
	case 1:
	//FriendshipTypeUnknown
	case 2:
		//FriendshipTypeConfirm
		/**
		 * 2. 友谊确认
		 */
		log.Printf("friend ship confirmed with%s", friendship.Contact().Name())
	case 3:
		//FriendshipTypeReceive
		/*
			1. 新的好友请求
			设置请求后，我们可以从request.hello中获得验证消息,
			并通过`request.accept（）`接受此请求
		*/
		if friendship.Hello() == viper.GetString("addFriendKeywords") {
			if err = friendship.Accept(); err != nil {
				ErrorFormat("添加好友失败", err)
			}
		} else {
			log.Printf("%s未能自动通过好友申请, 因为验证消息是%s", friendship.Contact().Name(), friendship.Hello())
		}
	case 4:
	//FriendshipTypeVerify
	default:
	}
	log.Printf("%s好友关系是: %s", friendship.Contact().Name(), friendship.Type())
}

/*
	@method onHeartbeat 获取机器人的心跳。
	@param {*} user
*/
func onHeartbeat(context *Context, data string) {
	log.Println("========================onHeartbeat👇========================")
	log.Printf("获取机器人的心跳: %s", data)
	//	没啥用
}

func onError(context *Context, err error) {
	ErrorFormat("机器人错误", err)
	return
}

func onMessage(context *Context, message *user.Message) {
	messages := encodeMessage(message)
	if message.Self() { // Bot 自己发的消息
		return
	}
	if message.Age() > 2*60*time.Second {
		log.Println("消息已丢弃，因为它太旧（超过2分钟）")
	}

	if message.Type() == schemas.MessageTypeText { // 文本消息
		if messages.Status { // 群聊状态
			//if strings.Contains(message.Text(),"add") {} // 添加自定义的操作
			if message.MentionSelf() {
				log.Printf("%s@我 %s", messages.UserName, strings.Replace(strings.Replace(message.Text(), "@", "", 1), viper.GetString("bot.name"), "", 1))
				DingMessage(fmt.Sprintf("%s @我 %s", messages.AutoInfo, strings.Replace(strings.Replace(message.Text(), "@", "", 1), viper.GetString("bot.name"), "", 1)))
			}
			if strings.Contains(message.Text(), "基于你的优异表现，+") {
				SayMsg(message, `
					我也要! [旺柴] 给你表演个才艺吧!
					《放鸽子》
				`)
			}
		}
		if strings.Contains("加群", message.Text()) {
			// 邀请进群
		}
	}
	// 打印所有的 TXT 类型的消息
	log.Printf("%s 说: %s", messages.AutoInfo, message.Text())
}

func encodeMessage(message *user.Message) MessageInfo {
	//var result []byte
	UserName := message.From().Name()
	UserID := message.From().ID()
	messages := MessageInfo{
		Status:   false,
		UserName: UserName,
		UserID:   UserID,
		AutoInfo: "用户ID: [" + UserID + "] 用户名称: [" + UserName + "]",
	}
	if message.Room() != nil {
		messages.Status = true
		messages.RoomID = message.Room().ID()
		messages.RoomName = strings.Replace(strings.Replace(message.Room().String(), "Room<", "", 1), ">", "", 1)
		messages.AutoInfo = "群聊ID: [" + messages.RoomID + "] 群聊名称: [" + messages.RoomName + "] " + messages.AutoInfo
	}
	//if result, err = json.Marshal(messages); err != nil {
	//	fmt.Println(err)
	//}
	//jsonStringData := string(result)
	// 这里保存json数据
	return messages
}

func main() {
	i := 0
	// 重试次数 10
	for i <= 10 {
		i++
		// 钉钉推送
		ViperRead()
		DingBotCheck()
		var bot = NewWechaty(WithPuppetOption(wp.Option{
			Token:    viper.GetString("wechaty.wechaty_puppet_service_token"),
			Endpoint: viper.GetString("wechaty.wechaty_puppet_endpoint"),
		}))
		log.Printf("Token:%s", viper.GetString("wechaty.wechaty_puppet_service_token"))
		log.Printf("Endpoint: %s", viper.GetString("wechaty.wechaty_puppet_endpoint"))

		//bot.OnScan(onScan).
		bot.OnLogin(onlogin).
			OnLogout(onLogout).
			OnMessage(onMessage).
			OnRoomInvite(onRoomInvite).
			OnRoomTopic(onRoomTopic).
			OnRoomJoin(onRoomJoin).
			OnRoomLeave(onRoomleave).
			OnFriendship(onFriendship).
			//OnHeartbeat(onHeartbeat).
			OnError(onError)
		//bot.DaemonStart()

		if err := bot.Start(); err != nil {
			// 重启Bot
			ErrorFormat("Bot 错误", err)
			if i > 10 {
				os.Exit(0)
			}
			log.Printf("正在重新启动程序, 当前重试次数: 第%v次", i)
			time.Sleep(10 * time.Second)
		} else {
			i = 0
			// Bot 守护程序
			var quitSig = make(chan os.Signal)
			signal.Notify(quitSig, os.Interrupt, os.Kill)
			select {
			case <-quitSig:
				ViperWrite()
				log.Fatal("程序退出!")
			}
		}
	}
}
