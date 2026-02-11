package res

import (
	"fmt"
	"log"
	"sync"

	"go.uber.org/zap"
)

// Resource 资源接口
type Resource interface {
	Name() string
	Close() error // 添加 Close 方法
}

var (
	// 存储所有注册的资源
	resources   = make(map[string]func() Resource)
	registerMux sync.RWMutex // 新增读写锁
)

// NewResource 创建资源实例
func NewResource(name string) Resource {
	registerMux.RLock()
	defer registerMux.RUnlock()

	if constructor, exists := resources[name]; exists {
		return constructor()
	}

	return nil
}

// RegisterResource 注册全局资源
func RegisterResource(name string, constructor func() Resource) {
	registerMux.Lock()
	defer registerMux.Unlock()

	if _, exists := resources[name]; exists {
		panic("resource type already registered: " + name)
	}

	resources[name] = constructor
}

// ResourceNames 所有注册的资源名称
func ResourceNames() []string {
	registerMux.Lock()
	defer registerMux.Unlock()

	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	return names
}

func InitResources() {
	// 初始化 MySQL 连接
	if err := RegisterMysqlDBToGlobal(); err != nil {
		log.Fatalf("Failed to register global resources: %v", err)
	}

	// 初始化 ...连接
}
func Print() {
	zap.L().Info("++++++++++++++ resource info begin: ++++++++++++++")
	for _, name := range ResourceNames() {
		zap.L().Info("-- ", zap.String("name", name))
	}
	zap.L().Info("++++++++++++++ resource info end: ++++++++++++++")
}

// 添加关闭所有资源的方法
func CloseAllResources() error {
	registerMux.Lock()
	defer registerMux.Unlock()

	var errs []error
	for name, constructor := range resources {
		if constructor == nil {
			continue
		}

		resource := constructor()
		if closer, ok := resource.(interface{ Close() error }); ok {
			zap.L().Info("Closing resource", zap.String("name", name))
			if err := closer.Close(); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to close %d resources", len(errs))
	}
	return nil
}
