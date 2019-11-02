// The jpeg command reads a PNG image from stdin
// and writes it as a JPEG to stdout.
package main

import (
  "fmt"
  "image"
  "image/jpeg"
  _ "image/png" // register png decoder
  "io"
  "os"
)

func main() {
  if err := toJPEG(os.Stdin, os.Stdout); err != nil {
    fmt.Fprintf(os.Stderr, "jpeg: %v\n", err)
    os.Exit(1)
  }
}

func toJPEG(in io.Reader, out io.Writer) error {
  img, kind, err := image.Decode(in)
  if err != nil {
    return err
  }
  fmt.Fprintln(os.Stderr, "Input format =", kind)
  return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
}

/*
// Usage

$ go build ./main.go
$ cat path/to/someimage.png | ./jpegconv > convimage.jpg
Input format = png
$ ls
jpegconv main.go convimage.jpg

*/
