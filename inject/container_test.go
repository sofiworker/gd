package inject

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
)

func TestProvideModuleEmptyName(t *testing.T) {
	container := New()
	module := NewModule("")
	err := container.ProvideModule(module)
	if err != ErrModuleNameEmpty {
		t.Errorf("want ErrModuleNameEmpty, got %v", err)
	}
}

func TestProvideModuleDuplicateName(t *testing.T) {
	container := New()
	module := NewModule("test")

	err := container.ProvideModule(module)
	if err != nil {
		t.Errorf("want nil error, got %v", err)
	}

	err = container.ProvideModule(module)
	if err != ErrModuleAlreadyExists {
		t.Errorf("want ErrModuleAlreadyExists, got %v", err)
	}
}

func TestProvideModuleWithParentNameNonExistentParent(t *testing.T) {
	container := New()
	module := NewModule("child")

	err := container.ProvideModuleWithParentName("nonexistent", module)
	if err != nil {
		t.Errorf("want nil error, got %v", err)
	}

	found := container.FindModuleByName("child", container.root)
	if found == nil {
		t.Error("module should be attached to root")
	}
}

func TestFindModuleByNameSuccess(t *testing.T) {
	container := New()
	module := NewModule("test")

	err := container.ProvideModule(module)
	if err != nil {
		t.Errorf("want nil error, got %v", err)
	}

	found := container.FindModuleByName("test", container.root)
	if found == nil {
		t.Error("module not found")
	}
}

func TestFindModuleByNameNotFound(t *testing.T) {
	container := New()
	found := container.FindModuleByName("nonexistent", container.root)
	if found != nil {
		t.Errorf("want nil, got %v", found)
	}
}

func TestProvideModuleWithMultiLevelHierarchy(t *testing.T) {
	container := New()

	// 创建多层级模块
	apiModule := NewModule("api")
	err := container.ProvideModule(apiModule)
	if err != nil {
		t.Fatalf("failed to provide api module: %v", err)
	}

	httpModule := NewModule("http")
	err = container.ProvideModuleWithParentName("api", httpModule)
	if err != nil {
		t.Fatalf("failed to provide http module: %v", err)
	}

	handlerModule := NewModule("handler")
	err = container.ProvideModuleWithParentName("http", handlerModule)
	if err != nil {
		t.Fatalf("failed to provide handler module: %v", err)
	}

	// 验证层级关系
	if container.root.modules[0].Name != "api" {
		t.Errorf("expected first level module to be 'api', got %s", container.root.modules[0].Name)
	}

	apiMod := container.FindModuleByName("api", container.root)
	if apiMod == nil {
		t.Fatal("api module not found")
	}

	httpMod := container.FindModuleByName("http", container.root)
	if httpMod == nil {
		t.Fatal("http module not found")
	}

	handlerMod := container.FindModuleByName("handler", container.root)
	if handlerMod == nil {
		t.Fatal("handler module not found")
	}

	// 打印模块结构
	t.Logf("Container structure:")
	printModuleTree(t, container.root, 0)
}

func TestProvideModuleWithSiblings(t *testing.T) {
	container := New()

	// 创建第一层级的多个模块
	apiModule := NewModule("api")
	err := container.ProvideModule(apiModule)
	if err != nil {
		t.Fatalf("failed to provide api module: %v", err)
	}

	grpcModule := NewModule("grpc")
	err = container.ProvideModule(grpcModule)
	if err != nil {
		t.Fatalf("failed to provide grpc module: %v", err)
	}

	// 为 api 模块添加多个子模块
	httpModule := NewModule("http")
	restModule := NewModule("rest")
	websocketModule := NewModule("websocket")

	err = container.ProvideModuleWithParentName("api", httpModule)
	if err != nil {
		t.Fatalf("failed to provide http module: %v", err)
	}

	err = container.ProvideModuleWithParentName("api", restModule)
	if err != nil {
		t.Fatalf("failed to provide rest module: %v", err)
	}

	err = container.ProvideModuleWithParentName("api", websocketModule)
	if err != nil {
		t.Fatalf("failed to provide websocket module: %v", err)
	}

	// 为 grpc 模块添加子模块
	serviceModule := NewModule("service")
	clientModule := NewModule("client")

	err = container.ProvideModuleWithParentName("grpc", serviceModule)
	if err != nil {
		t.Fatalf("failed to provide service module: %v", err)
	}

	err = container.ProvideModuleWithParentName("grpc", clientModule)
	if err != nil {
		t.Fatalf("failed to provide client module: %v", err)
	}

	// 验证根节点的子模块数量
	if len(container.root.modules) != 2 {
		t.Errorf("expected root to have 2 children, got %d", len(container.root.modules))
	}

	// 验证 api 模块的子模块数量
	apiMod := container.FindModuleByName("api", container.root)
	if len(apiMod.modules) != 3 {
		t.Errorf("expected api module to have 3 children, got %d", len(apiMod.modules))
	}

	// 验证 grpc 模块的子模块数量
	grpcMod := container.FindModuleByName("grpc", container.root)
	if len(grpcMod.modules) != 2 {
		t.Errorf("expected grpc module to have 2 children, got %d", len(grpcMod.modules))
	}

	// 打印完整的模块树结构
	t.Log("Container structure:")
	printModuleTree(t, container.root, 0)
}

func TestConcurrentModuleOperations(t *testing.T) {
	container := New()
	const numGoroutines = 10
	const numModulesPerGoroutine = 100

	// 创建根模块
	rootModules := []string{"api", "grpc", "db", "cache", "queue"}
	for _, name := range rootModules {
		if err := container.ProvideModule(NewModule(name)); err != nil {
			t.Fatalf("failed to create root module %s: %v", name, err)
		}
	}

	// 使用 WaitGroup 等待所有 goroutine 完成
	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines*2)

	// 并发提供模块
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < numModulesPerGoroutine; j++ {
				parentName := rootModules[j%len(rootModules)]
				moduleName := fmt.Sprintf("module_%d_%d", routineID, j)
				module := NewModule(moduleName)

				if err := container.ProvideModuleWithParentName(parentName, module); err != nil {
					errChan <- fmt.Errorf("routine %d failed to provide module %s: %v", routineID, moduleName, err)
					return
				}
			}
		}(i)
	}

	// 并发查找模块
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < numModulesPerGoroutine; j++ {
				randomParent := rootModules[rand.Intn(len(rootModules))]
				found := container.FindModuleByName(randomParent, container.root)
				if found == nil {
					errChan <- fmt.Errorf("routine %d failed to find module %s", routineID, randomParent)
					return
				}
			}
		}(i)
	}

	// 等待所有操作完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误发生
	for err := range errChan {
		t.Error(err)
	}

	// 验证最终的模块数量
	t.Log("Final container structure:")
	printModuleTree(t, container.root, 0)

	// 验证每个根模块的子模块数量
	for _, rootName := range rootModules {
		root := container.FindModuleByName(rootName, container.root)
		if root == nil {
			t.Errorf("root module %s not found", rootName)
			continue
		}
		expectedChildren := numGoroutines * (numModulesPerGoroutine / len(rootModules))
		if len(root.modules) != expectedChildren {
			t.Errorf("root module %s: expected %d children, got %d", rootName, expectedChildren, len(root.modules))
		}
	}
}

func printModuleTree(t *testing.T, module *Module, level int) {
	indent := strings.Repeat("  ", level)
	t.Logf("%s- %s", indent, module.Name)
	for _, child := range module.modules {
		printModuleTree(t, child, level+1)
	}
}

func BenchmarkModuleOperations(b *testing.B) {
	/**
	goos: windows
	goarch: amd64
	pkg: github.com/chuck1024/gd/v2/inject
	cpu: QEMU Virtual CPU version 2.5+
	BenchmarkModuleOperations
	BenchmarkModuleOperations/DeepHierarchyProvideAndFind
	BenchmarkModuleOperations/DeepHierarchyProvideAndFind-8         	       1	52894112500 ns/op
	BenchmarkModuleOperations/ConcurrentProvideAndFind
	BenchmarkModuleOperations/ConcurrentProvideAndFind-8            	       1	105963969800 ns/op
	PASS
	**/
	// 预定义测试规模
	const (
		numRootModules    = 10
		numChildrenPerMod = 10
		treeDepth         = 5
	)

	b.Run("DeepHierarchyProvideAndFind", func(b *testing.B) {
		container := New()

		// 重置定时器，不计入初始化时间
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 创建一个深层次的模块树
			moduleNames := createDeepModuleTree(b, container, numRootModules, numChildrenPerMod, treeDepth)

			// 随机查找一些模块
			b.StopTimer()
			targetNames := selectRandomModules(moduleNames, 100)
			b.StartTimer()

			for _, name := range targetNames {
				if found := container.FindModuleByName(name, container.root); found == nil {
					b.Fatalf("module %s not found", name)
				}
			}

			b.StopTimer()
			container = New() // 重置容器
			b.StartTimer()
		}
	})

	b.Run("ConcurrentProvideAndFind", func(b *testing.B) {
		container := New()
		var wg sync.WaitGroup

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 并发创建和查找模块
			wg.Add(numRootModules)
			for j := 0; j < numRootModules; j++ {
				go func(rootID int) {
					defer wg.Done()
					createModuleBranch(container, rootID, numChildrenPerMod, treeDepth)
				}(j)
			}
			wg.Wait()

			b.StopTimer()
			container = New()
			b.StartTimer()
		}
	})
}

// 创建深层次模块树，返回所有模块名称
func createDeepModuleTree(b *testing.B, container *Container, numRoots, numChildren, depth int) []string {
	var moduleNames []string

	// 创建根模块
	for i := 0; i < numRoots; i++ {
		rootName := fmt.Sprintf("root_%d", i)
		rootMod := NewModule(rootName)
		if err := container.ProvideModule(rootMod); err != nil {
			b.Fatalf("failed to provide root module: %v", err)
		}
		moduleNames = append(moduleNames, rootName)

		// 递归创建子模块
		createChildModules(b, container, rootName, numChildren, depth-1, &moduleNames)
	}

	return moduleNames
}

// 递归创建子模块
func createChildModules(b *testing.B, container *Container, parentName string, numChildren, depth int, moduleNames *[]string) {
	if depth <= 0 {
		return
	}

	for i := 0; i < numChildren; i++ {
		childName := fmt.Sprintf("%s_child_%d", parentName, i)
		childMod := NewModule(childName)
		if err := container.ProvideModuleWithParentName(parentName, childMod); err != nil {
			b.Fatalf("failed to provide child module: %v", err)
		}
		*moduleNames = append(*moduleNames, childName)

		// 递归创建下一层
		createChildModules(b, container, childName, numChildren, depth-1, moduleNames)
	}
}

// 为并发测试创建模块分支
func createModuleBranch(container *Container, rootID, numChildren, depth int) {
	rootName := fmt.Sprintf("root_%d", rootID)
	rootMod := NewModule(rootName)
	if err := container.ProvideModule(rootMod); err != nil {
		return
	}

	createBranchChildren(container, rootName, numChildren, depth-1)
}

// 创建分支子模块
func createBranchChildren(container *Container, parentName string, numChildren, depth int) {
	if depth <= 0 {
		return
	}
	for i := 0; i < numChildren; i++ {
		childName := fmt.Sprintf("%s_child_%d", parentName, i)
		childMod := NewModule(childName)
		if err := container.ProvideModuleWithParentName(parentName, childMod); err != nil {
			return
		}
		createBranchChildren(container, childName, numChildren, depth-1)
	}
}

// 从模块名列表中随机选择一些模块
func selectRandomModules(moduleNames []string, count int) []string {
	if count > len(moduleNames) {
		count = len(moduleNames)
	}

	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = moduleNames[rand.Intn(len(moduleNames))]
	}
	return selected
}

func BenchmarkContainerConcurrentAccess(b *testing.B) {
	// 预定义测试参数
	const (
		numReaders = 100
		numWriters = 10
		numModules = 1000
	)

	b.Run("ConcurrentReadWrite", func(b *testing.B) {
		container := New()
		var wg sync.WaitGroup

		// 预先创建一些模块
		for i := 0; i < 10; i++ {
			rootMod := NewModule(fmt.Sprintf("root_%d", i))
			if err := container.ProvideModule(rootMod); err != nil {
				b.Fatalf("failed to provide root module: %v", err)
			}
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			// 启动读取 goroutines
			wg.Add(numReaders)
			for r := 0; r < numReaders; r++ {
				go func() {
					defer wg.Done()
					for j := 0; j < 100; j++ {
						name := fmt.Sprintf("root_%d", j%10)
						_ = container.FindModuleByName(name, container.root)
					}
				}()
			}

			// 启动写入 goroutines
			wg.Add(numWriters)
			for w := 0; w < numWriters; w++ {
				go func(writerID int) {
					defer wg.Done()
					for j := 0; j < numModules/numWriters; j++ {
						modName := fmt.Sprintf("mod_%d_%d", writerID, j)
						parentName := fmt.Sprintf("root_%d", j%10)
						mod := NewModule(modName)
						_ = container.ProvideModuleWithParentName(parentName, mod)
					}
				}(w)
			}

			wg.Wait()
			b.StopTimer()
			container = New() // 重置容器
			// 重新创建根模块
			for i := 0; i < 10; i++ {
				rootMod := NewModule(fmt.Sprintf("root_%d", i))
				if err := container.ProvideModule(rootMod); err != nil {
					b.Fatalf("failed to provide root module: %v", err)
				}
			}
			b.StartTimer()
		}
	})

	/*
		goos: windows
		goarch: amd64
		pkg: github.com/chuck1024/gd/v2/inject
		cpu: QEMU Virtual CPU version 2.5+
		BenchmarkContainerConcurrentAccess
		BenchmarkContainerConcurrentAccess/ConcurrentReadsOnly
		BenchmarkContainerConcurrentAccess/ConcurrentReadsOnly-8         	     188	   7249627 ns/op
	*/
	b.Run("ConcurrentReadsOnly", func(b *testing.B) {
		container := New()
		var wg sync.WaitGroup

		// 预先创建模块树
		createTestModuleTree(container)

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			wg.Add(numReaders)
			for r := 0; r < numReaders; r++ {
				go func() {
					defer wg.Done()
					for j := 0; j < 1000; j++ {
						depth := j % 3
						name := fmt.Sprintf("module_%d_%d", depth, j%10)
						_ = container.FindModuleByName(name, container.root)
					}
				}()
			}
			wg.Wait()
		}
	})
}

// 创建测试用的模块树
func createTestModuleTree(container *Container) {
	// 创建三层深度的模块树
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			name := fmt.Sprintf("module_%d_%d", i, j)
			if i == 0 {
				module := NewModule(name)
				container.ProvideModule(module)
			} else {
				parentName := fmt.Sprintf("module_%d_%d", i-1, j)
				module := NewModule(name)
				container.ProvideModuleWithParentName(parentName, module)
			}
		}
	}
}
