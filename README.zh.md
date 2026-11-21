# Codie Bible

**Make your Codie as Bibel** —— 配对你的 codie 设备,让它朗读圣经:故事、诗篇、每日金句。它以耶稣的形象说话,也可换成修士,适合安静的场合。

Codie Bible 是 [CODE27 codie](https://code27.com) 伴读设备的桌面应用,用 [egoist/mygo](https://github.com/egoist/mygo) 构建 —— 纯 Go、无 cgo,原生窗口覆盖 macOS / Windows / Linux。

## 功能

- **阅读** —— 内置 13 个公有领域译本(KJV、ASV、WEB、和合本简繁、路德1912、RV1909、Louis Segond 1910、Almeida、俄语_synodal、口語訳、개역한글、Smith–Van Dyke),全文搜索、划线、笔记、耶稣的话红字显示。
- **灵修** —— 每日金句、司布真《晨昏灵修》(1866)、读经计划(一年通读、90天新约、将临期)。
- **祷告** —— ACTS 引导式祷告,每步配经文。忏悔的内容只留在这台电脑上;设备只会听到安慰经文的朗读。
- **故事** —— 圣经的四十个故事,从创世到新天新地,由设备讲给你听,支持断点续读。
- **节期** —— 从将临期到基督君王节的教会年历、节期与礼仪颜色、经课表的主日经课,并可将节日提醒推到设备。
- **音乐** —— 八首公有领域赞美诗与圣咏(巴赫、Sacred Harp、美军军乐队的圣诞颂歌)轻声背景播放。
- **自带默认形象** —— 两个形象默认穿着应用内置的 VRM(pixiv MIT 示例人形);可上传模型或克隆声音替换。
- **十一种界面语言** —— 界面语言与内置圣经语言一致:English、简体中文、繁體中文、Deutsch、Español、Français、Português、Русский、日本語、한국어、العربية(从右向左排版)。

配对向导四步:地址 → 六位配对码 → 选择形象 → 试读一节。设备说话时会自带动作;应用可设定动作幅度,也可为形象上传模型与克隆声音。

## 许可

- 圣经与来源:[internal/bible/data/LICENSES.md](internal/bible/data/LICENSES.md) —— 全部公有领域。
- 赞美诗:Wikimedia Commons 录音,公有领域为主;巴赫 BWV 147 为 CC0、BWV 140 为 CC BY-SA 2.0(MIT Concert Choir)—— 见 [internal/hymn](internal/hymn)。
- 司布真《晨昏灵修》(1866):公有领域;转录取自 dwahrens/morning-and-evening 数据文件(Apache-2.0 只约束打包,不约束文本)。
- 经课表:取自 [marmanold/Date-Lectionary](https://github.com/marmanold/Date-Lectionary)(BSD-2-Clause);经课本身只是圣经引用。
- 默认形象:pixiv [three-vrm](https://github.com/pixiv/three-vrm) 示例人形(MIT)。

## 构建

需要 Go ≥ 1.27 与 bun / Node ≥ 22.6(mygo 配置是 TypeScript)。

```sh
GOTOOLCHAIN=auto go tool mygo build            # 当前平台
go run ./cmd/genbible -only kjv                # 重新生成某一译本(需联网)
go test ./internal/...                         # 测试
```

生成的数据(圣经、灵修、经课表)已提交入仓,构建不需要联网。

## 许可

AGPL-3.0 —— 与共用其播放器管线的 codie-music-player 一致。见 [LICENSE](LICENSE)。
