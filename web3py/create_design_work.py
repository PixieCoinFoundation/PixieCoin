import sys
import time
import pprint
import my_constant
import getopt

from web3.providers import HTTPProvider
from web3 import Web3

def wait_for_receipt(w3, tx_hash, poll_interval):
   while True:
       tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
       if tx_receipt:
         return tx_receipt
       time.sleep(poll_interval)

if __name__ == '__main__':
  opts, args = getopt.getopt(sys.argv[1:], 'hn:d:k:i:p:c:w', ['help','contract_address=','name=', 'desc=', 'nickname=','icon=','http_provider=','wait'])

  # default http provider
  http_provider = 'http://localhost:8545'
  name,desc,nickname,icon,contract_address_str = '','','','',''
  wait = False

  for key, value in opts:
    # print('parameter',key,value)
    if key in ['-h', '--help']:
      print('upload design work info to ethereum')
      print('parameter:')
      print('-h\thelp')
      print('-n\tdesign work name')
      print('-d\tdesign work description')
      print('-nn\tdesigner nickname')
      print('-i\tdesign work icon info')
      print('-hp\thttp provider address')
      sys.exit(0)
    if key in ['-n', '--name']:
      name = value
    if key in ['-d', '--desc']:
      desc = value
    if key in ['-k','--nickname']:
      nickname = value
    if key in ['-i','--icon']:
      icon = value
    if key in ['-p','--http_provider']:
      http_provider = value
    if key in ['-c','--contract_address']:
      contract_address_str = value
    if key in ['-w','--wait']:
      wait = True
  
  if name=='' or nickname=='' or icon=='' or contract_address_str=='':
    print('empty parameter in name or nickname or icon or contract address')
    sys.exit(-1)

  w3 = Web3(HTTPProvider(http_provider))
  connected = w3.isConnected()

  if not connected:
    print('not connected to',http_provider)
    sys.exit(-1)

  contract_address = w3.toChecksumAddress(contract_address_str)

  ca = w3.eth.contract(
     address=contract_address,
     abi=my_constant.CONTRACT_ABI)

  coin_base = w3.eth.accounts[0]
  print('eth.accounts[0]',coin_base)

  upload_address = ca.functions.showUploadAddress().call(get_block_id=False)
  print('upload address',upload_address)

  job = ca.functions.createDesignWork(name,desc,nickname,icon)

  gas_estimate = job.estimateGas()
  print("Gas estimate:",gas_estimate)

  tx_hash = job.transact({'gas':gas_estimate,'from':upload_address})
  print("transaction hash:",tx_hash.hex())
  
  if wait:
    receipt = wait_for_receipt(w3, tx_hash, 1)
    print("Transaction receipt mined!",'\n',receipt)

    rp = dict(receipt)
    print("status:",rp['status'],"\n")