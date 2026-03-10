# Boggle Experiments: Updated Results

Re-run with corrected EEINSU die (was EEINSV). Also fixes a bug in
`cmd/sameness/sameness.go` where `evaluate()` always used `ClassicDice`
regardless of the `dice` parameter passed in.

## Playable Words (Dictionary Sizes)

| Metric | Classic | New (corrected) | Blog (old) |
|:---|---:|---:|---:|
| Playable words | 276,339 | 261,803 | 257,822* |
| Classic-exclusive | 14,663 | — | ~18,517 |
| New-exclusive | — | 127 | 134 |
| Classic-exclusive with F+K | 1,572 | — | 1,572 |

*Inferred from blog: 276,339 - 18,517 = 257,822.

The corrected die (U instead of V) adds ~4,000 more playable words to the
New Boggle dictionary, since U is substantially more common in English words
than V.

## Word Counts per Board (10,000 boards)

| Dice Set | Mean | Min | Max |
|:---|---:|---:|---:|
| Classic | 80.4 | 0 | 452 |
| New (corrected) | 90.2 | 0 | 510 |

## Sameness: Unique Words in 10-Game Sessions (1,000 sessions)

| Dice Set | Mean unique words | Min | Max |
|:---|---:|---:|---:|
| Classic | 771.9 | 332 | 1,433 |
| New (corrected) | 850.3 | 405 | 1,543 |

Note: the original `sameness.go` had a bug where `evaluate()` always used
`ClassicDice` inside the loop regardless of the `dice` parameter. With this
fix, New Boggle shows ~10% more unique words per session than Classic, but
the distributions still largely overlap — the blog post's core conclusion
("effectively the same diversity") remains directionally correct.

## Max Word Count Board (100,000 boards, Classic)

516 words:
```
l  e  s  h
r  u  t  a
n  o  l  g
m  s  e  a
```
