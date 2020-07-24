# OpenCC-go
golang 简体繁体中文互转


## go mod demo

(```)
    package main

    import (
        "fmt"

        ccgo "github.com/ApesPlan/OpenCC-go"
    )

    func main() {
        // 简体转繁体
        str := "刘德华"
        s2t, err := ccgo.New("s2t")
        if err != nil {
            fmt.Printf("Error: %s\n", err)
        }
        out, err := s2t.Convert(str)
        if err != nil {
            fmt.Printf("Error: %s\n", err)
        }
        fmt.Println(out)

        // 繁体转简体
        // str := "劉德華"
        // t2s, err := ccgo.New("t2s")
        // if err != nil {
        // 	fmt.Printf("Error: %s\n", err)
        // }
        // out, err := t2s.Convert(str)
        // if err != nil {
        // 	fmt.Printf("Error: %s\n", err)
        // }
        // fmt.Println(out)
    }
(```)
