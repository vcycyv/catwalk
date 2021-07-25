package assembler

var (
	DrawerAss     DrawerAssembler
	ConnectionAss ConnectionAssembler
	DataSourceAss DataSourceAssembler
)

func init() {
	DrawerAss = NewDrawerAssembler()
	ConnectionAss = NewConnectionAssembler()
	DataSourceAss = NewDataSourceAssebler()
}
