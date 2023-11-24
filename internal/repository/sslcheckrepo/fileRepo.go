package sslcheckrepo

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"

	"github.com/ni3mm4nd/ssl-expiry-checker/internal/domain/sslcheck"
)

var fStore *fileStore
var fStoreConfig *fileStoreConfig

const DEFAULT_DATA_FILE_NAME = "sslchecks.json"

type fileStore struct {
	Mutex sync.Mutex
	Store []sslcheck.SSLCheck `json:"store"`
}

type fileStoreConfig struct {
	DataFileName string
}

func (f *fileStore) ReadAll() ([]sslcheck.SSLCheck, error) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	err := f.ReadFromFile()
	if err != nil {
		return []sslcheck.SSLCheck{}, err
	}

	return f.Store, nil
}

func (f *fileStore) WriteAll(checks []sslcheck.SSLCheck) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	f.Store = checks
	err := f.WriteToFile()
	return err
}

func (f *fileStore) Read(key string) (sslcheck.SSLCheck, error) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	for _, val := range f.Store {
		if val.TargetURL == key {
			return val, nil
		}
	}

	return sslcheck.SSLCheck{}, errors.New("not found")
}
func (f *fileStore) Write(val sslcheck.SSLCheck) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	for i, check := range f.Store {
		if check.TargetURL == val.TargetURL {
			f.Store[i] = val
			return f.WriteToFile()
		}
	}

	f.Store = append(f.Store, val)
	err := f.WriteToFile()
	return err
}

func (f *fileStore) Delete(key string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()

	for i, val := range f.Store {
		if val.TargetURL == key {
			f.Store = append(f.Store[:i], f.Store[i+1:]...)
			return f.WriteToFile()
		}
	}

	return errors.New("not found")
}

func (f *fileStore) ReadFromFile() error {
	log.Println("Reading from file")
	file, err := os.Open(fStoreConfig.DataFileName)
	if err != nil {
		return err
	}

	jsonData, err := io.ReadAll(file)
	log.Printf("Read %d bytes\n", len(jsonData))
	if err != nil {
		log.Println(err)
		return err
	}

	return json.Unmarshal(jsonData, &f.Store)
}

func (f *fileStore) WriteToFile() error {
	var file *os.File

	jsonData, err := json.MarshalIndent(f.Store, "", " ")
	if err != nil {
		return err
	}

	file, err = os.Create(fStoreConfig.DataFileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(jsonData)
	return err
}

func NewFileStore() sslcheck.SSLCheckRepository {
	fStore = &fileStore{
		Mutex: sync.Mutex{},
		Store: []sslcheck.SSLCheck{},
	}

	_, err := os.Stat(DEFAULT_DATA_FILE_NAME)
	if err != nil {
		f, err := os.Create(DEFAULT_DATA_FILE_NAME)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.Write([]byte("[]"))
	}

	fStoreConfig = &fileStoreConfig{
		DataFileName: DEFAULT_DATA_FILE_NAME,
	}

	return fStore
}

func NewFileStoreWithFileName(dbFileName string) sslcheck.SSLCheckRepository {
	fStore = &fileStore{
		Mutex: sync.Mutex{},
		Store: []sslcheck.SSLCheck{},
	}

	_, err := os.Stat(dbFileName)
	if err != nil {
		f, err := os.Create(dbFileName)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		f.Write([]byte("[]"))
	}

	fStoreConfig = &fileStoreConfig{
		DataFileName: dbFileName,
	}

	return fStore
}
