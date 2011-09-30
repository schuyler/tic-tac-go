Tic-Tac-Go
-=-=-=-=-=-

In 1997, as an assignment for an undergraduate-level "Introduction to
Artificial Intelligence" course, I implemented a tic-tac-toe playing program in
Lisp. I found that, with the tree search depth set sufficiently high, the game
could beat me (if I wasn't paying attention), but I could never beat it.

I have replicated this feat in the Go programming language, just because Go
looked cool, and I wanted an excuse to learn it. I have implemented the "AI" in
just about the least efficient way possible, purely as an excuse to learn how
to use Go's types, channels, and goroutines.

I will say that I am pretty impressed at the "expressiveness of Go" that Rob
Pike keeps banging on about. Consider the following gem from ttt.go:

    game := NewGame([2]Player{'X','O'})

In other words: This game is a new game, with 2 players, 'X' and 'O'. When you
can write code like that *and* compile it to a native binary, well, now you're
getting somewhere.

I've posted the code on Github because why the heck not. Consider it public
domain. Patches welcome.

SDE
9/30/2011
