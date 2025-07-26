package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"google_storage_upload/pkg/storage"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gsupload",
		Usage: "上传文件或目录到 Google Cloud Storage",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "bucket",
				Aliases:  []string{"b"},
				Usage:    "Google Cloud Storage 存储桶名称",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "credentials",
				Aliases: []string{"c"},
				Usage:   "Google Cloud 凭证文件路径 (可选，默认使用环境变量)",
			},
			&cli.IntFlag{
				Name:    "expiration",
				Aliases: []string{"e"},
				Usage:   "对象过期天数 (可选)",
				Value:   0,
			},
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "对象名称前缀 (可选)",
				Value:   "",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "file",
				Usage:     "上传单个文件",
				ArgsUsage: "FILE_PATH",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("请指定要上传的文件路径")
					}
					return uploadFile(c)
				},
			},
			{
				Name:      "dir",
				Usage:     "上传整个目录",
				ArgsUsage: "DIRECTORY_PATH",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("请指定要上传的目录路径")
					}
					return uploadDirectory(c)
				},
			},
			{
				Name:      "list",
				Usage:     "列出存储桶中的对象",
				ArgsUsage: "[PREFIX]",
				Action: func(c *cli.Context) error {
					prefix := ""
					if c.NArg() > 0 {
						prefix = c.Args().Get(0)
					}
					return listObjects(c, prefix)
				},
			},
			{
				Name:      "delete",
				Usage:     "删除存储桶中的对象",
				ArgsUsage: "OBJECT_NAME",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("请指定要删除的对象名称")
					}
					return deleteObject(c)
				},
			},
			{
				Name:      "info",
				Usage:     "获取对象的元数据",
				ArgsUsage: "OBJECT_NAME",
				Action: func(c *cli.Context) error {
					if c.NArg() < 1 {
						return fmt.Errorf("请指定要查询的对象名称")
					}
					return getObjectInfo(c)
				},
			},
			{
				Name:      "expire",
				Usage:     "设置对象的过期时间",
				ArgsUsage: "OBJECT_NAME DAYS",
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return fmt.Errorf("请指定对象名称和过期天数")
					}
					return setObjectExpiration(c)
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func uploadFile(c *cli.Context) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")
	expirationDays := c.Int("expiration")
	prefix := c.String("prefix")
	filePath := c.Args().Get(0)

	// 验证文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("无法访问文件 %s: %v", filePath, err)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("%s 是一个目录，请使用 'dir' 命令上传目录", filePath)
	}

	// 构建对象名称
	objectName := prefix
	if objectName != "" && !strings.HasSuffix(objectName, "/") {
		objectName += "/"
	}
	objectName += filepath.Base(filePath)

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 上传文件
	startTime := time.Now()
	result, err := uploader.UploadFile(filePath, objectName, expirationDays)
	if err != nil {
		return err
	}

	fmt.Printf("文件上传成功！\n")
	fmt.Printf("文件: %s\n", result.LocalPath)
	fmt.Printf("对象: %s\n", result.ObjectName)
	fmt.Printf("URL: %s\n", result.URL)
	fmt.Printf("大小: %d 字节\n", result.Size)
	fmt.Printf("类型: %s\n", result.MimeType)
	if result.Expiration != nil {
		fmt.Printf("过期时间: %s\n", result.Expiration.Format("2006-01-02 15:04:05"))
	}
	fmt.Printf("耗时: %v\n", time.Since(startTime))

	return nil
}

func uploadDirectory(c *cli.Context) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")
	expirationDays := c.Int("expiration")
	prefix := c.String("prefix")
	dirPath := c.Args().Get(0)

	// 验证目录是否存在
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		return fmt.Errorf("无法访问目录 %s: %v", dirPath, err)
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s 不是一个目录，请使用 'file' 命令上传文件", dirPath)
	}

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 上传目录
	startTime := time.Now()
	results, err := uploader.UploadDirectory(dirPath, prefix, expirationDays)
	if err != nil {
		return err
	}

	fmt.Printf("\n目录上传成功！\n")
	fmt.Printf("目录: %s\n", dirPath)
	fmt.Printf("前缀: %s\n", prefix)
	fmt.Printf("上传文件数: %d\n", len(results))
	if expirationDays > 0 {
		fmt.Printf("过期时间: %d 天后\n", expirationDays)
	}
	fmt.Printf("耗时: %v\n", time.Since(startTime))

	return nil
}

// listObjects 列出存储桶中的对象
func listObjects(c *cli.Context, prefix string) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 列出对象
	objects, err := uploader.ListObjects(prefix)
	if err != nil {
		return err
	}

	if len(objects) == 0 {
		fmt.Printf("存储桶 %s 中没有找到对象", bucketName)
		if prefix != "" {
			fmt.Printf("（前缀：%s）", prefix)
		}
		fmt.Println()
		return nil
	}

	fmt.Printf("存储桶 %s 中的对象", bucketName)
	if prefix != "" {
		fmt.Printf("（前缀：%s）", prefix)
	}
	fmt.Println(":")

	for i, obj := range objects {
		fmt.Printf("%d. %s\n", i+1, obj)
	}

	fmt.Printf("\n共 %d 个对象\n", len(objects))
	return nil
}

// deleteObject 删除存储桶中的对象
func deleteObject(c *cli.Context) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")
	objectName := c.Args().Get(0)

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 删除对象
	if err := uploader.DeleteObject(objectName); err != nil {
		return err
	}

	fmt.Printf("对象 %s 已成功删除\n", objectName)
	return nil
}

// getObjectInfo 获取对象的元数据
func getObjectInfo(c *cli.Context) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")
	objectName := c.Args().Get(0)

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 获取对象元数据
	attrs, err := uploader.GetObjectMetadata(objectName)
	if err != nil {
		return err
	}

	fmt.Printf("对象信息：\n")
	fmt.Printf("名称: %s\n", attrs.Name)
	fmt.Printf("存储桶: %s\n", attrs.Bucket)
	fmt.Printf("大小: %d 字节\n", attrs.Size)
	fmt.Printf("内容类型: %s\n", attrs.ContentType)
	fmt.Printf("创建时间: %s\n", attrs.Created.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", attrs.Updated.Format("2006-01-02 15:04:05"))
	fmt.Printf("生成: %d\n", attrs.Generation)
	fmt.Printf("元数据: %v\n", attrs.Metadata)

	// 检查是否设置了过期时间
	if expiration, ok := attrs.Metadata["expiration"]; ok {
		fmt.Printf("过期时间: %s\n", expiration)
	} else {
		fmt.Printf("过期时间: 未设置\n")
	}

	return nil
}

// setObjectExpiration 设置对象的过期时间
func setObjectExpiration(c *cli.Context) error {
	bucketName := c.String("bucket")
	credentialsFile := c.String("credentials")
	objectName := c.Args().Get(0)

	// 解析过期天数
	expirationDays := 0
	var err error
	expirationDaysStr := c.Args().Get(1)
	expirationDays, err = strconv.Atoi(expirationDaysStr)
	if err != nil {
		return fmt.Errorf("无效的过期天数: %s", expirationDaysStr)
	}

	if expirationDays <= 0 {
		return fmt.Errorf("过期天数必须大于 0")
	}

	// 创建上传器
	ctx := context.Background()
	uploader, err := storage.NewUploader(ctx, bucketName, credentialsFile)
	if err != nil {
		return err
	}
	defer uploader.Close()

	// 设置对象过期时间
	if err := uploader.SetObjectExpiration(objectName, expirationDays); err != nil {
		return err
	}

	fmt.Printf("对象 %s 的过期时间已设置为 %d 天后\n", objectName, expirationDays)
	return nil
}
