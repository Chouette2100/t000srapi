/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

Ver. 0.0.0

*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

type Config struct {
	SR_acct    string //	SHOWROOMのアカウント名
	SR_pswd    string //	SHOWROOMのパスワード
	MaxNoRooms int    // データを取得するルームの最大数
}


/*
	SHOWROOMのフォローしている配信者の一部あるいは全部について、ファンレベルの達成状況を調べるサンプルプログラム
	配信が行われていない場合も達成状況を調べることができます。

	$ cd ~/go/src/t000srapi
	$ go mod init					<== 注意：パッケージ部分のソースをダウンロードした場合はimport部分はそのままにしておいて
	$ go mod tidy					<== 	  go.modに“replace github.com/Chouette2100/srapi ../srapi”みたいなのを追加します。
	$ go build t000srapi.go
	$ cat config.yml 
	sr_acct: ${SRACCT}
	sr_pswd: ${SRPSWD}
	maxnorooms: 3					<== フォローしているルームのうち最初からここで指定した数のルームについてファンレベル達成状況を表示します。
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRACCT xxxxxxxxx
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRPSWD xxxxxxxxx
	$ ./t000srapi config.yml


*/
func main() {

	//	ログファイルを設定する。
	logfile := exsrapi.CreateLogfile()
	defer logfile.Close()

	if len(os.Args) != 2 {
		//      引数が足りない(設定ファイル名がない)
		log.Printf("usage:  %s NameOfConfigFile\n", os.Args[0])
		return
	}

	//	設定ファイルを読み込む
	var config Config
	exsrapi.LoadConfig(os.Args[1], &config)

	//	cookiejarがセットされたHTTPクライアントを作る
	client, jar, err := exsrapi.CreateNewClient(config.SR_acct)
	if err != nil {
		log.Printf("CreateNewClient() returned error %s\n", err.Error())
		return
	}
	//	すべての処理が終了したらcookiejarを保存する。
	defer jar.Save()

	//	SHOWROOMのサービスにログインし、ユーザIDを取得する。
	userid, status := exsrapi.LoginShowroom(client, config.SR_acct, config.SR_pswd)
	if status != 0 {
		log.Printf(" LoginShowroom returned status = %d\n", status)
		return
	}

	//	フォローしている配信者のリストを作成する。
	rooms, status := srapi.CrwlFollow(client, config.MaxNoRooms)
	if status != 0 {
		log.Printf(" CrwlFollow returned status = %d\n", status)
		return
	}

	//	配信者のリストから、ファンレベルの達成状況を調べる。
	roomafnls, status := exsrapi.GetActiveFanNextLevel( client, userid, rooms )

	//	フォローしている配信者のリストを表示する。
	for _, roomafnl := range roomafnls {
		fmt.Printf("********************************************************************************\n")
		if status != 0 {
			log.Printf("***** ApiActiveFanNextlevel() returned error. status=%d\n", status)
			return
		}
		fmt.Printf("Room_id=%s ( %s )\n", roomafnl.Room_id, roomafnl.Main_name)
		fmt.Printf("current level = %d\n", roomafnl.Afnl.Level)
		fmt.Printf("next level =    %d\n", roomafnl.Afnl.Next_level.Level)
		for _, c := range roomafnl.Afnl.Next_level.Conditions {
			fmt.Printf("%s\n", c.Label)
			for _, cd := range c.Condition_details {
				fmt.Printf("  %-12s (目標)%5d %-10s (実績)%5d %-10s\n", cd.Label, cd.Goal, cd.Unit, cd.Value, cd.Unit)
			}
		}
	}
}
