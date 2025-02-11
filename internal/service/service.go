package service

import "mime/multipart"

type StorageService interface {
	SaveFile(attachment *multipart.FileHeader) (string, error)
}
