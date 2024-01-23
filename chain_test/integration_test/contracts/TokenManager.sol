// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./DToken.sol";

contract TokenManager {

    event NewContractCreated(address indexed a);
    event TakeERC20Result(bool);
    event Something(uint256);

    function genToken() external {
        DToken d = new DToken("dtoken","dt",10000000000000000000000000,msg.sender);
        emit NewContractCreated(address(d));
    }

    function transferHT(address payable to) external payable {
        to.transfer(msg.value);
    }

    function transfer(DToken d, address to,uint256 value) external {
        bool  b = d.transfer(to,value);
        require(b,"transfer failed");
    }

    function takeERC20(DToken d,address from,uint256 value) external {
        require(d.allowance(from,address(this)) >= value, "not enough allowance");
        bool ok = d.transferFrom(from,address(this),value);
        emit TakeERC20Result(ok);
    }

    function staticCallOtherContract(DToken d) external returns (uint256) {
        (bool success, bytes memory returnData) = address(d).staticcall(abi.encodeWithSignature("balanceOf(address)",address(this)));
        require(success);
        uint256 bal = abi.decode(returnData,(uint256));
        emit Something(bal);    // just to modify some thing
        return bal;
    }
}
