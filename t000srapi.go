/*!
Copyright © 2022 chouette.21.00@gmail.com
Released under the MIT license
https://opensource.org/licenses/mit-license.php

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
	$ go mod init					<== 注意：パッケージ部分のソースをダウンロードした場合はimport部分は書き換えず、
	$ go mod tidy					<== 	  go.modに“replace github.com/Chouette2100/srapi ../srapi”みたいに追加します。
	$ go build t000srapi.go
	$ cat config.yml 
	sr_acct: ${SRACCT}				<== ログインアカウントを環境変数 SRACCT で与えます。ここに直接アカウントを書くこともできます。
	sr_pswd: ${SRPSWD}				<== ログインパスワードを環境変数 SRPSWD で与えます。ここに直接パスワードを書くこともできます。
	maxnorooms: 3					<== フォローしているルームのうち最初からここで指定した数のルームについてファンレベル達成状況を表示します。
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRACCT xxxxxxxxx
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRPSWD xxxxxxxxx
	$ ./t000srapi config.yml


	go mod init	で不具合があるときは go mod init t000srapi.go を試してください（ソースの位置（ディレクトリ、ディレクトリ構成）を検討する）
	パスワードが間違っていたなどでログインに失敗したときは、再ログインの前に *_cookies を削除してください。

	今後の課題
	・達成状況が必要な（あるいは不必要な）ルームを選択できるようにする。
	・現在のレベルからレベル１０（あるいはレベル１５）を達成するのに必要な視聴時間、ポイント、コメント数を表示する。

	Ver. 0.0.0
	Ver. 0.1.0 結果出力先の変更を容易にする。CreateLogfile()のインターフェース変更に対応する。
	Ver. 0.1.1 GetActiveFanNextLevel()実行後のエラー処理の位置ずれを直す。

*/
func main() {

	//	ログファイルを設定する。
	logfile := exsrapi.CreateLogfile("", fmt.Sprintf("%d", os.Getpid()))
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
	if status != 0 {
		log.Printf("***** ApiActiveFanNextlevel() returned error. status=%d\n", status)
		return
	}

	pfnc := log.Printf
	//	フォローしている配信者のファンレベル進捗状況を表示する。
	for _, roomafnl := range roomafnls {
		pfnc("********************************************************************************\n")
		pfnc("Room_id=%s ( %s )\n", roomafnl.Room_id, roomafnl.Main_name)
		pfnc("current level = %d\n", roomafnl.Afnl.Level)
		pfnc("next level =    %d\n", roomafnl.Afnl.Next_level.Level)
		for _, c := range roomafnl.Afnl.Next_level.Conditions {
			pfnc("%s\n", c.Label)
			for _, cd := range c.Condition_details {
				pfnc("  %-12s (目標)%5d %-10s (実績)%5d %-10s\n", cd.Label, cd.Goal, cd.Unit, cd.Value, cd.Unit)
			}
		}
	}
}
