# boggle

A TUI game of single-player boggle inspired by [wordshake.com/boggle](https://wordshake.com/boggle).

Play a hosted game: `ssh boggle.fly.dev`.

## Installation

```console
$ go install github.com/lukasschwab/boggle@latest
```

## Usage

See also `boggle --help`.

```
boggle - play an unsanctioned game of boggle.

This game uses the Collins Scrabble Words 2019 dictionary, a three-minute timer,
and requires words be at least four letters long.

Each game outputs a .boggle file describing the board and your performance. You
can replay a board by passing a .boggle file to this program:

    $ boggle -file past-game.boggle

Or by providing the "serialized" short-form description of the board, included
in the YAML frontmatter of each .boggle file:

    $ boggle -board Y3VkbnF1dG5kZHVybHllYXg=

The following options are available:

  -board string
        serialized board string
  -file string
        .boggle file to configure game
  -new
        use the "New" Boggle dice set (since 2008).
        See: https://www.bananagrammer.com/2013/10/the-boggle-cube-redesign-and-its-effect.html
  -skip
        skip interactive game
  -solve
        print all possible words on board after game
  -url string
        web URL of a public .boggle file to configure game
```
