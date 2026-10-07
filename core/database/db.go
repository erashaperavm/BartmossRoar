package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// DB 是对 badger.DB 的薄封装。
// 只负责生命周期和事务入口，不掺业务逻辑。
// 六个逻辑表共用同一个实例，靠 key 前缀区分。
type DB struct {
	badger *badger.DB
	path   string
}

// Options 保留可控参数，prototype 阶段直接用默认值即可。
type Options struct {
	Path string // 空字符串表示内存模式
	// 日志级别，prototype 建议 Error 或 Silent，别被 INFO 刷屏
	LogLevel badger.Logger
}

// Open 打开或创建数据库。
// 只允许一个进程持有同一路径，重复打开会返回错误。
func Open(opts Options) (*DB, error) {
	badgerOpts := badger.DefaultOptions(opts.Path)

	if opts.LogLevel != nil {
		badgerOpts = badgerOpts.WithLogger(opts.LogLevel)
	} else {
		// prototype 阶段默认静默，避免每次 txn 都打日志
		badgerOpts = badgerOpts.WithLogger(nil)
	}

	raw, err := badger.Open(badgerOpts)
	if err != nil {
		return nil, fmt.Errorf("badger open %q: %w", opts.Path, err)
	}

	return &DB{
		badger: raw,
		path:   opts.Path,
	}, nil
}

// Close 关闭数据库。
// 先停后台 compaction，再刷盘，最后释放文件锁。
func (db *DB) Close() error {
	if db.badger == nil {
		return nil
	}
	if err := db.badger.Close(); err != nil {
		return fmt.Errorf("badger close: %w", err)
	}
	db.badger = nil
	return nil
}

// Path 返回当前数据库路径，方便日志和调试。
func (db *DB) Path() string {
	return db.path
}

// ------------------------------------------------------------
// 事务入口
//
// 只暴露 View / Update 两个方法，业务层不直接碰 *badger.DB。
// 这样以后想换底层引擎，只需改这个文件。
// ------------------------------------------------------------

// View 只读事务，返回 error 会中止事务。
func (db *DB) View(fn func(txn *badger.Txn) error) error {
	if db.badger == nil {
		return errors.New("database: not open")
	}
	return db.badger.View(fn)
}

// Update 读写事务，返回 error 会回滚。
func (db *DB) Update(fn func(txn *badger.Txn) error) error {
	if db.badger == nil {
		return errors.New("database: not open")
	}
	return db.badger.Update(fn)
}

// ------------------------------------------------------------
// 通用工具
// ------------------------------------------------------------

// Set 写入一个带可选 TTL 的键值。
// ttl <= 0 表示永不过期。
func (db *DB) Set(key, value []byte, ttl time.Duration) error {
	return db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, value)
		if ttl > 0 {
			entry = entry.WithTTL(ttl)
		}
		return txn.SetEntry(entry)
	})
}

// Get 读取一个键。
// 不存在时返回 ErrKeyNotFound，调用方自行判断。
func (db *DB) Get(key []byte) ([]byte, error) {
	var out []byte

	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		out, err = item.ValueCopy(nil)
		return err
	})

	return out, err
}

// Delete 删除一个键。
func (db *DB) Delete(key []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Has 判断键是否存在。
func (db *DB) Has(key []byte) (bool, error) {
	exists := false
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil
		}
		if err != nil {
			return err
		}
		exists = true
		return nil
	})
	return exists, err
}

// ------------------------------------------------------------
// 前缀扫描
//
// 六个逻辑表靠前缀区分，所以统一提供一个扫描入口。
// 业务层只需要传前缀和回调。
// ------------------------------------------------------------

// ScanPrefix 遍历所有匹配前缀的键值。
// 回调返回 error 会中止扫描。
func (db *DB) ScanPrefix(prefix []byte, fn func(key, value []byte) error) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = prefix
		opts.PrefetchValues = false // 需要时再 ValueCopy，省内存

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.KeyCopy(nil)
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			if err := fn(key, val); err != nil {
				return err
			}
		}
		return nil
	})
}

// ------------------------------------------------------------
// 维护操作
//
// prototype 阶段不强制调用，但留着以后有用。
// ------------------------------------------------------------

// RunGC 手动触发 value log GC。
// Badger 会自动做一部分，但写多删多时手动跑一次可以回收空间。
func (db *DB) RunGC() error {
	if db.badger == nil {
		return errors.New("database: not open")
	}
	// 反复调用直到没有可回收的，或出错
	for {
		err := db.badger.RunValueLogGC(0.5)
		if err != nil {
			// ErrNoRewrite 表示没有可回收的，正常结束
			if errors.Is(err, badger.ErrNoRewrite) {
				return nil
			}
			return fmt.Errorf("value log gc: %w", err)
		}
	}
}

// Backup 把数据库备份到 writer。
// 以后要做导出、迁移时用得上。
func (db *DB) Backup(w interface{ Write([]byte) (int, error) }) error {
	if db.badger == nil {
		return errors.New("database: not open")
	}
	_, err := db.badger.Backup(w, 0)
	return err
}
