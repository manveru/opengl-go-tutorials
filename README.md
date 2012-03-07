# Tutorials for OpenGL in Go

Tutorials translated from "NeHe Productions":http://nehe.gamedev.net/ into Go with SDL.

## Dependencies

First we need to install the dependencies:

    go get github.com/banthar/gl github.com/banthar/glu github.com/banthar/Go-SDL/sdl

## Running the examples

You can use `go run`, like this:

    cd lesson01
    go run lesson01.go

Please note that starting with lesson06, you _have_ to cd into the directory
because we start using external data.
