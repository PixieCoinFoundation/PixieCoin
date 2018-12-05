/// @title Contracts for PixieGame. Holds all related game structs, events and base variables.
/// @author Pixie Foundation (http://pixiecoin.io/)

pragma solidity ^0.4.19;

contract Ownable {
  address public owner;

  event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);

  function Ownable() public {
    owner = msg.sender;
  }


  modifier onlyOwner() {
    require(msg.sender == owner);
    _;
  }


  function transferOwnership(address newOwner) public onlyOwner {
    require(newOwner != address(0));
    OwnershipTransferred(owner, newOwner);
    owner = newOwner;
  }
}

contract ERC721 {
    // Required methods
    function totalSupply() public view returns (uint256 total);
    function balanceOf(address _owner) public view returns (uint256 balance);
    function ownerOf(uint256 _tokenId) external view returns (address owner);
    function approve(address _to, uint256 _tokenId) external;
    function transfer(address _to, uint256 _tokenId) external;
    function transferFrom(address _from, address _to, uint256 _tokenId) external;

    // Events
    event Transfer(address from, address to, uint256 tokenId);
    event Approval(address owner, address approved, uint256 tokenId);

    // Optional
    //function name() public view returns (string name);
    //function symbol() public view returns (string symbol);
    function tokensOfOwner(address _owner) external view returns (uint256[] tokenIds);
    //function tokenMetadata(uint256 _tokenId, string _preferredTransport) public view returns (string infoUrl);
    function takeOwnership(uint256 _tokenId) public;
}

contract PixieBase {
  /// @dev create new design work event
  event NewDesignWork(uint64 workId,uint64 registerTime);
  
  /// @dev transfer design work event
  event TransferDesignWork(address indexed from,address indexed to,uint indexed designWorkId);
  
  /// @dev product clothes from design work event
  event ProductDesignWork(uint indexed designWorkId);

  /// @dev data structure of design work
  struct DesignWork {
    uint64 workId;
    uint64 registerTime;
  }

  /// @dev global design work data
  DesignWork[] designWorks;
  mapping (uint => address) public designWorkToOwner;
  mapping (address => uint) ownerDesignWorkCount;
  mapping (uint256 => address) public designWorkIndexToApproved;
}

contract PixieAccessControl is PixieBase,Ownable {
  /// @dev product clothes from design work cool down time
  uint64 productCooldownTime = 30 minutes;

  /// @dev upload design work address.only owner can change.
  address uploadDesignWorkAddress;

  function changeUploadAddress(address _uploadAddress) external onlyOwner {
    uploadDesignWorkAddress = _uploadAddress;
  }

  function changeCoolDownTime(uint64 _productCoolDownTime) external onlyOwner {
    productCooldownTime = _productCoolDownTime;
  }

  function showUploadAddress() external view returns (address) {
    return uploadDesignWorkAddress;
  }

  function showProductCoolDownTime() external view returns (uint) {
    return productCooldownTime;
  }
}

contract PixieInternal is PixieAccessControl,ERC721 {
    string public constant pname = "PixieGame";
    string public constant psymbol = "PG";
  
    function _owns(address _claimant, uint256 _tokenId) internal view returns (bool) {
        return designWorkToOwner[_tokenId] == _claimant;
    }

    function _transfer(address _from, address _to, uint256 _tokenId) internal {
        ownerDesignWorkCount[_to]++;
        designWorkToOwner[_tokenId] = _to;

        if (_from != address(0)) {
            ownerDesignWorkCount[_from]--;
            delete designWorkIndexToApproved[_tokenId];
        }

        Transfer(_from, _to, _tokenId);
    }

    function _approve(uint256 _tokenId, address _approved) internal {
        designWorkIndexToApproved[_tokenId] = _approved;
    }

    function _approvedFor(address _claimant, uint256 _tokenId) internal view returns (bool) {
        return designWorkIndexToApproved[_tokenId] == _claimant;
    }

    function _createDesignWork(uint64 _workId,uint64 _registerTime) internal {
      designWorks.push(DesignWork(_workId,_registerTime));
    
      NewDesignWork(_workId,_registerTime);
    }
}

contract PixieGame is PixieInternal {
  function PixieGame() public {
    uploadDesignWorkAddress = msg.sender;
  }

  function createDesignWork(uint64 _workId,uint64 _registerTime) external {
    require(msg.sender == uploadDesignWorkAddress);

    _createDesignWork(_workId,_registerTime);
  }

  function batchCreateDesignWork(uint64[] _workIds,uint64[] _registerTimes) external {
    require(msg.sender == uploadDesignWorkAddress);
    require(_workIds.length == _registerTimes.length);

    uint256 i= 0;
    while (i < _workIds.length) {
      _createDesignWork(_workIds[i],_registerTimes[i]);
      i++;
    }
  }

  function balanceOf(address _owner) public view returns (uint256 count) {
      return ownerDesignWorkCount[_owner];
  }

  function transfer(
      address _to,
      uint256 _tokenId
  )
      external
  {
      require(_to != address(0));
      require(_to != address(this));
      require(_owns(msg.sender, _tokenId));

      _transfer(msg.sender, _to, _tokenId);
  }

  function approve(
      address _to,
      uint256 _tokenId
  )
      external
  {
      require(_owns(msg.sender, _tokenId));

      _approve(_tokenId, _to);

      Approval(msg.sender, _to, _tokenId);
  }

  function transferFrom(
      address _from,
      address _to,
      uint256 _tokenId
  )
      external
  {
      require(_to != address(0));
      require(_to != address(this));
      require(_approvedFor(msg.sender, _tokenId));
      require(_owns(_from, _tokenId));

      _transfer(_from, _to, _tokenId);
  }

  /// @dev no limit on total supply.just return current size of design work list.
  function totalSupply() public view returns (uint) {
      return designWorks.length;
  }

  function ownerOf(uint256 _tokenId)
      external
      view
      returns (address owner)
  {
      owner = designWorkToOwner[_tokenId];

      require(owner != address(0));
  }

  function tokensOfOwner(address _owner) external view returns(uint256[] ownerTokens) {
      uint256 tokenCount = balanceOf(_owner);

      if (tokenCount == 0) {
          return new uint256[](0);
      } else {
          uint256[] memory result = new uint256[](tokenCount);
          uint256 totalDesignWorks = totalSupply();
          uint256 resultIndex = 0;

          uint256 designWorkId;

          for (designWorkId = 1; designWorkId <= totalDesignWorks; designWorkId++) {
            if (designWorkToOwner[designWorkId] == _owner) {
                  result[resultIndex] = designWorkId;
                  resultIndex++;
              }
          }

          return result;
      }
  }

  function takeOwnership(uint256 _tokenId) public {
    require(designWorkIndexToApproved[_tokenId] == msg.sender);
    address owner = designWorkToOwner[_tokenId];
    _transfer(owner, msg.sender, _tokenId);
  }
}