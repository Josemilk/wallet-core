// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Factory} from "./Factory.sol";
import {Pool} from "./Pool.sol";

contract Router is ReentrancyGuard {
    using SafeERC20 for IERC20;

    Factory public immutable factory;
    error InvalidAddress();
    error DeadlineExpired();
    error PoolNotFound();

    constructor(address factoryAddress) {
        if (factoryAddress == address(0)) revert InvalidAddress();
        factory = Factory(factoryAddress);
    }

    function getAmountOut(address tokenIn, address tokenOut, uint256 amountIn) external view returns (uint256) {
        address pool = factory.getPool(tokenIn, tokenOut);
        if (pool == address(0)) revert PoolNotFound();
        return Pool(pool).quote(amountIn, tokenIn);
    }

    function swapExactTokensForTokens(
        address tokenIn,
        address tokenOut,
        uint256 amountIn,
        uint256 minAmountOut,
        address recipient,
        uint256 deadline
    ) external nonReentrant returns (uint256 amountOut) {
        if (block.timestamp > deadline) revert DeadlineExpired();
        if (recipient == address(0) || amountIn == 0) revert InvalidAddress();
        address pool = factory.getPool(tokenIn, tokenOut);
        if (pool == address(0)) revert PoolNotFound();

        IERC20(tokenIn).safeTransferFrom(msg.sender, pool, amountIn);
        amountOut = Pool(pool).swap(tokenIn, amountIn, minAmountOut, recipient);
    }
}
