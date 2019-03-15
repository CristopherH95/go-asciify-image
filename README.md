# go-asciify-image
Basic image to ascii art converter in Go

This is a simple program which takes in an image and produces an ascii art version.
In order to prevent massive ascii art pictures, it resizes images whose width or height is over 200px.

## build
Build to an executable with
```
go build
```

## usage
Run with
```
./asciify <file-path>
```
