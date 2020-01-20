pragma solidity 0.5.12;

contract Exchange {
    function ethToTokenTransferInput(uint256 tokensBought, uint256 deadline, address recipient) public payable returns (uint256);
    function tokenToEthTransferInput(uint256 tokensSold, uint256 minEth, uint256 deadline, address recipient) public returns (uint256);
    function tokenToTokenTransferInput(uint256 tokensSold, uint256 minTokensBought, uint256 minEthBought, uint256 deadline, address recipient, address tokenAddress) public returns (uint256);
    function approve(address spender, uint256 value) public returns (bool);
}

contract Factory {
    function getExchange(address tokenAddress) public view returns (address);
}

contract Token {
    function balanceOf(address addr) public view returns (uint256);
    function transfer(address recipient, uint amount) public;
    function approve(address spender, uint tokens) public returns (bool);
    function allowance(address tokenOwner, address spender) public view returns (uint);
}

contract IncognitoUniSwap {

    address _factory;
    address _owner;

    // fallback function which allows transfer eth.
    function() external payable {}

    constructor(address factoryAddress) public {
        _owner = msg.sender;
        _factory = factoryAddress;
    }

    modifier isOwner {
        require(msg.sender == _owner, "unauthorized");
        _;
    }

    // verifyProof verifies burnt proof from txid
    function verifyProof(string memory txId) internal returns (bool) {
        return true;
    }

    /**
     * exchange verifies burnt proof, and start swap ETH/token to ETH/token
     * @param txId is temporary used until proof's type is specified.
     * @param fromToken token that will be sold, if token=0x1 then it is ETH
     * @param toToken token that will be bought and be transferred back to the smart contract.
     * After swapping success, amount and type will be sent to `Porting` contract.
     */
    function exchange(string memory txId, address fromToken, address toToken, uint256 swapAmount) public {
        require(fromToken != toToken);

        // verify proof from txId
        // TODO: currently set txId, change to proof type later.
        require(verifyProof(txId));

        uint256 amount = 0;
        if (fromToken == address(0x1)) { // eth to token
            amount = eth2Token(toToken, swapAmount);
        } else if (toToken == address(0x1)) { // token to eth
            amount = token2Eth(fromToken, swapAmount);
        } else { // token to token
            amount = token2Token(fromToken, toToken, swapAmount);
        }
        // TODO: call Porting smart contract function to mint privacy token.
    }

    /**
     * transferToken transfers token to recipient
     * @param tokenAddress token's address is used to transfer token.
     * @param recipient is address that receives token.
     * @param amount number of token which is transferred to recipient.
     */
    function transferToken(address tokenAddress, address recipient, uint amount) public isOwner {
        uint256 balance = Token(tokenAddress).balanceOf(address(this));
        require(balance >= amount && amount > 0);
        Token(tokenAddress).transfer(recipient, amount);
    }

    /**
     * transferEth transfers eth to recipient
     * @param recipient is address that receives eth
     * @param amount number of eth which is transferred to recipient.
     */
    function transferEth(address payable recipient, uint256 amount) public isOwner {
        require(address(this).balance >= amount, "insufficience funds");
        recipient.transfer(amount);
    }

    /**
     * tokenBalance gets token's balance of this smart contract
     *  @param tokenAddress is used to get exchange address.
     */
    function tokenBalance(address tokenAddress) public view returns (uint256) {
        return Token(tokenAddress).balanceOf(address(this));
    }

    /**
     * getExchangeAddress returns token's exchange address.
     * @param tokenAddress is address of specific token. eg: DAI in rinkeby: 0x2448eE2641d78CC42D7AD76498917359D961A783
     */
    function getExchangeAddress(address tokenAddress) public view returns (address) {
        address exchangeAddress = Factory(_factory).getExchange(tokenAddress);
        require(exchangeAddress != address(0x0), "exchange not found");
        return exchangeAddress;
    }

    /**
     * approve gives access to exchangeAddress to transfer a particular token's amount of current smart contract.
     * @param tokenAddress address of token that is used for approving process.
     * @param exchangeAddress is the address that will be given authority to change balance of the smart contract.
     * @param amount allowed amount. The amount will be added into token.allowances field.
     */
    function approve(address tokenAddress, address exchangeAddress, uint256 amount) internal {
        Exchange(exchangeAddress).approve(exchangeAddress, amount);
        Token(tokenAddress).approve(exchangeAddress, amount);
    }

   /**
    * eth2Token swaps eth to an erc20 token.
    * @param tokenAddress target token.
    * @param ethAmount amount that is spent to swap token.
    */
    function eth2Token(address tokenAddress, uint256 ethAmount) internal returns (uint256) {
        // eth must be greater than 0.
        require(address(this).balance >= ethAmount, "eth must be greater than 0");
        // get exchange address
        address exchangeAddress = getExchangeAddress(tokenAddress);
        // swap
        uint256 amount = Exchange(exchangeAddress).ethToTokenTransferInput.value(ethAmount)(1, block.timestamp, address(this));
        require(amount > 0);

        // approve gives exchange permission to change token's amount of this smart contract.
        approve(tokenAddress, exchangeAddress, amount);
        return amount;
    }

    /**
     * token2Eth swaps token to eth.
     * @param tokenAddress address of token that will be used to swap into eth.
     * @param tokenAmount swapped token's amount.
     */
    function token2Eth(address tokenAddress, uint256 tokenAmount) internal returns (uint256) {
        // get exchange address
        address exchangeAddress = getExchangeAddress(tokenAddress);
        // call swap function
        uint256 amount = Exchange(exchangeAddress).tokenToEthTransferInput(tokenAmount, 1, block.timestamp, address(this));
        require(amount > 0);
        return amount;
    }

    /**
     * token2Token swaps token to another token.
     * @param soldTokenAddress token address that will be used to swap into another.
     * @param boughtTokenAddress token address that smart contract received token after swapped.
     * @param tokenSold token's amount that smart contract spend to exchange token from boughtTokenAddress.
     */
    function token2Token(address soldTokenAddress, address boughtTokenAddress, uint256 tokenSold) internal returns (uint256) {
        uint256 minTokenBought = 1;
        uint256 minEthBought = 1;
        address soldExchange = Factory(_factory).getExchange(soldTokenAddress);
        uint256 amount = Exchange(soldExchange).tokenToTokenTransferInput(tokenSold, minTokenBought, minEthBought, block.timestamp, address(this), boughtTokenAddress);
        address boughtExchange = Factory(_factory).getExchange(boughtTokenAddress);
        approve(boughtTokenAddress, boughtExchange, amount);
        return amount;
    }
}

