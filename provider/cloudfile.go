package provider

import "context"
import "io"

type ICloudFile interface {
	GetLink(remotePath string) string
	Put(ctx context.Context, remotePath string, local io.Reader, size int64) error
	PutFile(ctx context.Context, remotePath string, reader string) error
	DeleteFile(remotePath string) error
}