package main

import (
  "github.com/kr/pretty"
  "github.com/bmeg/arachne/aql"
  "os"
)

func main() {
  cli, err := aql.Connect("10.50.50.123:9090", false)
  if err != nil {
    panic(err)
  }

  v, err := cli.GetVertex(os.Args[1], os.Args[2])
  if err != nil {
    panic(err)
  }
  pretty.Println(v)


}
