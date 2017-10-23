{
  "language": "Solidity",
  "sources": {
    {{range $index, $file := .}}
    {{if $index}},{{end}}
    "{{basename $file}}": {
      "urls": [
        "{{$file}}"
      ]
    }
    {{end}}
  },
  "settings": {
    "outputSelection": {
      "*": {
        "*": [
          "abi",
          "ast",
          "evm.bytecode.object",
          "evm.bytecode.sourceMap",
          "evm.bytecode.linkReferences",
          "evm.deployedBytecode.object",
          "evm.deployedBytecode.sourceMap",
          "evm.deployedBytecode.linkReferences"
        ]
      }
    }
  }
}
