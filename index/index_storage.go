package index

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	storagego "github.com/supabase-community/storage-go"
)

type DocumentId struct {
	UserId     string
	Collection uuid.UUID
	DocId      uuid.UUID
	Filename   string
}

func (id DocumentId) path() string {
	return fmt.Sprintf("%s/%s/%s.pdf", id.UserId, id.Collection, id.DocId)
}

const bucket = "documents"

func (index Index) Upload(doc DocumentId, data []byte) error {

	index.Storage.CreateBucket(bucket, storagego.BucketOptions{
		Public:        false,
		FileSizeLimit: "50mb",
	})

	resp := index.Storage.UploadFile(bucket, doc.path(), bytes.NewReader(data))
	if resp.Error != "" {
		return fmt.Errorf("could not upload file: %v", resp.Error)
	}

	return nil
}

func (index Index) Download(doc DocumentId) ([]byte, error) {
	return index.Storage.DownloadFile(bucket, doc.path())
}

func (index Index) Delete(doc DocumentId) error {
	resp := index.Storage.RemoveFile(bucket, []string{doc.path()})
	if resp.Error != "" {
		return fmt.Errorf("could not delete file: %v", resp.Error)
	}

	return nil
}
