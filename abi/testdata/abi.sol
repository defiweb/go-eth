pragma solidity >=0.7.0 <0.9.0;

contract Test {
    struct Struct {
        bytes32 A;
        bytes32 B;
    }
    
    constructor(uint256 a) { }
    event EventA(uint256 indexed a, uint256 b);
    event EventB(uint256 indexed a, uint256 b) anonymous;
    error ErrorA(uint256 a, uint256 b);
    function Foo(uint256 a) public returns (uint256) { return 0; }
    function Bar(Struct[2][2] memory a) public returns (uint256[2][2] memory) { return [[0, 0], [0, 0]]; }
    fallback() external { }
    receive() payable external { }
}
