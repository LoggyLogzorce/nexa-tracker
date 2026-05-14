package attachment

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"nexa-task-tracker/internal/core/task"
	"nexa-task-tracker/internal/core/user"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Upload(ctx context.Context, taskID uint, userID uuid.UUID, filename string, file io.Reader) (*Attachment, error)
	GetByID(ctx context.Context, id, taskID uint) (*Attachment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]AttachmentResponse, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]AttachmentResponse, error)
	Delete(ctx context.Context, id, taskID uint, userID uuid.UUID) error
}

type service struct {
	repo       Repository
	taskRepo   task.Repository
	userRepo   user.Repository
	uploadPath string
}

func NewService(repo Repository, taskRepo task.Repository, userRepo user.Repository, uploadPath string) Service {
	return &service{
		repo:       repo,
		taskRepo:   taskRepo,
		userRepo:   userRepo,
		uploadPath: uploadPath,
	}
}

func (s *service) Upload(ctx context.Context, taskID uint, userID uuid.UUID, filename string, file io.Reader) (*Attachment, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fileUID := uuid.New().String()
	ext := filepath.Ext(filename)
	storedName := fileUID + ext

	taskDir := filepath.Join(s.uploadPath, fmt.Sprintf("%d", taskID))
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		log.Printf("Failed to create task directory: %v", err)
		return nil, ErrCreateDirectory
	}

	filePath := filepath.Join(taskDir, storedName)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return nil, ErrWriteFile
	}
	defer dst.Close()

	// читаем первые 512 байт для MIME, не прерывая поток
	buf := make([]byte, 512)
	n, readErr := file.Read(buf)
	if readErr != nil && readErr != io.EOF {
		os.Remove(filePath)
		log.Printf("Failed to read file header: %v", readErr)
		return nil, ErrWriteFile
	}
	mimeType := detectMimeType(buf[:n])

	// пишем буфер + остаток потока
	_, err = dst.Write(buf[:n])
	if err != nil {
		os.Remove(filePath)
		log.Printf("Failed to write file header: %v", err)
		return nil, ErrWriteFile
	}
	written, err := io.Copy(dst, file)
	if err != nil {
		os.Remove(filePath)
		log.Printf("Failed to write file: %v", err)
		return nil, ErrWriteFile
	}
	fileSize := int64(n) + written

	attachment := &Attachment{
		TaskID:   taskID,
		UserID:   userID,
		Filename: filename,
		FilePath: filePath,
		FileSize: fileSize,
		MimeType: &mimeType,
	}

	if err := s.repo.Create(ctxT, attachment); err != nil {
		os.Remove(filePath)
		log.Printf("Failed to create attachment record: %v", err)
		return nil, ErrCreateAttachmentRecord
	}

	return attachment, nil
}

func (s *service) GetByID(ctx context.Context, id, taskID uint) (*Attachment, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAttachmentNotFound
		}
		return nil, err
	}

	if a.TaskID != taskID {
		return nil, ErrAttachmentNotFound
	}

	return a, nil
}

func (s *service) GetByTaskID(ctx context.Context, taskID uint) ([]AttachmentResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	attachments, err := s.repo.GetByTaskID(ctxT, taskID)
	if err != nil {
		return nil, err
	}

	if len(attachments) == 0 {
		return []AttachmentResponse{}, nil
	}

	userIDsMap := make(map[uuid.UUID]struct{})
	for _, a := range attachments {
		userIDsMap[a.UserID] = struct{}{}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]user.User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	response := make([]AttachmentResponse, len(attachments))
	for i, a := range attachments {
		response[i] = AttachmentResponse{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			TaskID:    a.TaskID,
			Filename:  a.Filename,
			FileSize:  a.FileSize,
			MimeType:  a.MimeType,
		}
		if u, ok := usersMap[a.UserID]; ok {
			response[i].User = AttachmentUserResponse{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
			}
		}
	}

	return response, nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]AttachmentResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tasks, err := s.taskRepo.GetByProjectID(ctxT, projectID, false)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []AttachmentResponse{}, nil
	}

	taskIDs := make([]uint, len(tasks))
	for i, t := range tasks {
		taskIDs[i] = t.ID
	}

	attachments, err := s.repo.GetByTaskIDs(ctxT, taskIDs)
	if err != nil {
		return nil, err
	}

	if len(attachments) == 0 {
		return []AttachmentResponse{}, nil
	}

	userIDsMap := make(map[uuid.UUID]struct{})
	for _, a := range attachments {
		userIDsMap[a.UserID] = struct{}{}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]user.User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	response := make([]AttachmentResponse, len(attachments))
	for i, a := range attachments {
		response[i] = AttachmentResponse{
			ID:        a.ID,
			CreatedAt: a.CreatedAt,
			TaskID:    a.TaskID,
			Filename:  a.Filename,
			FileSize:  a.FileSize,
			MimeType:  a.MimeType,
		}
		if u, ok := usersMap[a.UserID]; ok {
			response[i].User = AttachmentUserResponse{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
			}
		}
	}

	return response, nil
}

func (s *service) Delete(ctx context.Context, id, taskID uint, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAttachmentNotFound
		}
		return err
	}

	if a.TaskID != taskID {
		return ErrAttachmentNotFound
	}

	if a.UserID != userID {
		return ErrNotAttachmentOwner
	}

	if err := os.Remove(a.FilePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove file: %v", err)
		return ErrRemoveFile
	}

	return s.repo.Delete(ctxT, id)
}

func detectMimeType(data []byte) string {
	if len(data) < 512 {
		return "application/octet-stream"
	}
	return http.DetectContentType(data)
}
