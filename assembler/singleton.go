package assembler

var (
	DrawerAss     DrawerAssembler
	ConnectionAss ConnectionAssembler
	DataSourceAss DataSourceAssembler
	ServerAss     ServerAssembler
	ModelFileAss  ModelFileAssembler
	ModelAss      ModelAssembler
)

func init() {
	DrawerAss = NewDrawerAssembler()
	ConnectionAss = NewConnectionAssembler()
	DataSourceAss = NewDataSourceAssebler()
	ServerAss = NewServerAssembler()
	ModelFileAss = NewModelFileAssembler()
	ModelFileAss = ModelFileAssembler(NewModelAssembler())
}
