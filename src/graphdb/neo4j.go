package graphdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/fx"
)

type Neo4jGraphService struct {
	ctx     context.Context
	config  config.GraphDbConfig
	session neo4j.SessionWithContext
}

// CreateEdge implements GraphDbService.
func (n *Neo4jGraphService) CreateEdge(edge *EdgeInfo) error {
	panic("unimplemented")
}

// func cloneMap(source *map[string]interface{}) *map[string]interface{} {
// 	clone := make(map[string]interface{})
// 	for key, value := range *source {
// 		clone[key] = value
// 	}

// 	return &clone
// }

// CreateNode implements GraphDbService.
func (n *Neo4jGraphService) CreateNode(node *NodeInfo) error {
	// Define the Cypher query to create a node
	// 	query := `
	// CREATE (p:Person {name: $name, age: $age})
	// RETURN p.name, p.age
	// `
	// 	params := map[string]interface{}{
	// 		"name": "John Doe",
	// 		"age":  30,
	// 	}

	queryAttrs := make([]string, 0, len(*node.Attrs)+1)
	for key := range *node.Attrs {
		queryAttrs = append(queryAttrs, fmt.Sprintf("%s: $%s", key, key))
	}
	queryAttrs = append(queryAttrs, "_id: $_id")

	query := fmt.Sprintf(`
	CREATE (node: %s {%s} )
	RETURN node._id
	`, node.Label, strings.Join(queryAttrs, ","))
	(*node.Attrs)["_id"] = node.Id

	fmt.Printf("executing query %s\n", query)

	// Execute the query inside a transaction
	records, err := n.session.Run(n.ctx, query, *node.Attrs)
	if err != nil {
		return err
	}

	if records.Next(n.ctx) {
		id, found := records.Record().Get("node._id")
		if !found {
			panic("unable to create node")
		}
		fmt.Printf("Created node %s", id)
		return nil
	}

	return nil
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
	driver, err := neo4j.NewDriverWithContext(graphcfg.Uri, neo4j.BasicAuth(graphcfg.Username, graphcfg.Password, ""))
	if err != nil {
		log.Fatal("Error creating Neo4j driver: ", err)
	}
	//defer driver.Close()

	// Start a new session
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	//defer session.Close()

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			driver.Close(ctx)
			session.Close(ctx)
			fmt.Println("Application is stopping. closed neo4j driver and session")
			return nil
		},
	})

	return &Neo4jGraphService{
		ctx:     ctx,
		config:  *cfg.GraphDb,
		session: session,
	}
}
