package assembler

var (
	DrawerAss     DrawerAssembler
	ConnectionAss ConnectionAssembler
)

func init() {
	DrawerAss = NewDrawerAssembler()
	ConnectionAss = NewConnectionAssembler()
}
