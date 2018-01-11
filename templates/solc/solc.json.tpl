{
  "language": "Solidity",
  "sources": {
    {{range $index, $match := .}}
    {{if $index}},{{end}}
    "{{$match.Filename}}": {
      "content": {{$match.Content}}
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
