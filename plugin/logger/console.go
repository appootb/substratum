package logger

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/appootb/substratum/v2/logger"
	"github.com/appootb/substratum/v2/proto/go/common"
	"github.com/appootb/substratum/v2/util/jsonpb"
)

func Init() {
	if logger.Implementor() == nil {
		logger.RegisterImplementor(&Console{})
	}
}

type Console struct {
	level int32
}

func (log *Console) UpdateLevel(level logger.Level) {
	atomic.StoreInt32(&log.level, int32(level))
}

func (log *Console) Log(level logger.Level, md *common.Metadata, msg string, c logger.Content) {
	if int32(level) < atomic.LoadInt32(&log.level) {
		return
	}
	var (
		meta    []byte
		content []byte
	)
	if md != nil {
		meta, _ = jsonpb.Marshal(md)
	}
	if c != nil && len(c) > 0 {
		content, _ = json.Marshal(c)
	}
	fmt.Println(fmt.Sprintf("%v metadata: %v, %v: %v", level.String(), string(meta), msg, string(content)))
}
