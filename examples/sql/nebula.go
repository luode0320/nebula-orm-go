package sql

// NebulaSpaceName NebulaVertexName NebulaEdgeName 统一定义名称
var NebulaSpaceName = "space_luode"
var NebulaVertexName = "test_vertex"
var NebulaEdgeName = "test_edge"

const (
	// NebulaInitSql 初始化sql
	NebulaInitSql = `
		create space if not exists space_luode(vid_type=fixed_string(64));

		use space_luode;

		create tag if not exists test_vertex(
			chain_key string,
			parent_key string
		);

		create edge if not exists test_edge(test string);
	`
)
