# gopher-game

gopher-game is a simple game in a terminal.

It is an attempt to play with Go bindings to OpenGL and to make a game for my little daughter.
I was inspired by awesome game [zombies-on-ice](https://github.com/loov/zombies-on-ice).

![gopher-game](https://github.com/austinov/gopher-game/blob/assets/screenshot.gif)


In the beginning, as usual, run:
```
    $ go get github.com/austinov/gopher-game
    $ cd $GOPATH/github.com/austinov/gopher-game
    $ glide up
    $ go build
    $ ./gopher-game
```

### Download

Alfa version:

* [Linux 64bit](https://github.com/austinov/gopher-game/blob/assets/gopher-game.linux-64bit.tar.gz)

#### Prerequisites

- [go-gl/gl](https://github.com/go-gl/gl)
- [go-gl/glfw](https://github.com/go-gl/glfw)
- [go.geo](https://github.com/paulmach/go.geo)


#### TODO

- more realistic gravity
- enemies intelligence
