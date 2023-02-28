package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type curiouscatResponseT struct {
	UserData struct {
		Id       int    `json:"id"`
		Username string `json:"username"`
		Answers  int    `json:"answers"`
	} `json:"userData"`
	Posts []struct {
		Type string `json:"type"`
		Post struct {
			Id             int    `json:"id"`
			Timestamp      int    `json:"timestamp"`
			SecondsElapsed int    `json:"seconds_elapsed"`
			Comment        string `json:"comment"`
			Reply          string `json:"reply"`
		} `json:"post"`
	}
}

var (
	url    string
	client *http.Client
)

func main() {

	username := flag.String("username", "", "Nome do usuário no CuriousCat")
	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	fmt.Println("Usuário selecionado:", *username)
	fmt.Printf("- -\n")

	url = fmt.Sprintf("https://curiouscat.live/api/v2.1/profile?username=%s", *username)
	client = http.DefaultClient
	client.Timeout = time.Second * 5

	readPosts := 0
	done := false
	maxTimestamp := 0
	var err error

	for {

		if done {
			break
		}

		posts := getPosts(maxTimestamp)
		var curiouscatResponse *curiouscatResponseT
		err = json.Unmarshal(posts, &curiouscatResponse)

		if err != nil {
			fmt.Println("Erro fazendo parse da resposta da API:", err)
			done = true
			continue
		}

		if len(curiouscatResponse.Posts) == 0 {
			if readPosts == 0 {
				fmt.Printf("Puxa, o usuário %s não possui nenhuma pergunta publicada!\n", *username)
			} else {
				fmt.Printf("Nenhum post a mais para carregar!\n")
			}
			done = true
			continue
		}

		// A partir da primeira chamada para API, o max_timestamp da próxima requisição
		// precisa ser baseado no primeiro timestamp lido
		maxTimestamp = curiouscatResponse.Posts[0].Post.Timestamp - 1

		for _, post := range curiouscatResponse.Posts {
			readPosts++
			fmt.Printf("Data: %s\n", time.Unix(int64(post.Post.Timestamp), 0))
			fmt.Printf("Há: %s\n", time.Since(time.Now().Add(time.Second*time.Duration(post.Post.SecondsElapsed)*-1)).Truncate(time.Hour))
			fmt.Printf("Pergunta: %s\n", post.Post.Comment)
			fmt.Printf("Resposta: %s\n", post.Post.Reply)
			fmt.Printf("- -\n")
		}

	}

	fmt.Println("Encerrei.")
	os.Exit(0)

}

func getPosts(maxTimestamp int) []byte {

	currentUrl := url
	if maxTimestamp > 0 {
		currentUrl = fmt.Sprintf("%s&max_timestamp=%d", url, maxTimestamp)
	}

	response, err := client.Get(currentUrl)

	if err != nil {
		fmt.Println("Erro ao requisitar o serviço:", err)
		os.Exit(0)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Erro lendo a resposta da API:", err)
		os.Exit(0)
	}
	return body
}
