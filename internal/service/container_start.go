package service

import (
	"fmt"

	"algo-container-manager/internal/common"
	"algo-container-manager/internal/db"
	"algo-container-manager/internal/k8s"
	"algo-container-manager/internal/model"

	"github.com/google/uuid"
)

func (s *ContainerService) Start(req model.StartAlgorithmRequest) error {
	uid := uuid.New().String()

	s.prepareStartRequest(&req, uid)

	if err := k8s.CreateDeployment(s.clientset, req); err != nil {
		return fmt.Errorf("create deployment failed: %w", err)
	}

	if err := k8s.CreateService(s.clientset, req); err != nil {
		return fmt.Errorf("create service failed: %w", err)
	}

	record := s.buildDeployRecord(req, uid)

	if err := s.saveDeployRecord(record); err != nil {
		_ = k8s.DeleteAlgorithm(s.clientset, req.Namespace, req.DeploymentName)
		return fmt.Errorf("save deploy record failed: %w", err)
	}

	return nil
}

func (s *ContainerService) prepareStartRequest(req *model.StartAlgorithmRequest, uid string) {
	if req.Namespace == "" {
		req.Namespace = "default"
	}
	if req.DeploymentName == "" {
		req.DeploymentName = common.BuildDeploymentName(req.Name, req.Version, uid)
	}
	if req.ServiceName == "" {
		req.ServiceName = common.BuildServiceName(req.DeploymentName)
	}
	if req.Port == 0 {
		req.Port = 80
	}
	if req.CPU == "" {
		req.CPU = "500m"
	}
	if req.Memory == "" {
		req.Memory = "512Mi"
	}
	if req.Replicas == 0 {
		req.Replicas = 1
	}
}

func (s *ContainerService) buildDeployRecord(req model.StartAlgorithmRequest, uid string) model.DeployRecord {
	return model.DeployRecord{
		UUID:               uid,
		AlgorithmVersionID: req.AlgorithmVersionID,
		ImageID:            req.ImageID,
		Namespace:          req.Namespace,
		K8sDeploymentName:  req.DeploymentName,
		K8sServiceName:     req.ServiceName,
		DeployStatus:       "deploying",
		AccessEndpoint:     "",
	}
}

func (s *ContainerService) saveDeployRecord(record model.DeployRecord) error {
	return db.DB.Create(&record).Error
}
