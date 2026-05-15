package repository

import "log_parser3000/internal/domain"

type Repository interface {
	SaveParsedLog(parsed *domain.ParsedLog) (int, error)
	SaveLogError(filename string, errMsg string) (int, error)

	GetLogByID(id int) (*domain.Log, error)
	GetNodeByID(id int) (*domain.Node, error)
	GetNodesByLogID(logID int) ([]domain.Node, error)
	GetPortsByNodeID(nodeID int) ([]domain.Port, error)
	GetTopologyByLogID(logID int) ([]domain.TopologyGroup, error)
}
