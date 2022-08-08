# t000srapi

	SHOWROOMのフォローしている配信者の一部あるいは全部について、ファンレベルの達成状況を調べるサンプルプログラム
	配信が行われていない場合も達成状況を調べることができます。
	
	以下はソースを作成する（コピペ）するところからの説明です。
	githubからソースをダウンロードして使用する場合は次の記事をご参照ください。
		[【Unix/Linux】Githubにあるサンプルプログラムの実行方法](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/220e38)
		[【Windows】Githubにあるサンプルプログラムの実行方法](https://zenn.dev/chouette2100/books/d8c28f8ff426b7/viewer/e27fc9)

```
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
	sr_acct: ${SRACCT}				<== ログインアカウントを環境変数 SRACCT で与えます。ここに直接アカウントを書くこともできます。
	sr_pswd: ${SRPSWD}				<== ログインパスワードを環境変数 SRPSWD で与えます。ここに直接パスワードを書くこともできます。
	maxnorooms: 3					<== フォローしているルームのうち最初からここで指定した数のルームについてファンレベル達成状況を表示します。
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRACCT xxxxxxxxx
	$ export SRACCT=xxxxxxxx		<== SHOWROOMのアカウント名		Cシェルの場合は setenv SRPSWD xxxxxxxxx
	$ ./t000srapi config.yml
```

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
