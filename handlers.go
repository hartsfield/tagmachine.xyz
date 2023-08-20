package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// home is displays the main page
func home(w http.ResponseWriter, r *http.Request) {
	var v viewData
	v.Order = "chron"
	if len(postDBChron) < 20 {
		v.Stream = postDBChron[:]
	} else {
		v.Stream = postDBChron[:20]
	}
	exeTmpl(w, r, &v, "main.tmpl")
}
func pageInOrder(db []*post, r *http.Request, count int, v *viewData) map[string]string {
	var bb bytes.Buffer
	var nextCount string
	params, err := url.ParseQuery(strings.Split(r.RequestURI, "?")[1])
	if err != nil {
		log.Println(err)
	}
	if params["count"] == nil {
		params["count"] = append(params["count"], "0")
	}
	if params["count"][0] != "None" {
		count, err := strconv.Atoi(params["count"][0])
		if err != nil {
			log.Println(err)
		}
		if len(db) <= count+20 {
			v.Stream = db[count-(count-len(db)):]
			nextCount = "None"
		} else {
			v.Stream = db[count : count+20]
			nextCount = strconv.Itoa(count + 20)
		}
		err = templates.ExecuteTemplate(&bb, "stream.tmpl", v)
		if err != nil {
			log.Println(err)
		}
	}
	return map[string]string{
		"success":  "true",
		"template": bb.String(),
		"count":    nextCount,
	}
}

// getByChron returns 20 posts at a time in chronological order
func getByChron(w http.ResponseWriter, r *http.Request) {
	var count int = 20
	var v viewData
	v.Order = "chron"
	if len(strings.Split(r.RequestURI, "?")) > 1 {
		ajaxRes := pageInOrder(postDBChron, r, count, &v)
		ajaxResponse(w, ajaxRes)
	} else {
		v.Stream = postDBChron[count+(len(postDBChron)-count):]
		exeTmpl(w, r, &v, "main.tmpl")
	}
}

// getByRanked returns 20 posts at a time in ranked order.
func getByRanked(w http.ResponseWriter, r *http.Request) {
	var count int = 20
	var v viewData
	v.Order = "ranked"
	if len(strings.Split(r.RequestURI, "?")) > 1 {
		ajaxRes := pageInOrder(postDBRank, r, count, &v)
		ajaxResponse(w, ajaxRes)
	} else {
		if len(postDBRank) < count {
			v.Stream = postDBRank[:]
		} else {
			v.Stream = postDBRank[:count]
		}
		exeTmpl(w, r, &v, "main.tmpl")
		// v.Stream = postDBRank[count+(len(postDBRank)-count):]
		// exeTmpl(w, r, &v, "main.tmpl")
	}
	// var v viewData
	// v.Order = "ranked"
	// var count int = 20
	// if len(strings.Split(r.RequestURI, "?")) > 1 {
	// 	params, err := url.ParseQuery(strings.Split(r.RequestURI, "?")[1])
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	if params["count"] == nil {
	// 		params["count"] = append(params["count"], "0")
	// 	}
	// 	if params["count"][0] != "None" {
	// 		count, err = strconv.Atoi(params["count"][0])
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		var nextCount string
	// 		if len(postDBRank) < count {
	// 			v.Stream = postDBRank[count+(len(postDBRank)-count):]
	// 			nextCount = "None"
	// 		} else {
	// 			v.Stream = postDBRank[count+1 : count+20]
	// 			nextCount = strconv.Itoa(count + 20)
	// 		}
	// 		var bb bytes.Buffer
	// 		err = templates.ExecuteTemplate(&bb, "stream.tmpl", v)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		ajaxResponse(w, map[string]string{
	// 			"success":  "true",
	// 			"template": bb.String(),
	// 			"count":    nextCount,
	// 		})
	// 	}
	// } else {
	// 	if len(postDBRank) < count {
	// 		v.Stream = postDBRank[:]
	// 	} else {
	// 		v.Stream = postDBRank[:count]
	// 	}
	// 	exeTmpl(w, r, &v, "main.tmpl")
	// }
}

// viewPost returns a single post, with replies.
func viewPost(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.RequestURI, "/")
	var p post
	rdb.HGetAll(rdx, parts[len(parts)-1]).Scan(&p)
	if len(p.Id) == 11 {
		getAllChidren(&p, "RANK")
	} else {
		p.BodyText = "This post was automatically deleted."
	}
	var v viewData
	v.Stream = nil
	v.Stream = append(v.Stream, &p)
	v.ViewType = "post"
	exeTmpl(w, r, &v, "post.tmpl")
}

func syncApk(apkFile *os.File) {

}

// handleForm verifies a users submissions and then adds it to the database.
func handleForm(w http.ResponseWriter, r *http.Request) {
	mr, err := r.MultipartReader()
	if err != nil {
		log.Println(err)
	}
	var bodyText string
	var parent string
	var tempFileName string
	var mediaType string
	for {
		part, err_part := mr.NextPart()
		if err_part == io.EOF {
			break
		}
		if part.FormName() == "myFile" {
			fileBytes, err := io.ReadAll(io.LimitReader(part, 10<<20))
			if err != nil {
				log.Println(err)
			}
			mt := http.DetectContentType(fileBytes)
			if mt != "image/jpeg" && mt != "image/png" && mt != "video/mp4" && mt != "video/webm" && mt != "image/gif" {
				ajaxResponse(w, map[string]string{
					"success": "false",
					"replyID": "",
					"error":   "png - jpg - gif - webm - mp4 only",
				})
				return
			}
			mediaType = strings.Split(mt, "/")[0]
			var fileExtension string
			switch mt {
			case "image/png":
				fileExtension = "png"
			case "image/jpeg":
				fileExtension = "jpg"
			case "image/gif":
				fileExtension = "gif"
			case "video/mp4":
				fileExtension = "mp4"
			case "video/webm":
				fileExtension = "webm"

			}
			tempFile, err := os.CreateTemp("public/temp-images", "u-*."+fileExtension)
			if err != nil {
				log.Println(err)
			}
			defer tempFile.Close()
			tempFileName = tempFile.Name()

			tempFile.Write(fileBytes)
		}
		if part.FormName() == "myText" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			bodyText = buf.String()
		}
		if part.FormName() == "parent" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			parent = buf.String()
		}
	}

	parentExists, err := rdb.Exists(rdx, parent).Result()
	if err != nil {
		log.Println(err)
	}

	if parentExists == 0 && parent != "root" {
		ajaxResponse(w, map[string]string{
			"success": "false",
			"replyID": "",
		})
		return
	}

	if len(bodyText) < 5 || len(bodyText) > 1000 {
		ajaxResponse(w, map[string]string{"success": "false"})
		return
	}

	var data post
	data.Id = genPostID(10)
	data.TS = time.Now()
	data.FTS = data.TS.Format("2006-01-02 03:04:05 pm")
	data.Parent = parent
	data.BodyText = bodyText
	data.MediaType = mediaType
	rdb.HSet(
		rdx, data.Id,
		"bodytext", bodyText,
		"id", data.Id,
		"ts", data.TS,
		"fts", data.FTS,
		"parent", parent,
		"childCount", "0",
		"media", tempFileName,
		"mediaType", mediaType,
	)
	if data.Parent != "root" {
		rdb.ZAdd(rdx, parent+":CHILDREN:CHRON", redis.Z{Score: float64(time.Now().UnixMilli()), Member: data.Id})
		rdb.ZAdd(rdx, parent+":CHILDREN:RANK", redis.Z{Score: 0, Member: data.Id})
		bubbleUp(&data)
	} else {
		rdb.ZAdd(rdx, "ANON:POSTS:CHRON", redis.Z{Score: float64(time.Now().UnixMilli()), Member: data.Id})
		rdb.ZAdd(rdx, "ANON:POSTS:RANK", redis.Z{Score: 0, Member: data.Id})
		popLast()
	}
	ajaxResponse(w, map[string]string{
		"success":   "true",
		"replyID":   data.Id,
		"timestamp": data.FTS,
	})
	beginCache()
}
