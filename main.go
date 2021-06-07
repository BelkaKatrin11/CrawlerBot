package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {
	// парсинг RSS
	items, err := parser.ParseRSS()
	if err != nil {
		log.Print(err)
		return
	}

	// создаем клиент для подключения к базе
	cl, err := esdb.NewClient()
	if err != nil {
		log.Println(err)
		return
	}

	db := esdb.NewDatabase(cl)
	/////// убрать
	err = db.GetInfo()
	if err != nil {
		log.Println(err)
		return
	}
	///////

	// создаем индкс и маппинг
	err = db.PutIndex()
	if err != nil {
		log.Println(err)
		return
	}

	// функция нахождения дупликатов статей
	// возвращает список не скачанных статей
	notDownloaded, err := db.FindDuplicates(items...)
	if err != nil {
		log.Println(err)
		return
	}

	repoItems := make([]models.RepoItem, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// в цикле скачиваем статьи
	for _, item := range notDownloaded {
		wg.Add(1)

		it := item

		// скачиваем многопоточно - разбиваем на горутины (потоки)
		go func() {
			defer wg.Done()

			// СКАЧИВАЕТ И ПАРСИТ СТАТЬЮ
			article, err := parser.ParseArticle(it.Link)
			if err != nil {
				return
			}

			mu.Lock()
			repoItems = append(repoItems, models.FillRepoItem(it, article))
			mu.Unlock()
		}()


		wg.Wait()
	}

	// печатает новые скачанные ссылки статей
	for _, art := range repoItems {
		fmt.Println("Downloaded one more article: ", art.Link)
	}

	// запихиваем новые статьи в базу эластик
	err = db.FillDatabaseWithItems(repoItems...)
	if err != nil {
		log.Println(err)
		return
	}

	//-----------------

	//функция полнотекстового поиска - ищем по слову, передаваемому в функу
	err = db.SearchByKey("бесы")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.TermAggregationByField("title")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.CardinalityAggregationByField("title")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.CardinalityAggregationByField("pubdate")
	if err != nil {
		log.Println(err)
		return
	}

	err = db.DateHistogramAggregation()
	if err != nil {
		log.Println(err)
		return
	}
}
