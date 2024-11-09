package graphdb

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nedvisol/go-connectdots/config"
	"github.com/nedvisol/go-connectdots/util"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.uber.org/fx"
)

type Neo4jGraphService struct {
	ctx    context.Context
	config config.GraphDbConfig
	driver neo4j.DriverWithContext
}

func (n *Neo4jGraphService) getSession(ctx context.Context) neo4j.SessionWithContext {
	return n.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
}

// CreateEdge implements GraphDbService.
func (n *Neo4jGraphService) UpdateEdge(edge *EdgeInfo, allowUpsert bool) error {
	queryAttrs := make([]string, 0, len(*edge.Attrs)+1)
	for key := range *edge.Attrs {
		queryAttrs = append(queryAttrs, fmt.Sprintf("edge.%s = $%s", key, key))
	}

	mergeOrMatch := util.Ternary(allowUpsert, "MERGE", "MATCH")

	query := fmt.Sprintf(`
	MATCH (left:%s { _id: $left_id })
	MATCH (right:%s { _id: $right_id })
	%s (left)-[edge:%s {_id: $edge_id}]->(right)
	SET %s
	RETURN edge._id
	`,
		edge.Left.Label,
		edge.Right.Label,
		mergeOrMatch,
		edge.Label,
		strings.Join(queryAttrs, ","))

	(*edge.Attrs)["left_id"] = edge.Left.Id
	(*edge.Attrs)["right_id"] = edge.Right.Id
	(*edge.Attrs)["edge_id"] = edge.Id

	//fmt.Printf("executing query %s\n", query)

	// Execute the query inside a transaction
	session := n.getSession(n.ctx)
	defer session.Close(n.ctx)
	records, err := session.Run(n.ctx, query, *edge.Attrs)
	if err != nil {
		return err
	}

	if records.Next(n.ctx) {
		id, found := records.Record().Get("edge._id")
		if !found {
			fmt.Printf("error updating edge %s\n", id)
			panic("unable to update edge")
		}
		//fmt.Printf("Updated node %s\n", id)
		return nil
	}

	return nil
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
	session := n.getSession(n.ctx)
	defer session.Close(n.ctx)
	records, err := session.Run(n.ctx, query, *node.Attrs)
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
func (n *Neo4jGraphService) UpdateNode(node *NodeInfo, allowUpsert bool) error {
	queryAttrs := make([]string, 0, len(*node.Attrs)+1)
	for key := range *node.Attrs {
		queryAttrs = append(queryAttrs, fmt.Sprintf("node.%s = $%s", key, key))
	}

	mergeOrMatch := util.Ternary(allowUpsert, "MERGE", "MATCH")

	query := fmt.Sprintf(`
	%s (node: %s {_id : $_id})
	SET %s
	RETURN node._id
	`, mergeOrMatch, node.Label, strings.Join(queryAttrs, ","))
	(*node.Attrs)["_id"] = node.Id

	//fmt.Printf("executing query %s\n", query)

	// Execute the query inside a transaction
	session := n.getSession(n.ctx)
	defer session.Close(n.ctx)
	records, err := session.Run(n.ctx, query, *node.Attrs)
	if err != nil {
		return err
	}

	if records.Next(n.ctx) {
		id, found := records.Record().Get("node._id")
		if !found {
			fmt.Printf("error updating node %s\n", id)
			panic("unable to update node")
		}
		//fmt.Printf("Updated node %s\n", id)
		return nil
	}

	return nil
}

func NewNeo4jGraphService(lifecycle fx.Lifecycle, ctx context.Context, cfg *config.Config) GraphDbService {
	graphcfg := cfg.GraphDb
	driver, err := neo4j.NewDriverWithContext(graphcfg.Uri, neo4j.BasicAuth(graphcfg.Username, graphcfg.Password, ""))
	if err != nil {
		log.Fatal("Error creating Neo4j driver: ", err)
	}
	//defer driver.Close()

	// Start a new session
	//session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	//defer session.Close()

	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			driver.Close(ctx)
			fmt.Println("Application is stopping. closed neo4j driver and session")
			return nil
		},
	})

	return &Neo4jGraphService{
		ctx:    ctx,
		config: *cfg.GraphDb,
		driver: driver,
	}
}
