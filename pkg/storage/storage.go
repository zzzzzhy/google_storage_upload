package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Uploader 定义了上传文件到 Google Cloud Storage 的接口
type Uploader struct {
	client     *storage.Client
	bucketName string
	ctx        context.Context
}

// UploadResult 表示上传结果
type UploadResult struct {
	LocalPath  string
	ObjectName string
	URL        string
	Size       int64
	MimeType   string
	Expiration *time.Time
}

// NewUploader 创建一个新的 Uploader 实例
func NewUploader(ctx context.Context, bucketName, credentialsFile string) (*Uploader, error) {
	var client *storage.Client
	var err error

	if credentialsFile != "" {
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	} else {
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("创建 Storage 客户端失败: %v", err)
	}

	// 不再检查存储桶是否存在，直接返回 Uploader
	return &Uploader{
		client:     client,
		bucketName: bucketName,
		ctx:        ctx,
	}, nil
}

// Close 关闭 Uploader 的客户端连接
func (u *Uploader) Close() error {
	return u.client.Close()
}

// UploadFile 上传单个文件到 Google Cloud Storage
func (u *Uploader) UploadFile(filePath, objectName string, expirationDays int) (*UploadResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败 %s: %v", filePath, err)
	}

	// 如果没有指定对象名称，使用文件名
	if objectName == "" {
		objectName = filepath.Base(filePath)
	}

	// 确保对象名称不以 / 开头
	objectName = strings.TrimPrefix(objectName, "/")

	// 检测文件的 MIME 类型
	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 创建对象写入器
	bucket := u.client.Bucket(u.bucketName)
	obj := bucket.Object(objectName)
	wc := obj.NewWriter(u.ctx)
	wc.ContentType = mimeType

	// 复制文件内容到对象
	if _, err = io.Copy(wc, file); err != nil {
		return nil, fmt.Errorf("上传文件 %s 失败: %v", filePath, err)
	}

	// 关闭写入器
	if err := wc.Close(); err != nil {
		return nil, fmt.Errorf("完成上传 %s 失败: %v", filePath, err)
	}

	// 创建上传结果
	result := &UploadResult{
		LocalPath:  filePath,
		ObjectName: objectName,
		URL:        fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.bucketName, objectName),
		Size:       fileInfo.Size(),
		MimeType:   mimeType,
	}

	// 如果指定了过期时间，设置对象的生命周期
	if expirationDays > 0 {
		expirationTime := time.Now().AddDate(0, 0, expirationDays)
		result.Expiration = &expirationTime

		// 设置对象元数据
		objAttrs := storage.ObjectAttrsToUpdate{
			Metadata: map[string]string{
				"expiration": expirationTime.Format(time.RFC3339),
			},
		}
		if _, err := obj.Update(u.ctx, objAttrs); err != nil {
			return result, fmt.Errorf("设置对象 %s 过期时间失败: %v", objectName, err)
		}
	}

	return result, nil
}

// UploadDirectory 递归上传目录到 Google Cloud Storage
func (u *Uploader) UploadDirectory(dirPath, objectPrefix string, expirationDays int) ([]*UploadResult, error) {
	var results []*UploadResult

	// 确保目录路径以 / 结尾
	if !strings.HasSuffix(dirPath, string(os.PathSeparator)) {
		dirPath += string(os.PathSeparator)
	}

	// 确保对象前缀不以 / 开头，但以 / 结尾（如果不为空）
	if objectPrefix != "" {
		objectPrefix = strings.TrimPrefix(objectPrefix, "/")
		if !strings.HasSuffix(objectPrefix, "/") {
			objectPrefix += "/"
		}
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身
		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("计算相对路径失败 %s: %v", path, err)
		}

		// 构建对象名称
		objectName := objectPrefix + strings.ReplaceAll(relPath, "\\", "/")

		// 上传文件
		result, err := u.UploadFile(path, objectName, expirationDays)
		if err != nil {
			return err
		}

		results = append(results, result)
		fmt.Printf("已上传: %s -> %s\n", path, result.URL)

		return nil
	})

	if err != nil {
		return results, fmt.Errorf("上传目录 %s 失败: %v", dirPath, err)
	}

	return results, nil
}

// ListObjects 列出存储桶中的对象
func (u *Uploader) ListObjects(prefix string) ([]string, error) {
	var objects []string

	it := u.client.Bucket(u.bucketName).Objects(u.ctx, &storage.Query{
		Prefix: prefix,
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return objects, fmt.Errorf("列出对象失败: %v", err)
		}
		objects = append(objects, attrs.Name)
	}

	return objects, nil
}

// DeleteObject 删除存储桶中的对象
func (u *Uploader) DeleteObject(objectName string) error {
	obj := u.client.Bucket(u.bucketName).Object(objectName)
	if err := obj.Delete(u.ctx); err != nil {
		return fmt.Errorf("删除对象 %s 失败: %v", objectName, err)
	}
	return nil
}

// GetObjectMetadata 获取对象的元数据
func (u *Uploader) GetObjectMetadata(objectName string) (*storage.ObjectAttrs, error) {
	attrs, err := u.client.Bucket(u.bucketName).Object(objectName).Attrs(u.ctx)
	if err != nil {
		return nil, fmt.Errorf("获取对象 %s 的元数据失败: %v", objectName, err)
	}
	return attrs, nil
}

// SetObjectExpiration 设置对象的过期时间
func (u *Uploader) SetObjectExpiration(objectName string, expirationDays int) error {
	if expirationDays <= 0 {
		return fmt.Errorf("过期天数必须大于 0")
	}

	expirationTime := time.Now().AddDate(0, 0, expirationDays)
	objAttrs := storage.ObjectAttrsToUpdate{
		Metadata: map[string]string{
			"expiration": expirationTime.Format(time.RFC3339),
		},
	}

	obj := u.client.Bucket(u.bucketName).Object(objectName)
	if _, err := obj.Update(u.ctx, objAttrs); err != nil {
		return fmt.Errorf("设置对象 %s 过期时间失败: %v", objectName, err)
	}

	return nil
}
