package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Parm struct {
	Url        string
	Prefix     string
	MacAddress string
}

func ffmpeg(path string, prefix string, url string) {
	if err := ensureDir(path + "/" + prefix); err != nil {
		fmt.Println("Directory creation failed with error: " + err.Error())
		os.Exit(1)
	}
	ffmpegCommand := "ffmpeg -fflags nobuffer -rtsp_transport tcp -i \"" + url + "\"  -vsync 0 -copyts  -vcodec copy -movflags frag_keyframe+empty_moov -an -hls_flags delete_segments+append_list -f segment -segment_list_flags live -segment_time 4 -segment_list_size 3 -segment_format mpegts -segment_list " + path + "/" + prefix + "/index.m3u8 -segment_list_type m3u8 -segment_list_entry_prefix /stream/ " + path + "/" + prefix + "/%d.ts"

	Shellout(ffmpegCommand)

	// fmt.Println(ffmpegCommand)

}

func PostFileBodyPath(ip string, path string, prefix string, macAdd string) {

	var count int = 0

	for {
		if CountFileInFolder(path+"/"+prefix) > 3 {
			tsFileCount := strconv.Itoa(count)
			fr := path + "/" + prefix + "/" + tsFileCount + ".ts"

			if Exists(fr) {
				now := time.Now()
				formatted2 := now.Format("2006-01-02T15:04:05")

				url := "http://" + ip + urlPart + formatted2 + "/4"
				// fmt.Println("URL:>", url)

				bodyText := "{\"filePath\":" + "\"" + path + "/" + prefix + "/" + tsFileCount + ".ts" + "\"}"

				var jsonStr = []byte(bodyText)
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
				req.Header.Set("Authorization", macAdd)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Accept-Encoding", "gzip, deflate, br")
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
				fmt.Println("---------------------------------------")
				fmt.Println("response Status:", resp.Status, fr)
				fmt.Println("response Headers:", resp.Header)
				body, _ := io.ReadAll(resp.Body)
				fmt.Println("response Body:", string(body))
				fmt.Println("---------------------------------------")

				os.Remove(fr)
				count++
			} else {
				fmt.Printf("------------>>>>>>>>> file not found %s \n", fr)
				time.Sleep(time.Second * 1)
			}

		} else {
			time.Sleep(time.Second * 3)

		}

	}

}

var serverIp string = "192.168.1.1"
var mainpath string = "/path/video"
var urlPart string = ""

func main() {

	os.RemoveAll(mainpath)
	time.Sleep(time.Second * 1)
	ensureDir(mainpath)
	time.Sleep(time.Second * 1)

	parm := []Parm{
		{Prefix: "a", Url: "rtsp://admin:admin@192.168.1.1:554/p1", MacAddress: "12341A1FF2BF"},
		{Prefix: "c", Url: "rtsp://admin:admin@192.168.1.2:554/p2", MacAddress: "12341A1FF299"},
	}

	sliceLength := len(parm)

	wg.Add(sliceLength)
	fmt.Printf("Running %d Camera...\n", len(parm))

	for i := 0; i < sliceLength; i++ {
		go ffmpeg(mainpath, parm[i].Prefix, parm[i].Url)
		go PostFileBodyPath(serverIp, mainpath, parm[i].Prefix, parm[i].MacAddress)
	}
	wg.Wait()
	fmt.Println("Finished for loop")

}
