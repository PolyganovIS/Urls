package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура для представления контента сайта
type SiteContent struct {
	Title   string
	Date    time.Time
	Content string
}

// Функция для имитации долгой работы, загрузки контента сайта
func DownloadSiteContent(ctx context.Context, url string) SiteContent {
	// Генерация случайного времени ожидания от 5 до 10 секунд
	rand.Seed(time.Now().UnixNano())
	sleepTime := time.Duration(rand.Intn(6)+5) * time.Second
	// Ожидание случайного времени
	select {
	case <-time.After(sleepTime):
	case <-ctx.Done():
		return SiteContent{}
	}
	// Возвращаем имитацию контента сайта
	return SiteContent{
		Title:   "Заголовок " + url,
		Date:    time.Now(),
		Content: "Содержание сайта " + url,
	}
}

// Функция для параллельного скачиваниn контента сайтов
func ParallelDownload(ctx context.Context, urls <-chan string, numWorkers int) map[string]SiteContent {
	mu := sync.Mutex{}
	m := make(map[string]SiteContent)
	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case url, ok := <-urls:
					if !ok {
						return
					}
					content := DownloadSiteContent(ctx, url)
					mu.Lock()
					m[url] = content
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	wg.Wait()
	return m
}

func main() {
	urls := make(chan string)
	// Запускаем функцию для параллельного скачивания контента
	go func() {
		urls <- "http://example.com"
		urls <- "http://example.org"
		urls <- "http://examp1e.net"
		close(urls)
	}()

	// Запускаем параллельное скачивание с максимальным количеством воркеров
	ctx := context.Background()
	result := ParallelDownload(ctx, urls, 3)
	// Выводим результат
	for url, content := range result {
		fmt.Printf("Caйт: %s\nЗаголовок: %s\nДата: %s\nСодержание: %s\n\n", url, content.Title, content.Date.Format(time.RFC3339), content.Content)
	}
}
