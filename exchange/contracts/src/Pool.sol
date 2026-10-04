// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {ReentrancyGuard} from "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import {Pausable} from "@openzeppelin/contracts/utils/Pausable.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

contract Pool is ERC20, ReentrancyGuard, Pausable, Ownable {
    using SafeERC20 for IERC20;
    uint256 public constant MINIMUM_LIQUIDITY = 1_000;
    uint256 public constant BPS = 10_000;
    uint256 public immutable swapFeeBps;
    address public immutable token0;
    address public immutable token1;
    address public immutable router;
    uint112 private reserve0;
    uint112 private reserve1;

    error InvalidToken();
    error InvalidAmount();
    error InvalidRecipient();
    error InvalidOutput();
    error InsufficientLiquidity();
    error OnlyRouter();
    error InvalidFee();

    event Swap(address indexed sender, address indexed recipient, address indexed tokenIn, uint256 amountIn, uint256 amountOut);
    event LiquidityAdded(address indexed provider, uint256 amount0, uint256 amount1, uint256 shares);
    event LiquidityRemoved(address indexed provider, uint256 amount0, uint256 amount1, uint256 shares);
    event ReservesSynced(uint112 reserve0, uint112 reserve1);

    constructor(address _token0, address _token1, address _router, uint256 _swapFeeBps, address _owner)
        ERC20("Exchange LP", "XLP") Ownable(_owner)
    {
        if (_token0 == address(0) || _token1 == address(0) || _token0 == _token1 || _router == address(0) || _owner == address(0)) revert InvalidToken();
        if (_swapFeeBps > 100) revert InvalidFee();
        token0 = _token0; token1 = _token1; router = _router; swapFeeBps = _swapFeeBps;
    }

    modifier onlyRouter() { if (msg.sender != router) revert OnlyRouter(); _; }

    function getReserves() external view returns (uint112, uint112) { return (reserve0, reserve1); }

    function quote(uint256 amountIn, address tokenIn) public view returns (uint256 amountOut) {
        if (amountIn == 0) revert InvalidAmount();
        bool zeroForOne = tokenIn == token0;
        if (!zeroForOne && tokenIn != token1) revert InvalidToken();
        uint256 reserveIn = zeroForOne ? reserve0 : reserve1;
        uint256 reserveOut = zeroForOne ? reserve1 : reserve0;
        if (reserveIn == 0 || reserveOut == 0) revert InsufficientLiquidity();
        uint256 net = amountIn * (BPS - swapFeeBps);
        return (net * reserveOut) / (reserveIn * BPS + net);
    }

    function addLiquidity(uint256 amount0, uint256 amount1, address recipient)
        external nonReentrant whenNotPaused returns (uint256 shares)
    {
        if (amount0 == 0 || amount1 == 0 || recipient == address(0)) revert InvalidAmount();
        IERC20(token0).safeTransferFrom(msg.sender, address(this), amount0);
        IERC20(token1).safeTransferFrom(msg.sender, address(this), amount1);
        uint256 supply = totalSupply();
        if (supply == 0) {
            shares = _sqrt(amount0 * amount1);
            if (shares <= MINIMUM_LIQUIDITY) revert InvalidAmount();
            _mint(address(0x000000000000000000000000000000000000dEaD), MINIMUM_LIQUIDITY);
            shares -= MINIMUM_LIQUIDITY;
        } else {
            uint256 s0 = amount0 * supply / reserve0;
            uint256 s1 = amount1 * supply / reserve1;
            shares = s0 < s1 ? s0 : s1;
        }
        if (shares == 0) revert InvalidAmount();
        _mint(recipient, shares);
        _updateReserves(uint256(reserve0) + amount0, uint256(reserve1) + amount1);
        emit LiquidityAdded(recipient, amount0, amount1, shares);
    }

    function removeLiquidity(uint256 shares, address recipient)
        external nonReentrant whenNotPaused returns (uint256 amount0, uint256 amount1)
    {
        if (shares == 0 || balanceOf(msg.sender) < shares || recipient == address(0)) revert InvalidAmount();
        uint256 supply = totalSupply();
        amount0 = uint256(reserve0) * shares / supply;
        amount1 = uint256(reserve1) * shares / supply;
        if (amount0 == 0 || amount1 == 0) revert InvalidAmount();
        _burn(msg.sender, shares);
        _updateReserves(uint256(reserve0) - amount0, uint256(reserve1) - amount1);
        IERC20(token0).safeTransfer(recipient, amount0);
        IERC20(token1).safeTransfer(recipient, amount1);
        emit LiquidityRemoved(msg.sender, amount0, amount1, shares);
    }

    function swap(address tokenIn, uint256 amountIn, uint256 minAmountOut, address recipient)
        external nonReentrant whenNotPaused onlyRouter returns (uint256 amountOut)
    {
        if (amountIn == 0 || recipient == address(0)) revert InvalidAmount();
        bool zeroForOne = tokenIn == token0;
        if (!zeroForOne && tokenIn != token1) revert InvalidToken();
        uint256 reserveOut = zeroForOne ? reserve1 : reserve0;
        amountOut = quote(amountIn, tokenIn);
        if (amountOut < minAmountOut || amountOut >= reserveOut) revert InvalidOutput();
        address tokenOut = zeroForOne ? token1 : token0;
        IERC20(tokenOut).safeTransfer(recipient, amountOut);
        if (zeroForOne) _updateReserves(uint256(reserve0) + amountIn, uint256(reserve1) - amountOut);
        else _updateReserves(uint256(reserve0) - amountOut, uint256(reserve1) + amountIn);
        emit Swap(msg.sender, recipient, tokenIn, amountIn, amountOut);
    }

    function pause() external onlyOwner { _pause(); }
    function unpause() external onlyOwner { _unpause(); }
    function sync() external onlyOwner { _updateReserves(IERC20(token0).balanceOf(address(this)), IERC20(token1).balanceOf(address(this))); }

    function _updateReserves(uint256 r0, uint256 r1) private {
        if (r0 > type(uint112).max || r1 > type(uint112).max) revert InvalidAmount();
        reserve0 = uint112(r0); reserve1 = uint112(r1);
        emit ReservesSynced(reserve0, reserve1);
    }

    function _sqrt(uint256 y) private pure returns (uint256 z) {
        if (y > 3) { z = y; uint256 x = y / 2 + 1; while (x < z) { z = x; x = (y / x + x) / 2; } }
        else if (y != 0) z = 1;
    }
}
