package main

import (
	"net/http"
	"io/ioutil"
	"html/template"
	"os"
	"io"
	"time"
	"fmt"
	"log"
	"github.com/julienschmidt/httprouter"
)

func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles(VIDEO_DIR +"upload.html")
	fmt.Println("testPageHandler")
	t.Execute(w, nil)
}

func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid 

	log.Println(vl);
	video, err := os.Open(vl)
	if err != nil {
		log.Printf("Open Error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Interna Server Error")
		return 
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()
}

func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	log.Println("uploadHandler")
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return 
	}

	file, _, err := r.FormFile("file") //<form>
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return 
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "File")
		return 
	}

	fn := p.ByName("vid-id")
	log.Println("uploadHandler ", fn)
	err = ioutil.WriteFile(VIDEO_DIR + fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Inter Error")
		return 
	}

	w.WriteHeader(http.StatusCreated) //
	io.WriteString(w, "Uploaded successful")
}