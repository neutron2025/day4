package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

func main() {
	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		var user User
		// 读取请求体到bytes.Buffer
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		// 打印请求体内容
		fmt.Println(buf.String())

		// 重置buf的读取指针
		buf.Reset()

		// 重置请求体读取器的位置
		r.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))

		// 检查 Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Invalid Content-Type", http.StatusUnsupportedMediaType)
			return
		}

		// 再次将buf的内容复制到r.Body中，以便可以解码JSON
		r.Body = io.NopCloser(buf)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			fmt.Println("Decode error:", err)
			return
		}

		fmt.Fprintf(w, "%s %s is %d years old!", user.Firstname, user.Lastname, user.Age)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		peter := User{
			Firstname: "John",
			Lastname:  "Doe",
			Age:       25,
		}

		json.NewEncoder(w).Encode(peter)
	})

	http.ListenAndServe(":8080", nil)
}

// $ go run json.go

// $ curl -s -XPOST -d'{"firstname":"Elon","lastname":"Musk","age":48}' http://localhost:8080/decode
// Elon Musk is 48 years old!

// $ curl -s http://localhost:8080/encode
// {"firstname":"John","lastname":"Doe","age":25}
