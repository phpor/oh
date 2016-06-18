# A surprisingly powerful Unix shell

Oh is a Unix shell. If you've used other Unix shells, oh should feel
familiar.

![gif](img/oh.gif)

Where oh diverges from traditional Unix shells is its programming
language features.

At its core, oh is a heavily modified dialect of the Scheme programming
language, complete with first-class continuations and proper tail
recursion. Like early Scheme implementations, oh exposes environments
as first-class values. Oh extends environments to allow both public and
private members and uses these extended first-class environments as the
basis for its prototype-based object system.

Written in Go, oh is also a concurrent programming language. It exposes
channels, in addition to pipes, as first-class values. As oh uses the
same syntax for code and data, channels and pipes can, in many cases, be
used interchangeably. This homoiconic nature also allows oh to support
fexprs which, in turn, allow oh to be easily extended.
In fact, much of oh (currently, 529 of 6047 
lines of code) is written in oh.

For a detailed comparison to other Unix shells see: [Comparing oh to other Unix Shells](https://htmlpreview.github.io/?https://raw.githubusercontent.com/michaelmacinnis/oh/master/doc/comparison.html)

## Installing

With Go 1.5 or greater installed,

    go get github.com/michaelmacinnis/oh

(Oh compiles and runs, but should be considered experimental, on Windows.)

## Using

For more detail see: [Using oh](doc/manual.md)

## License

[MIT](LICENSE)

