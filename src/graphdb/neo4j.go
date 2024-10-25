package graphdb

import (
	"context"
	"fmt"
	"log"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/fx"
)

type Neo4jGraphService struct {
	config config.GraphDbConfig
}

// CreateEdge implements GraphDbService.
func (n *Neo4jGraphService) CreateEdge(edge *EdgeInfo) error {
	panic("unimplemented")
}

// CreateNode implements GraphDbService.
func (n *Neo4jGraphService) CreateNode(node *NodeInfo) error {
	panic("unimplemented")
}

// DeleteNode implements GraphDbService.
func (n *Neo4jGraphService) DeleteNode(node *NodeInfo) error {
	panic("unimplemented")
}

// UpdateNode implements GraphDbService.
func (n *Neo4jGraphService) UpdateNode(node *NodeInfo) error {
	panic("unimplemented")
}

func NewNeo4jGraphService(lifecycle fx.Lifecycle, ctx context.Context, cfg *config.Config) GraphDbService {
	graphcfg := cfg.GraphDb
	driver, err := neo4j.NewDriver(graphcfg.Uri, neo4j.BasicAuth(graphcfg.Username, graphcfg.Password, ""))
	if err != nil {
		log.Fatal("Error creating Neo4j driver: ", err)
	}
	//defer driver.Close()

	// Start a new session
	session, _ := driver.NewSession(neo4j.SessionConfig{})
	//defer session.Close()

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			driver.Close()
			session.Close()
			fmt.Println("Application is stopping. closed neo4j driver and session")
			return nil
		},
	})

	return &Neo4jGraphService{
		config: *cfg.GraphDb,
	}
}
