package service

import (
	"log_parser3000/internal/domain"
	"log_parser3000/internal/repository"
)

type Service struct {
	repo   repository.Repository
	parser interface {
		ParseFromDataPath(relPath string) (int, error)
	}
}

func New(repo repository.Repository, parser interface {
	ParseFromDataPath(relPath string) (int, error)
}) *Service {
	return &Service{repo: repo, parser: parser}
}

func (s *Service) ParseLog(archivePath string) (int, error) {
	return s.parser.ParseFromDataPath(archivePath)
}

func (s *Service) GetLog(id int) (*domain.Log, error) {
	return s.repo.GetLogByID(id)
}

func (s *Service) GetNode(id int) (*domain.Node, error) {
	return s.repo.GetNodeByID(id)
}

func (s *Service) GetPorts(nodeID int) ([]domain.Port, error) {
	return s.repo.GetPortsByNodeID(nodeID)
}

func (s *Service) GetTopology(logID int) (*domain.TopologyResponse, error) {
	groups, err := s.repo.GetTopologyByLogID(logID)
	if err != nil {
		return nil, err
	}
	nodes, err := s.repo.GetNodesByLogID(logID)
	if err != nil {
		return nil, err
	}
	return &domain.TopologyResponse{
		Nodes:  nodes,
		Groups: groups,
	}, nil
}
