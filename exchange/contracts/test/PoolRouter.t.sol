// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import {Test} from "forge-std/Test.sol";
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {Factory} from "../src/Factory.sol";
import {Router} from "../src/Router.sol";
import {Pool} from "../src/Pool.sol";

contract TestToken is ERC20 {
    constructor(string memory n, string memory s) ERC20(n, s) {}
    function mint(address to, uint256 amount) external { _mint(to, amount); }
}

contract PoolRouterTest is Test {
    TestToken tokenA;
    TestToken tokenB;
    Factory factory;
    Router router;
    Pool pool;

    address owner = address(0xA11CE);
    address lp = address(0xB0B);
    address trader = address(0xC0DE);

    function setUp() public {
        tokenA = new TestToken("Token A", "A");
        tokenB = new TestToken("Token B", "B");
        factory = new Factory(owner, 30);
        router = new Router(address(factory));
        vm.prank(owner);
        factory.setRouter(address(router));
        vm.prank(owner);
        address poolAddress = factory.createPool(address(tokenA), address(tokenB));
        pool = Pool(poolAddress);

        tokenA.mint(lp, 100 ether);
        tokenB.mint(lp, 100 ether);
        tokenA.mint(trader, 10 ether);
        tokenB.mint(trader, 10 ether);
    }

    function testLiquidityAndSwap() public {
        vm.startPrank(lp);
        tokenA.approve(address(pool), type(uint256).max);
        tokenB.approve(address(pool), type(uint256).max);
        pool.addLiquidity(50 ether, 50 ether, lp);
        vm.stopPrank();

        vm.startPrank(trader);
        tokenA.approve(address(router), type(uint256).max);
        uint256 quoted = router.getAmountOut(address(tokenA), address(tokenB), 1 ether);
        uint256 before = tokenB.balanceOf(trader);
        router.swapExactTokensForTokens(
            address(tokenA), address(tokenB), 1 ether, quoted * 99 / 100, trader, block.timestamp + 1 hours
        );
        uint256 received = tokenB.balanceOf(trader) - before;
        vm.stopPrank();

        assertGt(received, 0);
        assertGe(received, quoted * 99 / 100);
    }

    function testOnlyFactoryOwnerCreatesPools() public {
        vm.prank(trader);
        vm.expectRevert();
        factory.createPool(address(0x1), address(0x2));
    }
}
