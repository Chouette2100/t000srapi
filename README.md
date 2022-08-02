# t000srapi

	SHOWROOMのフォローしている配信者の一部あるいは全部について、ファンレベルの達成状況を調べるサンプルプログラム
	配信が行われていない場合も達成状況を調べることができます。
```
	$ cd ~/go/src/t000srapi
	$ vi t000srapi.go				<== このソースを作成する。
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
```

	go mod init	で不具合があるときは go mod init t000srapi.go を試してください（ソースの位置（ディレクトリ、ディレクトリ構成）を検討する）
	パスワードが間違っていたなどでログインに失敗したときは、再ログインの前に *_cookies を削除してください。

	今後の課題
	・達成状況が必要な（あるいは不必要な）ルームを選択できるようにする。
	・現在のレベルからレベル１０（あるいはレベル１５）を達成するのに必要な視聴時間、ポイント、コメント数を表示する。

	Ver. 0.0.0
	Ver. 0.1.0 結果出力先の変更を容易にする。CreateLogfile()のインターフェース変更に対応する。
