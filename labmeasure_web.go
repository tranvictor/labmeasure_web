package main

import (
	"github.com/tranvictor/labmeasure"
	// "golang.org/x/net/websocket"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var templates = template.Must(template.ParseFiles(
	"tmpl/measure/body.html", "tmpl/measure/body/result.html"))

func display(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl+".html", data)
	if err != nil {
		fmt.Printf("%s", err)
	}
}

func trapError(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}

func handleUploadGET(w http.ResponseWriter) {
	histories, err := LoadHistories()
	if err != nil {
		fmt.Printf("%s", err)
	}
	display(w, "body", StatPresenter{nil, "", histories})
}

func makeTimestamp() string {
	mili := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("%d", mili)
}

func getFileFromPart(part *multipart.Part, stamp string) (string, error) {
	file := "uploads/" + stamp + ".json"
	dst, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, part)
	if err != nil {
		return "", err
	}
	return file, nil
}

func getParam(part *multipart.Part) (string, float32, error) {
	name := part.FormName()
	fmt.Printf("Read name %s", name)
	if name == "submit" {
		return name, float32(0), nil
	}
	buffer := make([]byte, 1024)
	n, err := part.Read(buffer)
	fmt.Printf("Read %s", string(buffer[:n]))
	if err != nil {
		return "", 0.0, err
	}
	value, err := strconv.ParseFloat(string(buffer[:n]), 32)
	if err != nil {
		return "", 0.0, err
	}
	return name, float32(value), nil
}

func getFormParas(stamp string, r *http.Request) ([]string, float32, float32, error) {
	reader, err := r.MultipartReader()
	files := make([]string, 0)
	preThreshold := float32(0.0)
	recThreshold := float32(0.0)
	if err != nil {
		return []string{""}, 0.0, 0.0, err
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		if part.FileName() != "" {
			file, err := getFileFromPart(part, stamp)
			if err != nil {
				return []string{""}, 0.0, 0.0, err
			}
			files = append(files, file)
		} else {
			name, value, err := getParam(part)
			if err != nil {
				return []string{""}, 0.0, 0.0, err
			}
			if name == "precision-threshold" {
				preThreshold = value
			} else if name == "recall-threshold" {
				recThreshold = value
			}
		}
	}

	return files, float32(preThreshold), float32(recThreshold), nil
}

func handleUploadPOST(w http.ResponseWriter, r *http.Request) {
	timestamp := makeTimestamp()
	files, preThreshold, recThreshold, err := getFormParas(timestamp, r)
	if trapError(w, err) {
		return
	}

	for _, file := range files {
		conf := labmeasure.Config{
			"resources/diffbot_v2.json",
			file,
			preThreshold,
			recThreshold,
		}
		stat := labmeasure.Measure(conf)
		history := History{}
		history.RecordFromStat(stat)
		history.ID = timestamp
		history.WriteToJson(fmt.Sprintf("results/%s.json", timestamp))
		b, _ := json.MarshalIndent(stat.IncorrectRecords, "", "  ")
		err := ioutil.WriteFile(
			fmt.Sprintf("incorrects/%s.json", timestamp), b, 0666)

		if err != nil {
			fmt.Printf("%q", err)
		}

		histories, err := LoadHistories()
		if err != nil {
			fmt.Printf("%q", err)
		}

		display(
			w, "body",
			StatPresenter{&stat, timestamp, histories},
		)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//GET displays the upload form.
	case "GET":
		handleUploadGET(w)
	//POST takes the uploaded file and process
	case "POST":
		handleUploadPOST(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func downloadIncorrectHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	filename := paths[len(paths)-1]
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"incorrect_%s.json\"", filename))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	f, err := ioutil.ReadFile(fmt.Sprintf("incorrects/%s.json", filename))
	if trapError(w, err) {
		return
	}
	io.Copy(w, bytes.NewReader(f))
}

func downloadUploadHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	filename := paths[len(paths)-1]
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=\"uploaded_%s.json\"", filename))
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	f, err := ioutil.ReadFile(fmt.Sprintf("uploads/%s.json", filename))
	if trapError(w, err) {
		return
	}
	io.Copy(w, bytes.NewReader(f))
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	filename := paths[len(paths)-1]
	history, err := LoadFromJson(filename + ".json")
	if err != nil {
		fmt.Printf("%s", err)
	}
	display(w, "result", HistoryPresenter{history})
}

func main() {
	http.HandleFunc("/measure/body", uploadHandler)
	http.HandleFunc("/measure/body/incorrects/", downloadIncorrectHandler)
	http.HandleFunc("/measure/body/uploads/", downloadUploadHandler)
	http.HandleFunc("/measure/body/results/", resultHandler)
	// http.HandleFunc("/measure/body/results/", resultHandler)
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.ListenAndServe(":5000", nil)
}
