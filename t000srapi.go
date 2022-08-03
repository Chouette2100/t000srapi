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
	"strconv"
	"strings"

	"github.com/gnue/go-disp_width"

	"github.com/Chouette2100/exsrapi"
	"github.com/Chouette2100/srapi"
)

/*
	SHOWROOMのフォローしている配信者の一部あるいは全部について、ファンレベルの達成状況を調べるサンプルプログラム
	配信が行われていない場合も達成状況を調べることができます。

	$ cd ~/go/src/t000srapi
	$ vi t000srapi.go				<== このソースを作成する。
	$ go mod init					<== 注意：パッケージ部分のソースをダウンロードした場合はimport部分は書き換えず、
	$ go mod tidy					<== 	  go.modに“replace github.com/Chouette2100/srapi ../srapi”みたいに追加します。
	$ go build t000srapi.go
	$ cat config.yml
	target:
	- 視聴時間,    10, 10, 10, 15,  15,  15,  15,  30,   30,   30
	- 無料ギフト,   0, 10, 40, 49, 396, 495, 495, 990, 2475, 4950
	- コメント数,   0,  0,  0,  0,   0,   1,   1,   1,    1,    1
	sr_acct: ${SRACCT}
	sr_pswd: ${SRPSWD}
	maxnorooms: 3
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRACCT xxxxxxxxx
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRPSWD xxxxxxxxx
	$ ./t000srapi config.yml


	go mod init	で不具合があるときは go mod init t000srapi.go を試してください（ソースの位置（ディレクトリ、ディレクトリ構成）を検討する）
	パスワードが間違っていたなどでログインに失敗したときは、再ログインの前に *_cookies を削除してください。

	今後の課題
	・達成状況が必要な（あるいは不必要な）ルームを選択できるようにする。
	・現在のレベルからレベル１０（あるいはレベル１５）を達成するのに必要な視聴時間、ポイント、コメント数を表示する。
		(制約はありますが、レベル１０までの必要数を表示できるようにしてあります Ver.0.2.1)

	Ver. 0.0.0
	Ver. 0.1.0 結果出力先の変更を容易にする。CreateLogfile()のインターフェース変更に対応する。
	Ver. 0.1.1 GetActiveFanNextLevel()実行後のエラー処理の位置ずれを直す。
	Ver. 0.2.1 レベル10までに必要な視聴時間、ポイント、コメント数を表示する。
	Ver. 0.2.2 レベル0のときはtarget[cd.Label][9] - cd.Valueを必要な視聴時間として表示する(target[cd.Label][roomafnl.Afnl.Level-1]は存在しない)

*/

//	全角を含む文字列の表示幅を指定した（半角）文字数にするため必要なスペースを追加する。
func SetWidthConst(str string, width int) string {
	n := disp_width.Measure(str)
	if n > width {
		return str
	}
	return str + strings.Repeat(" ", width-n)
}

type Config struct {
	SR_acct    string //	SHOWROOMのアカウント名
	SR_pswd    string //	SHOWROOMのパスワード
	MaxNoRooms int    // データを取得するルームの最大数
	Target     []string
}

//	目標を算出する
func setTarget(config *Config) (target map[string][]int, err error) {

	target = map[string][]int{}

	if len(config.Target) != 3 {
		return nil, fmt.Errorf("設定ファイルtargetには三つの要素が必要です。")
	}

	labels := map[string]bool{"視聴時間": true, "無料ギフト": true, "コメント数": true}

	for i := 0; i < 3; i++ {
		sTgt := strings.Split(config.Target[i], ",")
		if len(sTgt) != 11 {
			return nil, fmt.Errorf("設定ファイルのtargetの一つの要素はラベルと10個の目標値をカンマ区切りで書く必要があります。")
		}
		for j := 0; j < 10; j++ {
			sTgt[j+1] = strings.Replace(sTgt[j+1], " ", "", -1)
			itgt, err := strconv.Atoi(sTgt[j+1])
			if err != nil {
				return nil, fmt.Errorf("設定ファイルのtargetの目標値には数値を書く必要があります。\"%s\"", sTgt[j+1])
			}
			if _, ok := labels[sTgt[0]]; !ok {
				return nil, fmt.Errorf("設定ファイルのtargetのラベルには\"視聴時間\"、\"無料ギフト\"、\"コメント\"のいずれかを書く必要があります。")
			}
			target[sTgt[0]] = append(target[sTgt[0]], itgt)
		}
		for j := 1; j < 10; j++ {
			target[sTgt[0]][j] += target[sTgt[0]][j-1]
		}
	}
	//	log.Printf("%+v\n", target)
	return target, nil
}

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
	status := exsrapi.LoadConfig(os.Args[1], &config)
	if status != 0 {
		log.Printf("LoadConfig error: %d\n", status)
	}

	target, err := setTarget(&config)
	if err != nil {
		log.Printf("setTarget error: %s\n", err)
		return
	}

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
	roomafnls, status := exsrapi.GetActiveFanNextLevel(client, userid, rooms)
	if status != 0 {
		log.Printf("***** ApiActiveFanNextlevel() returned error. status=%d\n", status)
		return
	}

	pfnc := fmt.Printf
	//	フォローしている配信者のファンレベル進捗状況を表示する。
	for _, roomafnl := range roomafnls {
		pfnc("********************************************************************************\n")
		pfnc("Room_id=%s ( %s )\n", roomafnl.Room_id, roomafnl.Main_name)
		pfnc("current level = %d\n", roomafnl.Afnl.Level)
		pfnc("next level =    %d\n", roomafnl.Afnl.Next_level.Level)
		for _, c := range roomafnl.Afnl.Next_level.Conditions {
			pfnc("%s\n", c.Label)
			for _, cd := range c.Condition_details {
				ok := true
				dt := 0
				if _, ok = target[cd.Label]; ok {
					if roomafnl.Afnl.Level > 0 {
						dt = target[cd.Label][9] - target[cd.Label][roomafnl.Afnl.Level-1] - cd.Value
					} else {
						dt = target[cd.Label][9] - cd.Value
					}
				}
				if ok && roomafnl.Afnl.Level <= 9 && dt > 0 {
					pfnc("  %s (目標)%5d %s (実績)%5d %s (Lv10まであと)%5d %s\n",
						SetWidthConst(cd.Label, 10), cd.Goal, SetWidthConst(cd.Unit, 8), cd.Value, SetWidthConst(cd.Unit, 8), dt, SetWidthConst(cd.Unit, 8))
				} else {
					pfnc("  %s (目標)%5d %s (実績)%5d %s\n",
						SetWidthConst(cd.Label, 10), cd.Goal, SetWidthConst(cd.Unit, 8), cd.Value, SetWidthConst(cd.Unit, 8))
				}
			}
		}
	}
}
