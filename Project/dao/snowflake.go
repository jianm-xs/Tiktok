// snowflake 分布式唯一ID生成
// 创建人：刘伟欢
// 创建时间：2022-5-24

package dao

import (
	"Project/config"
	"errors"
	"sync"
	"time"
)

type SnowflakeIdWorker struct {
	mutex         sync.Mutex // 互斥锁
	lastTimestamp int64      // 上次生成ID的时间戳
	workerId      int64      // 工作节点ID（0~31）
	dataCenterId  int64      // 数据中心机房ID（0~31）
	sequence      int64      // 序列号（0~4095）
}

const (
	epoch              int64 = 1640966400000                                  // 设置起始时间(时间戳/毫秒)：2022-01-01 00:00:00，有效期69年
	workerIdBits       int64 = 5                                              // 机器id所占的位数
	dataCenterIdBits   int64 = 5                                              // 数据标识id所占的位数
	timestampBits      int64 = 41                                             // 时间戳所占的位数
	maxWorkerId        int64 = -1 ^ (-1 << workerIdBits)                      // 支持的最大机器id，结果是31
	maxDataCenterId    int64 = -1 ^ (-1 << dataCenterIdBits)                  // 支持的最大数据标识id，结果是31
	maxTimestamp       int64 = -1 ^ (-1 << timestampBits)                     // 支持的最大时间戳
	sequenceBits       int64 = 12                                             // 序列在id中占的位数
	workerIdShift      int64 = sequenceBits                                   // 机器ID向左移12位
	datacenterIdShift  int64 = sequenceBits + workerIdBits                    // 数据标识id向左移17位(12+5)
	timestampLeftShift int64 = sequenceBits + workerIdBits + dataCenterIdBits // 时间截向左移22位(5+5+12)
	sequenceMask       int64 = -1 ^ (-1 << sequenceBits)                      // 生成序列的掩码，这里为4095(0b111111111111)
)

// 注册 ID 生成器
var RegisterIdWorker *SnowflakeIdWorker

// 视频 ID 生成器
var VideoIdWorker *SnowflakeIdWorker

// 评论 ID 生成器
var CommentIdWorker *SnowflakeIdWorker

// 粉丝 ID 生成器
var FollowIdWorker *SnowflakeIdWorker

// 点赞 ID 生成器
var FavoriteIdWorker *SnowflakeIdWorker

// 根据 workerId ， dataCenterId ,创建 ID 生成器
func createWorker(wId int64, dId int64) (*SnowflakeIdWorker, error) {
	if wId < 0 || wId > maxWorkerId {
		return nil, errors.New("Worker ID excess of quantity")
	}
	if dId < 0 || dId > maxDataCenterId {
		return nil, errors.New("Datacenter ID excess of quantity")
	}
	// 生成一个新节点
	return &SnowflakeIdWorker{
		lastTimestamp: 0,
		workerId:      wId,
		dataCenterId:  dId,
		sequence:      0,
	}, nil
}

// 初始化所有 ID 生成器
func InitIdWorker() error {
	//	创建 注册 ID 生成器
	registerIdWorker, err := createWorker(config.SnowFlakeCfg.WorkerId, config.SnowFlakeCfg.DateCenterId)
	RegisterIdWorker = registerIdWorker
	if err != nil { // 创建失败
		return err
	}

	//	创建 视频 ID 生成器
	videoIdWorker, err := createWorker(config.SnowFlakeCfg.WorkerId, config.SnowFlakeCfg.DateCenterId+1)
	VideoIdWorker = videoIdWorker
	if err != nil { // 创建失败
		return err
	}

	//	创建 评论 ID 生成器
	commentIdWorker, err := createWorker(config.SnowFlakeCfg.WorkerId, config.SnowFlakeCfg.DateCenterId+2)
	CommentIdWorker = commentIdWorker
	if err != nil {
		return err
	}

	//	创建 粉丝 ID 生成器
	followIdWorker, err := createWorker(config.SnowFlakeCfg.WorkerId, config.SnowFlakeCfg.DateCenterId+3)
	FollowIdWorker = followIdWorker
	if err != nil {
		return err
	}

	//	创建 点赞 ID 生成器
	favoriteIdWorker, err := createWorker(config.SnowFlakeCfg.WorkerId, config.SnowFlakeCfg.DateCenterId+4)
	FavoriteIdWorker = favoriteIdWorker

	return err
}

// 获取 ID ，ID 生成异常，返回-1与错误信息
func (w *SnowflakeIdWorker) NextId() (int64, error) {
	// 保障线程安全 加锁
	w.mutex.Lock()
	// 生成完成后 解锁
	defer w.mutex.Unlock()
	// 获取生成时的时间戳并转为毫秒
	now := time.Now().UnixNano() / 1e6
	// 如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过这个时候应当抛出异常
	if now < w.lastTimestamp {
		// 异常，生成 ID 失败
		return -1, errors.New("Clock moved backwards")
	}
	// 同一时间戳下生成 ID ，增加序列号
	if w.lastTimestamp == now {
		w.sequence = (w.sequence + 1) & sequenceMask
		if w.sequence == 0 {
			// 当前序列号超出12 bit
			// 阻塞到下一个毫秒，直到获得新的时间戳
			for now <= w.lastTimestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		// 当前时间与工作节点上一次生成ID的时间不一致 则需要重置工作节点生成ID的序号
		w.sequence = 0
	}
	// 将机器上一次生成ID的时间更新为当前时间
	w.lastTimestamp = now
	ID := int64((now-epoch)<<timestampLeftShift | w.dataCenterId<<datacenterIdShift | (w.workerId << workerIdShift) | w.sequence)
	return ID, nil
}
