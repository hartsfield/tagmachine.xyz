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

// pageInOrder gets a page from the db passed to it. This is the code for
// pagification.
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
	}
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

// submitRoot verifies a users submissions and then adds it to the database.
func submitRoot(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		return
	}
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
			tempFile, err := os.CreateTemp("public/temp", "u-*."+fileExtension)
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
			log.Println("bt:", bodyText)
		}
		if part.FormName() == "parent" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			parent = buf.String()
		}
	}

	if parent != "root" {
		ajaxResponse(w, map[string]string{
			"success": "false",
			"replyID": "",
			"error":   "thread no longer exists",
		})
		return
	}

	if len(bodyText) < 5 || len(bodyText) > 1000 {
		ajaxResponse(w, map[string]string{"success": "false"})
		return
	}

	var data *post = &post{
		Id:        genPostID(10),
		TS:        time.Now(),
		Parent:    parent,
		BodyText:  bodyText,
		MediaType: mediaType,
	}
	data.FTS = data.TS.Format("2006-01-02 03:04:05 pm")
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
	rdb.ZAdd(rdx, "ANON:POSTS:CHRON", redis.Z{Score: float64(time.Now().UnixMilli()), Member: data.Id})
	rdb.ZAdd(rdx, "ANON:POSTS:RANK", redis.Z{Score: 0, Member: data.Id})
	popLast()
	ajaxResponse(w, map[string]string{
		"success":   "true",
		"replyID":   data.Id,
		"timestamp": data.FTS,
	})
	beginCache()
}

// submitReply verifies a users submissions and then adds it to the database.
func submitReply(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		return
	}
	data, err := marshalPostData(r)
	if err != nil {
		log.Println(err)
	}
	parentExists, err := rdb.Exists(rdx, data.Parent).Result()
	if err != nil {
		log.Println(err)
	}
	if parentExists == 0 {
		ajaxResponse(w, map[string]string{
			"success":   "false",
			"replyID":   "",
			"timestamp": data.FTS,
		})
		return
	}
	if len(data.BodyText) < 5 || len(data.BodyText) > 1000 {
		ajaxResponse(w, map[string]string{"success": "false"})
		return
	}
	data.Id = genPostID(10)
	data.TS = time.Now()
	data.FTS = data.TS.Format("2006-01-02 03:04:05 pm")
	data.MediaType = "none"
	rdb.HSet(
		rdx, data.Id,
		"name", data.Author,
		"title", data.Title,
		"bodytext", data.BodyText,
		"id", data.Id,
		"ts", data.TS,
		"fts", data.FTS,
		"parent", data.Parent,
		"childCount", "0",
		"mediaType", "none",
	)
	rdb.ZAdd(rdx, data.Parent+":CHILDREN:CHRON", redis.Z{Score: float64(time.Now().UnixMilli()), Member: data.Id})
	rdb.ZAdd(rdx, data.Parent+":CHILDREN:RANK", redis.Z{Score: 0, Member: data.Id})
	bubbleUp(data)
	ajaxResponse(w, map[string]string{
		"success":   "true",
		"replyID":   data.Id,
		"timestamp": data.FTS,
	})
	beginCache()
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	// Check if the user is logged in. You can't post wothout being logged
	c := r.Context().Value(ctxkey)
	if a, ok := c.(*credentials); !ok || !a.IsLoggedIn {
		ajaxResponse(w, map[string]string{
			"success": "false",
			"replyID": "",
			"error":   "You must be logged in for that",
		})
		return false
	}
	return true
}
