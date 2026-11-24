# Codie Bible

![Codie Bible](docs/banner.png)

**Make your Codie as Bibel** — codie をペアリングすると、聖書を声に出して読み上げます。物語、詩篇、その日の聖句。イエスの姿で語り、静かな場には修道士の姿にもできます。

Codie Bible は [CODE27 の codie](https://code27.com) 伴読デバイス向けのデスクトップアプリで、[egoist/mygo](https://github.com/egoist/mygo) で作っています。純粋 Go、cgo なし。macOS、Windows、Linux でネイティブなウィンドウを出します。

English: [README.md](README.md) · 中文: [README.zh.md](README.zh.md)

## できること

- **閲読** — パブリックドメインの聖書が十三訳（KJV、ASV、WEB、和合本 简体/繁體、Luther 1912、Reina-Valera 1909、Louis Segond 1910、Almeida、シノダル訳、口語訳、개역한글、Smith–Van Dyke）。どの訳でも検索、ハイライト、メモ、キリストの言葉は赤字。
- **霊想** — その日の聖句、スパージョンの『朝夕の霊想』（1866）、読経計画（一年で通読、新約聖書を 90 日、待降節）。
- **祈り** — ACTS の手順に沿った導きの祈り。各段階に聖句がつきます。告白の内容はこのマシンにだけ残り、デバイスが聞くのは慰めの聖句だけです。
- **物語** — 創世から新しい天と地まで、聖書の物語が四十。デバイスが声に出して語り、途中から続けられます。
- **節期** — 待降節からキリスト王節までの教会暦、祭日と礼拝色、改訂共同聖書日課の主日聖句、デバイスへの祭日の告知。
- **音楽** — パブリックドメインの讃美歌とコラールが八曲（バッハ、Sacred Harp、米軍楽隊のキャロル）。背景で静かに流れます。
- **姿と声** — イエスはアプリ自身の姿、修道士は祭服の僧侶（どちらも VRoid Hub の VRM）。初めてのペアリングで、同梱の祝福の言葉からイエスの声を静かにクローンします。モデルや声は後から入れ替えられます。
- **十一の言語** — 画面は同梱の聖書と同じ言語を話します。English、简体中文、繁體中文、Deutsch、Español、Français、Português、Русский、日本語、한국어、العربية（右から左へ組みます）。

ペアリングは四段階です。アドレス、六桁のコード、姿、一節の音読。デバイスは話すとき自分の動きをつけ、アプリはその強さを決めます。読者自身のモデルとクローンした声を姿に渡すこともできます。

## ライセンス

- 聖書とその出典: [internal/bible/data/LICENSES.md](internal/bible/data/LICENSES.md) — すべてパブリックドメイン。
- 讃美歌: Wikimedia Commons の録音。バッハ BWV 147（CC0）と BWV 140（CC BY-SA 2.0、MIT Concert Choir）を除きパブリックドメイン。[internal/hymn](internal/hymn) を見てください。
- スパージョン『朝夕の霊想』（1866）: パブリックドメイン。翻刻は dwahrens/morning-and-evening のデータファイル（Apache-2.0 は包装にかかり、本文にはかからない）。
- 聖書日課表: [marmanold/Date-Lectionary](https://github.com/marmanold/Date-Lectionary)（BSD-2-Clause）。日課そのものは聖句の引用です。
- イエスの姿: [VRoid Hub](https://hub.vroid.com) の "Jesus"（model 307680147834872455）。そこに記された条件に従います。
- イエスの声の種: 合成した祝福の言葉（KJV の文、システムの声）。読者自身のデバイス上でクローンし、実在の人の音声としては配布しません。
- 修道士の標準の姿: [VRoid Hub](https://hub.vroid.com) の祭服の僧侶（model 8631357567214777894、"アバター_神父服モデル"、ALL・nee・t）。CODE27 の判断で同梱しています。そのモデルの条件は利用を制限しているので、合わないときは自分のモデルを渡してください。pixiv の [three-vrm](https://github.com/pixiv/three-vrm) サンプル人形（MIT）が控えです。

## ビルド

Go 1.27 以上と、bun または Node 22.6 以上が要ります（mygo の設定は TypeScript）。

```sh
GOTOOLCHAIN=auto go tool mygo build            # このプラットフォーム
go run ./cmd/genbible -only kjv                # 一つの訳を再生成（ネットワーク）
go test ./internal/...                         # テスト
```

生成されたデータ（聖書、霊想、聖書日課）はリポジトリに入っています。ビルドにネットワークは要りません。

## ライセンス

AGPL-3.0 — 再生の仕組みを共有する codie-music-player と同じです。[LICENSE](LICENSE) を見てください。
