package chart

import (
	"math"
	"testing"
)

func TestSankeyColumnAssignment(t *testing.T) {
	// A → B → D
	// A → C → D
	nodes := make([]sankeyLayoutNode, 4)
	links := []SankeyLink{
		{Source: 0, Target: 1, Value: 10},
		{Source: 0, Target: 2, Value: 5},
		{Source: 1, Target: 3, Value: 10},
		{Source: 2, Target: 3, Value: 5},
	}

	numCols := sankeyAssignColumns(nodes, links)
	if numCols != 3 {
		t.Fatalf("numCols = %d, want 3", numCols)
	}
	if nodes[0].Col != 0 {
		t.Errorf("node 0 col = %d, want 0", nodes[0].Col)
	}
	if nodes[1].Col != 1 {
		t.Errorf("node 1 col = %d, want 1", nodes[1].Col)
	}
	if nodes[2].Col != 1 {
		t.Errorf("node 2 col = %d, want 1", nodes[2].Col)
	}
	if nodes[3].Col != 2 {
		t.Errorf("node 3 col = %d, want 2", nodes[3].Col)
	}
}

func TestSankeyVerticalLayout(t *testing.T) {
	nodes := make([]sankeyLayoutNode, 3)
	links := []SankeyLink{
		{Source: 0, Target: 1, Value: 20},
		{Source: 0, Target: 2, Value: 10},
	}

	sankeyAssignColumns(nodes, links)
	sankeyComputeFlows(nodes, links)

	// Node 0: outflow=30, col=0.
	// Node 1: inflow=20, col=1.
	// Node 2: inflow=10, col=1.
	if nodes[0].OutFlow != 30 {
		t.Errorf("node 0 outflow = %g, want 30", nodes[0].OutFlow)
	}
	if nodes[1].InFlow != 20 {
		t.Errorf("node 1 inflow = %g, want 20", nodes[1].InFlow)
	}
	if nodes[2].InFlow != 10 {
		t.Errorf("node 2 inflow = %g, want 10", nodes[2].InFlow)
	}

	sankeyLayoutNodes(nodes, 2, 20, 10,
		0, 200, 0, 300)

	// Node 1 should be taller than node 2 (2:1 ratio).
	if nodes[1].H <= nodes[2].H {
		t.Errorf("node 1 height %g should exceed node 2 height %g",
			nodes[1].H, nodes[2].H)
	}

	ratio := float64(nodes[1].H / nodes[2].H)
	if math.Abs(ratio-2.0) > 0.1 {
		t.Errorf("height ratio = %g, want ~2.0", ratio)
	}
}

func TestSankeyHitTestNode(t *testing.T) {
	sv := &sankeyView{
		nodes: []sankeyLayoutNode{
			{X: 10, Y: 10, W: 20, H: 50, Index: 0},
			{X: 170, Y: 20, W: 20, H: 40, Index: 1},
		},
	}

	kind, idx, ok := sv.hitTest(20, 30)
	if !ok || kind != "node" || idx != 0 {
		t.Errorf("hitTest(20,30) = (%s,%d,%v), want (node,0,true)",
			kind, idx, ok)
	}

	kind, idx, ok = sv.hitTest(180, 40)
	if !ok || kind != "node" || idx != 1 {
		t.Errorf("hitTest(180,40) = (%s,%d,%v), want (node,1,true)",
			kind, idx, ok)
	}
}

func TestSankeyHitTestLink(t *testing.T) {
	// Build a simple rectangular ribbon polygon.
	poly := []float32{
		50, 20, // top-left
		150, 20, // top-right
		150, 40, // bottom-right
		50, 40, // bottom-left
	}
	sv := &sankeyView{
		links: []sankeyLayoutLink{
			{Index: 0, Poly: poly},
		},
	}

	kind, idx, ok := sv.hitTest(100, 30)
	if !ok || kind != "link" || idx != 0 {
		t.Errorf("hitTest(100,30) = (%s,%d,%v), want (link,0,true)",
			kind, idx, ok)
	}
}

func TestSankeyHitTestMiss(t *testing.T) {
	sv := &sankeyView{
		nodes: []sankeyLayoutNode{
			{X: 10, Y: 10, W: 20, H: 50, Index: 0},
		},
		links: []sankeyLayoutLink{
			{Index: 0, Poly: []float32{
				50, 20, 150, 20, 150, 40, 50, 40,
			}},
		},
	}

	_, _, ok := sv.hitTest(0, 0)
	if ok {
		t.Error("hitTest(0,0) should miss")
	}

	_, _, ok = sv.hitTest(100, 100)
	if ok {
		t.Error("hitTest(100,100) should miss")
	}
}

func TestSankeyNaNInfSkipped(t *testing.T) {
	nodes := make([]sankeyLayoutNode, 3)
	links := []SankeyLink{
		{Source: 0, Target: 1, Value: 10},
		{Source: 0, Target: 2, Value: math.NaN()},
		{Source: 0, Target: 2, Value: math.Inf(1)},
		{Source: 0, Target: 2, Value: -5},
	}

	sankeyAssignColumns(nodes, links)
	sankeyComputeFlows(nodes, links)

	if nodes[0].OutFlow != 10 {
		t.Errorf("node 0 outflow = %g, want 10", nodes[0].OutFlow)
	}
	if nodes[2].InFlow != 0 {
		t.Errorf("node 2 inflow = %g, want 0", nodes[2].InFlow)
	}
}

func TestSankeySelfLoopSkipped(t *testing.T) {
	nodes := make([]sankeyLayoutNode, 2)
	links := []SankeyLink{
		{Source: 0, Target: 1, Value: 10},
		{Source: 0, Target: 0, Value: 99}, // self-loop
	}
	sankeyAssignColumns(nodes, links)
	sankeyComputeFlows(nodes, links)

	// Self-loop should not contribute to flow.
	if nodes[0].OutFlow != 10 {
		t.Errorf("node 0 outflow = %g, want 10", nodes[0].OutFlow)
	}
}

func TestSankeyValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     SankeyCfg
		wantErr bool
	}{
		{"empty nodes",
			SankeyCfg{}, true},
		{"empty links",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}},
			}, true},
		{"negative NodeWidth",
			SankeyCfg{
				Nodes:     []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links:     []SankeyLink{{Source: 0, Target: 1, Value: 1}},
				NodeWidth: -1,
			}, true},
		{"negative NodeGap",
			SankeyCfg{
				Nodes:   []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links:   []SankeyLink{{Source: 0, Target: 1, Value: 1}},
				NodeGap: -1,
			}, true},
		{"negative link value",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links: []SankeyLink{{Source: 0, Target: 1, Value: -5}},
			}, true},
		{"link index out of range",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}},
				Links: []SankeyLink{{Source: 0, Target: 5, Value: 10}},
			}, true},
		{"self-loop",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links: []SankeyLink{{Source: 0, Target: 0, Value: 10}},
			}, true},
		{"cycle detected",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links: []SankeyLink{
					{Source: 0, Target: 1, Value: 10},
					{Source: 1, Target: 0, Value: 5},
				},
			}, true},
		{"valid",
			SankeyCfg{
				Nodes: []SankeyNode{{Label: "A"}, {Label: "B"}},
				Links: []SankeyLink{{Source: 0, Target: 1, Value: 10}},
			}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v",
					err, tt.wantErr)
			}
		})
	}
}

func TestPointInPolygon(t *testing.T) {
	// Simple square: (0,0) (10,0) (10,10) (0,10).
	poly := []float32{0, 0, 10, 0, 10, 10, 0, 10}

	if !pointInPolygon(5, 5, poly) {
		t.Error("center should be inside")
	}
	if pointInPolygon(15, 5, poly) {
		t.Error("outside right should miss")
	}
	if pointInPolygon(5, 15, poly) {
		t.Error("below should miss")
	}
}

func TestHasCycle(t *testing.T) {
	// DAG: 0→1→2
	if hasCycle(3, []SankeyLink{
		{Source: 0, Target: 1, Value: 1},
		{Source: 1, Target: 2, Value: 1},
	}) {
		t.Error("DAG incorrectly detected as cyclic")
	}

	// Cycle: 0→1→0
	if !hasCycle(2, []SankeyLink{
		{Source: 0, Target: 1, Value: 1},
		{Source: 1, Target: 0, Value: 1},
	}) {
		t.Error("cycle not detected")
	}
}
