# Codie Bible

**Make your Codie as Bibel** — pair your codie device, and it reads the Bible aloud: the stories, the psalms, the verses of the day. It speaks as Jesus, or as a monk for the quieter rooms.

Codie Bible is a desktop app for [CODE27's codie](https://code27.com) companion device, built with [egoist/mygo](https://github.com/egoist/mygo) — pure Go, no cgo, native windows on macOS, Windows and Linux.

## What it does

- **Read** — thirteen public-domain Bibles (KJV, ASV, WEB, 和合本 简体/繁體, Luther 1912, Reina-Valera 1909, Louis Segond 1910, Almeida, Синодальный, 口語訳, 개역한글, Smith–Van Dyke), search over any of them, highlights, notes, and the words of Christ in red.
- **Devotions** — the verse of the day, Spurgeon's *Morning & Evening* (1866), and reading plans: through the Bible in a year, the New Testament in 90 days, Advent.
- **Pray** — a guided prayer in the ACTS steps, each with its scriptures. What is confessed stays on this machine; the device only hears the comfort verses read aloud.
- **Stories** — forty of the Bible's stories, from the creation to the new creation, told aloud on the device with resume.
- **Seasons** — the church year from Advent to Christ the King, its feasts and colors, the Revised Common Lectionary's Sunday readings, and feast alarms on the device.
- **Music** — eight public-domain hymns and chorales (Bach, the Sacred Harp, the US service bands' carols) playing quietly in the background.
- **A figure and a voice of its own** — Jesus wears the app's own figure of him (a VRM brought from VRoid Hub), and at the first pairing the app quietly clones his voice from a blessing it ships; the monk wears pixiv's MIT-licensed sample humanoid. Upload a model or clone a voice to replace either.
- **Eleven languages** — the interface speaks what the Bibles read: English, 简体中文, 繁體中文, Deutsch, Español, Français, Português, Русский, 日本語, 한국어, العربية (laid out from the right).

Pairing walks four steps — the address, the six-figure code, the character, a verse read aloud. The device plays its own motions as it speaks; the app sets how strong they read, and may give a character a model and a cloned voice of the reader's own.

## Licenses

- The Bibles and their sources: [internal/bible/data/LICENSES.md](internal/bible/data/LICENSES.md) — all public domain.
- The hymns: Wikimedia Commons recordings, all public domain except Bach BWV 147 (CC0) and BWV 140 (CC BY-SA 2.0, MIT Concert Choir) — see [internal/hymn](internal/hymn).
- Spurgeon's *Morning and Evening* (1866): public domain; transcription via the dwahrens/morning-and-evening data file (Apache-2.0 covers the packaging, not the text).
- The lectionary table: from [marmanold/Date-Lectionary](https://github.com/marmanold/Date-Lectionary) (BSD-2-Clause); the readings themselves are scripture references.
- The figure of Jesus: a VRM titled "Jesus" from [VRoid Hub](https://hub.vroid.com) (model 307680147834872455), used by the terms the model carries there.
- The seed of Jesus's voice: a synthesized blessing (KJV words, a system voice); it is cloned on the reader's own device, never shipped as audio of a real person.
- The monk's default figure: pixiv's [three-vrm](https://github.com/pixiv/three-vrm) sample humanoid (MIT).

## Building

Needs Go ≥ 1.27 and bun or Node ≥ 22.6 (the mygo config is TypeScript).

```sh
GOTOOLCHAIN=auto go tool mygo build            # this platform
go run ./cmd/genbible -only kjv                # redo one Bible (network)
go test ./internal/...                         # the tests
```

The generated data (Bibles, devotional, lectionary) is committed: a build never needs the network.

## License

AGPL-3.0 — the same as the codie-music-player whose player plumbing this app shares. See [LICENSE](LICENSE).
