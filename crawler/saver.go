package crawler

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Saver is an interface for saving page graph
type Saver interface {
	// Save - saves a page graph to the file
	Save(filename string, nodes []*PageNode, edges []*PageEdge) (err error)

	// GetExtension - returns file extension for the saver
	GetExtension() string
}

// GEXFSaver - saver interface implementation for GEXF graph format
type GEXFSaver struct {
}

// GetExtension - implementation of Saver
func (s GEXFSaver) GetExtension() string {
	return ".gexf"
}

// Save - Saver interface implementation for GEXFSaver
func (s *GEXFSaver) Save(filename string, nodes []*PageNode, edges []*PageEdge) (err error) {
	file, err := os.Create(filename)

	if err != nil {
		return
	}

	writer := io.Writer(file)

	document := `<?xml version="1.0" encoding="UTF-8"?>
	<gexf xmlns="http://www.gexf.net/1.2draft" version="1.2">
		<meta lastmodifieddate="%s">
			<creator>WebCrawler</creator>
			<description>WebCrawler web graph xml</description>
		</meta>
		<graph mode="static" defaultedgetype="directed">
			<nodes>
%s
			</nodes>
			<edges>
%s
			</edges>
		</graph>
	</gexf>`

	document = fmt.Sprintf(document, time.Now().Format("2006-01-02"),
		s.getNodesString(nodes), s.getEdgesString(edges))

	writer.Write([]byte(document))

	return nil
}

func (s *GEXFSaver) getNodesString(nodes []*PageNode) (result string) {
	for _, node := range nodes {
		result += fmt.Sprintf("\t\t\t\t<node id='%s' label='%s' />\n", node.url.String(), node.url.String())
	}

	return
}

func (s *GEXFSaver) getEdgesString(edges []*PageEdge) (result string) {
	for i, edge := range edges {
		result += fmt.Sprintf("\t\t\t\t<edge id='%d' source='%s' target='%s' />\n",
			i, edge.source.url.String(), edge.destination.url.String())
	}

	return
}
