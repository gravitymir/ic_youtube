package main

//$$(
//'#page-manager > ytd-browse > ytd-two-column-browse-results-renderer #primary #contents ytd-rich-grid-row #contents ytd-rich-item-renderer #content ytd-rich-grid-media #dismissible #details #meta h3 a[href]'
//).map((e) => e.href.replace('https://www.youtube.com/watch?v=', '').replace(/\=\d{1,5}s/, '').replace(/\&t/, '')).reverse();

import (
	"bytes"
	"fmt"
	"github.com/kkdai/youtube/v2"
	"github.com/robfig/cron/v3"
	"golang.org/x/exp/slices"
	"log"
	"os"
	"os/exec"
	"strings"
	myYouTube "workspace/youtube"
)

func scanDir(dir string) (list []string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range entries {
		if folder.Name() == ".DS_Store" {
			continue
		}
		numId := strings.Split(folder.Name(), " ")
		list = append(list, numId[1])
	}
	return list
}

func spotLightFoldersWithoutVideo() {
	fmt.Println("spotLightFoldersWithoutVideo()")
	var yt *myYouTube.YouTube
	dir := "../media/"
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, folder := range entries {
		if folder.Name() == ".DS_Store" {
			continue
		}
		folderContent, err := os.ReadDir(dir + folder.Name() + "/")
		if err != nil {
			fmt.Println(err)
			continue
		}

		contentNameSlice := make([]string, 7)
		for _, content := range folderContent {
			contentNameSlice = append(contentNameSlice, content.Name())
		}
		if !slices.Contains(contentNameSlice, "audio_en.mp3") {
			if slices.Contains([]string{
				"000768 MX9mS6AwRI4",
				"000766 6_lw4iUPZAI",
				"000079 O-LYqHnsNcY",
				"000619 Sg0lJNgRPrQ",
				"000300 g9u9hPAjoC4",
				"000609 Lk7TGOEtWiE",
				"000054 Sl6ac-ChpLw",
				"000063 KuAFH9HgjPE",
				"000062 1bN-Soliu6k",
				"000436 2eb7w5pgXFU",
				"000430 4piOfAvgIwo",
				"000403 N_3V6_oTB8o",
				"000390 HWOY3W-7zIc",
				"000233 0BqBJxiR__8",
				"000372 reyFMjqtLos",
				"000363 j3qSoYAfjvw",
				"000358 BYOsCNKL4Cg",
				"000354 eq6BccWQ_Rk",
				"000367 wCQuXf_wQCk",
				"000277 SW-yUT2kZmM",
				"000239 HcvdBvfGt4s",
				"000237 mcAJkTYvSXA",
				"000236 JUD-0fJooSo",
				"000235 mTiepS11sws",
				"000080 YbgbMxnfX7w", //music
				"000167 jokUwt9D4Tg", //corals
				"000089 FybmR50cWO8",
			}, folder.Name()) {
				continue
			}

			fmt.Printf("%s don't have audio_en.mp3\n%s%s\n", folder.Name(),
				"https://www.youtube.com/watch?v=", strings.Split(folder.Name(), " ")[1])
		}

		if slices.Contains(contentNameSlice, "video.mp4") {
			if !slices.Contains(contentNameSlice, "audio.mp4") {
				ffmpegGetMP3FromVideo(folder.Name())
			}
		} else {
			fmt.Printf("try download %s", folder.Name())
			indexID := strings.Split(folder.Name(), " ")
			err := downloadVideo22(indexID[1], folder.Name())
			if err != nil {
				fmt.Println(err)
			}
			err = ffmpegGetMP3FromVideo(folder.Name())
			if err != nil {
				fmt.Println(err)
			}

			if slices.Contains(contentNameSlice, "video.mp4") {
				err = ffmpegGetLowQualityVideo(folder.Name())
				fmt.Println(yt.VideoAndAudio[0].URL)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println(yt.VideoAndAudio[0].URL)
			}

		}
		if slices.Contains(contentNameSlice, "details.json") &&
			slices.Contains(contentNameSlice, "subtitles.json") &&
			slices.Contains(contentNameSlice, "thumbnail.jpg") {
			continue
		}
		if yt, err = myYouTube.Init(folder.Name()); err != nil {
			fmt.Println(err)
			continue
		}
		fullPath := dir + folder.Name()
		if !slices.Contains(contentNameSlice, "details.json") {
			yt.SaveDetailsPretty(fullPath + "/details.json")
		}
		if !slices.Contains(contentNameSlice, "subtitles.json") {
			if len(yt.SubtitlesTracks) > 0 {
				err = yt.SaveSubtitlesJsonPretty(yt.SubtitlesTracks[0],
					fullPath+"/subtitles.json")
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if !slices.Contains(contentNameSlice, "thumbnail.jpg") {
			err = yt.SaveThumbnailJPG(fullPath + "/thumbnail.jpg")
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func checkNews() {
	fmt.Println("checkNews()")

	client := youtube.Client{}
	playList, err := client.GetPlaylist("PLcihjUVySO7oupvXPjcvLQZoTX_QnFE52")
	if err != nil {
		return
	}
	var newsSlice []string

	mediaList := scanDir("../media/")
	for _, v := range playList.Videos {
		if !slices.Contains(mediaList, v.ID) {
			newsSlice = append(newsSlice, v.ID)
		}
	}
	if len(newsSlice) == 0 {
		return
	}

	fmt.Printf("%v, %d\n", newsSlice, len(mediaList))
	nextNumFolder := len(mediaList) + 1
	countRepeatZero := 6 - len(fmt.Sprintf("%d", nextNumFolder))
	//	str := fmt.Sprintf("%s %s", fullNumber, id)
	for k, videoID := range newsSlice {
		fullPath := fmt.Sprintf("%s%s%d %s",
			"../media/",
			strings.Repeat("0", countRepeatZero),
			nextNumFolder+k,
			videoID)

		if err := os.Mkdir(fullPath, os.ModePerm); err != nil {
			fmt.Println(err)
			continue
		}
	}
	spotLightFoldersWithoutVideo()
}
func downloadVideo(ID, folderName, fileName, qualityOrTAG string) error {
	cmd := exec.Command("youtubedr", "download",
		"-q", qualityOrTAG,
		"-o", fmt.Sprintf("../media/%s/%s", folderName, fileName),
		ID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Printf("%s 720p video.mp4 download is done\n", folderName)
	return nil
}

// quality 720
func downloadVideo22(ID, folderName string) error {
	fmt.Println("downloadVideo22", ID, folderName)
	err := downloadVideo(ID, folderName, "video.mp4", "22")
	return err
}

//ffmpeg -i "../000775 blvHjq15zsY/video.mp4" -i "../000775 blvHjq15zsY/audio_ru.mp3"\
//-filter_complex "[0:a]volume=0.1[a1];[a:1]volume=1[a2];\
//[a1][a2]amerge=inputs=2[out]" -map 0:v -map "[out]" -c:v copy -c:a aac -b:a "../000775 blvHjq15zsY/3.mp4"
// ffmpeg -i "../000775 blvHjq15zsY/video.mp4" -filter:a "volume=0.1" "../000775 blvHjq15zsY/video222.mp4"
// ffmpeg -i VOCALS -i MUSIC -filter_complex amix=inputs=2:duration=longest:dropout_transition=0:weights="1 0.25":normalize=0 OUTPUT
// ffmpeg -i "../000775 blvHjq15zsY/video.mp4" -i "../000775 blvHjq15zsY/audio_ru.mp3" -c:v copy -c:a copy "../000775 blvHjq15zsY/output.mp4"

// quality 360
func downloadVideo18(ID, folderName string) error {
	err := downloadVideo(ID, folderName, "video_360.mp4", "18")
	return err
}
func ffmpegGetLowQualityVideo(folderName string) error {
	from := fmt.Sprintf("../media/%s/video.mp4", folderName)
	to := fmt.Sprintf("../media/%s/video_240.mp4", folderName)
	auxSlice := []string{
		"-i", from,
		"-threads", "0",
		"-preset", "ultrafast",
		"-s", "320x240",
		"-c:v", "libx264", to,
	}
	fmt.Printf("%s %s", "ffmpeg", strings.Join(auxSlice, " "))
	cmd := exec.Command("ffmpeg", auxSlice...)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: creating video_240.mp4 in %s is ready %s", folderName, err)
	} else {
		fmt.Printf("Creating video_240.mp4 in %s is ready.", folderName)
	}
	return err
}

func ffmpegMergeTranslate(folderName, VideoFileName, AudioFileName, outFileName string) error {
	//ffmpeg -i "../000775 blvHjq15zsY/video.mp4" -i "../000775 blvHjq15zsY/audio_ru.mp3"\
	//-filter_complex "[0:a]volume=0.1[a1];[1:a]volume=1[a2];\
	//[a1][a2]amerge=inputs=2,pan=stereo|FL<c0+c1|FR<c2+c3[out]" -map 0:v -map "[out]" -c:v copy -c:a aac -b:a 192k "../000775 blvHjq15zsY/3.mp4"
	if VideoFileName == "" {
		VideoFileName = fmt.Sprintf("../media/%s/video.mp4", folderName)
	}
	if AudioFileName == "" {
		AudioFileName = fmt.Sprintf("../media/%s/video_240.mp4", folderName)
	}

	auxSlice := []string{
		"-i", VideoFileName,
		"-i", AudioFileName,
		"-filter_complex",
		"[0:a]volume=0.17[a1];[1:a]volume=1[a2];[a1][a2]amerge=inputs=2[out]",
		"-map", "0:v", "-map", "[out]", "-c:v", "copy", "-c:a", "aac", outFileName,
	}
	fmt.Printf("%s %s", "ffmpeg", strings.Join(auxSlice, " "))
	cmd := exec.Command("ffmpeg", auxSlice...)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: ffmpegMergeTranslate in %s is ready %s", folderName, err)
	} else {
		fmt.Printf("ffmpegMergeTranslate in %s is ready.", folderName)
	}
	return err
}
func ffmpegGetMP3FromVideo(folderName string) error {
	//path := "/Users/gravitymir/Documents/golang/media/"
	//from := fmt.Sprintf("%s%s/video.mp4", path, folder)

	video := fmt.Sprintf("../media/%s/video.mp4", folderName)
	mp3 := fmt.Sprintf("../media/%s/audio.mp3", folderName)
	cmd := exec.Command("ffmpeg", "-i", video, mp3)

	err := cmd.Run()

	if err != nil {
		return err
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	//fmt.Printf("%q\n", out.String())

	fmt.Printf("<<<-----------\naudio.mp3 from video.mp4 ready\nffmpeg -i %s %s\n----------->>>\n",
		video, mp3)
	return err
}

func getInfo(ID string) error {
	cmd := exec.Command("youtubedr", "info", ID)

	err := cmd.Run()

	if err != nil {
		return err
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	fmt.Printf("%q\n", out.String())
	return nil
}

func downloadPlaylist(ID string) {
	cmd := exec.Command("youtubedr", "list", "PLcihjUVySO7oupvXPjcvLQZoTX_QnFE52")

	//cmd.Stdin = strings.NewReader("")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%q\n", out.String())
}

func run() error {
	c := cron.New()
	//@every 1h30m10s
	checkNews()
	//if entryID, err := c.AddFunc("@every 1h", checkNews); err != nil {
	//	fmt.Println("Cron error entryID: ", entryID)
	//	return err
	//}
	spotLightFoldersWithoutVideo()
	//if entryID, err := c.AddFunc("@every 1h", spotLightFoldersWithoutVideo); err != nil {
	//	fmt.Println("Cron error entryID: ", entryID)
	//	return err
	//}
	c.Start()
	fmt.Println(c.Entries())
	select {}
}

func main() {
	///opt/homebrew/Cellar/ffmpeg/
	//brew uninstall ffmpeg
	//brew update
	//brew upgrade
	//brew cleanup
	//brew install ffmpeg --force
	//brew link ffmpeg
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

//$$(
//'#page-manager > ytd-browse > ytd-two-column-browse-results-renderer #primary #contents ytd-rich-grid-row #contents ytd-rich-item-renderer #content ytd-rich-grid-media #dismissible #details #meta h3 a[href]'
//).map((e) => e.href.replace('https://www.youtube.com/watch?v=', '').replace(/\=\d{1,5}s/, '').replace(/\&t/, '')).reverse();
