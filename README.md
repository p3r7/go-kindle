# go-kindle

branch using [`go-text`](https://github.com/go-text), esp. the [`go-text/typesetting`](https://github.com/go-text/typesetting) package.

kinda works but won't run on kindle itself due to the kernel being too old (encountering [#24980](https://github.com/golang/go/issues/24980))...

## aknowledgments

[`text`](./text/) sub-package is a modified version of [`ebiten/v2/text/v2`](https://pkg.go.dev/github.com/hajimehoshi/ebiten/v2/text/v2) ([about](https://ebitengine.org/en/blog/textv2.html)), esp. the [`texti18n` example code](https://github.com/hajimehoshi/ebiten/blob/main/examples/texti18n/main.go). indeed, `go-text/typesetting` is pretty low-level and poorly documented so i had to take a lot of inspiration from `ebiten` to make it work...
