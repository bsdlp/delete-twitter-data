package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	TwitterApiKey       string `required:"true" split_words:"true"`
	TwitterApiSecret    string `required:"true" split_words:"true"`
	TwitterAccessToken  string `required:"true" split_words:"true"`
	TwitterAccessSecret string `required:"true" split_words:"true"`
	TwitterDataDir      string `required:"true" split_words:"true"`
}

type LikeItem struct {
	Like struct {
		TweetId     string `json:"tweetId"`
		FullText    string `json:"fullText"`
		ExpandedUrl string `json:"expandedUrl"`
	} `json:"like"`
}

func readLikes(likesFilePath string) ([]LikeItem, error) {
	fh, err := os.Open(likesFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening like.js: %s", err)
	}
	defer fh.Close()

	// skip the javascript stuff
	// window.YTD.like.part0 =
	_, err = fh.Seek(24, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking like.js: %s", err)
	}

	likes := []LikeItem{}
	err = json.NewDecoder(fh).Decode(&likes)
	if err != nil {
		return nil, fmt.Errorf("error decoding like.js: %s", err)
	}
	return likes, nil
}

func deleteLike(client *twitter.Client, item LikeItem) error {
	id, err := strconv.ParseInt(item.Like.TweetId, 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing tweet id '%s': %s", item.Like.TweetId, err)
	}

	_, _, err = client.Favorites.Destroy(&twitter.FavoriteDestroyParams{ID: id})
	if err != nil {
		return fmt.Errorf("error unliking tweet %s: %s", item.Like.TweetId, err)
	}
	return nil
}

func main() {
	flag.Parse()
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	likes, err := readLikes(filepath.Join(cfg.TwitterDataDir, "like.js"))
	if err != nil {
		log.Fatal(err)
	}

	config := oauth1.NewConfig(cfg.TwitterApiKey, cfg.TwitterApiSecret)
	token := oauth1.NewToken(cfg.TwitterAccessToken, cfg.TwitterAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	for _, like := range likes {
		err = deleteLike(client, like)
		if err != nil {
			log.Println(err)
		}
	}
}
