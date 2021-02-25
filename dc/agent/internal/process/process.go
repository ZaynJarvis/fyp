package process

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"sync"

	"github.com/zaynjarvis/fyp/dc/agent/internal/storage"

	"github.com/sirupsen/logrus"
	"github.com/zaynjarvis/fyp/dc/agent/internal/transport"
	"github.com/zaynjarvis/fyp/dc/api"
)

type Processor struct {
	mu    sync.RWMutex
	imgCh chan *api.ImageReport
	cfg   *api.CollectionConfig
	cloud transport.CollectionService
	st    storage.Storage
}

func New(imgCh chan *api.ImageReport, cloud transport.CollectionService) *Processor {
	return &Processor{
		imgCh: imgCh,
		cloud: cloud,
		st:    storage.New(storage.Config{}),
	}
}

func (p *Processor) Execute() {
	go p.updateCfg()
	for img := range p.imgCh {
		if p.getCfg() == nil || len(p.getCfg().Rules) == 0 {
			logrus.Debug("no configured rules, skip")
			continue
		}
		p.execute(img)
	}
}

func (p *Processor) execute(img *api.ImageReport) {
	// [keep buffer]
	var res map[string]interface{}
	if err := json.Unmarshal(img.GetResult(), &res); err != nil {
		logrus.Error("unmarshal err: ", err)
		return
	}
	if !p.execRuleEngine(res) {
		return
	}

	switch img.Type {
	case api.ContentType_Image_PNG:
		if err := p.st.Image(img.Id, "image/png", img.GetImg()); err != nil {
			logrus.Error(err)
		}
	case api.ContentType_Image_JPEG:
		if err := p.st.Image(img.Id, "image/jpeg", img.GetImg()); err != nil {
			logrus.Error(err)
		}
	}

	if err := p.st.Data(img.Id, res); err != nil {
		logrus.Error(err)
	}

	// TODO: text indexing

	p.cloud.SendNotification(&api.CollectionEvent{
		Type:    img.Type,
		Message: "OK",
	})
}

func (p *Processor) execRuleEngine(res map[string]interface{}) bool {
	rules := p.getCfg().Rules
	valid := false
Loop:
	for _, rule := range rules {
		f := res[rule.Field].(float64)
		switch rule.Op {
		case api.Rule_lt:
			logrus.Info("hit")
			operand, err := strconv.ParseFloat(rule.Operand, 64)
			if err != nil {
				logrus.Error(err)
			}
			if f < operand && rand.Float64() < rule.SampleRate {
				valid = true
				break Loop
			}
		default:
			logrus.Error("unknown op: ", rule.Op)
		}
	}
	return valid
}

func (p *Processor) getCfg() *api.CollectionConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.cfg
}

func (p *Processor) updateCfg() {
	for cfg := range p.cloud.RecvConfig() {
		// TODO: add rate limiter
		p.mu.Lock()
		p.cfg = cfg
		p.st.UpdateConfig(storage.Config{
			ObjStoreAddr:  cfg.ObjectStoragePath,
			DataStoreAddr: cfg.DocumentStoragePath,
			TextIndexAddr: cfg.TextIndexPath,
		})
		p.mu.Unlock()
	}
}
