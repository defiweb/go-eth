pragma solidity >=0.7.0 <0.9.0;

contract Test {
    // Enum
    enum Status {
        Inactive,
        Active,
        Paused
    }
    
    // Struct
    struct Struct {
        bytes32 A;
        bytes32 B;
        Status status;
    }
    
    // Custom Type
    type CustomUint is uint256;
    
    // Events
    event EventA(uint256 indexed a, string b);
    event EventB(uint256 indexed a, string indexed b);
    event EventC(uint256 indexed a, string b) anonymous;
    
    // Error
    error ErrorA(uint256 a, uint256 b);
    
    // Public Variable
    Struct public structField;
    
    // Mapping
    mapping(address => Struct) public structsMapping;
    
    // Array
    Struct[] public structsArray;
    
    // Constructor
    constructor(CustomUint a) {}
    
    // Functions
    function Foo(CustomUint a) public returns (CustomUint) { return a; }
    
    function Bar(Struct[2][2] memory a) public returns (uint8[2][2] memory) { return [[0, 0], [0, 0]]; }
    
    // Fallback and Receive functions
    fallback() external {}
    
    receive() payable external {}
}
