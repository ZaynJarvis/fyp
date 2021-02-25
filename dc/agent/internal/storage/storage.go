package storage

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ObjStoreAddr   string
	DocStorageAddr string
	TextIndexAddr  string
}

type Obj interface {
	Image(string, []byte) error
	Close()
}
type Doc interface {
	Doc(string, interface{}) error
	Close()
}
type Index interface {
	TextIndex(string, string) error
	Close()
}
type Storage interface {
	Obj
	Doc
	Index
	UpdateConfig(cfg Config)
}

type ComposedStorage struct {
	mu    sync.RWMutex
	obj   Obj
	data  Doc
	index Index
	cfg   Config
}

func New(cfg Config) Storage {
	c := &ComposedStorage{cfg: Config{}}
	c.UpdateConfig(cfg)
	return c
}

func (c *ComposedStorage) Image(id string, data []byte) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.obj.Image(id, data)
}

func (c *ComposedStorage) Close() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.UpdateConfig(Config{})
}

func (c *ComposedStorage) Doc(id string, result interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data.Doc(id, result)
}

func (c *ComposedStorage) TextIndex(id string, text string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.index.TextIndex(id, text)
}

func (c *ComposedStorage) UpdateConfig(cfg Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.updateObj(cfg.ObjStoreAddr)
	c.updateData(cfg.DocStorageAddr)
	c.updateIndex(cfg.TextIndexAddr)
}

func (c *ComposedStorage) updateObj(addr string) {
	if addr == c.cfg.ObjStoreAddr {
		return
	}
	if c.obj != nil {
		c.obj.Close()
	}
	if addr != "" {
		m, err := newMinIO(addr)
		if err != nil {
			logrus.Error(err)
			c.cfg.ObjStoreAddr = ""
			return
		}
		c.obj = m
		logrus.Info("object storage connected")
	}
	c.cfg.ObjStoreAddr = addr
}

func (c *ComposedStorage) updateData(addr string) {
	if addr == c.cfg.ObjStoreAddr {
		return
	}
	if c.data != nil {
		c.data.Close()
	}
	if addr != "" {
		m, err := newMongoDB(addr)
		if err != nil {
			logrus.Error(err)
			c.cfg.DocStorageAddr = ""
			return
		}
		c.data = m
		logrus.Info("document storage connected")
	}
	c.cfg.DocStorageAddr = addr
}

func (c *ComposedStorage) updateIndex(addr string) {
	if addr == c.cfg.TextIndexAddr {
		return
	}
	if c.index != nil {
		c.index.Close()
	}
	logrus.Error("index not supported")
}
