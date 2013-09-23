package mycelium

import (
	"log"
	"time"
)

type Worker struct {
	roboFilter *RobotFilter
}

func NewWorker() *Worker {
	return &Worker{
		roboFilter: NewRobotFilter(),
	}
}

func (self *Worker) GetPage(url string) (*Page, error) {
	resp, err := self.roboFilter.PoliteGet(url)

	if err != nil {
		return nil, err
	}

	return NewPageFromResponse(resp), nil
}

func (self *Worker) GetPages(urls []string, timeout time.Duration) []*Page {
	pageChan := make(chan *Page)

	for _, url := range urls {
		go func() {
			page, err := self.GetPage(url)

			if err != nil {
				log.Println(err)
				return
			}

			pageChan <- page
		}()
	}

	pages := make([]*Page, 0)
	timeoutChan := time.After(timeout)

	for {
		select {
		case page := <-pageChan:
			pages = append(pages, page)
		case <-timeoutChan:
			return pages
		default:
			if len(pages) == len(urls) {
				return pages
			}
		}
	}

	return pages
}