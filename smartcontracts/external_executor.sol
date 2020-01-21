pragma solidity >=0.5.0 <0.7.0;

contract ExternalExecutor {

    /**
     * executeFunction is used to execute function provided in bytes code from a given address.
     * @param target is target's address that is used to call function.
     * @param code is encoded from function and its params.
     * @returns result of call process.
     * @notice this is a payable function in case users want to trigger a payable function in target.
     */
    function executeFunction(address target, bytes memory code) public payable returns (bytes memory) {
        (bool success, bytes memory result) = executor.call.value(msg.value)(code);
        require(success);
        return result;
    }
}
