package infrastructure

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/vcycyv/catwalk/domain"
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

func (f *fileService) Save(fileName string, reader io.Reader) (string, error) {
	logrus.Debugf("about save %s to mongodb", fileName)
	objectID, err := f.bucket.UploadFromStream(fileName, reader)
	if err != nil {
		logrus.Error("failed to upload content")
		return "", err
	}
	logrus.Debugf("saving %s to mongodb succeeded", fileName)

	return objectID.Hex(), nil
}

func (f *fileService) DirectContentToWriter(fileID string, writer io.Writer) error {
	id, _ := primitive.ObjectIDFromHex(fileID)
	_, err := f.bucket.DownloadToStream(id, writer)
	if err != nil {
		logrus.Error("failed to get file content: " + fileID)
	}
	return err
}

func (f *fileService) GetContent(fileID string) (io.Reader, error) {
	id, _ := primitive.ObjectIDFromHex(fileID)
	rtnVal, err := f.bucket.OpenDownloadStream(id)
	if err != nil {
		logrus.Error("failed to open download stream: " + fileID)
		return nil, err
	}
	return rtnVal, nil
}

func (f *fileService) Delete(fileID string) error {
	id, _ := primitive.ObjectIDFromHex(fileID)
	return f.bucket.Delete(id)
}
