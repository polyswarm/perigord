// Contains directives to generate required go source, such as embedding
// templates

//go:generate go-bindata -nomemcopy -prefix templates -pkg templates -ignore templates/*.go -o templates/bindata.go templates/...
//go:generate abigen --sol migration/bindings/Migrations.sol --pkg bindings --out migration/bindings/Migrations.go

package perigord
