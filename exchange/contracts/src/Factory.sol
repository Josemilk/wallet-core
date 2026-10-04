// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Pool} from "./Pool.sol";

contract Factory is Ownable {
    uint256 public constant DEFAULT_FEE_BPS = 30;
    address public router;
    uint256 public immutable feeBps;

    mapping(address => mapping(address => address)) public getPool;
    address[] public allPools;

    error InvalidAddress();
    error PoolExists();
    error RouterAlreadySet();
    error InvalidFee();

    event RouterSet(address indexed router);
    event PoolCreated(address indexed token0, address indexed token1, address pool, uint256 feeBps);

    constructor(address initialOwner, uint256 _feeBps) Ownable(initialOwner) {
        if (initialOwner == address(0)) revert InvalidAddress();
        if (_feeBps > 100) revert InvalidFee();
        feeBps = _feeBps;
    }

    function setRouter(address _router) external onlyOwner {
        if (_router == address(0)) revert InvalidAddress();
        if (router != address(0)) revert RouterAlreadySet();
        router = _router;
        emit RouterSet(_router);
    }

    function createPool(address tokenA, address tokenB) external onlyOwner returns (address pool) {
        if (router == address(0) || tokenA == address(0) || tokenB == address(0) || tokenA == tokenB) revert InvalidAddress();
        (address a, address b) = tokenA < tokenB ? (tokenA, tokenB) : (tokenB, tokenA);
        if (getPool[a][b] != address(0)) revert PoolExists();

        bytes32 salt = keccak256(abi.encodePacked(a, b));
        pool = address(new Pool{salt: salt}(a, b, router, feeBps, owner()));
        getPool[a][b] = pool;
        getPool[b][a] = pool;
        allPools.push(pool);
        emit PoolCreated(a, b, pool, feeBps);
    }

    function poolsLength() external view returns (uint256) { return allPools.length; }
}
