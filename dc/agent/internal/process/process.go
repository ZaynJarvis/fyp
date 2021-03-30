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

//s := storage.New(storage.Config{
//	ObjStoreAddr:  "localhost:9000",
//	DocStorageAddr: "localhost:27017",
//})
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
		go p.execute(img)
	}
}

func (p *Processor) execute(img *api.ImageReport) {
	// can have keep buffer
	if img.Id == "" {
		logrus.Error("invalid ID")
		return
	}
	var res map[string]interface{}
	if err := json.Unmarshal(img.GetResult(), &res); err != nil {
		logrus.Error("unmarshal err: ", err)
		return
	}
	if !p.execRuleEngine(res) {
		return
	}

	if err := p.st.Image(img.Id, img.GetImg()); err != nil {
		logrus.Error(err)
	}

	if err := p.st.Doc(img.Id, res); err != nil {
		logrus.Error(err)
	}

	// TODO: text indexing

	p.cloud.SendNotification(&api.CollectionEvent{
		Type:    img.Type,
		Message: "OK",
	})
}

func (p *Processor) execRuleEngine(res map[string]interface{}) bool {
	logrus.Debugf("exec rule on result: %s", res)
	rules := p.getCfg().Rules
	valid := false
	collectedRule := "unknown"
Loop:
	for _, rule := range rules {
		switch rule.Op {
		case api.Rule_exist:
			if f := res[rule.Field]; f != nil && rand.Float64() < rule.SampleRate {
				collectedRule = rule.String()
				valid = true
				break Loop
			}
		case api.Rule_not_exist:
			if f := res[rule.Field]; f == nil && rand.Float64() < rule.SampleRate {
				collectedRule = rule.String()
				valid = true
				break Loop
			}
		case api.Rule_lt:
			f := res[rule.Field].(float64)
			operand, err := strconv.ParseFloat(rule.Operand, 64)
			if err != nil {
				logrus.Error(err)
			}
			if f < operand && rand.Float64() < rule.SampleRate {
				collectedRule = rule.String()
				valid = true
				break Loop
			}
		case api.Rule_gt:
			f := res[rule.Field].(float64)
			operand, err := strconv.ParseFloat(rule.Operand, 64)
			if err != nil {
				logrus.Error(err)
			}
			if f > operand && rand.Float64() < rule.SampleRate {
				collectedRule = rule.String()
				valid = true
				break Loop
			}
		default:
			logrus.Error("no support to op: ", rule.Op)
		}
	}
	if valid {
		res["text"] = "collected because of the rule [" + collectedRule + "]"
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
		// should add rate r
		p.mu.Lock()
		p.cfg = cfg
		logrus.Debugf("received config: %s", cfg.String())
		p.st.UpdateConfig(storage.Config{
			ObjStoreAddr:   cfg.ObjectStoragePath,
			DocStorageAddr: cfg.DocumentStoragePath,
			TextIndexAddr:  cfg.TextIndexPath,
		})
		logrus.Info("storage update config completed")
		p.mu.Unlock()
	}
}
