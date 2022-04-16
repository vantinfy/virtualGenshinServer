package main

import (
	"VirtualGenshinServer/csvs"
	"VirtualGenshinServer/manage"
	"VirtualGenshinServer/module"
	"fmt"
)

func main() {
	fmt.Println("--- service start ---")

	// 加载配置
	csvs.CheckLoadCsv()

	// 禁词公共管理模块协程
	go manage.GetBanWords().Run()

	// 测试阶段只对玩家线程增加wait
	module.Wait.Add(1)

	player := module.NewPlayer()
	go player.Run()

	module.Wait.Wait()

	// 监听玩家连接，启用协程处理
	//ticker := time.NewTicker(time.Second * 3)
	//for {
	//	select {
	//	case <-ticker.C:
	//		player := module.NewPlayer()
	//		go player.Run()
	//	}
	//}

	// todo 大部分的fmt错误改为error返回	实现log模块 打印等级debug info warn error	打印信息暂定: [2022/04/04 02:42:51 [I] [main.go:16] hello world]
	/* beego 等级参考
	LevelEmergency = iota 紧急
	LevelAlert        //1 警告
	LevelCritical	  //2 关键/重要的
	LevelError		  //3 错误
	LevelWarning	  //4 警告
	LevelNotice		  //5 注意
	LevelInformational//6 信息
	LevelDebug		  //7 调试
	*/
}
