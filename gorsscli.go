package main

import (
    "fmt"
    "bufio"
    "log"
    "os"
    "github.com/mmcdole/gofeed"
)

func GetFeeds() ([]string){
    file,err := os.Open("urls")
    if err != nil {
        log.Fatalf("unable to read file: ",err)
    }
    defer file.Close()

    var lines []string
    scanner:= bufio.NewScanner(file)
    for scanner.Scan(){
        lines=append(lines, scanner.Text())
    }
    if err := scanner.Err(); err != nil{
        log.Fatalf("Unable to get lines from file: ",err)
    }
    return lines
}

func ReadFeed(url string) *gofeed.Feed {
    fp := gofeed.NewParser()
    feed, _ := fp.ParseURL(url)
    return feed
}

func main(){
    rss_urls := GetFeeds()

    for _,url := range rss_urls {
        feed := ReadFeed(url)

        for _,item := range feed.Items {
            fmt.Println(item.Title)
        }
    }
}
