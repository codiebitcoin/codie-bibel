# Codie Bible — page design

Seven pages, a sidebar, a player bar for the hymns, a pill for the
device's reading. Every string the pages read goes through `i18n.T`, in
the eleven languages the bundled Bibles read. No Go file runs past 800
lines; each page splits where it would.

The app speaks as the pairing guide's welcome does: *Make your Codie as
Bibel* — the tagline stands under the app's name in the sidebar, on the
guide's first step, and on the workshop page.

```
┌──────────┬─────────────────────────────────────────────┐
│ Codie    │                                             │
│ Bible    │               the page                      │
│          │                                             │
│ Make your│                                             │
│ Codie as │                                             │
│ Bibel    │                                             │
│          ├─────────────────────────────────────────────┤
│ Read     │ ◉ device pill      hymns bar (when playing) │
│ Devotions│                                             │
│ Pray     │                                             │
│ Stories  │                                             │
│ Seasons  │                                             │
│ Music    │                                             │
│ Device   │                                             │
│ Settings │                                             │
└──────────┴─────────────────────────────────────────────┘
```

## Read

*Books grid → chapter rows → the chapter itself.*

- **Books grid**: two groups, Old Testament and New, each a wrapped row
  of book buttons (localized names, `Read` heading). Above them the
  translation picker (a button naming the current translation, dropping
  the list `bibles.List()` gives) and the search field (`Search`,
  `SearchField` label). ← → move between chapters once one is open.
- **Chapter rows**: the book's name and its chapters, numbered, 8 a row.
- **The chapter**: verse numbers as superscripts in the accent color;
  the words of Christ in red (`bible.InRed`); the reader's highlights
  behind their verses in the accent's wash. A tap on a verse opens its
  menu: **Listen** (planned into turns and queued to the device), **Add
  to Stories** (skipped for now — the stories are curated), **Copy**
  (verse and reference), **Highlight** (on/off), **Note** (a note typed
  for the verse, kept locally). Footer: ‹ previous chapter, the
  reference `Ref.String()` in the page's terms, next ›.
- **Search**: on Enter, the results list (reference + the verse, the
  match marked by `bible.Hit.Start`); a result opens its chapter.
- i18n keys: the books' names come from the bible package's tables, not
  the catalog; the rest are catalog keys.

## Devotions

Three cards, top to bottom:

- **Verse of the Day**: the reference and the verse in the current
  translation, `Listen` beside it.
- **Morning & Evening**: the day's reading for the hour (`Morning`/
  `Evening` switch), the verse line and Spurgeon's text; `Listen`.
- **Reading plans**: a row per plan, its name and today's passages;
  tapping a passage opens it in Read; a plan's progress lives in the
  store (`plan:<id>` → the day last read).

## Pray

The ACTS steps as a walk: a step shows its scriptures (`PrayStep.Refs`,
rendered from the current translation) and its prompt; the confession
step carries the privacy line — *"what you write here stays on this
machine"* — and a field. After the last step, the **comfort verses**:
the device reads one (`InsertDialogue` is not used; the queue reads the
verse and the device's persona says a word of comfort of its own — that
chat is the device's work, not the app's). **Prayer list**: the prayers
kept (added from the supplication step), each markable answered.

## Stories

The curated list, in the Bible's order: name and reference line, tap to
open — the passages as paragraphs, `Listen` queuing them to the device
with resume, and the progress bar of turns read.

## Seasons

- The season card: the season's name and color swatch, the feast of the
  day if it is one, and *N days to <feast>*.
- The year's feasts, in order, each with its date and its readings'
  reference line.
- **This week's readings**: the Sunday's four lessons (lectionary), each
  a reference that opens in Read, and `Listen` queuing all four.
- Buttons on a feast: *Let Codie read it* (queue) and *Remind me* (the
  device alarm).

## Music

The hymns list (title, artist, license mark) with the one playing
highlighted; the bar below the pages: play/pause, previous, next,
repeat, volume, and the scrubber. Opening the page queues nothing;
pressing play queues all hymns from the first.

## Device

As shipped in M3: the pairing (status line, `Pair a device`), the
character (Jesus / the monk, each with `Upload a model…` and `Upload a
voice…`), the action group (Subtle / Normal / Strong — what the device
itself exposes).

## Settings

Language (the eleven), text size (Default/Large/Larger), serif on/off,
music auto-play on/off, and the privacy note — *the notes, highlights,
prayers and confessions stay on this machine; nothing of them leaves
it* — with the translation default the Read page's picker sets.
