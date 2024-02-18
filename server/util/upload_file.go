package util

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"server/model"
	"time"
)

var ctx = context.Background()

func UploadFile(file multipart.File, object string, c *model.ClientUploader) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	c.Cl.Bucket(c.BucketName).Object(c.UploadPath + object).NewWriter(ctx)
	
	wc := c.Cl.Bucket(c.BucketName).Object(c.UploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}