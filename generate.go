// Contains directives to generate requried go source, such as embedding
// templates

//go:generate go-bindata -nomemcopy -prefix templates -pkg templates -ignore templates/*.go -o templates/bindata.go templates/...

package perigord
