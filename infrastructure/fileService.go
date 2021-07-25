package infrastructure

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/vcycyv/blog/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

type fileService struct {
	bucket gridfs.Bucket
}

func NewFileService(bucket gridfs.Bucket) domain.FileService {
	return &fileService{
		bucket: bucket,
	}
}

func (f *fileService) Save(reader io.Reader, fileName string) (string, error) {
	logrus.Debugf("about save %s to mongodb", fileName)
	objectID, err := f.bucket.UploadFromStream(fileName, reader)
	if err != nil {
		logrus.Error("failed to upload content")
		return "", err
	}
	logrus.Debugf("saving %s to mongodb succeeded", fileName)

	return objectID.Hex(), nil
}

func (f *fileService) GetContent(fileID string, writer io.Writer) error {
	id, _ := primitive.ObjectIDFromHex(fileID)
	_, err := f.bucket.DownloadToStream(id, writer)
	return err
}

func (f *fileService) Delete(fileID string) error {
	id, _ := primitive.ObjectIDFromHex(fileID)
	return f.bucket.Delete(id)
}
