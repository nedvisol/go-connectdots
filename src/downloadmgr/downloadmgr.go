package downloadmgr

import (
	"context"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/nedvisol/go-connectdots/cacheditem"
)

const LIMIT_PER_HOST = 2

type DownloadManager struct {
	options       *DownloadManagerOptions
	client        *http.Client
	downloadQueue *downloadQueue
}

type DownloadManagerOptions struct {
	cacheDir       *string
	cachedItemRepo cacheditem.CachedItemRepository
}

type downloadQueue struct {
	limitPerHost int
	queues       map[string]chan struct{} // Semaphore per host
	mutex        sync.Mutex
	wg           sync.WaitGroup
}

type DownloadCallback func(data []byte)

func (dm *DownloadManager) getHashKey(request *http.Request) string {
	val := fmt.Sprintf("%s %s", request.Method, request.URL)
	// Compute the SHA-512 hash
	hash := sha512.New()
	hash.Write([]byte(val))

	// Get the final hashed output
	hashBytes := hash.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashBytes)
}

func (dm *DownloadManager) processDownload(ctx context.Context, request *http.Request, callback DownloadCallback) {
	//check for cache
	requestHashKey := dm.getHashKey(request)
	cachedItem, err := dm.options.cachedItemRepo.FindByKey(ctx, requestHashKey)
	if err != nil {
		panic("unable to get cached item from cachedItemRepo")
	}

	if cachedItem != nil {
		cacheExpiredAt := time.Unix(cachedItem.ExpiresAtSec, 0)
		filePath := fmt.Sprintf("%s/%s", *dm.options.cacheDir, requestHashKey)
		if time.Now().Before(cacheExpiredAt) {
			//return cached item in a different thread
			data, err := os.ReadFile(filePath)
			if err != nil {
				panic(err)
			}
			go func() {
				callback(data)
			}()
			return
		} else {
			//cache expired, delete from repo and delete file
			dm.options.cachedItemRepo.DeleteByKey(ctx, requestHashKey)
			err := os.Remove(filePath)
			if err != nil {
				panic(err)
			}
		}
	}
	// download from host and return content

	resp, err := dm.client.Do(request)
	if err != nil {
		panic("http request error")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("http read from body error")

	}
	go func() {
		callback(body)
	}()
}

func (dm *DownloadManager) Download(ctx context.Context, request *http.Request, callback DownloadCallback) {
	host := request.URL.Host
	dq := dm.downloadQueue

	// Get or create a semaphore for this host

	dq.mutex.Lock()
	if _, exists := dq.queues[host]; !exists {
		dq.queues[host] = make(chan struct{}, dq.limitPerHost) // Limit to 2 concurrent connections per host
	}
	sem := dq.queues[host]
	dq.mutex.Unlock()

	dq.wg.Add(1)
	go dm.processRequest(ctx, request, callback, sem)

}

func (dm *DownloadManager) processRequest(ctx context.Context, request *http.Request,
	callback DownloadCallback, sem chan struct{}) {
	dq := dm.downloadQueue

	defer dq.wg.Done()

	// Acquire a slot in the semaphore
	sem <- struct{}{}
	defer func() { <-sem }() // Release slot after the download

	// Perform the download
	dm.processDownload(ctx, request, callback)
}

// Wait waits for all downloads to complete
func (dm *DownloadManager) Wait() {
	dm.downloadQueue.wg.Wait()
}

func NewDownloadManager(opts *DownloadManagerOptions) *DownloadManager {
	downloadQueue := &downloadQueue{
		limitPerHost: LIMIT_PER_HOST,
		queues:       make(map[string]chan struct{}),
	}

	return &DownloadManager{
		options:       opts,
		client:        &http.Client{},
		downloadQueue: downloadQueue,
	}
}
