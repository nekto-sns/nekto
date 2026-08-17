# nekto-server

本プロジェクトはAGPLv3のもとで公開されています。詳細は [LICENSE.md](https://github.com/nekto-sns/nekto-server/blob/main/LICENSE.md) をご覧ください。

SNSのAPIサーバー部です。

未完成なので、多くの機能はまだ実装されていません。

現在実装されている機能は、
- ユーザー名からユーザー情報を取得
- Scratch Authを使って認証

の二つです。

## セットアップ

ビルド済みバイナリはまだリリースできていないので、ソースコードからビルドします。

goが必要です。

まずリポジトリをローカルにcloneします

```bash
git clone https://github.com/nekto-sns/nekto-server.git
```

次に、cloneしたリポジトリ内に移動してサーバーをビルドします
```bash
cd nekto-server
go mod tidy
go build
```

ビルドが終わったら、その次にサーバーを起動するための設定をします。

.envまたは環境変数に設定します。

設定は[.env.template](https://github.com/nekto-sns/nekto-server/blob/main/.env.template)を参考にしてください。

最後に

```bash
./server
```

でサーバーを起動すると、`0.0.0.0:<設定したポート番号>`でサーバーが起動します。
