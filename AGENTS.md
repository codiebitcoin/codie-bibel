# AGENTS.md

## Tests

Only write a test for behaviour that can actually break. If the code under test could change without the test noticing, or the test would fail on a harmless rewording, don't write it.

**Write tests for:**

- What input does: clicks, keys, drags, gestures and scrolling, and the state or effect they lead to (opening a chapter when a book is clicked, Space playing the hymns, a verse's menu reading it aloud, praying through the ACTS steps, leaving a prayer marked answered, switching translation).
- Logic that doesn't need a view: parsing verse references, the liturgical calendar (Easter and the feasts from it), the lectionary's week, splitting a reading into sentences for the device, playback position and resume, search across a translation.
- A regression you just fixed. Say in the test's comment what went wrong.
- Frame budgets, in `perf_test.go`, so idle views don't start asking for frames they don't need.

**Don't write tests that:**

- Check that a list of labels shows on screen. That repeats the view code and breaks whenever the wording changes.
- Check default settings or the order of rows in a form.
- Test MyGo itself, for example that a spring moves from 0 to 1.
- Depend on the wall clock. `ui.NewTester` starts its clock at the current time, so anything with a phase taken from the clock passes or fails depending on when the test runs. Also, `FramePending` ignores frames waiting on `After` or `Painter.After`. The calendar and devotions take their date as an argument, so a test can pass its own.

**How:**

- Build views with `ui.NewTester` and set up the app with `newTestApp` (in `internal/app/helpers_test.go`). Find elements by their accessibility label, then `Click`, `RightClick`, `Press`/`Move`/`Release` or `Key`, and check the app's state afterwards.
- The device is an interface (`device.Device`): tests set a `fakeDevice` that records what was asked of it and answers what the test needs.
- See what comes back from a goroutine with `posts` and `runPost`.
- Look at visuals instead of asserting on them. The `MM_FRAMES` tests save PNGs into that directory and skip when it isn't set. Check pixels only for coverage, such as whether a layer still covers the window.
- Start each test with a one-sentence comment saying the behaviour it checks, in the style of the code.
- Run `GOTOOLCHAIN=auto go test ./internal/...` before committing, and run a new UI test with `-count=5`. Fix or delete a flaky test; don't retry it.

## Text and data

- Every string the reader sees goes through `i18n.T`/`Tf` with a `zh_cn.go` entry: `TestNoChineseInCode` rejects Chinese literals in code, and `TestCatalogHasEveryKey` rejects an untranslated key. Scripture, devotional and hymn text is data — those packages are excluded from the Chinese check.
- The Bibles, the devotional, the plans, the lectionary and the hymns are generated into `internal/*/data` by `cmd/genbible` and committed: a build never needs the network. Don't edit files in `data/` by hand.
