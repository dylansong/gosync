package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config 结构体用于解析YAML配置
type Config struct {
	SyncConfigs []SyncConfig `yaml:"sync_configs"`
}

type SyncConfig struct {
	Name       string   `yaml:"name"`
	SourceDir  string   `yaml:"source_dir"`
	TargetDirs []string `yaml:"target_dirs"`
}

const configExample = `sync_configs:
  - name: "sync1"
    source_dir: "/path/to/source1"
    target_dirs:
      - "/path/to/target1"
      - "/path/to/target2"
  - name: "sync2"
    source_dir: "/path/to/source2"
    target_dirs:
      - "/path/to/target3"
      - "/path/to/target4"
`

const helpText = `使用说明:
1. 创建配置文件:
   echo '
sync_configs:
  - name: "sync1"
    source_dir: "/path/to/source1"
    target_dirs:
      - "/path/to/target1"
      - "/path/to/target2"' > config.yaml

2. 或者使用以下命令创建示例配置:
   gosync -example > config.yaml

3. 运行同步:
   gosync              # 使用当前目录下的 config.yaml
   gosync -config 文件路径  # 指定配置文件路径
`

func main() {
	// 定义命令行参数
	configPath := flag.String("config", "config.yaml", "配置文件路径 (默认为当前目录下的 config.yaml)")
	showExample := flag.Bool("example", false, "显示配置文件示例")
	
	// 自定义 Usage 信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", helpText)
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n配置文件示例:\n%s\n", configExample)
	}
	
	flag.Parse()

	// 如果使用 -example 参数，显示配置文件示例后退出
	if *showExample {
		fmt.Println("配置文件示例:")
		fmt.Println(configExample)
		return
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取当前目录失败: %v", err)
	}

	// 如果配置文件路径是相对路径，则基于当前目录
	configFilePath := *configPath
	if !filepath.IsAbs(configFilePath) {
		configFilePath = filepath.Join(currentDir, configFilePath)
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在: %s", configFilePath)
	}

	// 读取配置文件
	config := readConfig(configFilePath)

	// 执行同步
	for _, syncConfig := range config.SyncConfigs {
		fmt.Printf("正在处理同步配置: %s\n", syncConfig.Name)
		syncFiles(syncConfig)
	}
}

// 读取配置文件
func readConfig(configPath string) Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	return config
}

// 同步文件
func syncFiles(config SyncConfig) {
	// 遍历源目录
	err := filepath.Walk(config.SourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(config.SourceDir, path)
		if err != nil {
			return err
		}

		// 对目录不做处理
		if info.IsDir() {
			return nil
		}

		// 复制到每个目标目录
		for _, targetDir := range config.TargetDirs {
			targetPath := filepath.Join(targetDir, relPath)

			// 创建目标目录
			err = os.MkdirAll(filepath.Dir(targetPath), 0755)
			if err != nil {
				return err
			}

			// 复制文件
			err = copyFile(path, targetPath)
			if err != nil {
				return err
			}
			fmt.Printf("已复制: %s -> %s\n", path, targetPath)
		}

		return nil
	})

	if err != nil {
		log.Printf("同步失败: %v\n", err)
	}
}

// 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
